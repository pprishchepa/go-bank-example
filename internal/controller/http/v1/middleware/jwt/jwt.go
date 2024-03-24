package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
)

func Authorize(secret string) gin.HandlerFunc {
	keyFunc := hmacKeyFunc([]byte(secret))

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			unauthorizedResp(c, errors.New("authorization header required"))
			return
		}

		scheme, tokenString := parseAuthHeader(authHeader)
		if scheme != "bearer" {
			unauthorizedResp(c, errors.New("unexpected auth scheme"))
			return
		}
		if tokenString == "" {
			unauthorizedResp(c, errors.New("empty token"))
			return
		}

		token, err := jwt.Parse(tokenString, keyFunc)
		if err != nil {
			unauthorizedResp(c, fmt.Errorf("parse token: %w", err))
			return
		}
		if !token.Valid {
			unauthorizedResp(c, fmt.Errorf("invalid token"))
			return
		}

		c.Next()
	}
}

func unauthorizedResp(c *gin.Context, err error) {
	log.Debug().Err(err).Msg("unauthorized request")
	c.JSON(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	c.Abort()
}

func parseAuthHeader(s string) (scheme, token string) {
	chunks := strings.Split(strings.Trim(s, " "), " ")
	if len(chunks) == 2 {
		scheme = strings.ToLower(chunks[0])
		token = chunks[1]
	}
	return
}

func hmacKeyFunc(hmacSecret []byte) jwt.Keyfunc {
	if len(hmacSecret) == 0 {
		return func(token *jwt.Token) (interface{}, error) {
			return nil, fmt.Errorf("no HMAC secret configured")
		}
	}

	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	}
}
