package routes

import (
	"github.com/gin-gonic/gin"
	"samsungvoicebe/config"
	"samsungvoicebe/controllers"
	"samsungvoicebe/services"
)

func AnalysisRoutes(router *gin.RouterGroup, cfg *config.Config, service *services.AnalysisService) {
	analysisController := controllers.NewAnalysisController(cfg, service)

	router.GET("/:user_id/games", analysisController.GetGameHistoryList)
	router.GET("/game/:game_id/move/:move_order", analysisController.GetAnalyzedMoveByOrder)
}
