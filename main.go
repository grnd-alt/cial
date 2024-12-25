package main

import (
	"backendsetup/m/config"
	"backendsetup/m/db"
	"backendsetup/m/routes"
	"backendsetup/m/services"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	queries := db.Init(config)

	oidcProvider := services.InitOIDC(config)
	verifier:= oidcProvider.Verifier(&oidc.Config{ClientID: config.OIDCClientID})


	engine := routes.Init(verifier, config, queries)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "0.0.0.0", config.AppPort),
		Handler: engine,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
