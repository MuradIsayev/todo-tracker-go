package main

/*
TODO: Add more features like sorting tasks by date, filtering tasks by date range, etc.
TODO: Maybee add unit tests?
*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/MuradIsayev/todo-tracker/constants"
	"github.com/MuradIsayev/todo-tracker/project"
	"github.com/MuradIsayev/todo-tracker/task"
	"github.com/olekukonko/tablewriter"
)

func startREPL(id string) {
	fmt.Println("Welcome to the Task Management CLI - Interactive Mode")
	reader := bufio.NewReader(os.Stdin)

	taskTable := tablewriter.NewWriter(os.Stdout)
	taskService := task.NewTaskService(id, taskTable)
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
	// markTaskDone := markCommand.Bool("done", false, "Mark task with status DONE")
	// markTaskInProgress := markCommand.Bool("in-progress", false, "Mark task with status IN_PROGRESS")
	// markTaskTodo := markCommand.Bool("todo", false, "Mark task with status TODO")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		fmt.Println("input:", input)
		input = strings.TrimSpace(input)

		if input == "exit" || input == "quit" {
			break
		}

		handleCommand(input, taskService, listCommand, addCommand, updateCommand, deleteCommand, markCommand, listDone, listInProgress, listTodo)
	}
}

func handleCommand(input string, taskService *task.TaskService, listCommand, addCommand, updateCommand, deleteCommand, markCommand *flag.FlagSet, listDone, listInProgress, listTodo *bool) {
	parts := strings.Fields(input)
	fmt.Println("parts:", parts)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "delete":
		listCommand.Parse(args)

		if listCommand.Parsed() {
			fmt.Println("delete command parsed", listCommand.Args())
		}
	case "list":
		listCommand.Parse(args)
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
	case "add":
		addCommand.Parse(args)
		if addCommand.Parsed() {
			if len(addCommand.Args()) != 1 {
				fmt.Println("USAGE: add \"task title\"")
				return
			}
			if err := taskService.CreateTask(addCommand.Args()[0]); err != nil {
				fmt.Println("Error:", err)
			}
		}
	case "switch-project":
		// Switch project logic
	default:
		fmt.Println("Unknown command:", command)
	}
}

func main() {
	projectTable := tablewriter.NewWriter(os.Stdout)

	projectService := project.NewProjectService(constants.PROJECT_FILE_NAME, projectTable)

	// Commands
	replCommand := flag.NewFlagSet("repl", flag.ExitOnError)
	listCommand := flag.NewFlagSet(constants.LIST, flag.ExitOnError)
	addCommand := flag.NewFlagSet(constants.ADD, flag.ExitOnError)
	updateCommand := flag.NewFlagSet(constants.UPDATE, flag.ExitOnError)
	deleteCommand := flag.NewFlagSet(constants.DELETE, flag.ExitOnError)
	markCommand := flag.NewFlagSet(constants.MARK, flag.ExitOnError)

	// Flags for mark command of Project
	markProjectCompleted := markCommand.Bool("completed", false, "Mark task with status COMPLETED")
	markProjectStarted := markCommand.Bool("started", false, "Mark task with status STARTED")
	markProjectNotStarted := markCommand.Bool("not-started", false, "Mark task with status NOT_STARTED")

	// Flags for list command of Project
	listProjectCompleted := listCommand.Bool("completed", false, "List projects with status COMPLETED")
	listProjectStarted := listCommand.Bool("started", false, "List projects with status STARTED")
	listProjectNotStarted := listCommand.Bool("not-started", false, "List projects with status NOT_STARTED")

	// Flags for delete command of Project
	deleteAll := deleteCommand.Bool("all", false, "Delete all projects")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'add', 'list', 'update', 'delete', or 'mark' commands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case constants.ADD:
		addCommand.Parse(os.Args[2:])
		if addCommand.Parsed() {
			if len(addCommand.Args()) != 1 {
				fmt.Println("USAGE: add \"project name\"")
				return
			}
			if err := projectService.CreateProject(addCommand.Args()[0]); err != nil {
				fmt.Println("Error:", err)
			}
		}
	case constants.LIST:
		listCommand.Parse(os.Args[2:])
		if listCommand.Parsed() {
			var statusFilter project.ProjectStatus
			switch {
			case *listProjectCompleted:
				statusFilter = project.COMPLETED
			case *listProjectStarted:
				statusFilter = project.STARTED
			case *listProjectNotStarted:
				statusFilter = project.NOT_STARTED
			default:
				statusFilter = -1
			}
			if err := projectService.ListProjects(statusFilter); err != nil {
				fmt.Println("Error:", err)
			}
		}
	case constants.UPDATE:
		updateCommand.Parse(os.Args[2:])
		// if updateCommand.Parsed() {
		// 	if len(updateCommand.Args()) != 2 {
		// 		fmt.Println("USAGE: update <task_id> \"new task name\"")
		// 		return
		// 	}
		// 	if err := taskService.UpdateTaskTitle(updateCommand.Args()[0], updateCommand.Args()[1]); err != nil {
		// 		fmt.Println("Error:", err)
		// 	}
		// }
		if updateCommand.Parsed() {
			if len(updateCommand.Args()) != 2 {
				fmt.Println("USAGE: update <project_id> \"new project name\"")
				return
			}
			if err := projectService.UpdateProjectName(updateCommand.Args()[0], updateCommand.Args()[1]); err != nil {
				fmt.Println("Error:", err)
			}
		}
	case constants.DELETE:
		deleteCommand.Parse(os.Args[2:])
		// if deleteCommand.Parsed() {
		// 	if len(deleteCommand.Args()) != 1 && !*deleteAll {
		// 		fmt.Println("USAGE: delete <task_id> | --all")
		// 		return
		// 	}
		//
		// 	if *deleteAll {
		// 		if err := taskService.DeleteAllTasks(); err != nil {
		// 			fmt.Println("Error:", err)
		// 		}
		// 	} else {
		// 		if err := taskService.DeleteTask(deleteCommand.Args()[0]); err != nil {
		// 			fmt.Println("Error:", err)
		// 		}
		// 	}
		//
		// }
		if deleteCommand.Parsed() {
			if len(deleteCommand.Args()) != 1 && !*deleteAll {
				fmt.Println("USAGE: delete <project_id> | --all")
				return
			}

			if *deleteAll {
				if err := projectService.DeleteAllProjects(); err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				if err := projectService.DeleteProject(deleteCommand.Args()[0]); err != nil {
					fmt.Println("Error:", err)
				}
			}

		}
	case constants.MARK:
		markCommand.Parse(os.Args[2:])
		// if markCommand.Parsed() {
		// 	if len(markCommand.Args()) != 1 {
		// 		fmt.Println("USAGE: mark --done | --in-progress | --todo <task_id> ")
		// 		return
		// 	}
		// 	var status task.TaskStatus
		// 	switch {
		// 	case *markTaskDone:
		// 		status = task.DONE
		// 	case *markTaskInProgress:
		// 		status = task.IN_PROGRESS
		// 	case *markTaskTodo:
		// 		status = task.TODO
		// 	default:
		// 		fmt.Println("You must specify a task status using --done | --in-progress | --todo")
		// 		return
		// 	}
		// 	if err := taskService.UpdateTaskStatus(markCommand.Args()[0], status); err != nil {
		// 		fmt.Println("Error:", err)
		// 	}
		// }

		if markCommand.Parsed() {
			if len(markCommand.Args()) != 1 {
				fmt.Println("USAGE: mark --completed | --started | --not-completed <project_id> ")
				return
			}
			var status project.ProjectStatus
			switch {
			case *markProjectCompleted:
				status = project.COMPLETED
			case *markProjectStarted:
				status = project.STARTED
			case *markProjectNotStarted:
				status = project.NOT_STARTED
			default:
				fmt.Println("You must specify a task status using --done | --in-progress | --todo")
				return
			}
			if err := projectService.UpdateProjectStatus(markCommand.Args()[0], status); err != nil {
				fmt.Println("Error:", err)
			}
		}
	case "repl":
		replCommand.Parse(os.Args[2:])
		if replCommand.Parsed() {
			if len(replCommand.Args()) != 1 {
				fmt.Println("USAGE: repl <project_id>")
				return
			}

			if isProjectExists := projectService.IsProjectExists(replCommand.Args()[0]); !isProjectExists {
				fmt.Println("Project with ID=", replCommand.Args()[0], "not found")
				os.Exit(1)
			}

			startREPL(replCommand.Args()[0])
		}
	default:
		fmt.Println("Expected 'add', 'list', 'update', 'delete', or 'mark' commands")
		os.Exit(1)
	}
}
