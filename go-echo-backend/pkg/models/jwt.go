package models

import (
	"github.com/golang-jwt/jwt"
)

// JwtClaims custom claims
type JwtClaims struct {
	ID string `json:"id"`

	jwt.StandardClaims
}
