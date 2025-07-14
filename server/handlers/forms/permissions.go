package forms

import (
	"backend/context"
	"backend/db"
	"backend/utility"
	"errors"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgerrcode"
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

func GrantPermission(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	formID := c.Param("formId")

	type Payload struct {
		Role  db.PermissionRole `json:"role" validate:"oneof=view respond comment analyze edit manage"`
		User  *string           `json:"user"`
		Group *string           `json:"group"`
	}

	payload := Payload{}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			utils.FromError(
				utils.ErrorBadRequest,
				errors.New("Failed to parse request payload."),
			),
		)
	}

	if err := utils.Validate.Struct(payload); err != nil {
		message := utils.FormatValidationErrors(err)
		return c.JSON(
			http.StatusUnprocessableEntity,
			utils.FromError(
				utils.ErrorBadRequest,
				errors.New(message),
			),
		)
	}

	permission, err := cc.Query.GrantPermission(
		*cc.DbCtx,
		db.GrantPermissionParams{
			UserID:      user.ID,
			FormID:      formID,
			TargetUser:  payload.User,
			TargetGroup: payload.Group,
			Role:        payload.Role,
		},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.CheckViolation {
				return c.JSON(
					http.StatusUnprocessableEntity,
					utils.FromError(
						utils.ErrorBadRequest,
						errors.New("Failed to process payload - either a user ID or a group ID must be specified, but not both."),
					),
				)
			}

			if pgErr.Hint == "forbidden" {
				return c.JSON(
					http.StatusForbidden,
					utils.FromError(utils.HttpErrorCode(pgErr.Hint), errors.New(pgErr.Message)),
				)
			}

			if pgErr.Hint == "not-found" {
				return c.JSON(
					http.StatusNotFound,
					utils.FromError(utils.HttpErrorCode(pgErr.Hint), errors.New(pgErr.Message)),
				)
			}
		}

		log.Error("failed to grant permission", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to grant permission.")),
		)
	}

	return c.JSON(http.StatusCreated, permission)
}

func RevokePermission(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	formID := c.Param("formId")
	permissionID := c.Param("permissionId")

	err := cc.Query.RevokePermission(
		*cc.DbCtx,
		db.RevokePermissionParams{
			UserID:       user.ID,
			FormID:       formID,
			PermissionID: permissionID,
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

		log.Error("failed to revoke permission", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to revoke permission.")),
		)
	}

	return c.NoContent(http.StatusNoContent)
}
