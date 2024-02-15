package utils

import (
	"auth/shared"
	"fmt"

	"context"

	"github.com/go-redis/redis/v8"
)

func GetTokenFromRedis(username string) (shared.UserToken, error) {
	// Get the token from Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.224:6379", // Redis server address
		Password: "",                   // No password set
		DB:       0,                    // Use default DB
	})

	// Context to use for the Redis operation
	ctxR := context.Background()

	var userToken shared.UserToken

	// Retrieve data from Redis by key
	val, err := rdb.HGetAll(ctxR, username).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println(username, "Key does not exist")
		} else {
			fmt.Println("Error on get value from Redis")
		}
		return userToken, err
	}

	token, ok := val["confirm-code"]
	if !ok {
		fmt.Println("Token not found in data for key:", username)
		return userToken, nil
	}

	expireTime, ok := val["expire-at"]
	if !ok {
		fmt.Println("expireAt not found in data for key:", username)
		return userToken, nil
	}

	userToken = shared.UserToken{
		Token:  token,
		Expire: expireTime,
	}

	return userToken, err
}

// func WriteWorkflowIDToRedis (username string, workflowID string) error {
// 	// Get the token from Redis
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "192.168.1.224:6379", // Redis server address
// 		Password: "",               // No password set
// 		DB:       0,                // Use default DB
// 	})

// 	defer func() {
// 		if err := rdb.Close(); err != nil {
// 			fmt.Println("Error on close Redis connection")
// 		}
// 	}()

// 	// Context to use for the Redis operation
// 	ctxR := context.Background()

// 	// Set the hash fields in Redis
// 	err := rdb.HSet(ctxR, username, map[string]interface{}{
// 		"workflow-id": workflowID,
// 	}).Err()

// 	if err != nil {
// 		fmt.Println("Error on write workflow ID to Redis")
// 		return err
// 	}

// 	return nil
// }

// func GetWorkflowIDFromRedis (username string) (string, error) {
// 	// Get the token from Redis
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "192.168.1.224:6379", // Redis server address
// 		Password: "",               // No password set
// 		DB:       0,                // Use default DB
// 	})

// 	defer func() {
// 		if err := rdb.Close(); err != nil {
// 			fmt.Println("Error on close Redis connection")
// 		}
// 	}()

// 	// Context to use for the Redis operation
// 	ctxR := context.Background()

// 	// Get the hash fields from Redis
// 	val, err := rdb.HGetAll(ctxR, username).Result()
// 	if err != nil {
// 		if err == redis.Nil {
// 			fmt.Println(username, "Key does not exist")
// 		} else {
// 			fmt.Println("Error on get value from Redis")
// 		}
// 		return "", err
// 	}

// 	workflowID, ok := val["workflow_id"]
// 	if !ok {
// 		fmt.Println("Confirm code not found in data for key:", username)
// 		return "", nil
// 	}

// 	return workflowID, nil
// }
