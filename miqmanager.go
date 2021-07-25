package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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

	const cutset = " \n "
	var availableDates []string
	doc.Find(".abc__m").Each(func(_ int, monthSelection *goquery.Selection) {
		month := strings.Trim(monthSelection.Find(".abc__m__title").Text(), cutset)

		monthSelection.Find(".abc__d__item").Each(func(_ int, dateSelection *goquery.Selection) {
			childSelection := dateSelection.ChildrenFiltered("div:first-child")
			childText := strings.Trim(childSelection.Text(), cutset)
			if childText == "" {
				return
			}

			class, exists := childSelection.Attr("class")
			if exists && class == "no" {
				return
			}

			dateParsed, err := time.Parse("2 January 2006", fmt.Sprintf("%s %s", childText, month))
			if err != nil {
				log.Printf("error parsing child %v", err)
				return
			}

			availableDates = append(availableDates, dateParsed.Format("2006-01-02"))
		})
	})

	return availableDates, nil
}
