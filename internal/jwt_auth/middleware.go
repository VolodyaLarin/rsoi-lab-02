package jwt_auth

import (
	"context"
	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type CustomClaimsExample struct {
	name string
}

func (c CustomClaimsExample) Valid() error {
	return nil
}

func AuthMiddleware() gin.HandlerFunc {
	jwksURL := "https://dev-y3hp3rut2qnex3gw.us.auth0.com/.well-known/jwks.json"
	options := keyfunc.Options{
		Ctx: context.Background(),
		RefreshErrorHandler: func(err error) {
			log.Printf("There was an error with the jwt.Keyfunc\nError: %s", err.Error())
		},
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	}

	// Create the JWKS from the resource at the given URL.
	jwks, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		log.Fatalf("Failed to create JWKS from resource at the given URL.\nError: %s", err.Error())
	}

	return func(c *gin.Context) {
		tokenString := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 20)

		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(tokenString, claims, jwks.Keyfunc)
		if err != nil {
			log.Warning("Auth warning", err)
		}

		username, ok := claims["name"].(string)
		log.Warning(username, ok)

		if !ok {
			c.JSON(403, gin.H{})
			return
		}

		c.Set("x-username", username)

		c.Next()
	}

}
