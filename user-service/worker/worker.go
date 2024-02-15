package main

import (
	"log"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"user-service/shared"
	"user-service/workflow"
	"user-service/activity"
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{
		HostPort: "192.168.1.224:7233",
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, shared.UserRegistrationTaskQueue, worker.Options{})
	w.RegisterWorkflow(workflow.UserRegWorkflow)
	w.RegisterActivity(activity.CreateUser)
	w.RegisterActivity(activity.SendGeneratedNumberToEmail)
	w.RegisterActivity(activity.CheckConfirmationCodeFromRedis)
	w.RegisterActivity(activity.AddUserDataToRedis)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}