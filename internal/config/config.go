package config

import (
	"log"
	"os"

	env "github.com/joho/godotenv"
)

type Config struct {
	GRPC_Port string
	HTTP_Port string
}

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

	return &Config{
		GRPC_Port: grpc_port,
		HTTP_Port: http_port,
	}
}
