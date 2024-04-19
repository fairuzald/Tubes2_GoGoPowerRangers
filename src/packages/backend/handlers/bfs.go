package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func BFSHandlers(source string, destination string, maxDepth int) ([][]string, error) {
	if source == destination {
		return [][]string{{source}}, nil
	}

	queue := [][]string{{source}}
	visited := make(map[string]struct{})
	paths := [][]string{}
	found := false

	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Iteration to get
	for len(queue) > 0 {
		queueSameDepth := [][]string{}

		for len(queue) > 0 {
			queueSameDepth = append(queueSameDepth, queue[0])
			queue = queue[1:]
		}

		for _, path := range queueSameDepth {
			currentNode := path[len(path)-1]

			if currentNode == destination {
				found = true
				mutex.Lock()
				paths = append(paths, path)
				mutex.Unlock()
			} else if !found {
				wg.Add(1)
				go func(currentNode string, path []string) {
					defer wg.Done()
					links, err := ScrapeLinksSync(currentNode)

					if err != nil {
						// Handle error
						fmt.Print(err)
						return // Return to avoid deadlock

					}

					mutex.Lock()
					defer mutex.Unlock()
					for _, link := range links {
						if _, ok := visited[link]; !ok {
							visited[link] = struct{}{}
							newPath := append([]string(nil), path...)
							newPath = append(newPath, link)
							queue = append(queue, newPath)
						}
					}
				}(currentNode, path)
			}
		}
		wg.Wait()
	}

	return paths, nil
}

func BFSHTTPHandler(c *gin.Context) {
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
	paths, err = BFSHandlers(reqBody.Source, reqBody.Destination, 6)

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
