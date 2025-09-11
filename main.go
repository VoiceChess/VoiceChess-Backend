package main

import (
	"log"

	"samsungvoicebe/config"
	"samsungvoicebe/db"
	"samsungvoicebe/middleware"
	"samsungvoicebe/repo"
	"samsungvoicebe/routes"
	"samsungvoicebe/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.IsValid() {
		log.Println("✅ Environment loaded successfully")
	} else {
		log.Fatal("❌ GEMINI_API_KEY not configured")
	}

	database, err := db.New()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer database.Close()

	log.Println("✅ Database connected successfully")

	gameplayRepo := repo.NewGameplayRepo(database)
	analysisRepo := repo.NewAnalysisRepo(database)
	userRepo := repo.NewUserRepo(database)

	analysisService := services.NewAnalysisService(analysisRepo)
	gameplayService := services.NewGameplayService(gameplayRepo, analysisService)
	userService := services.NewUserService(userRepo)

	gin.SetMode(cfg.GinMode)

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Samsung Voice Backend API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	chatApi := r.Group("/api/chat")
	routes.ChatRoutes(chatApi, cfg)

	chessApi := r.Group("/api/chess")
	routes.ChessRoutes(chessApi, cfg)

	gameplayApi := r.Group("/api/gameplay")
	routes.GameplayRoutes(gameplayApi, cfg, gameplayService)

	analysisApi := r.Group("/api/analysis")
	routes.AnalysisRoutes(analysisApi, cfg, analysisService)

	userApi := r.Group("/api/user")
	routes.UserRoutes(userApi, cfg, userService)

	log.Printf("Base URL: http://localhost:%s/\n", cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
