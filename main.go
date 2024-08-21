package main

/*
TODO: Separate Task struct and its methods, and use instance based approach
TODO: Handle CLI commands and flag gracefully
TODO: Improve display of tasks
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

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

func readTasksFromFile() ([]Task, error) {
	fileContent, err := os.ReadFile(DB_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, fmt.Errorf("cannot read file: %v", err)
	}

	var tasks []Task
	if len(fileContent) > 0 {
		if err := json.Unmarshal(fileContent, &tasks); err != nil {
			return nil, fmt.Errorf("cannot convert JSON to struct: %v", err)
		}
	}
	return tasks, nil
}

func getNextID(tasks []Task) int {
	if len(tasks) == 0 {
		return 1
	}
	return tasks[len(tasks)-1].Id + 1
}

func createTask(title string) error {
	tasks, err := readTasksFromFile()
	if err != nil {
		return err
	}

	task := Task{
		Id:        getNextID(tasks),
		Title:     title,
		Status:    TODO,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	tasks = append(tasks, task)
	if err := writeTasksToFile(tasks); err != nil {
		return err
	}

	fmt.Printf("Task added successfully (ID: %d)\n", task.Id)
	return nil
}

func findTaskById(tasks []Task, id int) (int, *Task, error) {
	for i, task := range tasks {
		if task.Id == id {
			return i, &task, nil
		}
	}
	return -1, nil, fmt.Errorf("task with ID=%d not found", id)
}

func updateTaskStatus(id, taskStatus string) error {
	var numberRegex = regexp.MustCompile(`^[0-9]+$`)

	if !numberRegex.MatchString(id) {
		return errors.New("ID must only contain digits")
	}

	taskId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("ID must only contain digits: %v", err)
	}

	tasks, err := readTasksFromFile()
	index, task, err := findTaskById(tasks, taskId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	switch taskStatus {
	case MARK_TO_DO:
		task.Status = TODO
	case MARK_IN_PROGRESS:
		task.Status = IN_PROGRESS
	case MARK_DONE:
		task.Status = DONE
	default:
		return errors.New("Invalid status")
	}

	tasks[index] = *task

	err = writeTasksToFile(tasks)
	if err != nil {
		return err
	}

	fmt.Printf("Task status updated successfully (ID: %v)\n", id)
	return nil
}

func updateTaskTitle(id, title string) error {
	var numberRegex = regexp.MustCompile(`^[0-9]+$`)

	if !numberRegex.MatchString(id) {
		return errors.New("ID must only contain digits")
	}

	taskId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("ID must only contain digits: %v", err)
	}

	tasks, err := readTasksFromFile()
	index, task, err := findTaskById(tasks, taskId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if title != "" {
		task.Title = title
	}

	task.UpdatedAt = time.Now().UTC()
	tasks[index] = *task

	err = writeTasksToFile(tasks)
	if err != nil {
		return err
	}

	fmt.Printf("Task title updated successfully (ID: %v)\n", id)
	return nil
}

func deleteTask(id string) error {
	tasks, err := readTasksFromFile()
	if err != nil {
		return err
	}

	var numberRegex = regexp.MustCompile(`^[0-9]+$`)

	if !numberRegex.MatchString(id) {
		return errors.New("ID must only contain digits")
	}

	taskId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("ID must only contain digits: %v", err)
	}

	index, _, err := findTaskById(tasks, taskId)
	if err != nil {
		return err
	}

	tasks = append(tasks[:index], tasks[index+1:]...)
	if err := writeTasksToFile(tasks); err != nil {
		return err
	}

	fmt.Printf("Task deleted successfully (ID: %v)\n", id)
	return nil
}

func writeTasksToFile(tasks []Task) error {
	newTasks, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("Cannot convert Struct to JSON: %v", err)
	}

	if err := os.WriteFile(DB_FILE, newTasks, 0644); err != nil {
		return fmt.Errorf("cannot write to file: %v", err)
	}
	return nil
}

func listTasks() error {
	tasks, err := readTasksFromFile()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		fmt.Printf("ID: %d | Title: %s | Status: %s | CreatedAt: %s\n", task.Id, task.Title, task.Status, task.CreatedAt)
	}

	return nil
}

func main() {
	actionCommand := os.Args[1]

	switch actionCommand {
	case ADD:
		if len(os.Args) != 3 {
			fmt.Println("USAGE: add \"task to do\"")
			return
		}
		createTask(os.Args[2])
	case LIST:
		listTasks()

	case MARK_DONE, MARK_TO_DO, MARK_IN_PROGRESS:
		if len(os.Args) != 3 {
			fmt.Println("USAGE: mark-to-do | mark-in-progress | mark-done 1")
			return
		}
		taskId := os.Args[2]
		taskStatus := os.Args[1]
		updateTaskStatus(taskId, taskStatus)

	case UPDATE:
		if len(os.Args) != 4 {
			fmt.Println("USAGE: update 1 \"new title for 1st task\"")
			return
		}

		taskId := os.Args[2]
		taskTitle := os.Args[3]
		updateTaskTitle(taskId, taskTitle)

	case DELETE:
		if len(os.Args) != 3 {
			fmt.Println("USAGE: delete 1")
			return
		}
		taskId := os.Args[2]
		deleteTask(taskId)
	default:
		fmt.Println("Expected 'add', 'list', 'update', 'delete', 'mark-to-do | mark-in-progress | mark-done' commands")
	}
}
