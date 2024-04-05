package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AlgorithmParams struct {
	Title string
	Link  string
}

type Route struct {
	Path []AlgorithmParams
}

// func BFSAlgorithm(source AlgorithmParams, destination AlgorithmParams) ([]Route, int, error) {
// 	// Initialize queue for BFS
// 	queue := make([]AlgorithmParams, 0)
// 	queue = append(queue, source)

// 	// Initialize routes array to store all possible routes
// 	var routes []Route

// 	found := false

// 	// BFS traversal
// 	for len(queue) > 0 {
// 		currentNode := queue[0]
// 		queue = queue[1:]

// 		if currentNode.Link == destination.Link {
// 			found = true
// 			break
// 		}

// 		// Perform scraping to get links from current node
// 		links, err := handlers.ScrapingHandler(currentNode.Link)
// 		// linkMap := make(map[string]string)
// 		if err != nil {
// 			return nil, 0, err
// 		}

// 		// Find neighbors of currentNode
// 		// for _, link := range links {
// 		// 	if !visited[link["link"]] {
// 		// 		queue = append(queue, AlgorithmParams{
// 		// 			Title: link["title"],
// 		// 			Link:  link["link"],
// 		// 		})
// 		// 		visited[link["link"]] = true
// 		// 		parent[link["link"]] = currentNode.Link
// 		// 		linkMap[link["title"]] = link["link"]
// 		// 	}
// 		// }
// 		// storage[currentNode.Link] = linkMap
// 	}

// 	// If destination found, reconstruct all possible routes
// 	if found {
// 		current := destination.Link
// 		for current != source.Link {
// 			var path []string
// 			path = append(path, current)
// 			temp := parent[current]
// 			for temp != source.Link {
// 				path = append(path, temp)
// 				temp = parent[temp]
// 			}
// 			path = append(path, source.Link)

// 			// Reverse the path
// 			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
// 				path[i], path[j] = path[j], path[i]
// 			}

// 			routes = append(routes, Route{Path: path})

// 			current = parent[current]
// 		}

// 		// Reverse the routes array to maintain the order
// 		for i, j := 0, len(routes)-1; i < j; i, j = i+1, j-1 {
// 			routes[i], routes[j] = routes[j], routes[i]
// 		}
// 	}

// 	return routes, len(routes), nil

// }

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
