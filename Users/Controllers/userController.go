package Controllers

import (
	"TestApp/Models"
	"TestApp/Users/Services"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserController struct {
	Service *Services.UserService
}

func NewUsersController(service *Services.UserService) *UserController {
	return &UserController{Service: service}
}

func (controller *UserController) CreateUser(ctx echo.Context) error {
	var user Models.Users

	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	err := controller.Service.CreateUser(&user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error(), "hai": "hai"})
	}

	return ctx.JSON(http.StatusCreated, user)
}

func (controller *UserController) Login(ctx echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	token, err := controller.Service.Login(req.Username, req.Password)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
