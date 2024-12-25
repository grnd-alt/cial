package middleware

import (
	"errors"
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

func ProtectedMiddleware(verifier *oidc.IDTokenVerifier) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearer, err := getBearer(ctx.GetHeader("Authorization"))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		token, err := verifier.Verify(ctx, bearer)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		var claims Claims
		token.Claims(&claims)
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
