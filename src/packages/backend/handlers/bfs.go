package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BFSHandlers(source string, destination string, maxDepth int) ([][]string, error) {
	if source == destination {
		return [][]string{{source}}, nil
	}

	queue := [][]string{{source}}
	visited := make(map[string]bool)
	paths := [][]string{}
	found := false

	for len(queue) > 0 && len(paths) < maxDepth && !found {
		currentURL := queue[0]
		queue = queue[1:]

		currentPath := currentURL
		currentNode := currentURL[len(currentURL)-1]

		if currentNode == destination {
			if !found {
				found = true
			}
			paths = append(paths, currentPath)
		} else if !found {
			links, err := ScrapperHandlerLink(currentNode)
			if err != nil {
				return nil, fmt.Errorf("error while processing %s: %s", currentNode, err)
			}

			for _, link := range links {
				if !visited[link] && len(currentPath) < maxDepth {
					visited[link] = true
					newPath := append(currentPath, link)
					queue = append(queue, newPath)
				}
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

	paths, err := BFSHandlers(reqBody.Source, reqBody.Destination, 6)
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
