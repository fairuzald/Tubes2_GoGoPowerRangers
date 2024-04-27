package handlers

import (
	"fmt"
	"sync"
)

func IDSConcurrentHadlers(source string, destination string, maxDepth int) ([][]string, int, error) {
	if source == destination {
		return [][]string{{source}}, 0, nil
	} else {
		depth := 1
		for depth <= maxDepth {
			fmt.Println("Depth", depth)
			resultPaths, counter, err := DFSConcurrentHelper(source, destination, depth)
			if len(resultPaths) > 0 {
				return resultPaths, counter, err
			} else {
				depth++
			}
		}
	}
	return [][]string{}, 0, nil
}

func DFSConcurrentHelper(source string, destination string, maxDepth int) ([][]string, int, error) {
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

	semaphore := make(chan struct{}, 250)
	var mu sync.Mutex
	var mu1 sync.Mutex
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for len(stack) > 0 {
		maxDepthProcessed := stack[len(stack)-1].depth
		for i := 0; i < len(stack) && stack[len(stack)-1].depth == maxDepthProcessed; i++ {
			semaphore <- struct{}{}
			counter++
			wg.Add(1)
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if current.link == destination {
				mu.Lock()
				found = true
				mu.Unlock()
				wg.Done()
			} else if current.depth < maxDepth {
				// Launch goroutine
				go func(currentLink string, currentDepth int) {
					defer wg.Done()

					var links []string
					var err error
					links, ok := GetLinksFromCache(currentLink)
					if !ok || links == nil || len(links) == 0 {
						links, err = ScrapperHandlerLinkBuffer(currentLink)
						// fmt.Println(links)
						if err != nil {
							errCh <- fmt.Errorf("error while processing %s: %s", currentLink, err)
							return
						}
						SetLinksToCache(currentLink, links)
					}

					mu1.Lock()
					for _, link := range links {
						if _, ok := closestDistance[link]; !ok || currentDepth <= closestDistance[link] {
							closestDistance[link] = currentDepth
							if _, ok := parentsVisited[link]; !ok {
								parentsVisited[link] = make(map[string]struct{})
							}
							if _, ok := parentsVisited[link][currentLink]; !ok {
								parents[link] = append(parents[link], currentLink)
								parentsVisited[link][currentLink] = struct{}{}
								stack = append(stack, struct {
									link  string
									depth int
								}{link: link, depth: currentDepth + 1})
							}
						}
					}
					mu1.Unlock()
				}(current.link, current.depth)
			} else {
				wg.Done()
			}
			<-semaphore
		}
		wg.Wait()
	}

	close(errCh)
	for err := range errCh {
		return nil, 0, err
	}

	if found {
		paths = reconstructPath(&parents, source, destination)
	}
	return paths, counter, nil
}
