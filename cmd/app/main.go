package main

import (
	"context"
	"fmt"

	"github.com/catinapoke/go-microservice-example/internal/domain"
	"github.com/catinapoke/go-microservice-example/internal/handlers/errorshandler"
	goodscreate "github.com/catinapoke/go-microservice-example/internal/handlers/goodsCreate"
	goodsremove "github.com/catinapoke/go-microservice-example/internal/handlers/goodsRemove"
	goodsreprioritize "github.com/catinapoke/go-microservice-example/internal/handlers/goodsReprioritize"
	goodsupdate "github.com/catinapoke/go-microservice-example/internal/handlers/goodsUpdate"
	"github.com/catinapoke/go-microservice-example/internal/handlers/goodslist"
	"github.com/catinapoke/go-microservice-example/internal/repository/postgres"
	srvwrapper "github.com/catinapoke/go-microservice-example/utils/srwwrapper"
	"github.com/catinapoke/go-microservice-example/utils/tx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const (
	Port        = ":8080"
	DatabaseUrl = "postgres://user:password@postgresql:5432/example?sslmode=disable"
)

func main() {
	// Database
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, DatabaseUrl)
	if err != nil {
		log.Fatal(fmt.Errorf("connect to db: %w", err))
	}
	defer pool.Close()

	provider := tx.New(pool)
	repo := postgres.New(provider)

	// Service
	// repo := localgoods.New()
	model := domain.New(repo)

	creator := goodscreate.Handler{Model: model}
	updater := goodsupdate.Handler{Model: model}
	remover := goodsremove.Handler{Model: model}
	listing := goodslist.Handler{Model: model}
	prioritizer := goodsreprioritize.Handler{Model: model}

	e := echo.New()

	e.Logger.SetLevel(log.INFO)

	e.POST("/good/create", srvwrapper.New(creator.Handle).ServeHTTP)
	e.PATCH("/good/update", srvwrapper.New(updater.Handle).ServeHTTP)
	e.DELETE("/good/remove", srvwrapper.New(remover.Handle).ServeHTTP)
	e.GET("/good/list", srvwrapper.New(listing.Handle).ServeHTTP)
	e.PATCH("/good/reprioritize", srvwrapper.New(prioritizer.Handle).ServeHTTP)

	// Root level middleware
	e.Use(middleware.Recover())
	e.Use(errorshandler.HandleError)
	e.Use(middleware.Logger())

	e.Logger.Fatal(e.Start(Port))
}
