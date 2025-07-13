package forms

import (
	"backend/context"
	"backend/db"
	"backend/utility"
	"errors"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
)

func ListPermissions(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	formID := c.Param("formId")

	permissions, err := cc.Query.ListFormPermissions(
		*cc.DbCtx,
		db.ListFormPermissionsParams{
			FormID: formID,
			UserID: user.ID,
		},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Hint == "forbidden" {
			return c.JSON(
				http.StatusForbidden,
				utils.FromError(utils.ErrorForbidden, errors.New(pgErr.Message)),
			)
		}

		log.Error("failed to fetch form permissions", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to retrieve form permissions.")),
		)
	}

	return c.JSON(http.StatusOK, utils.EmptyArrayIfNull(permissions))
}
