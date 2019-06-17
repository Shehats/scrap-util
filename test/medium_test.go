package test

import (
	"testing"

	"github.com/scrap-util/src/scrapper"
	"github.com/stretchr/testify/assert"
)

func TestMediumFetch(t *testing.T) {
	mediumScrapper := scrapper.NewMediumScrapper("https://medium.com")
	if mediumScrapper == nil {
		t.Error("Instatiation Error")
	}
	var relevence scrapper.Relevence
	relevence = mediumScrapper.Get("/search?q=machine-learning")
	for k, publication := range relevence.Publications {
		assert.NotEmpty(t, k, "The slug shouldn't be nil")
		assert.NotEmpty(t, publication.Title, "Publication title shouldn't be empty", publication.Title)
		assert.NotEqual(t, publication.Uri, "", "Publication uri shouldn't be empty")
		assert.NotEqual(t, publication.Body, "", "Publication Content shouldn't be empty")

	}
}
