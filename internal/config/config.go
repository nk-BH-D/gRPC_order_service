package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func Loader() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("%v, used default value", err)
	}

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	return &Config{Port: port}
}
