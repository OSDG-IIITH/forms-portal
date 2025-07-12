package handlers

import (
	"backend/handlers/auth"

	"github.com/labstack/echo/v4"
)

func RegisterAll(router *echo.Group) {
	router.GET("/auth/login", auth.Login)
	router.GET("/auth/login/callback", auth.Callback)
}
