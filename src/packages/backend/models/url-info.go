package models

type WikipediaSearchResultURLInfo struct {
	Image struct {
		URL string `json:"source"`
	} `json:"thumbnail"`
	PageID int    `json:"pageid"`
	Title  string `json:"title"`
	URL    string `json:"url"`
	Terms  struct {
		Description []string `json:"description"`
	} `json:"terms"`
}

type WikipediaResponseURLInfo struct {
	BatchComplete string `json:"batchcomplete"`
	Query         struct {
		Pages map[string]WikipediaSearchResultURLInfo `json:"pages"`
	} `json:"query"`
}

type URLInfoBody struct {
	URL string `json:"url"`
}
