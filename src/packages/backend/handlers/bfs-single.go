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

func BFSHandlersSingleBackup(source string, destination string, maxDepth int) ([][]string, int, error) {
	if source == destination {
		return [][]string{{source}}, 1, nil
	}

	queue := [][]string{{source}}
	paths := [][]string{}
	visited := make(map[string]struct{})
	counter := 0

	var wg sync.WaitGroup
	var mutex sync.Mutex
	// Semaphore to limit the number of concurrent goroutines
	var sema = make(chan struct{}, 250)

	// Iteration to get
	for len(queue) > 0 && (maxDepth == 0 || len(queue[0]) <= maxDepth) {
		// Check if maxDepth is set and the current depth is greater than maxDepth

		queueSameDepth := [][]string{}

		for len(queue) > 0 {
			queueSameDepth = append(queueSameDepth, queue[0])
			queue = queue[1:]
		}

		fmt.Println("Depth", len(queueSameDepth[0]))

		for _, path := range queueSameDepth {
			currentNode := path[len(path)-1]
			counter++

			// Check if the current node is the destination
			if currentNode == destination {
				paths = append(paths, path)
				return paths, counter, nil
			} else {
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
							// Handle error and return to avoid deadlock
							fmt.Println(err)
							return
						}
						SetLinksToCache(currentNode, links2) // Update cache with fetched links
						links = links2
					}

					mutex.Lock()
					defer mutex.Unlock()
					for _, link := range links {
						if _, ok := visited[link]; !ok {
							visited[link] = struct{}{}
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
