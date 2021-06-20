package main

import (
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/joho/godotenv"
)

type Config struct {
	MIQPortalUrl     string
	SlackApiUrl      string
	SlackApiToken    string
	SlackChannelName string
}

func NewConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env is loaded")
	}

	return Config{
		MIQPortalUrl:     os.Getenv("MIQ_PORTAL_URL"),
		SlackApiUrl:      os.Getenv("SLACK_API_URL"),
		SlackApiToken:    os.Getenv("SLACK_API_TOKEN"),
		SlackChannelName: os.Getenv("SLACK_CHANNEL_NAME"),
	}
}

func main() {
	config := NewConfig()
	log.Printf("%s %s", config.MIQPortalUrl, config.SlackChannelName)

	main := NewMiqChecker(
		NewMiqManager(config.MIQPortalUrl),
		NewSlackManager(config.SlackApiUrl, config.SlackApiToken),
	)

	scheduler := gocron.NewScheduler(time.UTC)
	if _, err := scheduler.Every(10).Seconds().Do(func() {
		main.checkMiqPortal(config)
	}); err != nil {
		log.Fatalf("failed to schedule a job: %v", err)
	}
	scheduler.SingletonMode()
	scheduler.StartBlocking()
}
