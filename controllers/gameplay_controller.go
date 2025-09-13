package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"samsungvoicebe/config"
	"samsungvoicebe/models"
	"samsungvoicebe/services"
)

type GameplayController struct {
	Config  *config.Config
	Service *services.GameplayService
}

func NewGameplayController(cfg *config.Config, service *services.GameplayService) *GameplayController {
	return &GameplayController{
		Config:  cfg,
		Service: service,
	}
}

func (gc *GameplayController) PlayerMove(c *gin.Context) {
	gameID := c.Param("game_id")

	var req models.PlayerMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("GameplayController-PlayerMove-JsonBinding", err)
		return
	}

	var botMove models.BotMove
	var err error

	if gameID != "" {
		botMove, err = gc.Service.PlayerMove(&gameID, req.Fen, req.Move, req.BotLevel)
	} else {
		botMove, err = gc.Service.PlayerMove(nil, req.Fen, req.Move, req.BotLevel)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("GameplayController-PlayerMove-PlayerMove", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": botMove,
	})
}

func (gc *GameplayController) CreateGame(c *gin.Context) {
	userID := c.Param("user_id")

	gameID, err := gc.Service.CreateGame(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("GameplayController-CreateGame-CreateGame", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"game_id": gameID,
		},
	})
}

func (gc *GameplayController) GetHint(c *gin.Context) {
	var req models.HintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("GameplayController-GetHint-JsonBinding", err)
		return
	}

	hint, err := gc.Service.GetHint(req.Fen)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("GameplayController-GetHint-GetHint", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"hint": hint,
		},
	})
}

func (gc *GameplayController) PlayerMoveByVoiceTranscription(c *gin.Context) {
	var req models.PlayerMoveByTranscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("GameplayController-PlayerMoveByVoiceTranscription-JsonBinding", err)
		return
	}

	playerMove, err := gc.Service.PlayerMoveByVoiceTranscription(req.Fen, req.Transcription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("GameplayController-PlayerMoveByVoiceTranscription-PlayerMoveByVoiceTranscription", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": playerMove,
	})
}
