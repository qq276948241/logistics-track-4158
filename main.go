package main

import (
	"log"
	"logistics-tracker/config"
	"logistics-tracker/handlers"
	"logistics-tracker/middleware"
	"logistics-tracker/models"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	models.InitDB(cfg)

	middleware.InitAuth(cfg)

	r := gin.Default()

	r.POST("/api/login", handlers.Login)

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/waybills", handlers.CreateWaybill)
		api.GET("/waybills/:tracking_number/trackings", handlers.GetWaybillTrackings)
	}

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
