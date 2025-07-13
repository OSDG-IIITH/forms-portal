package forms

import (
	"backend/context"
	"backend/db"
	"backend/utility"
	"errors"
	"net/http"
	"sort"

	"github.com/charmbracelet/log"
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

	query := Query{"", "modified", "desc", 20, 0}

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

	var role db.NullPermissionRole
	if query.Role != "" {
		err := role.PermissionRole.Scan(query.Role)
		if err == nil {
			role.Valid = true
		}
	} else {
		role.Valid = false
	}

	forms, err := cc.Query.ListForms(
		*cc.DbCtx,
		db.ListFormsParams{
			OwnerID:   user.ID,
			Role:      role,
			LimitVal:  query.Limit,
			OffsetVal: query.Offset,
		},
	)
	if err != nil {
		log.Error("failed to fetch forms", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to fetch forms.")),
		)
	}

	sortForms(forms, query.Sort, query.Order)

	total, err := cc.Query.CountForms(
		*cc.DbCtx,
		db.CountFormsParams{
			OwnerID: user.ID,
			Role:    role,
		},
	)
	if err != nil {
		log.Error("failed to fetch forms", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to count forms.")),
		)
	}

	return c.JSON(http.StatusOK, utils.HttpResponse{
		Data: map[string]interface{}{
			"data": utils.EmptyArrayIfNull(forms),
			"pagination": map[string]int64{
				"offset": int64(query.Offset),
				"limit":  int64(query.Limit),
				"total":  total,
			},
		},
	})
}

func sortForms(forms []db.Form, sortBy, order string) {
	sort.Slice(forms, func(i, j int) bool {
		switch sortBy {
		case "modified":
			if order == "asc" {
				return forms[i].Modified.Time.Before(forms[j].Modified.Time)
			} else {
				return forms[i].Modified.Time.After(forms[j].Modified.Time)
			}
		case "title":
			if order == "asc" {
				return forms[i].Title < forms[j].Title
			} else {
				return forms[i].Title > forms[j].Title
			}
		default:
			return false
		}
	})
}
