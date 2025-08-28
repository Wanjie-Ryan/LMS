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

// jwt claim is data stored inside a jwt token
type CustomJWTClaims struct {
	ID                   uint `json:"id"`
	jwt.RegisteredClaims      // comes from the jwt library and already includes fields like IssuedAt, ExpiresAt, issuer/subject (optional)
}

func GenerateJWT(user models.User) (*string, *string, error) {

	// creating a CustomJwtClaims struct. Also, jwt.RegisteredClaims is a struct from the jwt library that can accept various data
	userClaims := CustomJWTClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48)),
		},
	}

	// creates a new JWT object
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	signedAccessToken, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return nil, nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomJWTClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48)),
		},
	})

	signedRefreshToken, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, nil, err
	}

	return &signedAccessToken, &signedRefreshToken, nil

}

// verifies a JWT string and extracts the claims
func ParseJWT(signedAccessToken string) (*CustomJWTClaims, error) {

	// the signedAccessToken is the first argument, followed by an empty calims struct that will be populated, the 3rd args is a function that returns the secret key.
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
