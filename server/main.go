package main

import (
	"backend/context"
	"backend/db"
	"backend/docs/openapi"
	"backend/handlers"
	"backend/middleware"
	"backend/utility"
	"context"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func main() {
	utils.SetupLogger()
	utils.LoadConfig()
	utils.SetupTLS()
	utils.LoadValidator()

	ctx := context.Background()
	conn, err := pgxpool.New(ctx, utils.Config.DatabaseUri)
	if err != nil {
		log.Error(
			"could not connect to database",
			"error", err.Error(), "url", utils.Config.DatabaseUri,
		)

		os.Exit(1)
	}

	defer conn.Close()
	q := db.New(conn)

	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	server.HTTPErrorHandler = utils.ErrorHandler
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := dbcontext.New(c, conn, &ctx, q)
			return next(cc)
		}
	})

	base := server.Group(utils.Config.BaseUrl)
	router := base.Group("/api")
	docs := router.Group("/docs")

	middleware.RegisterAll(server, router)
	handlers.RegisterAll(router)
	openapi.RegisterRoutes(docs)

	intf := utils.Config.Host
	port := utils.Config.Port
	addr := intf + ":" + port

	log.Info("listening for requests", "intf", intf, "port", port)
	if err := server.Start(addr); err != http.ErrServerClosed {
		log.Fatal("could not start server", "error", err)
	}
}
