package middleware

import (
	"backendsetup/m/db/sql/dbgen"
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	Sub      string `json:"sub"`
	Username string `json:"preferred_username"`
	Email    string `json:"email"`
	Expiry   int64  `json:"exp"`
}

func getBearer(header string) (string, error) {
	vals := strings.Split(header, " ")
	if len(vals) != 2 {
		return "", errors.New("Invalid Authorization header")
	}
	if vals[0] != "Bearer" {
		return "", errors.New("Invalid Authorization header")
	}
	return vals[1], nil
}

func ProtectedMiddleware(verifier *oidc.IDTokenVerifier, mode string, queries *dbgen.Queries) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// if mode == "development" {
		// 	ctx.Set("claims", Claims{
		// 		Sub:      "test",
		// 		Username: "test",
		// 		Email:    "test@test.com",
		// 	})
		// 	return
		// }
		bearer, err := getBearer(ctx.GetHeader("Authorization"))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		token, err := verifier.Verify(ctx, bearer)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		var claims Claims
		err = token.Claims(&claims)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.Set("claims", claims)
		go func(userID string) {
			err := queries.SetLastLogin(context.Background(), userID)
			if err != nil {
				log.Printf("failed to update last_login for user %s: %v", userID, err)
			}
		}(claims.Sub)
		ctx.Next()
	}
}
