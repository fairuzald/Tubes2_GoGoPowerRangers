package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func IDSHadlers(source string, destination string, maxDepth int) ([][]string, int, error) {
	if source == destination {
		return [][]string{{source}}, 0, nil
	} else {
		depth := 1
		for depth <= maxDepth {
			fmt.Println("Depth", depth)
			resultPaths, counter, err := DFSHelper(source, destination, depth)
			if len(resultPaths) > 0 {
				return resultPaths, counter, err
			} else {
				depth++
			}
		}
	}
	return [][]string{}, 0, nil
}

func DFSHelper(source string, destination string, maxDepth int) ([][]string, int, error) {
	stack := [][]string{{source}}
	paths := [][]string{}
	counter := 0

	for len(stack) > 0 {
		currentPath := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		currentNode := currentPath[len(currentPath)-1]
		counter++

		if currentNode == destination {
			paths = append(paths, currentPath)
		} else if len(currentPath) <= maxDepth {
			var links []string
			var err error
			links, ok := GetLinksFromCache(currentNode)
			if !ok || links == nil || len(links) == 0 {
				links, err = ScrapperHandlerLinkBuffer(currentNode)
				if err != nil {
					return nil, counter, fmt.Errorf("error while processing %s: %s", currentNode, err)
				}
				SetLinksToCache(currentNode, links)
			}

			for _, link := range links {
				if !isInArray(link, currentPath) {
					newPath := append([]string(nil), currentPath...)
					newPath = append(newPath, link)
					stack = append(stack, newPath)
				}
			}
		}
	}
	return paths, counter, nil
}

func isInArray(item string, array []string) bool {
	for _, value := range array {
		if value == item {
			return true
		}
	}
	return false
}

func IDSHTTPHandler(c *gin.Context) {
	type ReqBody struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	}

	startTime := time.Now()

	var reqBody ReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var paths [][]string
	var err error
	var count int
	paths, count, err = IDSHadlers(reqBody.Source, reqBody.Destination, 6)

	endTime := time.Now()
	runtime := endTime.Sub(startTime).Seconds()

	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to perform IDS algorithm"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Data received",
		"paths":        paths,
		"runtime":      runtime,
		"articleCount": count,
	})
}
