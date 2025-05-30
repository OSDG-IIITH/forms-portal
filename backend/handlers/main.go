package handlers

import (
	"backend/handlers/users"

	"github.com/labstack/echo/v4"
)

func RegisterAll(router *echo.Group) {
	router.GET("/users", users.List)
}
