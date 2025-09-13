package routes

import (
	"github.com/gin-gonic/gin"
	"samsungvoicebe/config"
	"samsungvoicebe/controllers"
	"samsungvoicebe/services"
)

func GameplayRoutes(router *gin.RouterGroup, cfg *config.Config, service *services.GameplayService) {
	gameplayController := controllers.NewGameplayController(cfg, service)

	router.POST("/game/:game_id/move", gameplayController.PlayerMove)
	router.POST("/:user_id/game", gameplayController.CreateGame)
	router.POST("/hint", gameplayController.GetHint)
	router.POST("/move-by-voice", gameplayController.PlayerMoveByVoiceTranscription)
	router.POST("/game/move", gameplayController.PlayerMove)

}
