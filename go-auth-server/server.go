package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Any("/login", func(ctx *gin.Context) {
		err := createToken(ctx)
		if err != nil {
			log.Println("Error authenticating request : ", err)
			ctx.Status(http.StatusUnauthorized)
		}
	})

	router.Any("/", func(ctx *gin.Context) {
		log.Println("Authentication request")
		log.Printf("Headers : %v\n", ctx.Request.Header)
		err := validateRequest(ctx)
		if err != nil {
			log.Println("Error authenticating request : ", err)
			ctx.Status(http.StatusUnauthorized)
		}
	})

	log.Println("Starting server at port 8080")
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
