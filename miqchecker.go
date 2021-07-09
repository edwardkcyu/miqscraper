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

func (m *MiqChecker) prepareSlackMessage(availableDates []string) (string, error) {
	sort.Strings(availableDates)

	formattedAvailableDates := make([]string, len(availableDates))
	for i, availableDate := range availableDates {
		date, err := time.Parse("2006-01-02", availableDate)
		if err != nil {
			return "", errors.Wrap(err, "error parsing date")
		}
		dateString := date.Format("02/01 (Mon)")
		formattedAvailableDates[i] = dateString
	}

	icon := ":no_entry_sign: Nothing available :cry:"
	hasAvailableDates := len(availableDates) > 0
	dateContents := strings.Join(formattedAvailableDates, ", ")
	if hasAvailableDates {
		hasTargetDates := strings.Contains(dateContents, "Tue") && strings.Contains(dateContents, "/09")
		if hasTargetDates {
			icon = ":white_check_mark:"
		} else {
			icon = ":eyes:"
		}
	}

	text := fmt.Sprintf(`%s %s`, icon, dateContents)

	return text, nil
}

func (m *MiqChecker) checkMiqPortal(slackChannelName string) error {
	availableDates, err := m.miqManager.fetchAvailableDates()
	if err != nil {
		return errors.Wrap(err, "failed to fetch available date: %v")
	}
	fmt.Println(availableDates)

	text, err := m.prepareSlackMessage(availableDates)
	if err != nil {
		return errors.Wrap(err, "error preparing slack message")
	}

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
