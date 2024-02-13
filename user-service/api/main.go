package main

import (
	"context"
	"fmt"
	"log"
	"user-service/shared"
	"user-service/workflow"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

func main() {
	r := gin.Default()

	// Initialize Temporal client
	c, err := client.Dial(client.Options{
		HostPort: "192.168.1.224:7233",
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	

	// Endpoint to start the workflow
	r.POST("/api/v1/user/register", func(cttx *gin.Context) {
		// Extract request parameters
		var reg shared.RegistrationBody

		if err := cttx.BindJSON(&reg); err != nil {
			cttx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		WorkflowID := "busOrderWorkflow" + uuid.New().String()

		// Start the BusOrderWorkflow using the Temporal client
		workflowOptions := client.StartWorkflowOptions{
			ID:        WorkflowID,               // Provide a unique workflow ID
			TaskQueue: shared.UserRegistrationTaskQueue, // Provide the appropriate task queue name
		}

		workflowFuture, _ := c.ExecuteWorkflow(context.Background(), workflowOptions, workflow.UserRegWorkflow, reg)
		
		fmt.Println(workflowFuture.Get(context.Background(), nil))

		// var PerformOrderResponse shared.PerformOrderResponse // Assuming WorkflowResponse is a struct type
		// err = workflowFuture.Get(context.Background(), &PerformOrderResponse)
		// fmt.Println(err)
		// if err != nil {
		// 	// Handle the error
		// 	// For example, log the error or return it as an HTTP response
		// 	cttx.JSON(500, gin.H{"error": err.Error()})
		// 	return
		// }
		
		// fmt.Println(PerformOrderResponse)

		// // Return workflow execution ID
		// cttx.JSON(200, gin.H{"OrderID": PerformOrderResponse.OrderID})
	})

	r.POST("/bus/payment", func(cttx *gin.Context) {

		var cnfSignal shared.ConfirmationSignal

		if err := cttx.BindJSON(&cnfSignal); err != nil {
			cttx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		// Handle payment logic here for GET request
		// You can start a new workflow or perform any other actions related to payment processing
		// For demonstration, we just return a success message
		signalName := "emailConfirmation"
		err = c.SignalWorkflow(context.Background(), cnfSignal.WorkflowID, "", signalName, cnfSignal.ConfirmationCode)
		if err != nil {
			cttx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		cttx.JSON(200, gin.H{"message": "Conirmation code sent successfully"})
	})

	// Run the server
	if err := r.Run(":2705"); err != nil {
		panic(err)
	}
}
