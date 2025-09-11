package controllers

import (
	"github.com/gin-gonic/gin"
	"samsungvoicebe/config"
	"samsungvoicebe/models"
	"samsungvoicebe/services"
)

type UserController struct {
	Config  *config.Config
	Service *services.UserService
}

func NewUserController(cfg *config.Config, service *services.UserService) *UserController {
	return &UserController{
		Config:  cfg,
		Service: service,
	}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	err := uc.Service.CreateUser(req.UserID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User created successfully"})
}
