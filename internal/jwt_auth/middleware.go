package jwt_auth

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"strings"
)

type CustomClaimsExample struct {
	name string
}

func (c CustomClaimsExample) Valid() error {
	return nil
}

func AuthMiddleware(c *gin.Context) {
	tokenString := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 20)

	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("LS6VIBs6bmAYkcwIwsDz2B1eW6512-GOCH_gXBkuWxeXQRSDTyFnl8UBGA4ZhepS"), nil
	})

	username, ok := claims["name"].(string)
	log.Warning(username, ok)

	if !ok {
		c.JSON(403, gin.H{})
		return
	}

	c.Set("x-username", username)

	c.Next()
}
