package main

import (
	"fmt"
	"log"
	"net/http"

	"bytes"
	"encoding/json"
	"gw/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Define a route for the gateway
	router.POST("/reg", func(ctx *gin.Context) {
		// Forward the request to the downstream service

		var reg models.RegistrationBody

		if err := ctx.BindJSON(&reg); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"Invalid Json Data": err.Error()})
			return
		}

		fmt.Println(reg)

		reqBody, err := json.Marshal(reg)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Forward the registration request to the user service
		regUrl := "http://localhost:2705/api/v1/user/register"
		resp, err := http.Post(regUrl, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()
	})

	router.POST("/login", func(ctx *gin.Context) {
		// Forward the request to the downstream service

		var login models.LoginBody

		if err := ctx.BindJSON(&login); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"Invalid Json Data": err.Error()})
			return
		}

		fmt.Println(login)

		reqBody, err := json.Marshal(login)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Forward the registration request to the user service
		loginUrl := "http://localhost:2705/api/v1/user/login"
		resp, err := http.Post(loginUrl, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()
	})

	// Run the HTTP server
	if err := router.Run(":2701"); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
