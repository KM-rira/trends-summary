package models

type Website struct {
	ToGetWebsiteName string `json:"toGetWebsiteName"`
	URL              string `json:"url"`
	ScrapingTag      string `json:"scrapingTag"`
	RequestText      string `json:"requestText"`
}
