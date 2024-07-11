package models

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
	"github.com/thaitanloi365/go-utils"
)

// Log keys for headers
const (
	RefErrorIDKey = "__RefErrorIDKey"
	UserIDKey     = "__UserIDKey"
	LanguageKey   = "Language"
	TimeZoneKey   = "Timezone"
)

func (c *CustomContext) GetRequestID() string {
	var reqID = c.Request().Header.Get(logger.HeaderRequestIDKey)
	if reqID == "" {
		reqID = c.Response().Header().Get(logger.HeaderRequestIDKey)
	}

	return reqID
}

func (c *CustomContext) GetBearerToken() string {
	var bearerToken = c.Request().Header.Get(echo.HeaderAuthorization)
	bearerToken = strings.TrimPrefix(bearerToken, "Bearer ")

	return bearerToken
}

func (c *CustomContext) SetRefErrorID(id string) {
	c.Set(RefErrorIDKey, id)
}

func (c *CustomContext) GetRefErrorID() string {
	if id, ok := c.Get(RefErrorIDKey).(string); ok {
		return id
	}

	return ""
}

func (c *CustomContext) GetLoggedUser() string {
	if id, ok := c.Get(RefErrorIDKey).(string); ok {
		return id
	}

	return ""
}

func (c *CustomContext) GetRequestInfo() *logger.RequestInfo {
	var user User
	c.GetUserFromContext(&user)

	var start = time.Now()
	var req = c.Request()
	var res = c.Response()

	var reqID = req.Header.Get(logger.HeaderRequestIDKey)
	if reqID == "" {
		reqID = res.Header().Get(logger.HeaderRequestIDKey)
	}

	var refErrorID = c.GetRefErrorID()
	var info = &logger.RequestInfo{
		ReqID:      reqID,
		UserID:     "",
		StatusCode: res.Status,
		Size:       res.Size,
		Method:     req.Method,
		Host:       req.Host,
		URI:        req.RequestURI,
		Latency:    time.Since(start),
		RemoteIP:   c.RealIP(),
		UserAgent:  req.UserAgent(),
		RefErrorID: refErrorID,
	}

	if req.URL != nil {
		info.URL = req.URL.String()
	}

	if uri, err := url.PathUnescape(req.RequestURI); err == nil {
		info.URI = uri
	}

	if user.ID != "" {
		info.UserID = user.ID
	}

	if len(req.Header) > 0 {
		info.Header = helper.FormatHeaders(req.Header)
	}
	return info
}

func (c *CustomContext) GetHttpRequestInfo(req *http.Request, res *echo.Response) *logger.RequestInfo {

	var reqID = req.Header.Get(logger.HeaderRequestIDKey)
	if reqID == "" {
		reqID = res.Header().Get(logger.HeaderRequestIDKey)
	}

	var refErrorID = c.GetRefErrorID()
	var info = &logger.RequestInfo{
		ReqID:      reqID,
		StatusCode: res.Status,
		Size:       res.Size,
		Method:     req.Method,
		Host:       req.Host,
		URI:        req.RequestURI,
		RemoteIP:   c.RealIP(),
		UserAgent:  req.UserAgent(),
		RefErrorID: refErrorID,
		UserID:     req.Header.Get(logger.UserIDKey),
	}

	if req.URL != nil {
		info.URL = req.URL.String()
	}

	if uri, err := url.PathUnescape(req.RequestURI); err == nil {
		info.URI = uri
	}

	if len(req.Header) > 0 {
		info.Header = helper.FormatHeaders(req.Header)
	}

	if claims, err := c.GetJwtClaimsInfo(); err == nil {
		info.UserID = claims.GetUserID()
		info.UserRole = string(claims.GetRole())
		info.GhostID = string(claims.GetGhostID())
	}
	return info
}

func (c *CustomContext) GetUserFromContext(in interface{}) error {
	if !utils.IsPtr(in) {
		return errors.New("input must be a pointer")
	}

	var rawUser = c.Get("user")
	if rawUser == nil {
		return errors.New("user is not found in context")
	}

	reflect.ValueOf(in).Elem().Set(reflect.ValueOf(rawUser).Elem())
	return nil
}

func (c *CustomContext) IsGhostLogin() bool {
	var value = c.Get("is_ghost")
	if value == nil {
		return false
	}

	if isAdmin, ok := value.(bool); ok {
		return isAdmin
	}
	return false
}

func (c *CustomContext) IsGhostLoggedIn() bool {
	var value = c.Get("is_ghost_logged_in")
	if value == nil {
		return false
	}

	if isAdminLoggedIn, ok := value.(bool); ok {
		return isAdminLoggedIn
	}
	return false
}

