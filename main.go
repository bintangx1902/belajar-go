package main

import (
	"TestApp/Configs"
	"TestApp/Items"
	"TestApp/Notes"
	"TestApp/Users"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
)

func landingPage(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func main() {
	if err := Configs.ConnectDB(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/", landingPage)
	e.GET("", landingPage)

	apiGroup := e.Group("/api")
	Notes.NotesRouter(apiGroup, Configs.DB)
	Users.UserRouter(apiGroup, Configs.DB)
	Items.ItemsRouter(apiGroup, Configs.DB)
	e.Logger.Fatal(e.Start(":8000"))

}
