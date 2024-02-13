package utils

import (
	"net/smtp"
	"user-service/shared"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"

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
