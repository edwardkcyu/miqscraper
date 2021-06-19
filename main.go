package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/joho/godotenv"
)

type Config struct {
	MIQPortalUrl     string
	SlackApiUrl      string
	SlackApiToken    string
	SlackChannelName string
	Cron             string
}

func NewConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env is loaded", err)
	}

	cron := os.Getenv("CRON")
	if cron == "" {
		cron = "*/1 * * * *"
	}

	return Config{
		MIQPortalUrl:     os.Getenv("MIQ_PORTAL_URL"),
		SlackApiUrl:      os.Getenv("SLACK_API_URL"),
		SlackApiToken:    os.Getenv("SLACK_API_TOKEN"),
		SlackChannelName: os.Getenv("SLACK_CHANNEL_NAME"),
		Cron:             cron,
	}
}

func main() {
	config := NewConfig()

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Cron(config.Cron).Do(func() {
		checkMiqPortal(config)
	})
	scheduler.StartBlocking()
}

func checkMiqPortal(config Config) {
	miqManager := NewMiqManager(config.MIQPortalUrl)
	availableDates, err := miqManager.fetchAvailableDates()
	if err != nil {
		log.Fatalf("failed to fetch available date: %v", err)
	}
	fmt.Println(availableDates)

	icon := ":no_entry_sign: Nothing available :cry:"
	if len(availableDates) > 0 {
		icon = ":white_check_mark:"
	}

	text := fmt.Sprintf(
		` %s %s`,
		icon,
		strings.Join(availableDates, ","),
	)

	slackManager := NewSlackManager(config.SlackApiUrl, config.SlackApiToken)
	if err := slackManager.SendMessage(config.SlackChannelName, text); err != nil {
		log.Fatalf("failed to send slack message: %v", err)
	}
	log.Println("Done checking MIQ")
}
