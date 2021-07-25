package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type CheckerTask struct {
	lastThreadId       string
	lastText           string
	slackChannel       string
	slackTargetChannel string
	miqManager         *MiqManager
	slackManager       *SlackManager
}

func NewCheckerTask(miqManager *MiqManager, slackManager *SlackManager, slackChannel string, slackTargetChannel string) *CheckerTask {
	return &CheckerTask{
		miqManager:         miqManager,
		slackManager:       slackManager,
		slackChannel:       slackChannel,
		slackTargetChannel: slackTargetChannel,
	}
}

func (m *CheckerTask) prepareSlackMessage(availableDates []string) (string, bool, error) {
	sort.Strings(availableDates)

	formattedAvailableDates := make([]string, len(availableDates))
	for i, availableDate := range availableDates {
		date, err := time.Parse("2 January 2006", availableDate)
		if err != nil {
			return "", false, errors.Wrap(err, "error parsing date")
		}
		dateString := date.Format("02/01 (Mon)")
		formattedAvailableDates[i] = dateString
	}

	icon := ":eyes: Nothing available :cry:"
	hasAvailableDates := len(availableDates) > 0
	dateContents := strings.Join(formattedAvailableDates, ", ")
	var hasTargetDates bool
	if hasAvailableDates {
		hasTargetDates = strings.Contains(dateContents, "/09 (Tue)")
		if hasTargetDates {
			icon = ":white_check_mark:"
		} else {
			icon = ":eyes:"
		}
	}

	text := fmt.Sprintf(`%s %s`, icon, dateContents)

	return text, hasTargetDates, nil
}

func (m *CheckerTask) checkMiqPortal() error {
	availableDates, err := m.miqManager.fetchAvailableDates()
	if err != nil {
		return errors.Wrap(err, "failed to fetch available date")
	}
	log.Println(availableDates)

	text, hasTargetDates, err := m.prepareSlackMessage(availableDates)
	if err != nil {
		return errors.Wrap(err, "error preparing slack message")
	}

	isTextDifferent := text != m.lastText
	if isTextDifferent {
		m.lastThreadId = ""
		m.lastText = text
	}

	threadId, err := m.slackManager.SendMessage(m.slackChannel, text, m.lastThreadId)
	if err != nil {
		return errors.Wrap(err, "failed to send slack channel")
	}
	if isTextDifferent {
		m.lastThreadId = threadId
	}

	if hasTargetDates {
		if _, err := m.slackManager.SendMessage(m.slackTargetChannel, text, m.lastThreadId); err != nil {
			return errors.Wrap(err, "failed to send to slack target channel")
		}
	}

	return nil
}
