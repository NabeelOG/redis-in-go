package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"redis-learn/handlers"
	redisDB "redis-learn/redis"
)

func main() {
	// 1. Connect to Redis
	client := redisDB.NewClient()
	defer client.Close()

	// 2. Create handler (injects Redis client into all handlers)
	h := handlers.NewHandler(client)

	// 3. Create Gin router
	r := gin.Default()

	// 4. Register routes
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
