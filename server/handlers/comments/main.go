package comments

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

func ListComments(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	formID := c.Param("formId")

	comments, err := cc.Query.ListComments(
		*cc.DbCtx,
		db.ListCommentsParams{
			FormID: formID,
			UserID: user.ID,
		},
	)

	if err != nil {
		log.Error("failed to fetch comments", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to fetch comments.")),
		)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": utils.EmptyArrayIfNull(comments),
	})
}

func CreateComment(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	formID := c.Param("formId")

	type Payload struct {
		Body    string  `json:"body" validate:"required"`
		Element string  `json:"element" validate:"required"`
		Parent  *string `json:"parent" validate:"omitempty,ulid"`
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

	comment, err := cc.Query.CreateComment(
		*cc.DbCtx,
		db.CreateCommentParams{
			FormID:  formID,
			UserID:  user.ID,
			Body:    payload.Body,
			Element: payload.Element,
			Parent:  payload.Parent,
		},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Hint == "forbidden" {
			return c.JSON(
				http.StatusForbidden,
				utils.FromError(
					utils.HttpErrorCode(pgErr.Hint),
					errors.New(pgErr.Message),
				),
			)
		}

		log.Error("failed to create form", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to create form.")),
		)
	}

	return c.JSON(http.StatusCreated, comment)
}

func UpdateComment(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	formID := c.Param("formId")
	commentID := c.Param("commentId")

	type Payload struct {
		Body  *string `json:"body"`
		State *string `json:"state" validate:"omitempty,oneof=hidden visible"`
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

	var state db.NullCommentState
	if payload.State != nil {
		err := state.CommentState.Scan(*payload.State)
		if err == nil {
			state.Valid = true
		}
	} else {
		state.Valid = false
	}

	comment, err := cc.Query.UpdateComment(
		*cc.DbCtx,
		db.UpdateCommentParams{
			ID:     commentID,
			FormID: formID,
			UserID: user.ID,
			Body:   payload.Body,
			State:  state,
		},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Hint == "forbidden" {
			return c.JSON(
				http.StatusForbidden,
				utils.FromError(
					utils.HttpErrorCode(pgErr.Hint),
					errors.New(pgErr.Message),
				),
			)
		}

		log.Error("failed to update comment", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to update comment.")),
		)
	}

	return c.JSON(http.StatusOK, comment)
}

func DeleteComment(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	formID := c.Param("formId")
	commentID := c.Param("commentId")

	err := cc.Query.DeleteComment(
		*cc.DbCtx,
		db.DeleteCommentParams{
			ID:     commentID,
			FormID: formID,
			UserID: user.ID,
		},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Hint == "forbidden" {
			return c.JSON(
				http.StatusForbidden,
				utils.FromError(
					utils.HttpErrorCode(pgErr.Hint),
					errors.New(pgErr.Message),
				),
			)
		}

		log.Error("failed to delete comment", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to delete comment.")),
		)
	}

	return c.NoContent(http.StatusNoContent)
}
