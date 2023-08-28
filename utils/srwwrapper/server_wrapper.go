package srvwrapper

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Wrapper[Req Validator, Res any] struct {
	fn func(ctx context.Context, req Req) (Res, error)
}

type Validator interface {
	Validate() error
}

func New[Req Validator, Res any](fn func(ctx context.Context, req Req) (Res, error)) *Wrapper[Req, Res] {
	return &Wrapper[Req, Res]{
		fn: fn,
	}
}

func (w *Wrapper[Req, Res]) ServeHTTP(c echo.Context) error {
	var req Req

	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	// Avoid restriction of binding Query parameters (only for GET/DELETE methods)
	methodName := c.Request().Method
	if methodName != "GET" && methodName != "DELETE" {
		if err := (&echo.DefaultBinder{}).BindQueryParams(c, &req); err != nil {
			return err
		}
	}

	reqValidation, ok := any(req).(Validator)
	if ok {
		errValidation := reqValidation.Validate()
		if errValidation != nil {
			c.Logger().Errorf("validation error: %+v", req)
			return c.String(http.StatusBadRequest, "bad request")
		}
	}

	resp, err := w.fn(c.Request().Context(), req)
	if err != nil {
		log.Printf("executor fail: %s", err)
		return c.String(http.StatusInternalServerError, "exec handler")
	}

	return c.JSON(http.StatusOK, resp)
}
