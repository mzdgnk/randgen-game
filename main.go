package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.GET("/", helloworld)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

func helloworld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World")
}

