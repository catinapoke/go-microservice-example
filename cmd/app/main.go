package main

import (
	"github.com/catinapoke/go-microservice-example/internal/handlers/errorshandler"
	goodscreate "github.com/catinapoke/go-microservice-example/internal/handlers/goodsCreate"
	goodsremove "github.com/catinapoke/go-microservice-example/internal/handlers/goodsRemove"
	goodsreprioritize "github.com/catinapoke/go-microservice-example/internal/handlers/goodsReprioritize"
	goodsupdate "github.com/catinapoke/go-microservice-example/internal/handlers/goodsUpdate"
	"github.com/catinapoke/go-microservice-example/internal/handlers/goodslist"
	srvwrapper "github.com/catinapoke/go-microservice-example/utils/srwwrapper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const port = ":8080"

func main() {
	creator := goodscreate.Handler{}
	updater := goodsupdate.Handler{}
	remover := goodsremove.Handler{}
	listing := goodslist.Handler{}
	prioritizer := goodsreprioritize.Handler{}

	e := echo.New()

	e.POST("/good/create", srvwrapper.New(creator.Handle).ServeHTTP)
	e.PATCH("/good/update", srvwrapper.New(updater.Handle).ServeHTTP)
	e.DELETE("/good/remove", srvwrapper.New(remover.Handle).ServeHTTP)
	e.GET("/good/list", srvwrapper.New(listing.Handle).ServeHTTP)
	e.PATCH("/good/reprioritize", srvwrapper.New(prioritizer.Handle).ServeHTTP)

	// Root level middleware
	e.Use(middleware.Recover())
	e.Use(errorshandler.HandleError)
	e.Use(middleware.Logger())

	e.Logger.Fatal(e.Start(port))
}
