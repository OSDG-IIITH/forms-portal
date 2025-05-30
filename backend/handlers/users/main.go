package users

import (
	"backend/context"
	"backend/utility"

	"errors"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
)

func List(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	users, err := cc.Query.ListUsers(*cc.DbCtx)

	if err != nil {
		log.Error(
			"failed to fetch users list from database",
			"error", err.Error(),
		)

		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(
				utils.ErrorInternal,
				errors.New("Failed to fetch users list from database."),
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		utils.HttpResponse{
			Data: utils.EmptyArrayIfNull(users),
		},
	)
}
