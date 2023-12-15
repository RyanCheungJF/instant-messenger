package main

type Message struct {
	// capital letters for public fields
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
