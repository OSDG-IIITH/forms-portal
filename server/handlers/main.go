package handlers

import (
	"backend/handlers/auth"
	"backend/handlers/forms"
	"backend/handlers/users"
	"backend/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterAll(router *echo.Group) {
	router.GET("/auth/login", auth.Login)
	router.GET("/auth/login/callback", auth.Callback)
	router.GET("/auth/logout", auth.Logout)
	router.GET("/auth/info", middleware.Auth(auth.Info))

	router.GET("/users/:userId", middleware.Auth(users.GetUser))

	router.GET("/forms", middleware.Auth(forms.ListForms))
	router.POST("/forms", middleware.Auth(forms.CreateForm))
	router.GET("/forms/:handle/:slug", middleware.Auth(forms.ResolveForm))
	router.GET("/forms/:formId", middleware.Auth(forms.GetForm))
	router.PATCH("/forms/:formId", middleware.Auth(forms.UpdateForm))
	router.DELETE("/forms/:formId", middleware.Auth(forms.DeleteForm))
}
