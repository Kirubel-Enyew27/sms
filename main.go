package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/twilio/twilio-go"
	v2010 "github.com/twilio/twilio-go/rest/api/v2010"
)

// Send SMS using Twilio
func sendSMS(to, message string) error {
	// Validate configuration
	accountSID := os.Getenv("ACCOUNT_SID")
	authToken := os.Getenv("AUTH_TOKEN")
	fromNumber := os.Getenv("FROM_NUMBER")

	if accountSID == "" || authToken == "" || fromNumber == "" {
		return fmt.Errorf("Twilio configuration is incomplete")
	}

	// Initialize Twilio client
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	// Create SMS message
	params := &v2010.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(fromNumber)
	params.SetBody(message)

	// Send SMS
	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("could not send SMS: %v", err)
	}

	return nil
}

// Function to be triggered by cron job
func scheduledSMS() {
	toNumber := os.Getenv("TO_NUMBER")
	message := os.Getenv("MESSAGE")

	err := sendSMS(toNumber, message)
	if err != nil {
		fmt.Println("Error sending SMS:", err)
	} else {
		fmt.Println("SMS sent successfully!")
	}
}

// Main function
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Set up cron job to run every 30 seconds
	c := cron.New()
	_, err = c.AddFunc("@every 6h", scheduledSMS) // Runs the scheduledSMS function every 6 hour
	if err != nil {
		log.Fatalf("Error setting up cron job: %v", err)
	}
	c.Start()

	// Keep the main function running to allow cron jobs to trigger
	select {} // Block forever
}
