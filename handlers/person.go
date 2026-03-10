package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"redis-learn/models"
	"redis-learn/postgres"
	redisDB "redis-learn/redis"
)

// Handler struct holds the Redis client — handlers need it to talk to Redis
type Handler struct {
	RedisClient *redis.Client
	DB          *gorm.DB
}

// Constructor — called once in main.go when wiring things up
func NewHandler(client *redis.Client, db *gorm.DB) *Handler {
	return &Handler{
		RedisClient: client,
		DB:          db,
	}
}

// POST /items
// Reads JSON body → saves to Redis → returns the created person
func (h *Handler) CreatePerson(c *gin.Context) {
	var person models.Person

	// ShouldBindJSON reads the request body and maps it to the Person struct
	// If the body is malformed JSON, it returns an error
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Validate required fields
	if person.ID == "" || person.Name == "" || person.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id, name, and email are required"})
		return
	}

	// Stamp timestamps
	person.CreatedAt = time.Now()
	person.UpdatedAt = time.Now()

	//save to postgre
	if err := postgres.CreatePerson(h.DB, person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create a person"})
		return
	}

	//save to redis
	if err := redisDB.CreatePerson(redisDB.Ctx, h.RedisClient, person); err != nil {
		// dont fail the request if Redis fails
		// PostgreSQL has the data, Redis  is just cache
		fmt.Printf("Warning: failed to cache person in Redis: %v\n", err)
	}

	c.JSON(http.StatusCreated, person)
}

// GET /items/:id
// Reads :id from the URL → fetches from Redis → returns the person
func (h *Handler) GetPerson(c *gin.Context) {
	id := c.Param("id") // extracts the :id segment from the URL path

	//check redis first
	person, err := redisDB.GetPerson(redisDB.Ctx, h.RedisClient, id)
	if err == nil {
		// found in Redis, return Immediately
		fmt.Println("Cache HIT - Returning from redis")
		c.JSON(http.StatusOK, person)
		return
	}

	// 2. Not in Redis (cache miss), go to PostgreSQL
	fmt.Println("Cache MISS - going to PostgreSQL")
	person, err = postgres.GetPerson(h.DB, id)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve person"})
		return
	}

	// store in redis for the next time
	if err := redisDB.CreatePerson(redisDB.Ctx, h.RedisClient, *person); err != nil {
		fmt.Printf("Warning: failed to cache person in Redis: %v\n", err)
	}

	c.JSON(http.StatusOK, person)
}

// PUT /items/:id
// Reads :id from URL + new data from body → updates Redis → returns updated person
func (h *Handler) UpdatePerson(c *gin.Context) {
	id := c.Param("id")

	// 1. Check if person exists in PostgreSQL
	existing, err := postgres.GetPerson(h.DB, id)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve person"})
		return
	}

	// 2. Bind incoming changes
	var updates models.Person
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to retieve person"})
		return
	}

	// Merge changes

	if updates.Name != "" {
		existing.Name = updates.Name
	}
	if updates.Email != "" {
		existing.Email = updates.Email
	}
	if updates.Age != 0 {
		existing.Age = updates.Age
	}
	existing.UpdatedAt = time.Now()

	// update postgres first (source of truth)
	if err := postgres.UpdatePerson(h.DB, existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update person"})
		return
	}

	// update Redis cache
	if err := redisDB.CreatePerson(redisDB.Ctx, h.RedisClient, *existing); err != nil {
		fmt.Printf("Warning: failed to update cache: %v\n", err)
	}

	c.JSON(http.StatusOK, existing)
}

// DELETE /items/:id
// Reads :id from URL → deletes from Redis → returns 204 No Content
func (h *Handler) DeletePerson(c *gin.Context) {
	id := c.Param("id")

	// 1. Delete from PostgreSQL first(source of truth)
	if err := postgres.DeletePerson(h.DB, id); err != nil {
		if err.Error() == "person with ID "+id+" does not exist" {
			c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete person"})
		return
	}

	// 2. Delete from Redis cache
	if err := redisDB.DeletePerson(redisDB.Ctx, h.RedisClient, id); err != nil {
		fmt.Printf("Warning: failed to delete from cache: %v\n", err)
	}

	c.Status(http.StatusNoContent)
}

// GET /items
// fetches from reddis and return the user
func (h *Handler) GetAll(c *gin.Context) {

	// 1. Try Redis First
	people, err := redisDB.ListAllPeople(redisDB.Ctx, h.RedisClient)
	if err == nil && len(people) > 0 {
		fmt.Println("Cache HIT - returning from Redis")
		c.JSON(http.StatusOK, people)
		return
	}

	// 2.Cache miss - go to PostgrSQL
	fmt.Println("Cache MISS - going to Postgre")
	people, err = postgres.ListAllPeople(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve people"})
		return
	}

	if len(people) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no people found"})
		return
	}

	// 3. Store each person in redis for the next time
	for _, person := range people {
		if err := redisDB.CreatePerson(redisDB.Ctx, h.RedisClient, *person); err != nil {
			fmt.Printf("Warning: failed to cache person %s: %v\n", person.ID, err)
		}
	}

	c.JSON(http.StatusOK, people)
}
