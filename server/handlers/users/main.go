package users

import (
	"backend/context"
	"backend/utility"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetUser(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	userId := c.Param("userId")

	user, err := cc.Query.GetUserById(*cc.DbCtx, userId)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			utils.FromError(utils.ErrorNotFound, errors.New("User not found")),
		)
	}

	return c.JSON(http.StatusOK, user)
}