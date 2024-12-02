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
	_ = godotenv.Load()

	// 환경변수 읽기
	AWS_ACCESS_KEY = os.Getenv("AWS_ACCESS_KEY")
	AWS_SECRET_KEY = os.Getenv("AWS_SECRET_KEY")
	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")

	if AWS_ACCESS_KEY == "" || AWS_SECRET_KEY == "" || JWT_SECRET_KEY == "" {
		log.Fatal("Required environment variables are not set")
	}

	fmt.Println("AWS_ACCESS_KEY:", AWS_ACCESS_KEY)
	fmt.Println("AWS_SECRET_KEY:", AWS_SECRET_KEY)
}
