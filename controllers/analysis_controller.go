package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"samsungvoicebe/config"
	"samsungvoicebe/middleware"
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
	firebaseUID, _ := c.Get(middleware.FirebaseUIDKey)
	if uid, ok := firebaseUID.(string); ok && uid != userID {
		c.JSON(403, gin.H{"error": "user_id does not match authenticated user"})
		return
	}

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

	firebaseUID, _ := c.Get(middleware.FirebaseUIDKey)
	userID, ok := firebaseUID.(string)
	if !ok || userID == "" {
		c.JSON(401, gin.H{"error": "missing authenticated user"})
		return
	}

	analyzedMove, err := ac.Service.GetAnalyzedMoveByOrder(moveOrder, gameID, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": analyzedMove})
}

func (ac *AnalysisController) GetFenFromPicture(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "Image file is required"})
		return
	}

	imageFile, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to open image file"})
		return
	}
	defer imageFile.Close()

	imageData := make([]byte, file.Size)
	_, err = imageFile.Read(imageData)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read image file"})
		return
	}

	fen, err := ac.Service.GetFenFromPicture(imageData)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": fen})
}
