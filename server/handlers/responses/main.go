package responses

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

func ListResponses(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	formID := c.Param("formId")

	type Query struct {
		Status db.ResponseStatus `query:"status" validate:"oneof=draft completed edited"`
		Sort   string            `query:"sort" validate:"oneof=started submitted"`
		Order  string            `query:"order" validate:"oneof=asc desc"`
		Limit  int32             `query:"limit" validate:"gte=1,lte=100"`
		Offset int32             `query:"offset" validate:"gte=0"`
	}

	query := Query{
		Status: "completed",
		Sort:   "submitted",
		Order:  "desc",
		Limit:  20,
	}

	err := c.Bind(&query)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			utils.FromError(
				utils.ErrorBadRequest,
				errors.New("Failed to parse request payload."),
			),
		)
	}

	err = utils.Validate.Struct(query)
	if err != nil {
		message := utils.FormatValidationErrors(err)
		return c.JSON(
			http.StatusUnprocessableEntity,
			utils.FromError(
				utils.ErrorBadRequest,
				errors.New(message),
			),
		)
	}

	responses, err := cc.Query.ListResponses(
		*cc.DbCtx,
		db.ListResponsesParams{
			FormID:    formID,
			UserID:    user.ID,
			Status:    query.Status,
			SortBy:    query.Sort,
			OrderBy:   query.Order,
			LimitVal:  query.Limit,
			OffsetVal: query.Offset,
		},
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Hint == "forbidden" {
			return c.JSON(
				http.StatusForbidden,
				utils.FromError(utils.HttpErrorCode(pgErr.Hint), errors.New(pgErr.Message)),
			)
		}

		log.Error("failed to fetch responses", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to fetch responses.")),
		)
	}

	total, err := cc.Query.CountResponses(
		*cc.DbCtx,
		db.CountResponsesParams{
			FormID: formID,
			UserID: user.ID,
			Status: query.Status,
		},
	)
	if err != nil {
		log.Error("failed to fetch responses", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to count responses.")),
		)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": utils.EmptyArrayIfNull(responses),
		"pagination": map[string]int64{
			"offset": int64(query.Offset),
			"limit":  int64(query.Limit),
			"total":  total,
		},
	})
}

func StartResponse(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	formID := c.Param("formId")

	response, err := cc.Query.StartResponse(
		*cc.DbCtx,
		db.StartResponseParams{
			FormID: formID,
			UserID: user.ID,
		},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && (pgErr.Hint == "form-closed" || pgErr.Hint == "forbidden") {
			return c.JSON(
				http.StatusForbidden,
				utils.FromError(utils.HttpErrorCode(pgErr.Hint), errors.New(pgErr.Message)),
			)
		}

		log.Error("failed to start response", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to start response.")),
		)
	}

	return c.JSON(http.StatusCreated, response)
}

func GetResponse(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	formID := c.Param("formId")
	responseID := c.Param("responseId")

	response, err := cc.Query.GetResponse(
		*cc.DbCtx,
		db.GetResponseParams{
			ID:     responseID,
			FormID: formID,
			UserID: user.ID,
		},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Hint == "forbidden" {
			return c.JSON(
				http.StatusForbidden,
				utils.FromError(utils.HttpErrorCode(pgErr.Hint), errors.New(pgErr.Message)),
			)
		}

		log.Error("failed to fetch response", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to fetch response.")),
		)
	}

	return c.JSON(http.StatusOK, response)
}

func GetAnswers(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	formID := c.Param("formId")
	responseID := c.Param("responseId")

	answers, err := cc.Query.GetAnswers(
		*cc.DbCtx,
		db.GetAnswersParams{
			ID:     responseID,
			FormID: formID,
			UserID: user.ID,
		},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Hint == "forbidden" {
			return c.JSON(
				http.StatusForbidden,
				utils.FromError(utils.HttpErrorCode(pgErr.Hint), errors.New(pgErr.Message)),
			)
		}

		log.Error("failed to fetch answers", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to fetch answers.")),
		)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": utils.EmptyArrayIfNull(answers),
	})
}

func SaveAnswer(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	formID := c.Param("formId")
	responseID := c.Param("responseId")

	type Payload struct {
		Question string `json:"question" validate:"required"`
		Value    string `json:"value" validate:"required"`
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

	answer, err := cc.Query.SaveAnswer(
		*cc.DbCtx,
		db.SaveAnswerParams{
			ID:       responseID,
			FormID:   formID,
			UserID:   user.ID,
			Question: payload.Question,
			Value:    payload.Value,
		},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Hint == "forbidden" {
			return c.JSON(
				http.StatusForbidden,
				utils.FromError(utils.HttpErrorCode(pgErr.Hint), errors.New(pgErr.Message)),
			)
		}

		log.Error("failed to save answer", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to save answer.")),
		)
	}

	return c.JSON(http.StatusOK, answer)
}

func SubmitResponse(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	formID := c.Param("formId")
	responseID := c.Param("responseId")

	response, err := cc.Query.SubmitResponse(
		*cc.DbCtx,
		db.SubmitResponseParams{
			ID:     responseID,
			FormID: formID,
			UserID: user.ID,
		},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && (pgErr.Hint == "form-closed" || pgErr.Hint == "forbidden") {
			return c.JSON(
				http.StatusForbidden,
				utils.FromError(utils.HttpErrorCode(pgErr.Hint), errors.New(pgErr.Message)),
			)
		}

		log.Error("failed to submit response", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to submit response.")),
		)
	}

	return c.JSON(http.StatusCreated, response)
}
