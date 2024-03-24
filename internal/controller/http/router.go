package http

import (
	"net/http"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/pprishchepa/go-bank-example/internal/config"
	httpv1 "github.com/pprishchepa/go-bank-example/internal/controller/http/v1"
	"github.com/pprishchepa/go-bank-example/internal/controller/http/v1/middleware/jwt"
)

func NewRouter(conf config.Config, wallet *httpv1.WalletRoutes) http.Handler {
	gin.SetMode(gin.ReleaseMode)

	e := gin.New()
	e.ContextWithFallback = true
	e.HandleMethodNotAllowed = true

	e.Use(gin.Recovery())
	e.Use(logger.SetLogger())

	e.GET("/", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	v1 := e.Group("/api/v1", jwt.Authorize(conf.Auth.JWTSecret))
	{
		wallet.RegisterRoutes(v1)
	}

	return e
}
