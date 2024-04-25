package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func IDSHadlers(source string, destination string, maxDepth int) ([][]string, int, error) {
	if source == destination {
		return [][]string{{source}}, 0, nil
	} else {
		cache := make(map[string][]string)
		depth := 1
		for depth <= maxDepth {
			fmt.Println("Depth", depth)
			resultPaths, counter, err := DFSHelper(source, destination, depth, &cache)
			if len(resultPaths) > 0 {
				return resultPaths, counter, err
			} else {
				depth++
			}
		}
	}
	return [][]string{}, 0, nil
}

func DFSHelper(source string, destination string, maxDepth int, cache *map[string][]string) ([][]string, int, error) {
	stack := [][]string{{source}}
	paths := [][]string{}
	counter := 0

	for len(stack) > 0 {
		currentPath := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		currentNode := currentPath[len(currentPath)-1]
		counter++

		if currentNode == destination {
			paths = append(paths, currentPath)
		} else if len(currentPath) <= maxDepth {
			var links []string
			var err error
			if value, ok := (*cache)[currentNode]; ok {
				links = value
			} else {
				links, err = ScrapperHandlerLinkBuffer(currentNode)

				if err != nil {
					return nil, counter, fmt.Errorf("error while processing %s: %s", currentNode, err)
				}
				(*cache)[currentNode] = links
			}

			for _, link := range links {
				if !isInArray(link, currentPath) {
					// Create a new path by appending the link to the current path
					newPath := append([]string(nil), currentPath...)
					newPath = append(newPath, link)
					// Add the new path to the stack for further exploration
					stack = append(stack, newPath)
				}

			}
		}
	}
	return paths, counter, nil
}

func isInArray(item string, array []string) bool {
	for _, value := range array {
		if value == item {
			return true
		}
	}
	return false
}

func IDSHTTPHandler(c *gin.Context) {
	type ReqBody struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	}

	startTime := time.Now() // Record the start time

	var reqBody ReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var paths [][]string
	var err error
	var count int
	// Measure runtime and set tes accordingly
	paths, count, err = BFSHandlers(reqBody.Source, reqBody.Destination, 6)

	// Calculate runtime
	endTime := time.Now()
	runtime := endTime.Sub(startTime).Seconds()

	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to perform BFS algorithm"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Data diterima",
		"paths":        paths,
		"runtime":      runtime,
		"articleCount": count,
	})
}
