package handler

import (
	// "fmt"
	"log/slog"
	"wog-server/domain"
	"wog-server/internal/logger"
	"wog-server/internal/service/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (r *AuthHandler) Register(c *gin.Context) {
	res := domain.NewAPIResponse()

	var reqUser domain.SignUpRequest
	if err := c.ShouldBindJSON(&reqUser); err != nil {
		res.AddError(1, "Некорректные данные. "+"details: "+ err.Error())
		res.Success = false
		c.JSON(400, res)
		return
	}

	user := domain.User(reqUser)

	userID, err := r.authService.Register(user)
	if err != nil {
		logger.Log.Error("Ошибка при создании пользователя", slog.String("err", err.Error()))
		res.AddError(1, "Ошибка при создании пользователя")
		res.Success = false
		c.JSON(400, res)
		return
	}
	
	res.Result = domain.SignUpResponse{
		ID: userID,
		Name: user.Name,
		Email: user.Email,
	}
	c.JSON(200, res)
}


func (r *AuthHandler) Login(c *gin.Context) {
	res := domain.NewAPIResponse()

	var reqUser domain.SignInRequest
	if err := c.ShouldBindJSON(&reqUser); err != nil {
		logger.Log.Error("Некорректные данные при логине", slog.String("err", err.Error()))
		res.AddError(1, "Некорректные данные запроса")
		res.Success = false
		c.JSON(400, res)
		return
	}

	if (reqUser.Email == "" && reqUser.Name == "") {
		res.AddError(2, "Должно быть указано хотя бы 1 поле из name, email")
		res.Success = false
		c.JSON(400, res)
		return
	}

	user := domain.User(reqUser)

	resUser, err := r.authService.Login(user)
	if err != nil {
		res.AddError(1, err.Error())
		res.Success = false
		c.JSON(400, res)
		return
	}

	res.Result = resUser
	c.JSON(200, res)
}

