package resource

import (
	"fmt"

	"github.com/labstack/echo"
)

type ErrBody struct {
	Message string `json:"message"`
	ErrID   string `json:"errid"`
}

func (b ErrBody) Error() string {
	return fmt.Sprintf("[%s] %s", b.ErrID, b.Message)
}

func newDefaultErrorResponse(c echo.Context, code int, msg string) error {
	return c.JSON(code, ErrBody{Message: msg, ErrID: "0000"})
}

func newErrorResponse(c echo.Context, code int, msg string, pgerrcode string) error {
	return c.JSON(code, ErrBody{Message: msg, ErrID: pgerrcode})
}
