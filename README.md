# Redis Learning with Go

## Phase 1: Basic Connection
- Simple Go program connecting to Redis in Docker
- Set and get key-value pairs
- Error handling for connections and missing keys

### How to run:
1. Start Redis: `docker run --name my-redis -p 6379:6379 -d redis:alpine`
2. Run: `go run main.go`

### What I learned:
- Connecting to Redis using go-redis
- Basic SET/GET operations
- Context usage
- Error handling with redis.Nil