func (c CustomContext) GetJwtClaims() (JwtClaims, error) {
	user, ok := c.Get("user_token").(*jwt.Token)
	if !ok {
		return JwtClaims{}, eris.New("User not authenticated")
	}
	var claims = user.Claims.(*JwtClaims)
	var jwtClaims JwtClaims

	var err = claims.Valid()
	if err != nil {
		if e, ok := err.(*echo.HTTPError); ok {

			if e.Code == http.StatusBadRequest {
				return jwtClaims, errs.ErrTokenMissing
			}

			if e.Code == http.StatusUnauthorized {
				return jwtClaims, errs.ErrTokenInvalid
			}
		}

		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return jwtClaims, errs.ErrTokenInvalid
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return jwtClaims, errs.ErrTokenExpired
			} else {
				return jwtClaims, errs.ErrTokenInvalid
			}
		}

		return jwtClaims, errs.ErrTokenInvalid
	}

	jwtClaims.ID = claims.ID
	jwtClaims.Audience = claims.Audience
	jwtClaims.Issuer = claims.Issuer
	jwtClaims.Subject = claims.Subject

	return jwtClaims, nil
}

func (c CustomContext) GetJwtClaimsInfo() (JwtClaimsInfo, error) {
	jwtClaims, err := c.GetJwtClaims()

	var info = JwtClaimsInfo{
		role:   enums.Role(jwtClaims.Audience),
		userID: jwtClaims.ID,
	}

	return info, err
}

func (c CustomContext) Success(i interface{}) error {
	var code = http.StatusOK

	if v, ok := i.(string); ok {
		return c.JSON(code, M{
			"message": v,
		})
	}

	var rt = reflect.TypeOf(i)
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		return c.JSON(code, M{
			"records": i,
		})
	}

	return c.JSON(code, i)
}

func (c CustomContext) BindAndValidate(i interface{}) error {
	var err = (&echo.DefaultBinder{}).BindPathParams(c, i)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = c.Bind(i)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = c.Validate(i)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return nil
}

func (c CustomContext) GetPathParamUint(tag string, fallbackValue ...uint) uint {
	v, err := strconv.ParseUint(c.Param(tag), 10, 32)
	if err != nil && len(fallbackValue) > 0 {
		return fallbackValue[0]
	}

	return uint(v)
}

func (c CustomContext) GetPathParamString(tag string, fallbackValue ...string) string {
	v := c.Param(tag)
	if v == "" && len(fallbackValue) > 0 {
		return fallbackValue[0]
	}
	return v

}

func (c CustomContext) GetQueryParamString(tag string, fallbackValue ...string) string {
	v := c.QueryParam(tag)
	if v == "" && len(fallbackValue) > 0 {
		return fallbackValue[0]
	}
	return v
}

func (c CustomContext) GetQueryParamInt(tag string, fallbackValue ...int) int {
	v, err := strconv.Atoi(c.QueryParam(tag))
	if err != nil && len(fallbackValue) > 0 {
		return fallbackValue[0]
	}
	return v
}

func (c CustomContext) GetQueryParamBool(tag string, fallbackValue ...bool) bool {
	v, err := strconv.ParseBool(c.QueryParam(tag))
	if err != nil && len(fallbackValue) > 0 {
		return fallbackValue[0]
	}
	return v
}

func (c CustomContext) GetQueryParamBoolPtr(tag string) *bool {
	v, err := strconv.ParseBool(c.QueryParam(tag))
	if err != nil {
		return nil
	}
	return &v
}

func (c CustomContext) GetQueryParamSliceString(tag string, separator ...string) []string {
	var sep = ","
	if len(separator) > 0 {
		sep = separator[0]
	}
	if v := c.QueryParam(tag); v != "" {
		return strings.Split(c.QueryParam(tag), sep)
	}

	return []string{}
}

func (c CustomContext) GetQuerySortSliceString(tag string, defaultValue string, separator ...string) [][]string {
	var sep = ","
	if len(separator) > 0 {
		sep = separator[0]
	}

	var listStatus = [][]string{}
	var v = c.QueryParam(tag)
	if v == "" {
		v = defaultValue
	}

	var segments = strings.Split(v, sep)

	for _, s := range segments {
		var segs = strings.Split(s, ":")
		if len(segs) == 2 {
			listStatus = append(listStatus, segs)
		}
	}

	return listStatus
}

func (c CustomContext) GetTimeZone() enums.Timezone {
	timezoneCode := enums.Timezone(c.Request().Header.Get(TimeZoneKey))
	return timezoneCode.DefaultIfInvalid()
}

func (c CustomContext) GetLanguage() enums.LanguageCode {
	language := enums.LanguageCode(c.Request().Header.Get(LanguageKey))
	return language.DefaultIfInvalid()
}
