package main

import (
	"fmt"
	"log"
	"net/http"

	"backendsetup/m/config"
	"backendsetup/m/db"
	"backendsetup/m/routes"
	"backendsetup/m/services"

	"github.com/coreos/go-oidc/v3/oidc"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	queries := db.Init(config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName)

	oidcProvider := services.InitOIDC(config)
	verifier := oidcProvider.Verifier(&oidc.Config{ClientID: config.OIDCClientID})

	engine := routes.Init(verifier, config, queries)

	fmt.Printf("listening on: %s:%d\n", "0.0.0.0", config.AppPort)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "0.0.0.0", config.AppPort),
		Handler: engine,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
