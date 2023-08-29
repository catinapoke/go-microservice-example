package main

import (
	"context"
	"fmt"

	"github.com/catinapoke/go-microservice-example/internal/config"
	"github.com/catinapoke/go-microservice-example/internal/domain"
	goodscreate "github.com/catinapoke/go-microservice-example/internal/handlers/goodsCreate"
	goodslist "github.com/catinapoke/go-microservice-example/internal/handlers/goodsList"
	goodsremove "github.com/catinapoke/go-microservice-example/internal/handlers/goodsRemove"
	goodsreprioritize "github.com/catinapoke/go-microservice-example/internal/handlers/goodsReprioritize"
	goodsupdate "github.com/catinapoke/go-microservice-example/internal/handlers/goodsUpdate"
	"github.com/catinapoke/go-microservice-example/internal/repository/natslogs"
	"github.com/catinapoke/go-microservice-example/internal/repository/postgres"
	"github.com/catinapoke/go-microservice-example/internal/repository/rediscache"
	"github.com/catinapoke/go-microservice-example/utils/logger"
	srvwrapper "github.com/catinapoke/go-microservice-example/utils/srwwrapper"
	"github.com/catinapoke/go-microservice-example/utils/tx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

var (
	Log echo.Logger
)

func main() {
	err := config.Init()

	if err != nil {
		log.Fatal(fmt.Errorf("can't unmarshal config: %w", err))
	}

	// Database
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, config.AppConfig.PostgresUrl)
	if err != nil {
		log.Fatal(fmt.Errorf("connect to db: %w", err))
	}
	defer pool.Close()

	provider := tx.New(pool)
	repo := postgres.New(provider)

	// Redis
	opt, err := redis.ParseURL(config.AppConfig.Redis.Url)
	if err != nil {
		log.Fatal(fmt.Errorf("connect to redis: %w", err))
	}

	client := redis.NewClient(opt)
	defer client.Close()

	_, err = client.Ping(ctx).Result()

	if err != nil {
		log.Fatal(fmt.Errorf("ping redis: %w", err))
	}

	rds := rediscache.New(client, repo, config.AppConfig.Redis.Key)

	// Nats
	conn, err := nats.Connect(config.AppConfig.Nats.Url)
	if err != nil {
		log.Fatal(fmt.Errorf("connect to nats: %w", err))
	}
	defer conn.Close()

	natsClient := natslogs.New(conn, rds, config.AppConfig.Nats.Subject)

	// Service
	model := domain.New(natsClient, provider)

	creator := goodscreate.Handler{Model: model}
	updater := goodsupdate.Handler{Model: model}
	remover := goodsremove.Handler{Model: model}
	listing := goodslist.Handler{Model: model}
	prioritizer := goodsreprioritize.Handler{Model: model}

	e := echo.New()

	e.Logger.SetLevel(log.INFO)
	logger.Log = e.Logger

	e.POST("/good/create", srvwrapper.New(creator.Handle).ServeHTTP)
	e.PATCH("/good/update", srvwrapper.New(updater.Handle).ServeHTTP)
	e.DELETE("/good/remove", srvwrapper.New(remover.Handle).ServeHTTP)
	e.GET("/good/list", srvwrapper.New(listing.Handle).ServeHTTP)
	e.PATCH("/good/reprioritize", srvwrapper.New(prioritizer.Handle).ServeHTTP)

	// Root level middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.Logger.Fatal(e.Start(config.AppConfig.ServicePort))
}
