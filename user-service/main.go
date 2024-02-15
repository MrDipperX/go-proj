package main

import (
	"github.com/go-redis/redis/v8"
	"fmt"
	"context"
)

func GetWorkflowIDFromRedis(username string) (map[string]string, error) {
    // Context to use for the Redis operation
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.1168:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	defer func() {
		if err := rdb.Close(); err != nil {
			fmt.Println("Error on close Redis connection")
		}
	}()

    ctxR := context.Background()

    // Get the hash fields from Redis
    val, err := rdb.HGetAll(ctxR, username).Result()
    if err != nil {
        fmt.Println("Error on get workflow ID from Redis:", err)
        return nil, err
    }

    // Print the dictionary retrieved from Redis
    fmt.Println("Dictionary from Redis:", val)

    return val, nil
}

func main() {
    // Replace "aaa" with the actual username/key
    dict, err := GetWorkflowIDFromRedis("aaa")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Access fields from the dictionary
    workflowID := dict["workflow_id"]
    fmt.Println("WorkflowID from Redis:", workflowID)
}