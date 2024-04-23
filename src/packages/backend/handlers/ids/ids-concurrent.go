package main

import (
	"backend/handlers"
	"fmt"
	"sync"
)

func IDSCon(source string, destination string, maxDepth int) ([][]string, error) {
	// Check if source is the same as destination
	if source == destination {
		return [][]string{{source}}, nil
	}

	// Channel to collect results from goroutines
	resultCh := make(chan [][]string)

	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	depth := 1
	for depth <= maxDepth {
		wg.Add(1)
		go func(d int) {
			defer wg.Done()
			resultPaths, _ := DFSCon(source, destination, d)
			resultCh <- resultPaths
		}(depth)
		depth++
	}

	// Start a goroutine to collect results and close the channel once all results are received
	go func() {
		wg.Wait() // Wait for all DFS calls to finish
		close(resultCh) // Close the channel after all results are received
	}()

	// Collect results from the channel
	paths := [][]string{}
	for result := range resultCh {
		if len(result) > 0 {
			paths = append(paths, result...)
		}
	}

	return paths, nil
}

func DFSCon(source string, destination string, maxDepth int) ([][]string, error) {
	stack := [][]string{{source}}
	visited := make(map[string]struct{})
	var mu sync.Mutex // Mutex to synchronize access to the visited map
	paths := [][]string{}

	for len(stack) > 0 {
		currentPath := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		currentNode := currentPath[len(currentPath)-1]

		if currentNode == destination {
			paths = append(paths, currentPath)
			continue
		}

		if len(currentPath) > maxDepth {
			continue
		}

		links, err := handlers.ScrapeLinksSync(currentNode)
		if err != nil {
			return nil, fmt.Errorf("error while processing %s: %s", currentNode, err)
		}

		var wg sync.WaitGroup
		for _, link := range links {
			wg.Add(1)
			go func(l string) {
				defer wg.Done()
				mu.Lock()
				if _, ok := visited[l]; !ok {
					visited[l] = struct{}{}
					newPath := append(currentPath, l)
					stack = append(stack, newPath)
				}
				mu.Unlock()
			}(link)
		}
		wg.Wait() // Wait for all links to be processed before moving to the next node
	}

	return paths, nil
}