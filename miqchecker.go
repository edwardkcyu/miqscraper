package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type MiqChecker struct {
	lastThreadId string
	lastText     string
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

	sort.Strings(availableDates)

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
	}

	text := fmt.Sprintf(`%s %s`, icon, strings.Join(formattedAvailableDates, ","))
	isTextDifferent := text != m.lastText
	if isTextDifferent {
		m.lastThreadId = ""
		m.lastText = text
	}

	threadId, err := m.slackManager.SendMessage(slackChannelName, text, m.lastThreadId)
	if err != nil {
		return errors.Wrap(err, "failed to send slack message: %v")
	}

	if isTextDifferent {
		m.lastThreadId = threadId
	}

	return nil
}
