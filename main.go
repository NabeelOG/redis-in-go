package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"redis-learn/config"
	"redis-learn/handlers"
	"redis-learn/postgres"
	redisDB "redis-learn/redis"
)

func main() {

	// 1. Load Config
	cfg := config.Load()

	// 2. Connect to Redis
	redisClient := redisDB.NewClient()
	defer redisClient.Close()

	// 3. Connect to PostgreSQL
	pgClient := postgres.NewClient(cfg)

	// 4. Create handler (injects Redis client into all handlers)
	h := handlers.NewHandler(redisClient, pgClient)

	// 5. Create Gin router
	r := gin.Default()

	// 6. Register routes
	r.POST("/items", h.CreatePerson)
	r.GET("/items/:id", h.GetPerson)
	r.PUT("/items/:id", h.UpdatePerson)
	r.DELETE("/items/:id", h.DeletePerson)
	r.GET("/allitems", h.GetAll)

	// 5. Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
