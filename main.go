package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	ADD    string = "add"
	UPDATE string = "update"
	DELETE string = "delete"
	LIST   string = "list"
)

const dbFile = "task.json"

func (ts TaskStatus) String() string {
	switch ts {
	case TODO:
		return "TODO"
	case IN_PROGRESS:
		return "IN_PROGRESS"
	case DONE:
		return "DONE"
	default:
		return "UNKNOWN"
	}
}

func readTasksFromFile() []Task {
	var tasks []Task
	fileContent, err := os.ReadFile(dbFile)
	if err != nil {
		// If the file doesn't exist, return an empty task list
		if os.IsNotExist(err) {
			return tasks
		}
		panic("Cannot read tasks file")
	}

	// Handle empty file case
	if len(fileContent) == 0 {
		return tasks
	}
	err = json.Unmarshal(fileContent, &tasks)
	if err != nil {
		panic("Cannot convert JSON to Struct")
	}
	return tasks
}

func getNextID() int {
	tasks := readTasksFromFile()

	if len(tasks) == 0 {
		return 1
	}

	return tasks[len(tasks)-1].Id + 1
}

func createTask(title string) {
	tasks := readTasksFromFile()

	task := Task{
		Id:        getNextID(),
		Title:     title,
		Status:    TODO,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tasks = append(tasks, task)
	newTasks, err := json.Marshal(tasks)
	if err != nil {
		panic("Cannot covert Struct to JSON")
	}

	os.WriteFile(dbFile, newTasks, 0644)
}
func listTasks() {
	tasks := readTasksFromFile()

	for _, task := range tasks {
		fmt.Printf("ID: %d | Title: %s | Status: %s | CreatedAt: %s\n", task.Id, task.Title, task.Status.String(), task.CreatedAt)
	}

}

func main() {
	// if len(os.Args) < 3 {
	// 	fmt.Println("Usage: <command> <title>")
	// 	return
	// }
	actionCommand := os.Args[1]

	switch actionCommand {
	case ADD:
		createTask(os.Args[2])
	case LIST:
		listTasks()
	}
}
