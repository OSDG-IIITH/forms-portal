package middleware

import (
	"backend/context"
	"backend/utility"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(utils.SessionCookieName)
		if err != nil {
			return c.JSON(
				http.StatusUnauthorized,
				utils.FromError(utils.ErrorUnauthorized, errors.New("Not logged in.")),
			)
		}

		session, err := utils.ValidateSession(cookie.Value)
		if err != nil {
			return c.JSON(
				http.StatusUnauthorized,
				utils.FromError(utils.ErrorUnauthorized, errors.New("Invalid session.")),
			)
		}

		cc := c.(*dbcontext.Context)
		user, err := cc.Query.GetUserById(*cc.DbCtx, session.ID)
		if err != nil {
			return c.JSON(
				http.StatusUnauthorized,
				utils.FromError(utils.ErrorUnauthorized, errors.New("Invalid user.")),
			)
		}

		c.Set("user", user)

		return next(c)
	}
}
