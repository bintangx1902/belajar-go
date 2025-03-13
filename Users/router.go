package Users

import (
	"TestApp/Users/Controllers"
	"TestApp/Users/Repositories"
	"TestApp/Users/Services"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func UserRouter(route *echo.Group, db *gorm.DB) {
	repo := Repositories.NewUsersRepository(db)
	service := Services.NewUsersService(repo)
	controller := Controllers.NewUsersController(service)

	userGroup := route.Group("/users")
	userGroup.POST("", controller.CreateUser)
	userGroup.POST("/login", controller.Login)

}
