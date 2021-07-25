package main

import (
	"fmt"
	"log"
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
	req, err := http.NewRequest(http.MethodGet, m.url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create miq portal request")
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36")
	req.Header.Add("cache-control", "no-store")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to access MIQ portal")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read dom document")
	}

	availableDates := []string{}
	doc.Find(".abc__m").Each(func(_ int, monthSelection *goquery.Selection) {
		month := strings.Trim(monthSelection.Find(".abc__m__title").Text(), " \n")

		monthSelection.Find(".abc__d__item").Each(func(_ int, dateSelection *goquery.Selection) {
			childSelection := dateSelection.Children()
			isAvailable, exists := childSelection.Attr("class")
			if !exists {
				return
			}

			if isAvailable == "no" {
				return
			}

			date := childSelection.Text()
			availableDates = append(availableDates, fmt.Sprintf("%s %s", date, month))
			log.Println(isAvailable, exists, childSelection.Text())
		})
	})

	return availableDates, nil
}
