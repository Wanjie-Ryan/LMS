package middleware

import (
	"errors"
	"strings"

	"github.com/Wanjie-Ryan/LMS/common"
	"github.com/Wanjie-Ryan/LMS/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AppMiddleware struct {
	DB *gorm.DB
}

func (a *AppMiddleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		// the line below adds an HTTP response header
		// its important for caching proxies like CDN edges, reverse proxies, or API gateways.
		// by adding Vary:Authorization, you're telling caching systems:: don't cache and reuse the same response for users with different Aauthorization headers.

		c.Response().Header().Add("Vary", "Authorization")

		authHeader := c.Request().Header.Get("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return common.SendUnauthorizedResponse(c, "Invalid Authorization")
		}

		authHeaderSplit := strings.Split(authHeader, " ")
		accessToken := authHeaderSplit[1]
		claims, err := common.ParseJWT(accessToken)
		if err != nil {
			return common.SendUnauthorizedResponse(c, err.Error())
		}

		if common.IsClaimExpired(claims) {
			return common.SendUnauthorizedResponse(c, "Token Expired")
		}

		var user models.User
		result := a.DB.First(&user, claims.ID)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return common.SendUnauthorizedResponse(c, "Malformed Token")
			}

			return common.SendInternalServerError(c, "error getting user")
		}

		c.Set("user", &user)
		return next(c)

	}
}
