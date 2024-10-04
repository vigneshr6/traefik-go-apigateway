package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("1234")

const Authorization = "Authorization"
const Bearer = "Bearer"

func extractToken(headers http.Header) (string, error) {
	authHeader, ok := headers[Authorization]
	if !ok {
		return "", fmt.Errorf("token not found")
	}
	auth := authHeader[0]
	if !strings.Contains(auth, Bearer) {
		return "", fmt.Errorf("token not found")
	}
	auth = auth[len(Bearer)+1:]
	log.Println("Found token : ", auth)
	return auth, nil
}

func validate(tokenString string) error {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check for verification errors
	if err != nil {
		log.Println("Error parsing existing token")
		return err
	}

	// Check if the token is valid
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	log.Printf("Token verified successfully. Claims: %+v\\n", token.Claims)
	// Return the verified token
	return nil
}

func createToken(ctx *gin.Context) error {
	username, _, ok := ctx.Request.BasicAuth()
	if !ok {
		return fmt.Errorf("basic auth not found")
	}
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"iss": "todo-app",                       // Issuer
		"aud": "user",                           // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		log.Println("Error signing token")
		return err
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
	log.Println("Token created successfully")
	return nil
}

func validateRequest(ctx *gin.Context) error {
	token, err := extractToken(ctx.Request.Header)
	if err != nil {
		return err
	}
	err = validate(token)
	if err != nil {
		return err
	}
	return nil
}
