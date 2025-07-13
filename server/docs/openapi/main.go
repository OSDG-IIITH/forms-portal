package openapi

import (
	_ "embed"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed index.html
var html string

//go:embed stoplight.min.js
var js string

//go:embed stoplight.min.css
var css string

//go:embed spec.yaml
var schema string

func RegisterRoutes(e *echo.Group) {
	e.GET("", func(c echo.Context) error {
		return c.HTML(http.StatusOK, html)
	})
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, html)
	})

	e.GET("/stoplight.min.js", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/javascript; charset=utf-8")
		return c.String(http.StatusOK, js)
	})
	e.GET("/stoplight.min.css", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/css; charset=utf-8")
		return c.String(http.StatusOK, css)
	})
	e.GET("/spec.yaml", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/yaml; charset=utf-8")
		return c.String(http.StatusOK, schema)
	})
}
