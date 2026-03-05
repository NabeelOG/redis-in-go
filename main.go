package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Person struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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

	//---DEMONSTRATE CRUD OPERATION---//

	//1. CREATE -- Store a person as JSON
	fmt.Println("CREATE: Adding a new person...")
	person1 := Person{
		ID:        "user:1001",
		Name:      "Nabeel Muhammed",
		Email:     "nmuhammed851@gmail.com",
		Age:       30,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = createPerson(ctx, client, person1)
	if err != nil {
		log.Fatalf("Failed to create person: %v", err)
	}
	fmt.Printf("Created: %+v\n\n", person1)

	//2. READ -- Retrieve the Person
	fmt.Println("READ: Retrieving person")
	retrieved, err := getPerson(ctx, client, "user:1001")
	if err != nil {
		log.Fatalf("Failed to get person: %v", err)
	}
	fmt.Printf("Retrieved %+v\n\n", retrieved)

	//3. UPDATE -- Modify the person
	fmt.Println("UPDATE: Updating person...")
	retrieved.Age = 31
	retrieved.Email = "nabeel@gmail.com"
	retrieved.UpdatedAt = time.Now()

	err = updatePerson(ctx, client, retrieved)
	if err != nil {
		log.Fatalf("Failed to update person: %v", err)
	}
	fmt.Printf("Updated %+v\n\n", retrieved)

	//4. READ again to Verify Updates
	fmt.Println("READ: verifying updates...")
	updated, err := getPerson(ctx, client, "user:1001")
	if err != nil {
		log.Fatalf("Failed to get person: %v", err)
	}
	fmt.Printf("After update: %+v\n\n", updated)

	//5. CREATE another person
	fmt.Println("CREATE: Adding another person...")
	person2 := Person{
		ID:        "user:1002",
		Name:      "Nihal",
		Email:     "nihal@gmail.com",
		Age:       17,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = createPerson(ctx, client, person2)
	if err != nil {
		log.Fatalf("Failed to create person: %v", err)
	}
	fmt.Printf("Created %+v\n\n", person2)

	//6. List all people (using key patterns)
	fmt.Println("LIST: All people in Redis...")
	people, err := listAttPeople(ctx, client)
	if err != nil {
		log.Fatalf("Failed to list all people: %v", err)
	}
	for i, p := range people {
		fmt.Printf(" %d. %s (%s)\n", i+1, p.Name, p.Email)
	}
	fmt.Println()

	//7. DELETE a person
	fmt.Println("DELETE: Removing person with ID user:1002...")
	err = deletePerson(ctx, client, "user:1002")
	if err != nil {
		log.Fatalf("Failed to delete person: %v", err)
	}
	fmt.Println("Deleted Successfully")

	//8. Try to get deleted user
	fmt.Println("\n VERIFY: Trying to get the deleted person")
	deleted, err := getPerson(ctx, client, "user:1002")
	if err == redis.Nil {
		fmt.Println("Person not found (correctly deleted)")
	} else if err != nil {
		log.Fatalf("Error: %v", err)
	} else {
		fmt.Printf("Person still exists: %+v\n", deleted)
	}

	//9. Check if user:1001 still exists
	fmt.Println("\nVERIFY: Checking if user:1001 still exists...")
	exisits, err := client.Exists(ctx, "user:1001").Result()
	if err != nil {
		log.Fatalf("Error checking existence: %v", err)
	}
	if exisits == 1 {
		fmt.Println("user:1001 still exists")
	} else {
		fmt.Println("user:1001 was deleted")
	}

	// Close the connection when done
	defer client.Close()

	fmt.Println("\n✨ All operations completed successfully!")
}

// CREATE PERSON FUNC
func createPerson(ctx context.Context, client *redis.Client, person Person) error {
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

// RETRIEVE PERSON FUNC
func getPerson(ctx context.Context, client *redis.Client, id string) (*Person, error) {
	jsonData, err := client.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}
	var person Person
	err = json.Unmarshal([]byte(jsonData), &person)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal person: %v", err)
	}

	return &person, nil
}

// UPDATE PERSON FUNC
func updatePerson(ctx context.Context, client *redis.Client, person *Person) error {
	exists, err := client.Exists(ctx, person.ID).Result()
	if err != nil {
		return fmt.Errorf("failed to check existence: %v", err)
	}
	if exists == 0 {
		return fmt.Errorf("person with ID %s does not exist", person.ID)
	}

	person.UpdatedAt = time.Now()
	//Reusing create function(Set with same key overwrites)
	return createPerson(ctx, client, *person)
}

// LIST ALL PERSON FUNC(using key pattern)
func listAttPeople(ctx context.Context, client *redis.Client) ([]*Person, error) {
	keys, err := client.Keys(ctx, "user:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys %v", err)
	}

	var people []*Person
	for _, key := range keys {
		person, err := getPerson(ctx, client, key)
		if err != nil {
			fmt.Printf("Warning: Could not parse person with key %s: %v\n", key, err)
			continue
		}
		people = append(people, person)
	}
	return people, nil
}

// DELETE: Remove a person form Redis
func deletePerson(ctx context.Context, client *redis.Client, id string) error {
	//check if exists
	exists, err := client.Exists(ctx, id).Result()
	if err != nil {
		return fmt.Errorf("failed to check existence: %v", err)
	}
	if exists == 0 {
		return fmt.Errorf("person with ID %s does not exist", id)
	}

	//Delete key
	err = client.Del(ctx, id).Err()
	if err != nil {
		return fmt.Errorf("failed to delete %v", err)
	}
	return nil
}
