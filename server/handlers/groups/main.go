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

func GetGroup(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	groupId := c.Param("groupId")

	group, err := cc.Query.GetGroup(
		*cc.DbCtx,
		db.GetGroupParams{
			ID:     groupId,
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

		log.Error("failed to fetch group", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to fetch group.")),
		)
	}

	return c.JSON(http.StatusOK, group)
}

func UpdateGroup(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	groupId := c.Param("groupId")

	type Payload struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
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

	group, err := cc.Query.UpdateGroup(
		*cc.DbCtx,
		db.UpdateGroupParams{
			ID:          groupId,
			UserID:      user.ID,
			Name:        payload.Name,
			Description: payload.Description,
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

		log.Error("failed to update group", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to update group.")),
		)
	}

	return c.JSON(http.StatusOK, group)
}

func DeleteGroup(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	groupId := c.Param("groupId")

	err := cc.Query.DeleteGroup(
		*cc.DbCtx,
		db.DeleteGroupParams{
			ID:     groupId,
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

		log.Error("failed to delete group", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to delete group.")),
		)
	}

	return c.NoContent(http.StatusNoContent)
}

func UpdateGroupDomain(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	groupId := c.Param("groupId")

	type Payload struct {
		Domain string `json:"domain" validate:"required,fqdn"`
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

	err := cc.Query.UpdateGroupDomain(
		*cc.DbCtx,
		db.UpdateGroupDomainParams{
			ID:     groupId,
			UserID: user.ID,
			Domain: payload.Domain,
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

		log.Error("failed to update group domain", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to update group domain.")),
		)
	}

	return c.NoContent(http.StatusNoContent)
}

func AddGroupMember(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	groupId := c.Param("groupId")

	type Payload struct {
		Email string `json:"email" validate:"required,email"`
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

	err := cc.Query.AddGroupMember(
		*cc.DbCtx,
		db.AddGroupMemberParams{
			GroupID:    groupId,
			UserID:     user.ID,
			TargetUser: payload.Email,
		},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Hint == "not-found" {
			return c.JSON(
				http.StatusNotFound,
				utils.FromError(utils.HttpErrorCode(pgErr.Hint), errors.New(pgErr.Message)),
			)
		}

		if errors.As(err, &pgErr) && pgErr.Hint == "forbidden" {
			return c.JSON(
				http.StatusForbidden,
				utils.FromError(utils.ErrorForbidden, errors.New(pgErr.Message)),
			)
		}

		log.Error("failed to add group member", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to add group member.")),
		)
	}

	return c.NoContent(http.StatusNoContent)
}

func RemoveGroupMember(c echo.Context) error {
	cc := c.(*dbcontext.Context)
	user := c.Get("user").(db.User)

	groupId := c.Param("groupId")
	targetUserId := c.Param("userId")

	err := cc.Query.RemoveGroupMember(
		*cc.DbCtx,
		db.RemoveGroupMemberParams{
			GroupID:      groupId,
			UserID:       user.ID,
			TargetUserID: targetUserId,
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

		log.Error("failed to remove group member", "error", err)
		return c.JSON(
			http.StatusInternalServerError,
			utils.FromError(utils.ErrorInternal, errors.New("Failed to remove group member.")),
		)
	}

	return c.NoContent(http.StatusNoContent)
}
