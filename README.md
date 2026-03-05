# Redis Learning with Go - Phase 2

## Phase 2: CRUD Operations with JSON

Building on Phase 1's basic Redis connection, this phase implements full CRUD operations using Redis as a document store with JSON.

## Features

* Store Go structs as JSON in Redis
* Complete CRUD operations:

  * **Create**: Store person records with JSON marshaling
  * **Read**: Retrieve and parse JSON back to structs
  * **Update**: Modify existing records
  * **Delete**: Remove records
* List all records using Redis key patterns
* Proper error handling with contextual messages
* Existence checks before update and delete operations

## Code Structure

```go
// Person struct stored in Redis
type Person struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Age       int       `json:"age"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// CRUD Helper Functions
createPerson(ctx, client, person) // CREATE
getPerson(ctx, client, id)        // READ
updatePerson(ctx, client, person) // UPDATE
deletePerson(ctx, client, id)     // DELETE
listAllPeople(ctx, client)        // LIST
```

## How to Run

### Start Redis using Docker

If Redis is already created:

```bash
docker start my-redis
```

Or create a new Redis container:

```bash
docker run --name my-redis -p 6379:6379 -d redis:alpine
```

### Run the program

```bash
go run main.go
```

## Expected Output

```
Connected to Redis: PONG

CREATE: Adding a new person...
Created: {ID:user:1001 Name:John Doe Email:john@example.com Age:30 ...}

READ: Retrieving person...
Retrieved: {ID:user:1001 Name:John Doe Email:john@example.com Age:30 ...}

UPDATE: Updating person...
Updated: {ID:user:1001 Name:John Doe Email:john.doe@newemail.com Age:31 ...}

READ: Verifying update...
After update: {ID:user:1001 Name:John Doe Email:john.doe@newemail.com Age:31 ...}

CREATE: Adding another person...
Created: {ID:user:1002 Name:Jane Smith Email:jane@example.com Age:25 ...}

LIST: All people in Redis...
1. John Doe (john.doe@newemail.com)
2. Jane Smith (jane@example.com)

DELETE: Removing person with ID user:1002...
Deleted successfully

VERIFY: Trying to get deleted person...
Person not found (correctly deleted)

VERIFY: Checking if user:1001 still exists...
user:1001 still exists

All CRUD operations completed successfully.
```

## Key Concepts Learned

### 1. JSON Marshaling and Unmarshaling

```go
// Convert struct -> JSON for storage
jsonData, _ := json.Marshal(person)

// Convert JSON -> struct for usage
var person Person
json.Unmarshal(jsonData, &person)
```

### 2. Redis Key Patterns

* Using `user:1001`, `user:1002` as keys for grouping related records
* Using `Keys(ctx, "user:*")` to retrieve all user-related keys

### 3. Redis Commands Used

* **SET** – Store JSON data
* **GET** – Retrieve JSON data
* **DEL** – Delete records
* **EXISTS** – Check if a key exists
* **KEYS** – Find keys matching a pattern

### 4. Error Handling Patterns

* `redis.Nil` for missing keys
* `fmt.Errorf` for adding contextual error information
* `log.Fatalf` only used in `main()` for fatal errors

## File Structure

```
redis-learn/
├── main.go
├── go.mod
├── go.sum
└── README.md
```

## Dependencies

* go-redis/redis/v8 – Redis client for Go
* Redis server running via Docker

## Next Steps

Possible Phase 3 improvements:

* Add HTTP API layer with REST endpoints
* Implement Redis transactions
* Add data validation
* Use Redis Hash data type instead of JSON
* Add unit tests
* Implement connection pooling

## Key Takeaways

* Redis can act as a document store by storing JSON data.
* CRUD patterns can be implemented on top of a key-value store.
* Proper error propagation improves maintainability.
* Struct tags like `json:"fieldname"` control how Go structs are serialized.

## Learning Progress

```
Phase 1: Basic Connection
Phase 2: CRUD Operations
```