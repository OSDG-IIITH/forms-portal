package handlers

import (
	"backend/handlers/auth"
	"backend/handlers/comments"
	"backend/handlers/forms"
	"backend/handlers/groups"
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

	router.GET("/forms/:formId/permissions", middleware.Auth(forms.ListPermissions))
	router.POST("/forms/:formId/permissions", middleware.Auth(forms.GrantPermission))
	router.DELETE("/forms/:formId/permissions/:permissionId", middleware.Auth(forms.RevokePermission))

	router.GET("/forms/:formId/comments", middleware.Auth(comments.ListComments))
	router.POST("/forms/:formId/comments", middleware.Auth(comments.CreateComment))
	router.PATCH("/forms/:formId/comments/:commentId", middleware.Auth(comments.UpdateComment))
	router.DELETE("/forms/:formId/comments/:commentId", middleware.Auth(comments.DeleteComment))

	router.GET("/groups", middleware.Auth(groups.ListGroups))
	router.POST("/groups", middleware.Auth(groups.CreateGroup))

	router.GET("/groups/:groupId", middleware.Auth(groups.GetGroup))
	router.PATCH("/groups/:groupId", middleware.Auth(groups.UpdateGroup))
	router.DELETE("/groups/:groupId", middleware.Auth(groups.DeleteGroup))

	router.PUT("/groups/:groupId/domain", middleware.Auth(groups.UpdateGroupDomain))
	router.POST("/groups/:groupId/members", middleware.Auth(groups.AddGroupMember))
	router.DELETE("/groups/:groupId/members/:userId", middleware.Auth(groups.RemoveGroupMember))
}
