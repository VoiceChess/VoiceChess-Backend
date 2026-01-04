package routes

import (
	"samsungvoicebe/config"
	"samsungvoicebe/controllers"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(router *gin.RouterGroup, cfg *config.Config) {
	// V1 - Gemini
	// chatController := controllers.NewChatController(cfg)
	// router.POST("/gemini", chatController.Chat)

	// V2 - Azure OpenAI
	chatControllerV2 := controllers.NewChatControllerV2(cfg)
	router.POST("/gemini", chatControllerV2.ChatWithAI)
}
