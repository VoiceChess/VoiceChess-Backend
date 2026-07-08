package main

import (
	"context"
	"log"

	"samsungvoicebe/config"
	"samsungvoicebe/db"
	"samsungvoicebe/firebaseauth"
	"samsungvoicebe/middleware"
	"samsungvoicebe/repo"
	"samsungvoicebe/routes"
	"samsungvoicebe/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	log.Println("✅ Environment loaded")

	// db
	database, err := db.New()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer database.Close()

	log.Println("✅ Database connected successfully")

	firebaseAuthClient, err := firebaseauth.NewClient(context.Background())
	if err != nil {
		log.Fatal("❌ Failed to initialize Firebase auth:", err)
	}
	log.Println("✅ Firebase auth initialized successfully")

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
			"message": "VoiceChess Backend API",
			"version": "1.0.1",
			"status":  "running",
		})
	})

	api := r.Group("/api", middleware.FirebaseAuth(firebaseAuthClient))

	chatApi := api.Group("/chat")
	routes.ChatRoutes(chatApi, cfg)

	chessApi := api.Group("/chess")
	routes.ChessRoutes(chessApi, cfg)

	gameplayApi := api.Group("/gameplay")
	routes.GameplayRoutes(gameplayApi, cfg, gameplayService)

	analysisApi := api.Group("/analysis")
	routes.AnalysisRoutes(analysisApi, cfg, analysisService)

	userApi := api.Group("/user")
	routes.UserRoutes(userApi, cfg, userService)

	log.Printf("Base URL: http://localhost:%s/\n", cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
