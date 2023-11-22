package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/enzoforreal/momo-api/internal/api"
	"github.com/enzoforreal/momo-api/internal/config"
	"github.com/enzoforreal/momo-api/internal/momo"
	"github.com/enzoforreal/momo-api/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {

	logger.Init()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	fmt.Println("Token URL:", cfg.Momo.TokenURL)

	fmt.Println("Environment Variables:")
	fmt.Println("MOMO_CONSUMER_KEY:", os.Getenv("MOMO_CONSUMER_KEY"))
	fmt.Println("MOMO_CONSUMER_SECRET:", os.Getenv("MOMO_CONSUMER_SECRET"))
	fmt.Println("MOMO_TOKEN_URL:", os.Getenv("MOMO_TOKEN_URL"))
	fmt.Println("MOMO_CALLBACK_URL:", os.Getenv("MOMO_CALLBACK_URL"))
	fmt.Println("MOMO_API_ENDPOINT:", os.Getenv("MOMO_API_ENDPOINT"))

	momoClient := momo.NewClient(cfg)

	token, err := momoClient.GetOAuthToken()
	if err != nil {

		log.Fatal("Failed to get token: ", err)
	}

	log.Printf("Token: %s", token)

	server := &http.Server{
		Addr:    ":8080",
		Handler: api.NewRouter(),
	}

	fmt.Printf("Server listening on port %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
	}
}
