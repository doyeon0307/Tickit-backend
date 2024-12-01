package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	AWS_ACCESS_KEY string
	AWS_SECRET_KEY string
	JWT_SECRET_KEY string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AWS_ACCESS_KEY = os.Getenv("AWS_ACCESS_KEY")
	AWS_SECRET_KEY = os.Getenv("AWS_SECRET_KEY")
	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")

	fmt.Println("AWS_ACCESS_KEY:", AWS_ACCESS_KEY)
	fmt.Println("AWS_SECRET_KEY:", AWS_SECRET_KEY)
}
