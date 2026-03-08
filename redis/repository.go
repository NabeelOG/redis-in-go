package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"

	"redis-learn/models"
)

// CREATE
func CreatePerson(ctx context.Context, client *redis.Client, person models.Person) error {
	jsonData, err := json.Marshal(person)
	if err != nil {
		return fmt.Errorf("failed to marshal person: %v", err)
	}

	err = client.Set(ctx, person.ID, jsonData, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to store in Redis: %v", err)
	}

	return nil
}

// READ
func GetPerson(ctx context.Context, client *redis.Client, id string) (*models.Person, error) {
	jsonData, err := client.Get(ctx, id).Result()
	if err != nil {
		return nil, err // caller checks for redis.Nil
	}

	var person models.Person
	err = json.Unmarshal([]byte(jsonData), &person)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal person: %v", err)
	}

	return &person, nil
}

// UPDATE
func UpdatePerson(ctx context.Context, client *redis.Client, person *models.Person) error {
	exists, err := client.Exists(ctx, person.ID).Result()
	if err != nil {
		return fmt.Errorf("failed to check existence: %v", err)
	}
	if exists == 0 {
		return fmt.Errorf("person with ID %s does not exist", person.ID)
	}

	// Set with same key overwrites — reuse CreatePerson
	return CreatePerson(ctx, client, *person)
}

// DELETE
func DeletePerson(ctx context.Context, client *redis.Client, id string) error {
	exists, err := client.Exists(ctx, id).Result()
	if err != nil {
		return fmt.Errorf("failed to check existence: %v", err)
	}
	if exists == 0 {
		return fmt.Errorf("person with ID %s does not exist", id)
	}

	err = client.Del(ctx, id).Err()
	if err != nil {
		return fmt.Errorf("failed to delete: %v", err)
	}

	return nil
}

// LIST ALL
func ListAllPeople(ctx context.Context, client *redis.Client) ([]*models.Person, error) {
	keys, err := client.Keys(ctx, "user:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %v", err)
	}

	var people []*models.Person
	for _, key := range keys {
		person, err := GetPerson(ctx, client, key)
		if err != nil {
			fmt.Printf("Warning: could not parse key %s: %v\n", key, err)
			continue
		}
		people = append(people, person)
	}

	return people, nil
}
