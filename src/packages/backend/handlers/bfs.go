package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BFSHandlers(source string, destination string) ([][]string, error) {
	if source == destination {
		return [][]string{{source}}, nil
	}

	queue := [][]string{{source}}
	visited := map[string]bool{source: true}
	var paths [][]string

	for len(queue) > 0 {
		currentURLs := queue[0]
		queue = queue[1:]

		currentURL := currentURLs[len(currentURLs)-1]
		if currentURL == destination {
			paths = append(paths, currentURLs)
		}

		links, err := ScrapperHandlerLink(currentURL)
		if err != nil {
			return nil, fmt.Errorf("error while processing %s: %s", currentURL, err)
		}

		for _, link := range links {
			if !visited[link] {
				visited[link] = true
				newPath := append([]string{}, currentURLs...)
				newPath = append(newPath, link)
				queue = append(queue, newPath)
			}
		}
	}

	return paths, nil
}

func PostAlgorithmHandler(c *gin.Context) {

	type ReqBody struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	}

	var reqBody ReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	paths, err := BFSHandlers(reqBody.Source, reqBody.Destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to perform BFS algorithm"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data diterima",
		"paths":   paths,
		"depth":   len(paths[0]) - 2,
	})
}
