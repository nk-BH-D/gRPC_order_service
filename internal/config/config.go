package config

import (
	"fmt"
	"log"
	"os"

	env "github.com/joho/godotenv"
)

type Config struct {
	GRPC_Port string
	HTTP_Port string
	DB_URL    string
}

// []
func Loader() *Config {
	err := env.Load()
	if err != nil {
		log.Printf("%v, used default value", err)
	}

	grpc_port := os.Getenv("GRPC_PORT")
	if grpc_port == "" {
		grpc_port = "50051"
	}

	http_port := os.Getenv("HTTP_PORT")
	if http_port == "" {
		http_port = "8080"
	}

	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		postgres_port := os.Getenv("POSTGRES_PORT")
		if postgres_port == "" {
			postgres_port = "5432"
		}
		db_url = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			"db",
			postgres_port,
			os.Getenv("POSTGRES_DB"),
		)
	}

	return &Config{
		GRPC_Port: grpc_port,
		HTTP_Port: http_port,
		DB_URL:    db_url,
	}
}
