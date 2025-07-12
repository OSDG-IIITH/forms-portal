package dbcontext

import (
	"backend/db"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	DbConn *pgxpool.Pool
	DbCtx  *context.Context
	Query  *db.Queries
}

func New(ctx echo.Context, conn *pgxpool.Pool, dbctx *context.Context, query *db.Queries) *Context {
	return &Context{ctx, conn, dbctx, query}
}
