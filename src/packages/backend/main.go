package main

import (
	"backend/router"
)

func main() {
	r := router.SetupRouter()
	r.Run()
}
