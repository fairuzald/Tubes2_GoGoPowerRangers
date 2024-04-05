package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ScrapingHandler(url string) ([]map[string]string, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	// Find all links with href starting with "https://en.wikipedia.org/wiki/"
	links := make([]map[string]string, 0)

	// Create a map to store and identify unique links set
	linkMap := make(map[string]bool)

	// Find all links with href starting with "/wiki/" and without colon ":" (Not wiki link but only distractor like Wikipedia news, recent events, file image, etc. )
	doc.Find("a[href^='/wiki/']").Each(func(i int, s *goquery.Selection) {
		// Get the href attribute value
		link, _ := s.Attr("href")

		// Check if the link contains a colon ":"
		if !strings.Contains(link, ":") {
			// Check if the link is not already in the map
			if !linkMap[link] {
				// Add the link to the map and result array
				linkMap[link] = true
				links = append(links, map[string]string{
					"title": s.Text(),
					"link":  "https://en.wikipedia.org" + link,
				})
			}
		}
	})

	return links, nil
}
