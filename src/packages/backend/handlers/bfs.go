package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func BFSHandlers(source string, destination string, maxDepth int) ([][]string, int, error) {
	if source == destination {
		return [][]string{{source}}, 1, nil
	}

	queue := [][]string{{source}}
	paths := [][]string{}
	counter := 0
	found := false

	var wg sync.WaitGroup
	var mutex sync.Mutex
	// Semaphore to limit the number of concurrent goroutines
	var sema = make(chan struct{}, 250)

	// Iteration to get
	for len(queue) > 0 && !found && (maxDepth == 0 || len(queue[0]) <= maxDepth) {
		// Check if maxDepth is set and the current depth is greater than maxDepth

		queueSameDepth := [][]string{}

		for len(queue) > 0 {
			queueSameDepth = append(queueSameDepth, queue[0])
			queue = queue[1:]
		}

		fmt.Println("Depth", len(queueSameDepth[0])-1)

		for _, path := range queueSameDepth {
			currentNode := path[len(path)-1]

			// Check if the current node is the destination
			if currentNode == destination {
				found = true
				mutex.Lock()
				paths = append(paths, path)
				mutex.Unlock()
			} else if !found {
				wg.Add(1)
				// Acquire token semaphore
				sema <- struct{}{}
				go func(currentNode string, path []string) {
					defer wg.Done()
					// Release token semaphore
					defer func() { <-sema }()

					// fmt.Println("Goroutine started for node:", currentNode)

					// Get links from the current node
					links, ok := GetLinksFromCache(currentNode)
					if !ok || links == nil || len(links) == 0 {
						links2, err := ScrapperHandlerLinkBuffer(currentNode)
						if err != nil {
							// Handle error
							fmt.Println(err)
							return // Return to avoid deadlock
						}
						SetLinksToCache(currentNode, links2)
						links = links2
					}
					counter += len(links)

					mutex.Lock()
					defer mutex.Unlock()
					for _, link := range links {
						if !isInArray(link, path) {
							// Create a new path by appending the link to the current path
							newPath := append([]string(nil), path...)
							newPath = append(newPath, link)
							queue = append(queue, newPath)
						}
					}

					// fmt.Println("Goroutine finished for node:", currentNode)
				}(currentNode, path)

			}
		}
		wg.Wait()
	}

	return paths, counter, nil
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
