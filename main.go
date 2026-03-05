package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Create a context
	ctx := context.Background()

	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	// Test the connection
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Printf("Connected to Redis: %s\n", pong)

	// Set a key-value pair
	err = client.Set(ctx, "name", "Gopher", 0).Err()
	if err != nil {
		log.Fatalf("Could not set key: %v", err)
	}
	fmt.Println("Key 'name' set successfully")

	// Get the value
	val, err := client.Get(ctx, "name").Result()
	if err != nil {
		log.Fatalf("Could not get key: %v", err)
	}
	fmt.Printf("Value for key 'name': %s\n", val)

	// Try to get a non-existent key
	val2, err := client.Get(ctx, "nonexistent").Result()
	if err == redis.Nil {
		fmt.Println("Key 'nonexistent' does not exist")
	} else if err != nil {
		log.Fatalf("Error getting key: %v", err)
	} else {
		fmt.Printf("Value for key 'nonexistent': %s\n", val2)
	}

	// Close the connection when done
	defer client.Close()

	fmt.Println("\n✨ All operations completed successfully!")
}
