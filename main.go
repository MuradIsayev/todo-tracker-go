package main

/*
TODO: Separate Task struct and its methods, and use instance based approach
TODO: Handle CLI commands and flag gracefully
TODO: Improve display of tasks
*/

import (
	"encoding/json"
	"errors"
	"flag"
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

func updateTaskStatus(id string, taskStatus TaskStatus) error {
	var numberRegex = regexp.MustCompile(`^[0-9]+$`)

	if !numberRegex.MatchString(id) {
		return errors.New("ID must only contain digits")
	}

	taskId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("Failed to convert string to integer: %v", err)
	}

	tasks, err := readTasksFromFile()
	index, task, err := findTaskById(tasks, taskId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	task.Status = taskStatus
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
		return fmt.Errorf("Failed to convert string to integer: %v", err)
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
		return fmt.Errorf("Failed to convert string to integer: %v", err)
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

func listTasks(statusFilter TaskStatus) error {
	tasks, err := readTasksFromFile()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if statusFilter == -1 || task.Status == statusFilter {
			fmt.Printf("ID: %d | Title: %s | Status: %s | CreatedAt: %s\n", task.Id, task.Title, task.Status, task.CreatedAt)
		}
	}

	return nil
}

func main() {
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	listDone := listCommand.Bool("done", false, "List tasks with status DONE")
	listInProgress := listCommand.Bool("in-progress", false, "List tasks with status IN_PROGRESS")
	listTodo := listCommand.Bool("todo", false, "List tasks with status TODO")

	addCommand := flag.NewFlagSet("add", flag.ExitOnError)
	updateCommand := flag.NewFlagSet("update", flag.ExitOnError)
	deleteCommand := flag.NewFlagSet("delete", flag.ExitOnError)

	markCommand := flag.NewFlagSet("mark", flag.ExitOnError)
	markDone := markCommand.Bool("done", false, "Mark task with status DONE")
	markInProgress := markCommand.Bool("in-progress", false, "Mark task with status IN_PROGRESS")
	markTodo := markCommand.Bool("todo", false, "Mark task with status TODO")

	switch os.Args[1] {
	case "add":
		addCommand.Parse(os.Args[2:])
		if addCommand.Parsed() {
			if len(addCommand.Args()) != 1 {
				fmt.Println("USAGE: add \"task title\"")
				return
			}
			createTask(addCommand.Args()[0])
		}
	case "list":
		listCommand.Parse(os.Args[2:])
		if listCommand.Parsed() {
			var statusFilter TaskStatus
			switch {
			case *listDone:
				statusFilter = DONE
			case *listInProgress:
				statusFilter = IN_PROGRESS
			case *listTodo:
				statusFilter = TODO
			default:
				statusFilter = -1
			}
			listTasks(statusFilter)
		}
	case "update":
		updateCommand.Parse(os.Args[2:])
		if updateCommand.Parsed() {
			if len(updateCommand.Args()) != 2 {
				fmt.Println("USAGE: update <task ID> \"task title\"")
				return
			}
			updateTaskTitle(updateCommand.Args()[0], updateCommand.Args()[1])
		}
	case "delete":
		deleteCommand.Parse(os.Args[2:])
		if deleteCommand.Parsed() {
			if len(deleteCommand.Args()) != 1 {
				fmt.Println("USAGE: delete <task ID>")
				return
			}
			deleteTask(deleteCommand.Args()[0])
		}
	case "mark":
		markCommand.Parse(os.Args[2:])
		if markCommand.Parsed() {
			if len(markCommand.Args()) != 1 {
				fmt.Println("USAGE: mark --done | --in-progress | --todo <task ID>")
				return
			}
			var taskStatus TaskStatus
			switch {
			case *markDone:
				taskStatus = DONE
			case *markInProgress:
				taskStatus = IN_PROGRESS
			case *markTodo:
				taskStatus = TODO
			default:
				fmt.Println("Please specify a valid status with --done, --in-progress, or --todo")
				return
			}
			updateTaskStatus(markCommand.Args()[0], taskStatus)
		}
	default:
		fmt.Println("Expected 'add', 'list', 'update', 'delete', or 'mark' commands")
		os.Exit(1)
	}
}
