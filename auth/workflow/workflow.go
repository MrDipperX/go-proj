package workflow

import (
	"auth/shared"
	"auth/activity"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func UserAuthWorkflow(ctx workflow.Context, userData shared.UserData) error {
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    3,
	}

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		RetryPolicy:            retryPolicy,
	}

	// Get the logger for the workflow
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting user authentication workflow")

	// Set the activity options
	ctx = workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx, activity.GenerateTokenAndWriteToRedis, userData).Get(ctx, nil)
	if err != nil {
		logger.Error("Error on add user data to Redis")
		return err
	}

	logger.Info("User authentication workflow completed successfully")

	// Return the JWT token
	return nil
}