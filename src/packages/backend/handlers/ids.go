package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// IDSHadlers implements the Iterative Deepening Search (IDS) algorithm.
func IDSHadlers(source string, destination string, maxDepth int) ([][]string, int, error) {
	if source == destination { // If source and destination are the same
		return [][]string{{source}}, 0, nil
	} else {
		depth := 1
		for depth <= maxDepth {
			fmt.Println("Depth", depth)

			// Perform DFS search with the current depth
			resultPaths, counter, err := DFSHelper(source, destination, depth)

			if len(resultPaths) > 0 { // If paths are found within the depth limit
				return resultPaths, counter, err // Return the paths and counter
			} else {
				depth++ // Increment depth to continue searching with increased depth limit
			}
		}
	}
	return [][]string{}, 0, nil // Return empty paths if no paths are found within the depth limit
}

// DFSHelper performs Depth-First Search (DFS) traversal with a depth limit.
func DFSHelper(source string, destination string, maxDepth int) ([][]string, int, error) {
	closestDistance := make(map[string]int) // Map to store closest distances from the source
	closestDistance[source] = 0             // Initialize distance for the source node
	stack := []struct {
		link  string
		depth int
	}{{
		link:  source,
		depth: 0,
	}}
	parents := make(map[string][]string)                   // Map to store parent-child relationships
	parentsVisited := make(map[string]map[string]struct{}) // Map to track visited parents for each node
	paths := [][]string{}                                  // Slice to store paths found
	found, counter := false, 0                             // Flags and counter for tracking progress
	closestDistance := make(map[string]int)
	closestDistance[source] = 0
	stack := []struct {
		link  string
		depth int
	}{{
		link:  source,
		depth: 0,
	}}
	parents := make(map[string][]string)
	parentsVisited := make(map[string]map[string]struct{})
	paths := [][]string{}
	found, counter := false, 0

	for len(stack) > 0 {
		current := stack[len(stack)-1] // Get the top element from the stack
		stack = stack[:len(stack)-1]   // Pop the top element from the stack

		if current.link == destination { // If the current node is the destination
			found = true // Set the found flag to true
			continue
		} else if current.depth < maxDepth { // If depth limit not reached
			counter++ // Increment the counter for node processing

			// Fetch links for the current node if not cached
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if current.link == destination {
			found = true
			continue
		} else if current.depth < maxDepth {
			counter++
			var links []string
			var err error
			links, ok := GetLinksFromCache(current.link)
			links, ok := GetLinksFromCache(current.link)
			if !ok || links == nil || len(links) == 0 {
				links, err = ScrapperHandlerLinkBuffer(current.link)
				if err != nil {
					return nil, counter, fmt.Errorf("error while processing %s: %s", currentNode, err)
				}
				SetLinksToCache(current.link, links) // Cache the fetched links
			}

			// Process links and update closest distances, parent-child relationships, and stack
			for _, link := range links {
				if _, ok := closestDistance[link]; !ok || current.depth <= closestDistance[link] {
					closestDistance[link] = current.depth
					if _, ok := parentsVisited[link]; !ok {
						parentsVisited[link] = make(map[string]struct{})
					}
					if _, ok := parentsVisited[link][current.link]; !ok {
						parents[link] = append(parents[link], current.link)
						parentsVisited[link][current.link] = struct{}{}
						stack = append(stack, struct {
							link  string
							depth int
						}{link: link, depth: current.depth + 1}) // Push next nodes to the stack
					}
				}
			}
		}
	}

	if found { // If destination is found within depth limit
		paths = reconstructPath(&parents, source, destination) // Reconstruct paths from source to destination
	}
	return paths, counter, nil // Return paths, counter, and no error
}

// IDSHTTPHandler handles HTTP requests for the IDS algorithm.
				if _, ok := closestDistance[link]; !ok || current.depth <= closestDistance[link] {
					closestDistance[link] = current.depth
					if _, ok := parentsVisited[link]; !ok {
						parentsVisited[link] = make(map[string]struct{})
					}
					if _, ok := parentsVisited[link][current.link]; !ok {
						parents[link] = append(parents[link], current.link)
						parentsVisited[link][current.link] = struct{}{}
						stack = append(stack, struct {
							link  string
							depth int
						}{link: link, depth: current.depth + 1})
					}
				}
			}
		}
	}
	if found {
		paths = reconstructPath(&parents, source, destination)
	}
	return paths, counter, nil
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

	queryParams := c.Query("method")

	var paths [][]string
	var err error
	var count int

	// Measure runtime and set tes accordingly
	if queryParams == "single" {
		paths, count, err = IDSHadlersBackupSingle(reqBody.Source, reqBody.Destination, 6)
	} else {
		paths, count, err = IDSHadlers(reqBody.Source, reqBody.Destination, 6)
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
