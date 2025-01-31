package services

import (
	"backendsetup/m/config"
	"context"
	"fmt"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
)

func InitOIDC(conf *config.Config) *oidc.Provider {
	provider,err := oidc.NewProvider(context.Background(), conf.OIDCIssuer)
	for err != nil {
		time.Sleep(5 * time.Second)
		fmt.Println("retrying keycloak connection")
		fmt.Println(err)
		provider, err = oidc.NewProvider(context.Background(), conf.OIDCIssuer)
	}
	return provider
}
