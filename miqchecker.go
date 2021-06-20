package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type MiqChecker struct {
	threadId     string
	miqManager   *MiqManager
	slackManager *SlackManager
}

func NewMiqChecker(miqManager *MiqManager, slackManager *SlackManager) *MiqChecker {
	return &MiqChecker{
		miqManager:   miqManager,
		slackManager: slackManager,
	}
}

func (m *MiqChecker) checkMiqPortal(config Config) {
	availableDates, err := m.miqManager.fetchAvailableDates()
	if err != nil {
		log.Fatalf("failed to fetch available date: %v", err)
	}
	fmt.Println(availableDates)

	formattedAvailableDates := make([]string, len(availableDates))
	for i, availableDate := range availableDates {
		date, _ := time.Parse("2006-01-02", availableDate)
		dateString := date.Format("2006-01-02(Mon)")
		formattedAvailableDates[i] = dateString
	}

	icon := ":no_entry_sign: Nothing available :cry:"
	hasAvailableDates := len(availableDates) > 0
	if hasAvailableDates {
		icon = ":white_check_mark:"
		m.threadId = ""
	}

	text := fmt.Sprintf(`%s %s`, icon, strings.Join(formattedAvailableDates, ","))

	threadId, err := m.slackManager.SendMessage(config.SlackChannelName, text, m.threadId)
	if err != nil {
		log.Fatalf("failed to send slack message: %v", err)
	}

	if !hasAvailableDates && m.threadId == "" {
		m.threadId = threadId
	}
}
