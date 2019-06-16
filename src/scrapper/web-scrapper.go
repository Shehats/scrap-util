package scrapper

// WebScrapper is the class that scraps the web sites for data
type WebScrapper interface {
	Create() *WebScrapper
	Get(tags ...string) interface{}
}
