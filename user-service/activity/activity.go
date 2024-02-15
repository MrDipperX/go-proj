package activity

import (
	"context"
	"user-service/shared"
	"user-service/utils"

	"go.temporal.io/sdk/activity"
	"github.com/go-redis/redis/v8"
	// "encoding/json"

	"database/sql"

	_ "github.com/lib/pq"
)


func SendGeneratedNumberToEmail(ctx context.Context, reg shared.RegistrationBody) (generatedNumber string, err error) {

	logger := activity.GetLogger(ctx)
    logger.Info("Preparing to Generate six digit number")

	numbers, err := utils.GenerateSixDigitNumber()

	if err != nil {
		logger.Error("Error at Generate six digit number")
		return "None", err
	}

	config := shared.EmailConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "UserService",
		Password: "123123Qa",
		From:     "userservice391@gmail.com",
		To:       []string{reg.Email},
		Subject:  "Six-digit Number",
		Body:     numbers,
	}

	logger.Info("Preparing to Send Email")

	if err := utils.SendEmail(config); err != nil {
		logger.Error("Error on Send Email")
		return numbers, err
	}

	logger.Info("Email sent successfully!")

	return numbers, nil
}


func AddUserDataToRedis(ctx context.Context, redisUserData shared.RedisUserData) error {

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

	// Parse the JSON data into UserData struct


	// Set the hash fields in Redis
	err := rdb.HSet(context.Background(), redisUserData.Username, map[string]interface{}{
		"confirm-code": redisUserData.ConfirmCode,
		"expire-at":    redisUserData.TokenExpireAt, // Convert time to Unix timestamp
		"token":        redisUserData.Token,
		"workflow_id":  redisUserData.WorkFlowID,
	}).Err()
	if err != nil {
		logger.Error("Error on write user data to Redis")
		return err
	}

	logger.Info("User data written to Redis successfully!")

	return nil
}


func CheckConfirmationCodeFromRedis(ctx context.Context, reg shared.UserConfirmation) (bool, error) {
	
	logger := activity.GetLogger(ctx)
	logger.Info("Preparing to check confirmation code from Redis")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.224::6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	// Context to use for the Redis operation
	ctxR := context.Background()

	// Retrieve data from Redis by key
	val, err := rdb.HGetAll(ctxR, reg.Username).Result()
	if err != nil {
		if err == redis.Nil {
			logger.Error(reg.Username, "Key does not exist")
		} else {
			logger.Error("Error on get value from Redis")
		}
		return false, err
	}

	confirmCode, ok := val["confirm-code"]
	if !ok {
		logger.Error("Confirm code not found in data for key:", reg.Username)
		return false, nil
	}

	// Check if the confirmation code matches
	if confirmCode != reg.ConfirmCode {
		logger.Error("Confirmation code does not match")
		return false, nil
	}

	return true, nil
}


func CreateUser(ctx context.Context, userData shared.RegistrationBody) error {
	// fmt.Printf(
	// 	"\nSimulate failure to trigger compensation. ReferenceId: %s\n",
	// 	transferDetails.ReferenceID,
	// )
	logger := activity.GetLogger(ctx)

	connectionStr := "postgres://temporal:temporal@192.168.1.224:5441/trydb?sslmode=disable"

	conn, err := sql.Open("postgres", connectionStr)
	if err != nil {
		logger.Error("Error on open connection to DB")
		return err
	}

	defer conn.Close()

	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS Users (Name Varchar(255), Surname Varchar(255), Phone Varchar(255) UNIQUE, Email Varchar(255) UNIQUE, Login Varchar(255) UNIQUE, Password Varchar(255), Validated BOOLEAN);")
	if err != nil {
		logger.Error("Error on create table")
		return err
	}

	hashedPassword, err := utils.HashPassword(userData.Password)

	if err != nil {
		logger.Error("Error on hash password")
		return err
	}

	_, err = conn.Exec("INSERT INTO Users VALUES ($1, $2, $3, $4, $5, $6, $7);", userData.Name, userData.Surname, userData.Phone, userData.Email, userData.Login, hashedPassword, userData.Confirmed)
	if err != nil {
		logger.Error("Error on insert data to table")
		return err
	}

	return nil
}
