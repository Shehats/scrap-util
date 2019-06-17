package scrapper

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Workiva/go-datastructures/set"
)

const TitleClass = "graf graf--h3 graf-after--figure graf--trailing graf--title"

type Publication struct {
	Title    string
	Category string
	Uri      string
	Body     string
}

// Relevence is a storege of useful publications
type Relevence struct {
	Publications map[string]Publication
}

// MediumScrapper is an implementation of WebScrapper
type MediumScrapper struct {
	uri string
}

func NewMediumScrapper(uri string) *MediumScrapper {
	scrapper := new(MediumScrapper)
	scrapper.uri = uri
	return scrapper
}

// Get is the implementation fo MediumScrapper
func (scrapper *MediumScrapper) Get(tags ...string) Relevence {
	uri := scrapper.uri

	soFar := set.New()
	visited := set.New()

	for _, tag := range tags {
		uri = uri + tag
	}

	resp, err := http.Get(uri)

	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatal("Error Loading HTTP response body. ", err)
	}

	var publications map[string]Publication

	publications = make(map[string]Publication)

	doc.Find("h3").Each(func(index int, element *goquery.Selection) {
		if element.HasClass(TitleClass) {
			slug := slugify(element.Text())
			publications[slug] = Publication{Title: element.Text()}
			soFar.Add(slug)
		}
	})

	doc.Find("a").Each(func(index int, element *goquery.Selection) {
		alt, exists := element.Attr("href")
		if exists {
			// Get publication uri if exists
			for _, tag := range soFar.Flatten() {
				if strings.Contains(alt, tag.(string)) {
					if !strings.Contains(alt, "responses") && !visited.Exists(tag) {
						publications[tag.(string)] = Publication{Title: publications[tag.(string)].Title, Uri: alt}
						publications[tag.(string)] = getContent(publications[tag.(string)])
						visited.Add(tag.(string))
					}
				}
			}
		}
	})

	cleanPublications := make(map[string]Publication)

	// Filter publications for empty:
	for k, publication := range publications {
		if publication.Uri != "" {
			cleanPublications[k] = publication
		}
	}

	return Relevence{Publications: cleanPublications}
}

func getContent(publication Publication) Publication {
	content := ""
	resp, err := http.Get(publication.Uri)
	if err == nil {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err == nil {
			doc.Find("p").Each(func(index int, element *goquery.Selection) {
				content = content + element.Text() + "\n"
			})
		}
	}
	return Publication{Title: publication.Title, Uri: publication.Uri, Body: content}
}

func slugify(str string) string {
	cleanStr := regexp.MustCompile(`[^[a-z-A-Z0-9\s]+]*`)
	s := strings.ToLower(str)
	s = cleanStr.ReplaceAllString(s, " ")
	retStr := ""
	for _, word := range strings.Split(s, " ") {
		retStr = retStr + "-" + word
	}
	retStr = strings.TrimRight(retStr, "-")
	retStr = strings.TrimLeft(retStr, "-")
	esnureCleanStr := regexp.MustCompile(`[-]{2,}`)
	retStr = esnureCleanStr.ReplaceAllString(retStr, "-")
	return retStr
}
