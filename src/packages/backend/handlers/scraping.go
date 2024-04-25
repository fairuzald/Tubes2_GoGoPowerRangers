package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

func ScrapeLink(url string, wg *sync.WaitGroup, linkCh chan<- string, visitedLinks map[string]bool) error {
	defer wg.Done()

	// Create a transport with proxy settings if needed
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	res, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to request URL: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return fmt.Errorf("failed to parse HTML document: %w", err)
	}

	doc.Find("a[href^='/wiki/']").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists && !strings.Contains(link, ":") && !visitedLinks[link] {
			linkCh <- "https://en.wikipedia.org" + link
			visitedLinks[link] = true
		}
	})

	return nil
}

func ScrapeLinkWithRetry(url string, wg *sync.WaitGroup, linkCh chan<- string, visitedLinks map[string]bool) error {
	defer wg.Done()

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	retries := 3           // Number of retries
	var res *http.Response // Declare res outside the loop

	for i := 0; i < retries; i++ {
		var err error // Declare err inside the loop
		res, err = client.Get(url)
		if err != nil {
			return fmt.Errorf("failed to request URL: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {
			break // Successful request, exit retry loop
		}

		if res.StatusCode == http.StatusTooManyRequests {
			// Wait before retrying
			time.Sleep(time.Second * 5)
			continue // Retry the request
		}

		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return fmt.Errorf("failed to parse HTML document: %w", err)
	}

	doc.Find("a[href^='/wiki/']").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists && !strings.Contains(link, ":") && !visitedLinks[link] {
			linkCh <- "https://en.wikipedia.org" + link
			visitedLinks[link] = true
		}
	})

	return nil
}

func ScrapeLinksSync(url string) ([]string, error) {
	var wg sync.WaitGroup
	linkCh := make(chan string, 1000)
	var links []string
	visitedLinks := make(map[string]bool)

	wg.Add(1)
	go func() {
		defer close(linkCh)
		defer wg.Wait()
		err := ScrapeLinkWithRetry(url, &wg, linkCh, visitedLinks)
		if err != nil {
			fmt.Println("Error scraping links:", err)
		}
	}()

	for link := range linkCh {
		links = append(links, link)
	}

	return links, nil
}

func ScrapperHandlerLink(url string) ([]string, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	links := make([]string, 0)
	linkMap := make(map[string]bool)

	doc.Find("a[href^='/wiki/']").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists && !strings.Contains(link, ":") {
			if !linkMap[link] {
				linkMap[link] = true
				links = append(links, "https://en.wikipedia.org"+link)
			}
		}
	})

	return links, nil
}

func ScrapperHandlerLinkBuffer(url string) ([]string, error) {

	res, err := http.Get(url)
	// Request the HTML page.
	if err != nil {
		return nil, fmt.Errorf("failed to request URL: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document.
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML document: %w", err)
	}

	// Find all links with href starting with "/wiki/" and without colon ":".
	links := make([]string, 0)
	linkMap := make(map[string]bool)

	doc.Find("main #mw-content-text a[href^='/wiki/']").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists && !strings.Contains(link, ":") && !strings.Contains(link, "/Main_Page") {
			if !linkMap[link] {
				linkMap[link] = true
				var buffer bytes.Buffer
				buffer.WriteString("https://en.wikipedia.org")
				buffer.WriteString(link)
				links = append(links, buffer.String())
			}
		}
	})

	return links, nil
}

func ScrapingHandlerPost(c *gin.Context) {
	type ReqBody struct {
		Url string `json:"url"`
	}

	use := c.Query("use")

	var reqBody ReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var links []string
	var err error
	tes := "normal"
	if use == "buffer" {
		tes = "buffer"
		links, err = ScrapperHandlerLinkBuffer(reqBody.Url)
	} else if use == "sync" {
		tes = "sync"
		links, err = ScrapeLinksSync(reqBody.Url)
	} else {
		tes = "normal"
		links, err = ScrapperHandlerLink(reqBody.Url)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to perform BFS algorithm"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"method": tes,
		"links":  links,
		"count":  len(links),
	})
}
