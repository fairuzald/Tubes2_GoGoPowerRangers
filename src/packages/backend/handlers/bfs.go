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
		closestMutex sync.Mutex     // Mutex to ensure safe access to closestDistance map
		foundMutex   sync.Mutex     // Mutex to ensure safe access to found flag
		wg           sync.WaitGroup // WaitGroup to synchronize goroutines
	)
	closestDistance := make(map[string]int)                // Map to track closest distance to each URL from source
	closestDistance[source] = 0                            // Initial distance of source URL is 0
	queue := []string{source}                              // Queue to process URLs in BFS manner
	parents := make(map[string][]string)                   // Map to track parents of each URL in shortest paths
	parentsVisited := make(map[string]map[string]struct{}) // Map to track visited parents for each URL
	paths := [][]string{}                                  // Initialize paths as empty slice

	depth, found, counter := 0, false, 0 // Initialize depth, found flag, and counter for processed URLs

	semaphore := make(chan struct{}, 250) // Semaphore to limit concurrent goroutines

	// Loop until the queue is empty, destination is found, or maxDepth is reached (if set)
	for len(queue) > 0 && !found && (maxDepth == 0 || depth <= maxDepth) {
		depth++
		fmt.Println("Depth", depth)

		size := len(queue)

		// Iterate through each node at the current depth level
		for i := 0; i < size; i++ {
			semaphore <- struct{}{} // Acquire semaphore to control goroutine concurrency
			current := queue[i]
			counter++
			wg.Add(1) // Add to WaitGroup to track active goroutines

			// Launch a goroutine to process the current node concurrently
			go func(current string) {
				defer func() {
					<-semaphore // Release semaphore when the goroutine exits
					wg.Done()   // Mark WaitGroup as done to track completed goroutines
				}()

				// Check if the current node is the destination
				if current == destination && !found {
					foundMutex.Lock()
					found = true
					foundMutex.Unlock()
				} else if !found { // If not the destination, process the node
					// Fetch links for the current node if not cached
					links, ok := GetLinksFromCache(current)
					if !ok || links == nil || len(links) == 0 {
						links2, err := ScrapperHandlerLinkBuffer(current)
						if err != nil {
							// Handle error and return to avoid deadlock
							fmt.Println(err)
							return
						}
						SetLinksToCache(current, links2) // Update cache with fetched links
						links = links2
					}

					closestMutex.Lock() // Lock to ensure safe access to shared data structures
					for _, link := range links {
						// Update closest distance and queue if needed
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
					closestMutex.Unlock() // Unlock after updating shared data
				}
			}(current)
		}

		wg.Wait() // Wait for all goroutines in this depth level to finish processing

		queue = queue[size:] // Remove processed nodes from the queue
	}

	fmt.Println("Depth", depth, found)

	if found {
		// Reconstruct paths if destination is found
		paths = reconstructPath(&parents, source, destination)
	}

	return paths, counter, nil
}

// Reconstruct paths from destination to source using parent nodes
func reconstructPath(parents *map[string][]string, source string, destination string) [][]string {
	var paths [][]string
	var mu sync.Mutex // Mutex to secure access to the paths slice

	// Recursive function to build paths
	var buildPaths func(string, []string)
	buildPaths = func(node string, path []string) {
		path = append(path, node) // Append the current node to the path

		if node == source { // If the current node is the source node
			// Reverse the path to get the correct order
			reversePath := make([]string, len(path))
			for i, j := 0, len(path)-1; i < len(path); i, j = i+1, j-1 {
				reversePath[i] = path[j]
			}

			// Lock and update the paths slice using the mutex to avoid concurrent modification
			mu.Lock()
			defer mu.Unlock()
			paths = append(paths, reversePath) // Add the reversed path to the paths slice
		} else {
			var wg sync.WaitGroup

			// Iterate through parent nodes of the current node
			for _, parent := range (*parents)[node] {
				wg.Add(1)
				go func(parent string) {
					defer wg.Done()
					buildPaths(parent, path) // Recursively build paths for parent nodes
				}(parent)
			}

			wg.Wait() // Wait for all recursive calls to finish
		}
	}

	// Start building paths from the destination node with an empty path slice
	buildPaths(destination, []string{})

	return paths // Return the constructed paths from destination to source
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

	queryParams := c.Query("method")

	var paths [][]string
	var err error
	var count int

	// Measure runtime and set tes accordingly
	if queryParams == "single" {
		paths, count, err = BFSHandlersSingle(reqBody.Source, reqBody.Destination, 6)
	} else {
		paths, count, err = BFSHandlers(reqBody.Source, reqBody.Destination, 6)
	}

	// Calculate runtime
	endTime := time.Now()
	runtime := endTime.Sub(startTime).Seconds()

	if err != nil {
		fmt.Println(err)
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
