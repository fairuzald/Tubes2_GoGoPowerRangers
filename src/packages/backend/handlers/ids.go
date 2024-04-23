package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func IDSHadlers(source string, destination string, maxDepth int) ([][]string, error) {
	if source == destination {
		return [][]string{{source}}, nil
	} else {
		cache := make(map[string][]string)
		depth := 1
		found := false
		for depth <= maxDepth && !found {
			resultPaths, err := DFSHelper(source, destination, depth, &cache)
			if len(resultPaths) > 0 {
				found = true
				return resultPaths, err
			} else {
				depth++
			}
		}
	}
	return [][]string{}, nil
}

func DFSHelper(source string, destination string, maxDepth int, cache *map[string][]string) ([][]string, error) {
	stack := [][]string{{source}}
	var mu sync.Mutex // Mutex to synchronize access to the visited map
	var mu1 sync.Mutex
	paths := [][]string{}

	for len(stack) > 0 {
		currentPath := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		currentNode := currentPath[len(currentPath)-1]

		if currentNode == destination {
			mu1.Lock()
			paths = append(paths, currentPath)
			mu1.Unlock()
			continue
		}

		if len(currentPath) > maxDepth {
			continue
		}

		var links []string
		var err error

		if value, ok := (*cache)[currentNode]; ok {
			links = value
		} else {
			links, err = ScrapperHandlerLinkBuffer(currentNode)
			if err != nil {
				return nil, fmt.Errorf("error while processing %s: %s", currentNode, err)
			}
			(*cache)[currentNode] = links
		}

		var wg sync.WaitGroup
		for _, link := range links {
			wg.Add(1)
			go func(l string) {
				defer wg.Done()
				mu.Lock()
				// Check if the link has not been visited before
				if !isInArray(l, currentPath) {
					// Create a new path by appending the link to the current path
					newPath := append(currentPath, l)
					// Add the new path to the stack for further exploration
					stack = append(stack, newPath)
				}
				mu.Unlock()
			}(link)
		}
		wg.Wait() // Wait for all links to be processed before moving to the next node
	}

	return paths, nil
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

	// Measure runtime and set tes accordingly
	paths, err = IDSHadlers(reqBody.Source, reqBody.Destination, 6)

	// Calculate runtime
	endTime := time.Now()
	runtime := endTime.Sub(startTime).Seconds()

	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to perform BFS algorithm"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data diterima",
		"paths":   paths,
		"runtime": runtime,
	})
}

func isInArray(item string, array []string) bool {
	for _, value := range array {
		if value == item {
			return true
		}
	}
	return false
}
