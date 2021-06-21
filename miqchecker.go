package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
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

func (m *MiqChecker) checkMiqPortal(slackChannelName string) error {
	availableDates, err := m.miqManager.fetchAvailableDates()
	if err != nil {
		return errors.Wrap(err, "failed to fetch available date: %v")
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

	threadId, err := m.slackManager.SendMessage(slackChannelName, text, m.threadId)
	if err != nil {
		return errors.Wrap(err, "failed to send slack message: %v")
	}

	if !hasAvailableDates && m.threadId == "" {
		m.threadId = threadId
	}

	return nil
}
