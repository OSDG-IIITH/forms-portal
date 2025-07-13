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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func ListForms(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	type Query struct {
		Role   string `query:"role" validate:"omitempty,oneof=view respond comment analyze edit manage"`
		Sort   string `query:"sort" validate:"oneof=modified title"`
		Order  string `query:"order" validate:"oneof=asc desc"`
		Limit  int32  `query:"limit" validate:"gte=1,lte=100"`
		Offset int32  `query:"offset" validate:"gte=0"`
	}

	query := Query{Sort: "modified", Order: "desc", Limit: 20}

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

	var filterRole db.NullPermissionRole
	if query.Role != "" {
		err := filterRole.PermissionRole.Scan(query.Role)
		if err == nil {
			filterRole.Valid = true
		}
	} else {
		filterRole.Valid = false
	}

	forms, err := cc.Query.ListForms(
		*cc.DbCtx,
		db.ListFormsParams{
			UserID:     user.ID,
			FilterRole: filterRole,
			SortBy:     query.Sort,
			OrderBy:    query.Order,
			LimitVal:   query.Limit,
			OffsetVal:  query.Offset,
		},
	)
	if err != nil {
		log.Error("failed to fetch forms", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to fetch forms.")),
		)
	}

	total, err := cc.Query.CountForms(
		*cc.DbCtx,
		db.CountFormsParams{
			UserID:     user.ID,
			FilterRole: filterRole,
		},
	)
	if err != nil {
		log.Error("failed to fetch forms", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to count forms.")),
		)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": utils.EmptyArrayIfNull(forms),
		"pagination": map[string]int64{
			"offset": int64(query.Offset),
			"limit":  int64(query.Limit),
			"total":  total,
		},
	})
}

func CreateForm(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	type Payload struct {
		Title             string              `json:"title" validate:"required"`
		Slug              string              `json:"slug" validate:"required"`
		Description       *string             `json:"description"`
		Structure         string              `json:"structure" validate:"required"`
		Live              bool                `json:"live"`
		Opens             *pgtype.Timestamptz `json:"opens"`
		Closes            *pgtype.Timestamptz `json:"closes"`
		MaxResponses      *int32              `json:"max_responses"`
		IndividualLimit   int32               `json:"individual_limit" validate:"gte=1"`
		EditableResponses bool                `json:"editable_responses"`
	}

	payload := Payload{Live: false, IndividualLimit: 1, EditableResponses: false}

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

	form, err := cc.Query.CreateFormWithPermissions(
		*cc.DbCtx,
		db.CreateFormWithPermissionsParams{
			OwnerID:           user.ID,
			Slug:              payload.Slug,
			Title:             payload.Title,
			Description:       payload.Description,
			Structure:         payload.Structure,
			Live:              payload.Live,
			Opens:             payload.Opens,
			Closes:            payload.Closes,
			MaxResponses:      payload.MaxResponses,
			IndividualLimit:   payload.IndividualLimit,
			EditableResponses: payload.EditableResponses,
		},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.CheckViolation {
			return c.JSON(
				http.StatusUnprocessableEntity,
				utils.FromError(
					utils.ErrorBadRequest,
					errors.New("Failed to process payload - max_responses must be greater than individual_limit"),
				),
			)
		}

		log.Error("failed to create form", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to create form.")),
		)
	}

	return c.JSON(http.StatusCreated, form)
}

func ResolveForm(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	handle := c.Param("handle")
	slug := c.Param("slug")

	form, err := cc.Query.ResolveFormByHandleAndSlug(
		*cc.DbCtx,
		db.ResolveFormByHandleAndSlugParams{
			Handle: handle,
			Slug:   slug,
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

		log.Error("failed to resolve form", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to retrieve form.")),
		)
	}

	return c.JSON(http.StatusOK, form)
}

func GetForm(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	formID := c.Param("formId")

	form, err := cc.Query.GetFormByID(
		*cc.DbCtx,
		db.GetFormByIDParams{
			ID:     formID,
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

		log.Error("failed to fetch form", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to retrieve form.")),
		)
	}

	return c.JSON(http.StatusOK, form)
}
