package workflow

import (
	"user-service/activity"
	"user-service/shared"

	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func UserRegWorkflow(ctx workflow.Context, reg shared.RegistrationBody) (err error) {
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

	// Set the activity options
	ctx = workflow.WithActivityOptions(ctx, ao)

	var ConfirmCode string

	err = workflow.ExecuteActivity(ctx, activity.SendGeneratedNumberToEmail, reg).Get(ctx, &ConfirmCode)

	if err != nil {
		workflow.GetLogger(ctx).Error("SendGeneratedNumberToEmail Activity failed.", "Error", err)
		return err
	}

	// Redis write confirmation code
	redisUserData := shared.RedisUserData{
		Username:    reg.Login,
		ConfirmCode: ConfirmCode}

	err = workflow.ExecuteActivity(ctx, activity.AddUserDataToRedis, redisUserData).Get(ctx, nil)

	if err != nil {
		workflow.GetLogger(ctx).Error("ConfirmationCodeToRedis Activity failed.", "Error", err)
		return err
	}


	// Start a timer for n minutes
    timerFuture := workflow.NewTimer(ctx, time.Minute*3)
    workflow.GetLogger(ctx).Info("Start a timer")

    // Use selector to wait for timer expiration or payment completion
    selector := workflow.NewSelector(ctx)
    var timerStatus int8

	var signalData string

    selector.AddFuture(timerFuture, func(f workflow.Future) {
        // Timer expired, handle accordingly (e.g., cancel order)
        workflow.GetLogger(ctx).Info("Timer expired, email not verified in time. Cancelling registration.")
        timerStatus = 0
    }).AddReceive(workflow.GetSignalChannel(ctx, "email"), func(c workflow.ReceiveChannel, more bool) {
        // Handle the "paymentCompleted" signal
		c.Receive(ctx, &signalData)

        workflow.GetLogger(ctx).Info("Email verified successfully!")
        timerStatus = 1
    })

    // Wait for either timer expiration, payment completion, or workflow cancellation
    selector.Select(ctx)

	if timerStatus == 1{
		// Redis check confirmation code

		userConfirmation := shared.UserConfirmation{
			Username:    reg.Login,
			ConfirmCode: signalData}

		var confirmationStatus bool

		err = workflow.ExecuteActivity(ctx, activity.CheckConfirmationCodeFromRedis, userConfirmation).Get(ctx, &confirmationStatus)

		if err != nil {
			workflow.GetLogger(ctx).Error("CheckConfirmationCodeFromRedis Activity failed.", "Error", err)
			return err
		}

		if confirmationStatus {
			// Hash password
			err = workflow.ExecuteActivity(ctx, activity.CreateUser, reg).Get(ctx, nil)

			if err != nil {
				workflow.GetLogger(ctx).Error("CreateUser Activity failed.", "Error", err)
				return err
			}

			// Token generation here



		} else {
			workflow.GetLogger(ctx).Error("Confirmation code is incorrect")
			return err
		}
    }else if timerStatus == 0{
		return err
    }




	return nil

}
