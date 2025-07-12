package middleware

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
)

func requestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()

		start := time.Now()
		if err := next(c); err != nil {
			c.Error(err)
		}
		stop := time.Now()

		latency := stop.Sub(start).String()
		status := res.Status

		logger := log.Debug
		if status >= 500 {
			logger = log.Error
		}

		logger(
			"processed request",
			"status", status, "method", req.Method,
			"endpoint", c.Path(), "uri", req.RequestURI,
			"ip", c.RealIP(), "latency", latency,
			"start", start.Format(time.RFC3339),
		)

		return nil
	}
}
