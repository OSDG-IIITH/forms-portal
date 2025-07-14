package groups

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

func ListGroups(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	type Query struct {
		Sort   string `query:"sort" validate:"oneof=created name type"`
		Order  string `query:"order" validate:"oneof=asc desc"`
		Limit  int32  `query:"limit" validate:"gte=1,lte=100"`
		Offset int32  `query:"offset" validate:"gte=0"`
	}

	query := Query{Sort: "created", Order: "desc", Limit: 20}

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

	groups, err := cc.Query.ListGroups(
		*cc.DbCtx,
		db.ListGroupsParams{
			UserID:    user.ID,
			SortBy:    query.Sort,
			OrderBy:   query.Order,
			LimitVal:  query.Limit,
			OffsetVal: query.Offset,
		},
	)
	if err != nil {
		log.Error("failed to fetch groups", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to fetch groups.")),
		)
	}

	total, err := cc.Query.CountGroups(
		*cc.DbCtx,
		user.ID,
	)
	if err != nil {
		log.Error("failed to count groups", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to count groups.")),
		)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": utils.EmptyArrayIfNull(groups),
		"pagination": map[string]int64{
			"offset": int64(query.Offset),
			"limit":  int64(query.Limit),
			"total":  total,
		},
	})
}

func CreateGroup(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	type Payload struct {
		Name        string       `json:"name" validate:"required"`
		Description *string      `json:"description"`
		Type        db.GroupType `json:"type" validate:"required,oneof=domain list"`
		Domain      *string      `json:"domain" validate:"omitempty,fqdn"`
		Members     []string     `json:"members" validate:"dive,email"`
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

	group, err := cc.Query.CreateGroup(
		*cc.DbCtx,
		db.CreateGroupParams{
			OwnerID:     user.ID,
			Name:        payload.Name,
			Description: payload.Description,
			Type:        payload.Type,
			Domain:      payload.Domain,
			Members:     payload.Members,
		},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return c.JSON(
				http.StatusConflict,
				utils.FromError(
					utils.ErrorConflict,
					errors.New("A group with the same name already exists."),
				),
			)
		}

		if errors.As(err, &pgErr) && pgErr.Hint == "not-found" {
			return c.JSON(
				http.StatusNotFound,
				utils.FromError(
					utils.HttpErrorCode(pgErr.Hint),
					errors.New(pgErr.Message),
				),
			)
		}

		log.Error("failed to create group", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to create group.")),
		)
	}

	return c.JSON(http.StatusCreated, group)
}
