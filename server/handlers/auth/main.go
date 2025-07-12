package auth

import (
	"backend/context"
	"backend/db"
	"backend/utility"
	"encoding/json"
	"io"

	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
)

func getLoginUrl() string {
	return utils.Config.CasBaseUrl +
		"/login?service=" + url.QueryEscape(utils.Config.CasServiceUrl)
}

func getLogoutUrl() string {
	return utils.Config.CasBaseUrl + "/logout"
}

func getValidationUrl(ticket string) string {
	return utils.Config.CasBaseUrl + "/serviceValidate?service=" +
		url.QueryEscape(utils.Config.CasServiceUrl) +
		"&ticket=" + url.QueryEscape(ticket) + "&format=JSON"
}

type CasResponse struct {
	ServiceResponse struct {
		AuthenticationSuccess struct {
			Attributes struct {
				ID    []string `json:"uid"`
				Name  []string `json:"Name"`
				Email []string `json:"E-Mail"`
			} `json:"attributes"`
		} `json:"authenticationSuccess,omitempty"`
		AuthenticationFailure struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"authenticationFailure,omitempty"`
	} `json:"serviceResponse"`
}

func validateTicket(ticket string) (*CasResponse, error) {
	url := getValidationUrl(ticket)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data CasResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func Login(c echo.Context) error {
	return c.Redirect(http.StatusFound, getLoginUrl())
}

func Callback(c echo.Context) error {
	frontend, _ := url.ParseRequestURI(utils.Config.FrontendUrl)
	query := frontend.Query()

	ticket := c.QueryParam("ticket")
	response, err := validateTicket(ticket)
	if err != nil {
		log.Error("failed to validate cas ticket", "error", err.Error())

		query.Set("error_code", string(utils.ErrorInternal))
		query.Set("error_message", "An error occurred while logging you in. Please try again.")
		frontend.RawQuery = query.Encode()
		return c.Redirect(http.StatusFound, frontend.String())
	}

	if response.ServiceResponse.AuthenticationFailure.Code == "INVALID_TICKET" {
		return c.Redirect(http.StatusFound, getLoginUrl())
	}

	attributes := response.ServiceResponse.AuthenticationSuccess.Attributes
	if attributes.ID == nil || attributes.ID[0] == "" {
		log.Warn("no user id in cas response", "response", response)

		query.Set("error_code", string(utils.ErrorUnauthorized))
		query.Set(
			"error_message",
			"Failed to authorize due to the following CAS error: "+
				response.ServiceResponse.AuthenticationFailure.Code+".",
		)
		frontend.RawQuery = query.Encode()
		return c.Redirect(http.StatusFound, frontend.String())
	}

	cc := c.(*dbcontext.Context)
	user, err := cc.Query.EnsureUser(*cc.DbCtx, db.EnsureUserParams{
		Handle: attributes.ID[0], Email: attributes.Email[0], Name: attributes.Name[0],
	})

	if err != nil {
		log.Error("failed to upsert user in database", "error", err.Error())

		query.Set("error_code", string(utils.ErrorInternal))
		query.Set("error_message", "Failed to fetch user details.")
		frontend.RawQuery = query.Encode()
		return c.Redirect(http.StatusFound, frontend.String())
	}

	session := utils.CreateSession(user.ID, utils.DefaultSessionTtl)
	c.SetCookie(session.Cookie())

	return c.Redirect(http.StatusFound, frontend.String())
}

func Logout(c echo.Context) error {
	return c.Redirect(http.StatusFound, getLogoutUrl())
}

func Info(c echo.Context) error {
	user := c.Get("user").(db.User)
	return c.JSON(http.StatusOK, user)
}
