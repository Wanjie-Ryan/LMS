package common

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Wanjie-Ryan/LMS/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type CustomJWTClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.User) (*string, *string, error) {

	userClaims := CustomJWTClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	signedAccessToken, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return nil, nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomJWTClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})

	signedRefreshToken, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, nil, err
	}

	return &signedAccessToken, &signedRefreshToken, nil

}

func ParseJWT(signedAccessToken string) (*CustomJWTClaims, error) {

	parsedJwtAccessToken, err := jwt.ParseWithClaims(signedAccessToken, &CustomJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		fmt.Println("error parsing access token", err.Error())
		return nil, err
	} else if claims, ok := parsedJwtAccessToken.Claims.(*CustomJWTClaims); ok {

		return claims, nil
	} else {
		log.Default().Println("error parsing access token")
		return nil, errors.New("unknown claims Type, cannot proceed")
	}
}
func IsClaimExpired(claims *CustomJWTClaims) bool {
	currentTime := jwt.NewNumericDate(time.Now())
	return claims.ExpiresAt.Before(currentTime.Time)
}
