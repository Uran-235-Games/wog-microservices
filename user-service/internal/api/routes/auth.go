package routers

import (
	"wog-server/internal/api/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, h *handler.AuthHandler) {
	auth := rg.Group("/auth")
	auth.POST("/sign-up", h.Register)
	auth.POST("/sign-in", h.Login)
}
