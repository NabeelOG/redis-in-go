package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"redis-learn/models"
	redisDB "redis-learn/redis"
)

// Handler struct holds the Redis client — handlers need it to talk to Redis
type Handler struct {
	RedisClient *redis.Client
}

// Constructor — called once in main.go when wiring things up
func NewHandler(client *redis.Client) *Handler {
	return &Handler{RedisClient: client}
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

	if err := redisDB.CreatePerson(redisDB.Ctx, h.RedisClient, person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create person"})
		return
	}

	c.JSON(http.StatusCreated, person)
}

// GET /items/:id
// Reads :id from the URL → fetches from Redis → returns the person
func (h *Handler) GetPerson(c *gin.Context) {
	id := c.Param("id") // extracts the :id segment from the URL path

	person, err := redisDB.GetPerson(redisDB.Ctx, h.RedisClient, id)
	if err == redis.Nil {
		// redis.Nil means the key doesn't exist — that's a 404, not a server error
		c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve person"})
		return
	}

	c.JSON(http.StatusOK, person)
}

// PUT /items/:id
// Reads :id from URL + new data from body → updates Redis → returns updated person
func (h *Handler) UpdatePerson(c *gin.Context) {
	id := c.Param("id")

	// First check the person actually exists before trying to update
	existing, err := redisDB.GetPerson(redisDB.Ctx, h.RedisClient, id)
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve person"})
		return
	}

	// Bind the incoming changes to a temporary struct
	var updates models.Person
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Merge: only overwrite fields that were actually sent
	// ID and CreatedAt are preserved from the existing record
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

	if err := redisDB.UpdatePerson(redisDB.Ctx, h.RedisClient, existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update person"})
		return
	}

	c.JSON(http.StatusOK, existing)
}

// DELETE /items/:id
// Reads :id from URL → deletes from Redis → returns 204 No Content
func (h *Handler) DeletePerson(c *gin.Context) {
	id := c.Param("id")

	err := redisDB.DeletePerson(redisDB.Ctx, h.RedisClient, id)
	if err != nil && err.Error() == "person with ID "+id+" does not exist" {
		c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete person"})
		return
	}

	// 204 means success but no body to return
	c.Status(http.StatusNoContent)
}
