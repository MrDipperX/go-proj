package activity

import (
	"auth/shared"
	"auth/utils"
	"context"

	"github.com/go-redis/redis/v8"
	"go.temporal.io/sdk/activity"
)



func GenerateTokenAndWriteToRedis(ctx context.Context, userData shared.UserData)  error {

	logger := activity.GetLogger(ctx)
	logger.Info("Preparing to write user data to Redis")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.224:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	defer func() {
		if err := rdb.Close(); err != nil {
			logger.Error("Error on close Redis connection")
		}
	}()


	token, err := utils.GenerateUserToken(userData.Username, shared.SecretKey)
	if err != nil {
		logger.Error("Error on generate user token")
		return err
	}

	// Parse the JSON data into UserData struct


	// Set the hash fields in Redis
	err = rdb.HSet(context.Background(), userData.Username, map[string]interface{}{
		"confirm-code": "",
		"expire-at":    token.Expire, // Convert time to Unix timestamp
		"token":        token.Token,
	}).Err()
	
	if err != nil {
		logger.Error("Error on write user data to Redis")
		return err
	}

	logger.Info("User data written to Redis successfully!")

	return nil
}