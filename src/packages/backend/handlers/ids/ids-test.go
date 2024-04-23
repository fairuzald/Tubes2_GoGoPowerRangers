package main

import (
	"backend/handlers"
	"fmt"
)

func main() {
	paths1, _ := handlers.ScrapeLinksSync("https://en.wikipedia.org/wiki/Joko_Widodo")
	fmt.Println(paths1)
	paths, _ := IDSCon("https://en.wikipedia.org/wiki/Joko_Widodo", "https://en.wikipedia.org/wiki/Malay_language", 3)
	fmt.Println(paths)
}