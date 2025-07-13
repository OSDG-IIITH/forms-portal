package handlers

import (
	"backend/handlers/auth"
	"backend/handlers/forms"
	"backend/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterAll(router *echo.Group) {
	router.GET("/auth/login", auth.Login)
	router.GET("/auth/login/callback", auth.Callback)
	router.GET("/auth/logout", auth.Logout)
	router.GET("/auth/info", middleware.Auth(auth.Info))

	router.GET("/forms", middleware.Auth(forms.ListForms))
	router.POST("/forms", middleware.Auth(forms.CreateForm))
}
