package main

import "time"

type TaskStatus int

type Task struct {
	Id        int        `json:"id"`
	Title     string     `json:"title"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

const (
	TODO TaskStatus = iota
	IN_PROGRESS
	DONE
)
