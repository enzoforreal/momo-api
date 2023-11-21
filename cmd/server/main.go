package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/enzoforreal/momo-api/internal/api"
	"github.com/enzoforreal/momo-api/internal/config"
	"github.com/enzoforreal/momo-api/internal/momo"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

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
