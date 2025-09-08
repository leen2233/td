package main

type Task struct {
	ID int `json:"taskId"`
	Text string `json:"text"`
	Timestamp int64 `json:"timestamp"`
	Done bool `json:"done"`
}
