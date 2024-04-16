package server

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var secretKey = []byte("secret-key")

const (
	AccessTokenDuration = 7 * 24 * time.Hour
	CookieExpDuration   = AccessTokenDuration - 1*time.Minute
)

func CreateJWTToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(AccessTokenDuration).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SetTokenCookie(c echo.Context, accessToken string, expiration time.Time) {
	cookie := new(http.Cookie)
	cookie.Name = "chess.access_token"
	cookie.Value = accessToken
	cookie.Expires = expiration
	c.SetCookie(cookie)
}
