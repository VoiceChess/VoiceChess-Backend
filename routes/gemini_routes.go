package routes

import (
	"samsungvoicebe/config"
	"samsungvoicebe/controllers"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(router *gin.RouterGroup, cfg *config.Config) {
	chatControllerV2 := controllers.NewChatControllerV2(cfg)
	router.POST("/gemini", chatControllerV2.ChatWithAI)
}
