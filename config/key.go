package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	AWS_ACCESS_KEY string
	AWS_SECRET_KEY string
	JWT_SECRET_KEY string
)

func init() {
	_ = godotenv.Load()

	AWS_ACCESS_KEY = os.Getenv("AWS_ACCESS_KEY")
	AWS_SECRET_KEY = os.Getenv("AWS_SECRET_KEY")
	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")

	fmt.Println("AWS_ACCESS_KEY:", AWS_ACCESS_KEY)
	fmt.Println("AWS_SECRET_KEY:", AWS_SECRET_KEY)
}
