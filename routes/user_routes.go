package routes

import (
	"github.com/gin-gonic/gin"
	"samsungvoicebe/config"
	"samsungvoicebe/controllers"
	"samsungvoicebe/services"
)

func UserRoutes(router *gin.RouterGroup, cfg *config.Config, service *services.UserService) {
	userController := controllers.NewUserController(cfg, service)

	router.POST("/", userController.CreateUser)
}
