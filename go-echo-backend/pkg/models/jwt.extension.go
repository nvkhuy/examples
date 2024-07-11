package models

import (
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/golang-jwt/jwt"
)

// ValidateToken parse and validate
func (jwtClaims *JwtClaims) ValidateToken(secret, jwtToken string) error {
	// ValidateToken validate jwt token
	key := []byte(secret)
	token, err := jwt.ParseWithClaims(jwtToken, jwtClaims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil || !token.Valid {
		return errs.ErrTokenInvalid
	}

	return nil

}

func NewJwtClaimsInfo() *JwtClaimsInfo {
	return &JwtClaimsInfo{}
}

func (jwtClaims *JwtClaims) GetJwtClaimsInfo() JwtClaimsInfo {
	var parts = strings.Split(jwtClaims.Issuer, "|")
	var info = JwtClaimsInfo{
		userID: jwtClaims.ID,
		role:   enums.Role(jwtClaims.Audience),
	}

	if len(parts) == 2 {
		info.isGhost = strings.EqualFold(parts[0], "ghost")
		info.ghostID = parts[1]
	}

	return info
}

func (jwtClaims *JwtClaimsInfo) GetUserID(fallback ...string) string {
	var id = jwtClaims.userID
	if len(fallback) > 0 && id == "" {
		id = fallback[0]
	}

	return id
}

func (jwtClaims *JwtClaimsInfo) GetRole(fallback ...enums.Role) enums.Role {
	var id = jwtClaims.role
	if len(fallback) > 0 && id == "" {
		id = fallback[0]
	}

	return id
}

func (jwtClaims *JwtClaimsInfo) SetTimezone(tz enums.Timezone) *JwtClaimsInfo {
	jwtClaims.timezone = tz
	return jwtClaims
}

func (jwtClaims *JwtClaimsInfo) SetRole(role enums.Role) *JwtClaimsInfo {
	jwtClaims.role = role
	return jwtClaims
}

func (jwtClaims *JwtClaimsInfo) SetUserID(userID string) *JwtClaimsInfo {
	jwtClaims.userID = userID
	return jwtClaims
}

func (jwtClaims *JwtClaimsInfo) GetGhostID() string {
	return jwtClaims.ghostID
}

func (jwtClaims *JwtClaimsInfo) GetTimezone() enums.Timezone {
	return jwtClaims.timezone
}
