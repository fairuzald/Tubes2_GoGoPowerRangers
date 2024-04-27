package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ScrapperHandlerLinkBuffer(url string) ([]string, error) {

	res, err := http.Get(url)
	// Request the HTML page.
	if err != nil {
		return nil, fmt.Errorf("failed to request URL: %w", err)
	}
	defer res.Body.Close()

	// Load the HTML document.
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML document: %w", err)
	}

	// Find all links with href starting with "/wiki/" and without colon ":".
	links := make([]string, 0)
	linkMap := make(map[string]struct{})

	doc.Find("main #mw-content-text a[href^='/wiki/']").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists && !strings.Contains(link, ":") && !strings.Contains(link, "/Main_Page") {
			if _, ok := linkMap[link]; !ok {
				linkMap[link] = struct{}{}
				var buffer bytes.Buffer
				buffer.WriteString("https://en.wikipedia.org")
				buffer.WriteString(link)
				links = append(links, buffer.String())
			}
		}
	})

	return links, nil

}
