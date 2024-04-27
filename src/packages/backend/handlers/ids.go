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
			if !ok || links == nil || len(links) == 0 {
				links, err = ScrapperHandlerLinkBuffer(current.link)
				// fmt.Println(links)
				if err != nil {
					return nil, counter, fmt.Errorf("error while processing %s: %s", current.link, err)
				}
				SetLinksToCache(current.link, links)
			}

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
