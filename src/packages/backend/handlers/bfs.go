package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func BFSHandlers(source string, destination string, maxDepth int) ([][]string, error) {
	if source == destination {
		return [][]string{{source}}, nil
	}

	type Result struct {
		Links []string
		Err   error
	}

	queue := [][]string{{source}}
	visited := make(map[string]struct{})
	paths := [][]string{}
	found := false

	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Iteration to get
	for len(queue) > 0 {
		queueSameDepth := [][]string{}

		for len(queue) > 0 {
			queueSameDepth = append(queueSameDepth, queue[0])
			queue = queue[1:]
		}

		for _, path := range queueSameDepth {
			currentNode := path[len(path)-1]

			if currentNode == destination {
				found = true
				mutex.Lock()
				paths = append(paths, path)
				mutex.Unlock()
			} else if !found {
				wg.Add(1)
				go func(currentNode string, path []string) {
					defer wg.Done()
					links, err := ScrapeLinksSync(currentNode)
					result := Result{Links: links, Err: err}

					if result.Err != nil {
						// Handle error
						return
					}

					mutex.Lock()
					defer mutex.Unlock()
					for _, link := range result.Links {
						if _, ok := visited[link]; !ok {
							visited[link] = struct{}{}
							newPath := append([]string(nil), path...)
							newPath = append(newPath, link)
							queue = append(queue, newPath)
						}
					}
				}(currentNode, path)
			}
		}
		wg.Wait()
	}

	return paths, nil
}

func NormalBFS(source string, destination string, maxDepth int) ([][]string, error) {
	if source == destination {
		return [][]string{{source}}, nil
	}

	queue := [][]string{{source}}
	visited := make(map[string]struct{})
	paths := [][]string{}
	found := false

	for len(queue) > 0 {
		currentURL := queue[0]
		queue = queue[1:]

		currentPath := currentURL
		currentNode := currentURL[len(currentURL)-1]

		if currentNode == destination {
			found = true
			paths = append(paths, currentPath)
		} else if len(currentPath) < maxDepth && !found {
			links, err := ScrapeLinksSync(currentNode)
			if err != nil {
				return nil, fmt.Errorf("error while processing %s: %s", currentNode, err)
			}

			for _, link := range links {
				if _, ok := visited[link]; !ok {
					visited[link] = struct{}{}
					newPath := append(currentPath, link)
					queue = append(queue, newPath)
				}
			}
		}
	}

	return paths, nil
}

func NormalBFSBool(source string, destination string, maxDepth int) ([][]string, error) {
	if source == destination {
		return [][]string{{source}}, nil
	}

	queue := [][]string{{source}}
	visited := make(map[string]bool)
	paths := [][]string{}
	found := false

	for len(queue) > 0 {
		currentURL := queue[0]
		queue = queue[1:]

		currentPath := currentURL
		currentNode := currentURL[len(currentURL)-1]

		if currentNode == destination {
			found = true
			paths = append(paths, currentPath)
		} else if len(currentPath) < maxDepth && !found {
			links, err := ScrapeLinksSync(currentNode)
			if err != nil {
				return nil, fmt.Errorf("error while processing %s: %s", currentNode, err)
			}

			for _, link := range links {
				if !visited[link] {
					visited[link] = true
					newPath := append(currentPath, link)
					queue = append(queue, newPath)
				}
			}
		}
	}

	return paths, nil
}

type QueueNode struct {
	Path  []string
	Depth int
}

func NormalBFS2(source string, destination string, maxDepth int) ([][]string, error) {
	if source == destination {
		return [][]string{{source}}, nil
	}

	queue := []QueueNode{{Path: []string{source}, Depth: 0}}
	visited := make(map[string]struct{})
	paths := [][]string{}
	found := false

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		currentNode := current.Path[current.Depth]

		if currentNode == destination {
			found = true
			paths = append(paths, current.Path)
		} else if current.Depth < maxDepth && !found {
			links, err := ScrapeLinksSync(currentNode)
			if err != nil {
				return nil, fmt.Errorf("error while processing %s: %s", currentNode, err)
			}

			for _, link := range links {
				if _, ok := visited[link]; !ok {
					visited[link] = struct{}{}
					newPath := append(current.Path, link)
					queue = append(queue, QueueNode{Path: newPath, Depth: current.Depth + 1})
				}
			}
		}
	}

	return paths, nil
}

func NormalBFS2Bool(source string, destination string, maxDepth int) ([][]string, error) {
	if source == destination {
		return [][]string{{source}}, nil
	}

	queue := []QueueNode{{Path: []string{source}, Depth: 0}}
	visited := make(map[string]bool)
	paths := [][]string{}
	found := false

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		currentNode := current.Path[current.Depth]

		if currentNode == destination {
			found = true
			paths = append(paths, current.Path)
		} else if current.Depth < maxDepth && !found {
			links, err := ScrapeLinksSync(currentNode)
			if err != nil {
				return nil, fmt.Errorf("error while processing %s: %s", currentNode, err)
			}

			for _, link := range links {
				if !visited[link] {
					visited[link] = true
					newPath := append(current.Path, link)
					queue = append(queue, QueueNode{Path: newPath, Depth: current.Depth + 1})
				}
			}
		}
	}

	return paths, nil
}

func PostAlgorithmHandler(c *gin.Context) {
	type ReqBody struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	}

	use := c.Query("use")
	tes := ""
	startTime := time.Now() // Record the start time

	var reqBody ReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var paths [][]string
	var err error

	// Measure runtime and set tes accordingly
	if use == "optimized" {
		paths, err = BFSHandlers(reqBody.Source, reqBody.Destination, 6)
		tes = "optimized"
	} else if use == "normal2" {
		paths, err = NormalBFS2(reqBody.Source, reqBody.Destination, 6)
		tes = "normal2"
	} else if use == "normal2bool" {
		paths, err = NormalBFS2Bool(reqBody.Source, reqBody.Destination, 6)
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
