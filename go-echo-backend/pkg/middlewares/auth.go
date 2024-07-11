package middlewares

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"

	"github.com/rotisserie/eris"

	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/rsa"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const contextKey = "user_token"

func IsBasicAuth() echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		var cc = c.(*models.CustomContext)

		var u = cc.App.Config.AdminUserName
		var p = cc.App.Config.AdminUserPassword

		if subtle.ConstantTimeCompare([]byte(username), []byte(u)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(p)) == 1 {
			return true, nil
		}

		return false, nil
	})

}

// IsAuthorized Authorization middleware
func (m *Middleware) IsAuthorized() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:       &models.JwtClaims{},
		ContextKey:   contextKey,
		SigningKey:   []byte(m.App.Config.JWTSecret),
		ErrorHandler: transformError,
	})
}

func (m *Middleware) IsValidSignature() echo.MiddlewareFunc {
	var rsaClient = rsa.New(m.App.Config)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var signature = c.Request().Header.Get("Signature")

			var err = rsaClient.Validate(signature)
			if err != nil {
				return eris.Wrapf(errs.ErrSignatureInvalid, "err=%+v", err)
			}

			return next(c)
		}
	}

}

func (m *Middleware) IsBasicAuth() echo.MiddlewareFunc {
	return IsBasicAuth()
}

func (m *Middleware) CheckTokenExpiredAndAttachUserInfo() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var cc = c.(*models.CustomContext)

			jwtClaims, err := cc.GetJwtClaims()
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			var user models.User
			err = cc.App.DB.First(&user, "id = ?", jwtClaims.ID).Error
			if err != nil {
				return errs.ErrTokenInvalid
			}

			var isValid = jwtClaims.VerifyIssuer(user.TokenIssuer, true)
			if !isValid {
				if parts := strings.Split(jwtClaims.Issuer, "|"); len(parts) == 2 {
					isValid = subtle.ConstantTimeCompare([]byte(parts[0]), []byte(("ghost"))) != 0

				}

				c.Set("is_ghost", isValid)
			}

			if !isValid {
				return errs.ErrTokenInvalid
			}

			c.Set("user", &user)

			return next(c)
		}
	}

}

func (m *Middleware) CheckAccountStatus() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var cc = c.(*models.CustomContext)

			var user models.User
			err := cc.GetUserFromContext(user)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			if user.AccountStatus != enums.AccountStatusActive {
				return errs.ErrAccountStatusNotActive
			}

			return next(c)
		}
	}

}

func (m *Middleware) CheckRole(roles ...enums.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var cc = c.(*models.CustomContext)
			var isValidRole bool = false

			jwtClaims, err := cc.GetJwtClaims()
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			for _, role := range roles {
				if role.String() == jwtClaims.Audience {
					isValidRole = true
					break
				}
			}

			if !isValidRole {
				return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("This route only valid for rules: %v", roles))
			}
			return next(c)
		}
	}

}

func (m *Middleware) RequireRoles(next echo.HandlerFunc, roles ...string) echo.HandlerFunc {
	return func(c echo.Context) error {
		var cc = c.(*models.CustomContext)
		var isValidRole bool = false

		jwtClaims, err := cc.GetJwtClaims()
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		for _, role := range roles {
			if role == jwtClaims.Audience {
				isValidRole = true
				break
			}
		}

		if !isValidRole {
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("This route only valid for rules: %v", roles))
		}

		return next(c)
	}
}

func (m *Middleware) IsGrant() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var cc = c.(*models.CustomContext)
			var jwtClaims models.JwtClaims
			var err error
			jwtClaims, err = cc.GetJwtClaims()
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden, err)
			}
			//if !cc.App.Config.IsProd() { // not apply casbin on prod
			//	role := jwtClaims.Subject    // sub
			//	path := c.Request().URL.Path // obj
			//	method := c.Request().Method // act
			//
			//	var isGrant bool
			//	isGrant, err = cc.App.Enforcer.Enforce(role, path, method)
			//	if err != nil {
			//		return eris.Wrap(err, "IsGrant Error")
			//	}
			//	if isGrant == false {
			//		return errs.ErrExceedPermission
			//	}
			//}
			var user models.User
			err = cc.App.DB.Select("ID", "TokenIssuer").First(&user, "id = ?", jwtClaims.ID).Error
			if user.ID == "" || err != nil {
				msg := fmt.Sprintf("IsGrant - Not found user - userId:%s", jwtClaims.ID)
				return eris.Wrap(errs.ErrUserNotFound, msg)
			}
			if user.TokenIssuer == "" {
				msg := fmt.Sprintf("IsGrant - Empty Token Issuer - userId:%s", jwtClaims.ID)
				return eris.Wrap(errs.ErrTokenInvalid, msg)
			}
			return next(c)
		}
	}
}

func IsAuthorizedWithQueryToken(secretKey string) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:       &models.JwtClaims{},
		ContextKey:   contextKey,
		SigningKey:   []byte(secretKey),
		ErrorHandler: transformError,
		TokenLookup:  "query:token",
		SuccessHandler: func(c echo.Context) {
			var claims models.JwtClaims
			if err := claims.ValidateToken(secretKey, c.QueryParam("token")); err == nil {
				c.Request().Header.Set(logger.UserIDKey, claims.ID)
			}
		},
	})
}

func transformError(err error) error {
	if err == middleware.ErrJWTMissing {
		return errs.ErrTokenMissing
	}

	if e, ok := err.(*echo.HTTPError); ok {
		if e.Code == http.StatusBadRequest {
			return errs.ErrTokenMissing
		}

		if e.Code == http.StatusUnauthorized {
			return errs.ErrTokenInvalid
		}
	}

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return errs.ErrTokenInvalid
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return errs.ErrTokenExpired
		} else {
			return errs.ErrTokenInvalid
		}
	}

	return err
}
