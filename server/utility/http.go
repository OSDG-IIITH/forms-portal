package utils

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
)

type HttpErrorCode string

const (
	ErrorBadRequest   HttpErrorCode = "bad-request"
	ErrorUnauthorized HttpErrorCode = "unauthorized"
	ErrorForbidden    HttpErrorCode = "forbidden"
	ErrorNotFound     HttpErrorCode = "not-found"
	ErrorConflict     HttpErrorCode = "already-exists"
	ErrorInvalidRange HttpErrorCode = "invalid-range"
	ErrorRateLimited  HttpErrorCode = "rate-limited"
	ErrorInternal     HttpErrorCode = "internal-server-error"

	ErrorFormClosed HttpErrorCode = "form-closed"
)

type HttpError struct {
	Code    HttpErrorCode `json:"code"`
	Message string        `json:"message"`
}

type HttpResponse struct {
	Error *HttpError `json:"error,omitempty"`
}

func FromError(code HttpErrorCode, err error) HttpResponse {
	return HttpResponse{
		Error: &HttpError{
			Code:    code,
			Message: err.Error(),
		},
	}
}

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	errCode := ErrorInternal
	msg := "An error occurred on the server side while processing your request. Please report this error."
	if httpError, ok := err.(*echo.HTTPError); ok {
		code = httpError.Code
		if code != http.StatusInternalServerError {
			msg = httpError.Message.(string)
		}

		if code == http.StatusNotFound {
			errCode = ErrorNotFound
			msg = "Cannot " + c.Request().Method + " this endpoint."
		}
	}

	if code == http.StatusInternalServerError {
		log.Error(
			"captured uncaught",
			"error", err, "code", code,
		)
	}

	c.JSON(code, FromError(errCode, errors.New(msg)))
}

func SetupTLS() {
	pool := x509.NewCertPool()
	file := Config.TlsCertsPath

	if file == "" {
		return
	}

	certs, err := os.ReadFile(file)
	if err != nil {
		log.Warn("could not read ca certificates", "file", file)
		return
	}

	ok := pool.AppendCertsFromPEM(certs)
	if !ok {
		log.Warn("failed to parse ca certs", "file", file)
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: pool},
		},
	}
	*http.DefaultClient = *client
}
