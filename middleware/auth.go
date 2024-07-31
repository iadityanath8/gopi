package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("vUHIf2JKIZRvTL1M9Yd4A/nvqUsD+jw9wTeb+bk8r5M=")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	Claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	return token.SignedString(jwtKey)
}

func extractTokenFromHeader(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	fmt.Println("getting the token", bearerToken)

	if len(bearerToken) > 7 && bearerToken[:7] == "Bearer " {
		return bearerToken[7:]
	}
	return ""
}

func VerifyToken(c *gin.Context) {
	tokenString := extractTokenFromHeader(c)
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization failed"})
		c.Abort()
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		// fmt.Println("what the hell just happended", tokenString)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	c.Set("user", claims.Username)
	c.Next()
}

func GetCurrentUser(c *gin.Context) string {
	user, _ := c.Get("user")
	if username, ok := user.(string); ok {
		return username
	}
	return ""
}
