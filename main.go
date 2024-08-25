package main

/*
TODO: Improve display of tasks (use libraries like tablewriter to beautify visualization)
TODO: Reduce repetition of code and improve code quality (separation of concerns and DRY principle)
TODO: Display date in a more readable format (e.g., "2021-09-01T12:00:00Z" -> "2021-09-01 12:00:00")
TODO: Add more features like sorting tasks by date, filtering tasks by date range, etc.
TODO: Understand what the code does and how it works
TODO: Maybee add unit tests?
*/

import (
	"flag"
	"fmt"
	"os"

	"github.com/MuradIsayev/todo-tracker/constants"
	"github.com/MuradIsayev/todo-tracker/task"
)

func main() {
	taskService := task.NewTaskService(constants.DB_FILE)

	// Commands
	listCommand := flag.NewFlagSet(constants.LIST, flag.ExitOnError)
	addCommand := flag.NewFlagSet(constants.ADD, flag.ExitOnError)
	updateCommand := flag.NewFlagSet(constants.UPDATE, flag.ExitOnError)
	deleteCommand := flag.NewFlagSet(constants.DELETE, flag.ExitOnError)
	markCommand := flag.NewFlagSet(constants.MARK, flag.ExitOnError)

	// Flags for list command
	listDone := listCommand.Bool("done", false, "List tasks with status DONE")
	listInProgress := listCommand.Bool("in-progress", false, "List tasks with status IN_PROGRESS")
	listTodo := listCommand.Bool("todo", false, "List tasks with status TODO")

	// Flags for mark command
	markDone := markCommand.Bool("done", false, "Mark task with status DONE")
	markInProgress := markCommand.Bool("in-progress", false, "Mark task with status IN_PROGRESS")
	markTodo := markCommand.Bool("todo", false, "Mark task with status TODO")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'add', 'list', 'update', 'delete', or 'mark' commands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case constants.ADD:
		addCommand.Parse(os.Args[2:])
		if addCommand.Parsed() {
			if len(addCommand.Args()) != 1 {
				fmt.Println("USAGE: add \"task title\"")
				return
			}
			if err := taskService.CreateTask(addCommand.Args()[0]); err != nil {
				fmt.Println("Error:", err)
			}
		}
	case constants.LIST:
		listCommand.Parse(os.Args[2:])
		if listCommand.Parsed() {
			var statusFilter task.TaskStatus
			switch {
			case *listDone:
				statusFilter = task.DONE
			case *listInProgress:
				statusFilter = task.IN_PROGRESS
			case *listTodo:
				statusFilter = task.TODO
			default:
				statusFilter = -1
			}
			if err := taskService.ListTasks(statusFilter); err != nil {
				fmt.Println("Error:", err)
			}
		}
	case constants.UPDATE:
		updateCommand.Parse(os.Args[2:])
		if updateCommand.Parsed() {
			if len(updateCommand.Args()) != 2 {
				fmt.Println("USAGE: update <task_id> \"task title\"")
				return
			}
			if err := taskService.UpdateTaskTitle(updateCommand.Args()[0], updateCommand.Args()[1]); err != nil {
				fmt.Println("Error:", err)
			}
		}
	case constants.DELETE:
		deleteCommand.Parse(os.Args[2:])
		if deleteCommand.Parsed() {
			if len(deleteCommand.Args()) != 1 {
				fmt.Println("USAGE: delete <task_id>")
				return
			}
			if err := taskService.DeleteTask(deleteCommand.Args()[0]); err != nil {
				fmt.Println("Error:", err)
			}
		}
	case constants.MARK:
		markCommand.Parse(os.Args[2:])
		if markCommand.Parsed() {
			if len(markCommand.Args()) != 1 {
				fmt.Println("USAGE: mark --done | --in-progress | --todo <task_id> ")
				return
			}
			var status task.TaskStatus
			switch {
			case *markDone:
				status = task.DONE
			case *markInProgress:
				status = task.IN_PROGRESS
			case *markTodo:
				status = task.TODO
			default:
				fmt.Println("You must specify a task status using --done | --in-progress | --todo")
				return
			}
			if err := taskService.UpdateTaskStatus(markCommand.Args()[0], status); err != nil {
				fmt.Println("Error:", err)
			}
		}
	default:
		fmt.Println("Expected 'add', 'list', 'update', 'delete', or 'mark' commands")
		os.Exit(1)
	}
}
