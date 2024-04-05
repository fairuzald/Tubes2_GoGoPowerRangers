package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
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

type AlgorithmParams struct {
	Title string
	Link  string
}

type Route struct {
	Path []AlgorithmParams
}

func BFSAlgorithm(source AlgorithmParams, destination AlgorithmParams) ([]Route, int, error) {
	// Initialize queue for BFS
	queue := make([]AlgorithmParams, 0)
	queue = append(queue, source)

	// Initialize storage to prevent same link scrapping
	storage := make(map[string][]AlgorithmParams)
	storage[source.Link] = make()

	// Initialize routes array to store all possible routes
	var routes []Route

	found := false

	// BFS traversal
	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		if currentNode.Link == destination.Link {
			found = true
			break
		}

		// Perform scraping to get links from current node
		links, err := ScrapingHandler(currentNode.Link)
		linkMap := make(map[string]string)
		if err != nil {
			return nil, 0, err
		}

		// Find neighbors of currentNode
		for _, link := range links {
			if !visited[link["link"]] {
				queue = append(queue, AlgorithmParams{
					Title: link["title"],
					Link:  link["link"],
				})
				visited[link["link"]] = true
				parent[link["link"]] = currentNode.Link
				linkMap[link["title"]] = link["link"]
			}
		}
		storage[currentNode.Link] = linkMap
	}

	// If destination found, reconstruct all possible routes
	if found {
		current := destination.Link
		for current != source.Link {
			var path []string
			path = append(path, current)
			temp := parent[current]
			for temp != source.Link {
				path = append(path, temp)
				temp = parent[temp]
			}
			path = append(path, source.Link)

			// Reverse the path
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}

			routes = append(routes, Route{Path: path})

			current = parent[current]
		}

		// Reverse the routes array to maintain the order
		for i, j := 0, len(routes)-1; i < j; i, j = i+1, j-1 {
			routes[i], routes[j] = routes[j], routes[i]
		}
	}

	return routes, len(routes), nil

}

func PostAlgorithmHandler(c *gin.Context) {
	type AlgorithmRequest struct {
		URL string `json:"url"`
	}

	var reqBody AlgorithmRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	req := reqBody.URL
	links, err := ScrapingHandler(req)
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
