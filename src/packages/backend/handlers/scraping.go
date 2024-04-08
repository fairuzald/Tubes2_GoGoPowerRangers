package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

func ScrapperHandlerLink(url string) ([]string, error) {
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

	// Find all links with href starting with "/wiki/" and without colon ":"
	links := make([]string, 0)
	linkMap := make(map[string]bool)

	doc.Find("a[href^='/wiki/']").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		if !strings.Contains(link, ":") {
			if !linkMap[link] {
				linkMap[link] = true
				links = append(links, "https://en.wikipedia.org"+link)
			}
		}
	})

	return links, nil
}

func ScrapingHandlerPost(c *gin.Context) {
	type ReqBody struct {
		Url string `json:"url"`
	}

	var reqBody ReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	links, err := ScrapperHandlerLink(reqBody.Url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to perform BFS algorithm"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data diterima",
		"links":   links,
		"count":   len(links),
	})
}
