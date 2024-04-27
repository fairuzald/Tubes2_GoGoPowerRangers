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

	var (
		closestMutex sync.RWMutex // Changed to RWMutex
		foundMutex   sync.Mutex
		wg           sync.WaitGroup
	)
	closestDistance := make(map[string]int)
	closestDistance[source] = 0
	queue := []string{source}
	parents := make(map[string][]string)
	parentsVisited := make(map[string]map[string]struct{})
	paths := [][]string{}

	depth, found, counter := 0, false, 0

	semaphore := make(chan struct{}, 250) // Set maximum concurrent goroutines

	for len(queue) > 0 && !found && (maxDepth == 0 || depth <= maxDepth) {
		depth++
		fmt.Println("Depth", depth)

		size := len(queue)

		for i := 0; i < size; i++ {
			semaphore <- struct{}{} // Acquire semaphore
			current := queue[i]
			counter++
			wg.Add(1)
			go func(current string) {
				defer func() {
					<-semaphore // Release semaphore
					wg.Done()   // Ensure wg.Done() is called for every goroutine
				}()

				if current == destination && !found {
					foundMutex.Lock()
					found = true
					foundMutex.Unlock()
				} else if !found {
					links, ok := GetLinksFromCache(current)
					if !ok || links == nil || len(links) == 0 {
						links2, err := ScrapperHandlerLinkBuffer(current)
						if err != nil {
							// Handle error
							fmt.Println(err)
							return // Return to avoid deadlock
						}
						SetLinksToCache(current, links2)
						links = links2
					}

					closestMutex.Lock()
					for _, link := range links {
						if _, ok := closestDistance[link]; !ok || depth <= closestDistance[link] {
							closestDistance[link] = depth
							queue = append(queue, link)
							if _, ok := parentsVisited[link]; !ok {
								parentsVisited[link] = make(map[string]struct{})
							}
							if _, ok := parentsVisited[link][current]; !ok {
								parents[link] = append(parents[link], current)
								parentsVisited[link][current] = struct{}{}
							}
						}
					}
					closestMutex.Unlock()
				}
			}(current)
		}

		wg.Wait()

		queue = queue[size:]
	}

	fmt.Println("Depth", depth, found)

	if found {
		// Reconstruct paths if destination is found
		paths = reconstructPath(&parents, source, destination)
	}

	return paths, counter, nil
}

func reconstructPath(parents *map[string][]string, source string, destination string) [][]string {
	var paths [][]string
	var mu sync.Mutex // Mutex untuk mengamankan penggunaan slice paths

	// Recursive function to build paths
	var buildPaths func(string, []string)
	buildPaths = func(node string, path []string) {
		path = append(path, node)

		if node == source {
			// Reverse the path to get the correct order
			reversePath := make([]string, len(path))
			for i, j := 0, len(path)-1; i < len(path); i, j = i+1, j-1 {
				reversePath[i] = path[j]
			}

			// Lock dan update slice paths dengan menggunakan mutex
			mu.Lock()
			defer mu.Unlock()
			paths = append(paths, reversePath)
		} else {
			var wg sync.WaitGroup

			for _, parent := range (*parents)[node] {
				wg.Add(1)
				go func(parent string) {
					defer wg.Done()
					buildPaths(parent, path)
				}(parent)
			}

			wg.Wait()
		}
	}

	// Start building paths from the destination
	buildPaths(destination, []string{})

	return paths
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
