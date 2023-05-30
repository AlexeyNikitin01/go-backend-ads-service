package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ads/internal/app"
)

func NewHTTPServer(port string, a app.App) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	s := &http.Server{
		Addr: port,
		Handler: router,
	}

	AppRouter(router.Group("api/v1"), a)

	return s
}
