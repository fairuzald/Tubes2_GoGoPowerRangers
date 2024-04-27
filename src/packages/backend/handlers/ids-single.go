package handlers

import (
	"fmt"
)

func IDSHadlersBackupSingle(source string, destination string, maxDepth int) ([][]string, int, error) {
	if source == destination {
		return [][]string{{source}}, 0, nil
	} else {
		depth := 1
		for depth <= maxDepth {
			fmt.Println("Depth", depth)
			resultPaths, counter, err := DFSHelperBackupSingle(source, destination, depth)
			if len(resultPaths) > 0 {
				return resultPaths, counter, err
			} else {
				depth++
			}
		}
	}
	return [][]string{}, 0, nil
}

func DFSHelperBackupSingle(source string, destination string, maxDepth int) ([][]string, int, error) {
	stack := [][]string{{source}}
	visited := make(map[string]struct{})
	paths := [][]string{}
	counter := 0

	for len(stack) > 0 {
		currentPath := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		currentNode := currentPath[len(currentPath)-1]
		counter++

		if currentNode == destination {
			paths = append(paths, currentPath)
			return paths, counter, nil
		} else if len(currentPath) <= maxDepth {
			var links []string
			links, ok := GetLinksFromCache(currentNode)
			if !ok || links == nil || len(links) == 0 {
				links2, err := ScrapperHandlerLinkBuffer(currentNode)
				if err != nil {
					// Handle error and return to avoid deadlock
					fmt.Println(err)
					break
				}
				SetLinksToCache(currentNode, links2) // Update cache with fetched links
				links = links2
			}

			for _, link := range links {
				if _, ok := visited[link]; !ok {
					visited[link] = struct{}{}
					newPath := append([]string(nil), currentPath...)
					newPath = append(newPath, link)
					stack = append(stack, newPath)
				}
			}
		}
	}
	return paths, counter, nil
}
