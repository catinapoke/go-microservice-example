package errorshandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorMessage struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details struct{} `json:"details"`
}

func HandleError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract the credentials from HTTP request header and perform a security
		// check
		err := next(c)

		if err != nil {
			method := c.Request().Method
			if method == "PATCH" || method == "DELETE" {
				return echo.NewHTTPError(http.StatusNotFound, ErrorMessage{
					Code:    3,
					Message: "errors.good.notFound",
				})
			}

			return err
		}

		return nil
	}
}
