package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadDotEnv() error {
	path := os.Getenv("DOTENV_PATH")
	if path == "" {
		path = "./.env.local"
	}

	err := godotenv.Load(path)
	if err != nil {
		log.Println(fmt.Sprintf("%v: [%v] %v", "error", "loadDotEnv", err.Error()))
		return err
	}
	return nil
}
