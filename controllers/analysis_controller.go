package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"samsungvoicebe/config"
	"samsungvoicebe/services"
)

type AnalysisController struct {
	Config  *config.Config
	Service *services.AnalysisService
}

func NewAnalysisController(cfg *config.Config, service *services.AnalysisService) *AnalysisController {
	return &AnalysisController{
		Config:  cfg,
		Service: service,
	}
}

func (ac *AnalysisController) GetGameHistoryList(c *gin.Context) {
	userID := c.Param("user_id")
	gamesHistoryList, err := ac.Service.GetGameHistoryList(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": gamesHistoryList})
}

func (ac *AnalysisController) GetAnalyzedMoveByOrder(c *gin.Context) {
	gameID := c.Param("game_id")
	moveOrderParam := c.Param("move_order")

	moveOrder, err := strconv.Atoi(moveOrderParam)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid move order"})
		return
	}

	analyzedMove, err := ac.Service.GetAnalyzedMoveByOrder(moveOrder, gameID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": analyzedMove})
}
