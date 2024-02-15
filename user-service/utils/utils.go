package utils

import (
	"net/smtp"
	"user-service/shared"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/go-redis/redis/v8"
	"context"

)

func SendEmail(config shared.EmailConfig) error {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", config.To, config.Subject, config.Body))

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	if err := smtp.SendMail(addr, auth, config.From, config.To, msg); err != nil {
		return err
	}

	return nil
}


func GenerateSixDigitNumber() (string, error) {
	// Seed the random number generator with current time
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random six-digit number
	number := r.Intn(900000) + 100000

	// Convert the number to a string
	return fmt.Sprintf("%06d", number), nil
}



func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}


func WriteWorkflowIDToRedis (username string, workflowID string) error {
	// Get the token from Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.224:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	defer func() {
		if err := rdb.Close(); err != nil {
			fmt.Println("Error on close Redis connection")
		}
	}()

	// Context to use for the Redis operation
	ctxR := context.Background()

	// Set the hash fields in Redis
	err := rdb.HSet(ctxR, username, map[string]interface{}{
		"workflow-id": workflowID,
	}).Err()

	if err != nil {
		fmt.Println("Error on write workflow ID to Redis")
		return err
	}

	return nil
}


func GetWorkflowIDFromRedis (username string) (string, error) {
	// Get the token from Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.224:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	defer func() {
		if err := rdb.Close(); err != nil {
			fmt.Println("Error on close Redis connection")
		}
	}()

	// Context to use for the Redis operation
	ctxR := context.Background()

	// Get the hash fields from Redis
	val, err := rdb.HGetAll(ctxR, username).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println(username, "Key does not exist")
		} else {
			fmt.Println("Error on get value from Redis")
		}
		return "", err
	}

	workflowID, ok := val["workflow_id"]
	if !ok {
		fmt.Println("Confirm code not found in data for key:", username)
		return "", nil
	}

	return workflowID, nil
}