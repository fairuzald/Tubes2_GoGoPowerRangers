package handlers

import (
	"fmt"
	"sync"
)

func BFSHandlersSingle(source string, destination string, maxDepth int) ([][]string, int, error) {
	if source == destination {
		return [][]string{{source}}, 1, nil
	}

	var (
		visitedMutex sync.Mutex
		wg           sync.WaitGroup
	)
	queue := []string{source}
	parents := make(map[string]string)
	paths := [][]string{}
	visitedNode := make(map[string]struct{})

	depth, found, counter := 0, false, 0

	semaphore := make(chan struct{}, 250) // Set maximum concurrent goroutines
	done := make(chan struct{})           // Channel to signal completion

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

				select {
				case <-done: // Check if search is done
					return
				default:
				}

				if current == destination && !found {
					found = true
					// Send signal to stop other goroutines
					close(done)
					return
				} else if !found {
					links, err := ScrapperHandlerLinkBuffer(current)
					if err != nil {
						// Handle error
						fmt.Println(err)
						return // Return to avoid deadlock
					}
					visitedMutex.Lock()
					for _, link := range links {
						if _, ok1 := visitedNode[link]; !ok1 {
							visitedNode[link] = struct{}{}
							queue = append(queue, link)
							if parents[link] == "" {
								parents[link] = current
							}
						}
					}
					visitedMutex.Unlock()
				}
			}(current)
		}

		wg.Wait()

		queue = queue[size:]
	}

	if found {
		paths = reconstructPathSingle(&parents, source, destination)
	}

	return paths, counter, nil
}

func reconstructPathSingle(parents *map[string]string, source string, destination string) [][]string {
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
			buildPaths((*parents)[node], path)
		}
	}

	// Start building paths from the destination
	buildPaths(destination, []string{})

	return paths
}
