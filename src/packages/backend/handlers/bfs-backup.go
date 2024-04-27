package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func BFSHandlersBackup(source string, destination string, maxDepth int) ([][]string, int, error) {
	if source == destination {
		return [][]string{{source}}, 1, nil
	}

	queue := [][]string{{source}}
	paths := [][]string{}
	counter := 0
	found := false

	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
		mu    sync.Mutex
	)
	// Semaphore to limit the number of concurrent goroutines
	var sema = make(chan struct{}, 250)

	for len(queue) > 0 && !found && (maxDepth == 0 || len(queue[0]) <= maxDepth) {
		queueSameDepth := [][]string{}

		for len(queue) > 0 {
			queueSameDepth = append(queueSameDepth, queue[0])
			queue = queue[1:]
		}

		fmt.Println("Depth", len(queueSameDepth[0]))

		for _, path := range queueSameDepth {
			currentNode := path[len(path)-1]
			counter++

			wg.Add(1)
			sema <- struct{}{}
			go func(currentNode string, path []string) {
				defer wg.Done()
				defer func() { <-sema }()

				if currentNode == destination {
					mu.Lock()
					defer mu.Unlock()
					found = true
					paths = append(paths, path)
					return
				}

				handleLink(currentNode, path, &queue, &mutex)
			}(currentNode, path)
		}
		wg.Wait()
	}

	return paths, counter, nil
}

func handleLink(currentNode string, path []string, queue *[][]string, mutex *sync.Mutex) {
	links, ok := GetLinksFromCache(currentNode)
	if !ok || links == nil || len(links) == 0 {
		links2, err := ScrapperHandlerLinkBuffer(currentNode)
		if err != nil {
			fmt.Println(err)
			return
		}
		SetLinksToCache(currentNode, links2)
		links = links2
	}

	mutex.Lock()
	defer mutex.Unlock()
	for _, link := range links {
		if !isInArray(link, path) {
			newPath := append([]string(nil), path...)
			newPath = append(newPath, link)
			*queue = append(*queue, newPath)
		}
	}
}
func BFSHTTPHandlerBackup(c *gin.Context) {
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
		paths, count, err = BFSHandlersSingleBackup(reqBody.Source, reqBody.Destination, 6)
	} else {
		paths, count, err = BFSHandlersBackup(reqBody.Source, reqBody.Destination, 6)
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
