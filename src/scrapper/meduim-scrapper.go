package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Relevence is a storege of useful users and tags from the scrapper
type Relevence struct {
	users          []string
	publicationIds []string
}

// MediumScrapper is an implementation of WebScrapper
type MediumScrapper struct {
	uri string
}

// Get is the implementation fo MediumScrapper
// Returns tags and users
func (scrapper *MediumScrapper) Get(tags ...string) Relevence {
	uri := ""
	for _, tag := range tags {
		uri = uri + tag
	}

	resp, err := http.Get(uri)

	if err != nil {
		log.Fatal(err)
	}

	users := make([]string, 0)

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatal("Error Loading HTTP response body. ", err)
	}

	doc.Find("a").Each(func(index int, element *goquery.Selection) {
		alt, exists := element.Attr("href")
		if exists {
			for _, u := range strings.Split(alt, "@") {
				match, _ := regexp.MatchString("([a-zA-Z]+\\.*)+", u)
				if match && !strings.Contains(u, "http") && !strings.Contains(u, "https") && !strings.ContainsAny(u, "?#") {
					users = append(users, u)
					fmt.Println(u)
				}
			}
		}
	})
	return Relevence{users: users, publicationIds: make([]string, 0)}
}
