package main

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type MiqManager struct {
	url string
}

func NewMiqManager(url string) *MiqManager {
	return &MiqManager{
		url: url,
	}
}

func (m MiqManager) fetchAvailableDates() ([]string, error) {
	res, err := http.Get(m.url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to access MIQ portal")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read dom document")
	}

	availableDates, exists := doc.Find("#accommodation-calendar-home").Attr("data-arrival-dates")
	if !exists {
		return nil, errors.New("dom attribute #accommodation-calendar-home not found")
	}

	if availableDates == "" {
		return []string{}, nil
	}

	return strings.Split(availableDates, "_"), nil
}
