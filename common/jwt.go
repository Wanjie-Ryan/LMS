package common

import "github.com/golang-jwt/jwt/v5"

type CustomJWTClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}