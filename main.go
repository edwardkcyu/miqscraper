package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	res, err := http.Get("https://allocation.miq.govt.nz/portal/")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	availableDates, exists := doc.Find("#accommodation-calendar-home").Attr("data-arrival-dates")
	if !exists {
		log.Fatal("no available date")
	}

	fmt.Printf("Available: %s", strings.Split(availableDates, "_"))
}
