package handlers

import (
	"backend/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NormalBFS performs a breadth-first search (BFS) algorithm from a given source to a destination using empty struct.
func NormalBFS(source string, destination string, maxDepth int) ([][]string, error) {
	// Check if source is the same as destination
	if source == destination {
		return [][]string{{source}}, nil
	}

	// Initialize the queue with the source URL
	queue := [][]string{{source}}
	// Track visited URLs to avoid cycles
	visited := make(map[string]struct{})
	// Store all paths found from source to destination
	paths := [][]string{}
	// Flag to indicate if destination has been found
	found := false

	// Loop until the queue is empty
	for len(queue) > 0 {
		// Get the first URL in the queue
		currentURL := queue[0]
		queue = queue[1:]

		// Extract the current path and current node (last URL in the path)
		currentPath := currentURL
		currentNode := currentURL[len(currentURL)-1]

		// Check if the current node is the destination
		if currentNode == destination {
			found = true
			paths = append(paths, currentPath)
		} else if len(currentPath) < maxDepth && !found {
			// Explore links from the current node if depth limit is not reached and destination is not found
			links, err := ScrapeLinksSync(currentNode)
			if err != nil {
				return nil, fmt.Errorf("error while processing %s: %s", currentNode, err)
			}

			// Iterate through each link
			for _, link := range links {
				// Check if the link has not been visited before
				if _, ok := visited[link]; !ok {
					visited[link] = struct{}{}
					// Create a new path by appending the link to the current path
					newPath := append(currentPath, link)
					// Add the new path to the queue for further exploration
					queue = append(queue, newPath)
				}
			}
		}
	}

	return paths, nil
}

// NormalBFSBool performs a BFS algorithm similar to NormalBFS but with a boolean map for visited nodes.
func NormalBFSBool(source string, destination string, maxDepth int) ([][]string, error) {
	// Check if source is the same as destination
	if source == destination {
		return [][]string{{source}}, nil
	}

	// Initialize the queue with the source URL
	queue := [][]string{{source}}
	// Track visited URLs to avoid cycles
	visited := make(map[string]bool)
	// Store all paths found from source to destination
	paths := [][]string{}
	// Flag to indicate if destination has been found
	found := false

	// Loop until the queue is empty
	for len(queue) > 0 {
		// Get the first URL in the queue
		currentURL := queue[0]
		queue = queue[1:]

		// Extract the current path and current node (last URL in the path)
		currentPath := currentURL
		currentNode := currentURL[len(currentURL)-1]

		// Check if the current node is the destination
		if currentNode == destination {
			found = true
			paths = append(paths, currentPath)
		} else if len(currentPath) < maxDepth && !found {
			// Explore links from the current node if depth limit is not reached and destination is not found
			links, err := ScrapeLinksSync(currentNode)
			if err != nil {
				return nil, fmt.Errorf("error while processing %s: %s", currentNode, err)
			}

			// Iterate through each link
			for _, link := range links {
				// Check if the link has not been visited before
				if !visited[link] {
					visited[link] = true
					// Create a new path by appending the link to the current path
					newPath := append(currentPath, link)
					// Add the new path to the queue for further exploration
					queue = append(queue, newPath)
				}
			}
		}
	}

	return paths, nil
}

// NormalBFSDepth performs a BFS algorithm using models.QueueNode to track dept using empty struct.
func NormalBFSDepth(source string, destination string, maxDepth int) ([][]string, error) {
	// Check if source is the same as destination
	if source == destination {
		return [][]string{{source}}, nil
	}

	// Initialize the queue with the source URL and depth 0
	queue := []models.QueueNode{{Path: []string{source}, Depth: 0}}
	// Track visited URLs to avoid cycles
	visited := make(map[string]struct{})
	// Store all paths found from source to destination
	paths := [][]string{}
	// Flag to indicate if destination has been found
	found := false

	// Loop until the queue is empty
	for len(queue) > 0 {
		// Get the first node in the queue
		current := queue[0]
		queue = queue[1:]

		// Extract the current node from the path based on its depth
		currentNode := current.Path[current.Depth]

		// Check if the current node is the destination
		if currentNode == destination {
			found = true
			paths = append(paths, current.Path)
		} else if current.Depth < maxDepth && !found {
			// Explore links from the current node if depth limit is not reached and destination is not found
			links, err := ScrapeLinksSync(currentNode)
			if err != nil {
				return nil, fmt.Errorf("error while processing %s: %s", currentNode, err)
			}

			// Iterate through each link
			for _, link := range links {
				// Check if the link has not been visited before
				if _, ok := visited[link]; !ok {
					visited[link] = struct{}{}
					// Create a new path by appending the link to the current path
					newPath := append(current.Path, link)
					// Add the new path and increment the depth to the queue for further exploration
					queue = append(queue, models.QueueNode{Path: newPath, Depth: current.Depth + 1})
				}
			}
		}
	}

	return paths, nil
}

// NormalBFSDepthBool performs a BFS algorithm using models.QueueNode and a boolean map for visited nodes.
func NormalBFSDepthBool(source string, destination string, maxDepth int) ([][]string, error) {
	// Check if source is the same as destination
	if source == destination {
		return [][]string{{source}}, nil
	}

	// Initialize the queue with the source URL and depth 0
	queue := []models.QueueNode{{Path: []string{source}, Depth: 0}}
	// Track visited URLs to avoid cycles
	visited := make(map[string]bool)
	// Store all paths found from source to destination
	paths := [][]string{}
	// Flag to indicate if destination has been found
	found := false

	// Loop until the queue is empty
	for len(queue) > 0 {
		// Get the first node in the queue
		current := queue[0]
		queue = queue[1:]

		// Extract the current node from the path based on its depth
		currentNode := current.Path[current.Depth]

		// Check if the current node is the destination
		if currentNode == destination {
			found = true
			paths = append(paths, current.Path)
		} else if current.Depth < maxDepth && !found {
			// Explore links from the current node if depth limit is not reached and destination is not found
			links, err := ScrapeLinksSync(currentNode)
			if err != nil {
				return nil, fmt.Errorf("error while processing %s: %s", currentNode, err)
			}

			// Iterate through each link
			for _, link := range links {
				// Check if the link has not been visited before
				if !visited[link] {
					visited[link] = true
					// Create a new path by appending the link to the current path
					newPath := append(current.Path, link)
					// Add the new path and increment the depth to the queue for further exploration
					queue = append(queue, models.QueueNode{Path: newPath, Depth: current.Depth + 1})
				}
			}
		}
	}

	return paths, nil
}

func BFSTestHTTPHandler(c *gin.Context) {
	use := c.Query("use")
	tes := ""
	startTime := time.Now()
	var reqBody models.BFSReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var paths [][]string
	var err error

	// Measure runtime and set tes accordingly
	if use == "normal2" {
		paths, err = NormalBFSDepth(reqBody.Source, reqBody.Destination, 6)
		tes = "normal2"
	} else if use == "normal2bool" {
		paths, err = NormalBFSDepthBool(reqBody.Source, reqBody.Destination, 6)
		tes = "normal2bool"
	} else if use == "normalbool" {
		paths, err = NormalBFSBool(reqBody.Source, reqBody.Destination, 6)
		tes = "normalbool"
	} else {
		paths, err = NormalBFS(reqBody.Source, reqBody.Destination, 6)
		tes = "normal"
	}

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
		"Tes":     tes,
		"runtime": runtime,
	})
}
