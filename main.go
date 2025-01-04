package main

import (
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	v2010 "github.com/twilio/twilio-go/rest/api/v2010"
)

// Initialize configuration
func initConfig() {
	viper.SetConfigName("config") // Config file name (without extension)
	viper.SetConfigType("yaml")   // Config file type
	viper.AddConfigPath(".")      // Look in the current directory

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}

// Send SMS using Twilio
func sendSMS(to, message string) error {
	// Validate configuration
	accountSID := viper.GetString("twilio.account_sid")
	authToken := viper.GetString("twilio.auth_token")
	fromNumber := viper.GetString("twilio.from_number")

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
	err := sendSMS(viper.GetString("twilio.to_number"), viper.GetString("twilio.message"))
	if err != nil {
		fmt.Println("Error sending SMS:", err)
	} else {
		fmt.Println("SMS sent successfully!")
	}
}

// Main function
func main() {
	initConfig() // Load configuration

	// Set up cron job to run every 30 seconds
	c := cron.New()
	_, err := c.AddFunc("@every 6h", scheduledSMS) // Runs the scheduledSMS function every 6 hour
	if err != nil {
		log.Fatalf("Error setting up cron job: %v", err)
	}
	c.Start()

	// Keep the main function running to allow cron jobs to trigger
	select {} // Block forever
}
