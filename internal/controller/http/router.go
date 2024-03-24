package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpv1 "github.com/pprishchepa/go-bank-example/internal/controller/http/v1"
)

func NewRouter(wallet *httpv1.WalletRoutes) http.Handler {
	gin.SetMode(gin.ReleaseMode)

	e := gin.New()
	e.ContextWithFallback = true
	e.HandleMethodNotAllowed = true

	e.Use(gin.Recovery())

	e.GET("/", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	v1 := e.Group("/api/v1")
	{
		wallet.RegisterRoutes(v1)
	}

	return e
}
