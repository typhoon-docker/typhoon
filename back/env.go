package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() {
	env := os.Getenv("GO_ENV")

	if "" == env {
		env = "development"
	}
	log.Println("running in " + env + " mode")

	// Load .env
	godotenv.Load()

	// Load .env.{GO_ENV}
	godotenv.Load(".env." + env)

	// Load .env.local then .env.{GO_ENV}.local if GO_ENV != test
	if "test" != env {
		godotenv.Load(".env.local")
		godotenv.Load(".env." + env + ".local")
	}
}
