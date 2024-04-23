package main

import (
	"backend/handlers"
	"fmt"
)

func IDS(source string, destination string, maxDepth int) ([][]string, error) {
	// Check if source is the same as destination
	if source == destination {
		return [][]string{{source}}, nil
	} else {
		depth := 1
		found := false
		for depth <= maxDepth && !found {
			fmt.Println(depth)
			resultPaths, err := DFS(source, destination, depth)
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

func DFS(source string, destination string, maxDepth int) ([][]string, error) {
	// Initialize the stack with the source URL. The last element is top of the stack.
	stack := [][]string{{source}}
	// Track visited URLs to avoid cycles
	visited := make(map[string]struct{})
	// Store all paths found from source to destination
	paths := [][]string{}

	// Loop until the stack is empty
	for len(stack) > 0 {
		// Pop the top path from the stack
		currentPath := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		// Get the last URL in the current path
		currentNode := currentPath[len(currentPath)-1]
		fmt.Println(currentNode)

		// Check if the current node is the destination
		if currentNode == destination {
			fmt.Println(currentPath)
			paths = append(paths, currentPath)
			// If destination found, continue to explore other paths without the need to check depth or visit other links
			continue
		}

		// Check if current path exceeds maximum depth
		if len(currentPath) > maxDepth {
			continue // Skip exploring further if maximum depth reached
		}

		// Explore links from the current node
		links, err := handlers.ScrapeLinksSync(currentNode)
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
				// Add the new path to the stack for further exploration
				stack = append(stack, newPath)
			}
		}
	}

	return paths, nil
}
