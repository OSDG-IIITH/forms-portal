package middleware

import (
	"backend/utility"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterAll(srv *echo.Echo, api *echo.Group) {
	srv.Use(requestLogger)
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{utils.Config.FrontendUrl},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
		},
		AllowCredentials: true,
	}))
}
