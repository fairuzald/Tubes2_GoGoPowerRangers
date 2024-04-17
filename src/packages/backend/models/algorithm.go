package models

// BFSReqBody represents the request body for the BFS algorithm
type BFSReqBody struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

// QueueNode represents a node in the queue
type QueueNode struct {
	Path  []string
	Depth int
}
