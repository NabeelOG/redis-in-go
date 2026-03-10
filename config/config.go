package config

type Config struct {
	RedisAddr  string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}

func Load() *Config {
	return &Config{
		RedisAddr:  "localhost:6379",
		DBHost:     "localhost",
		DBUser:     "nabeel",
		DBPassword: "nabeel",
		DBName:     "microservice_db",
		DBPort:     "5432",
	}
}
