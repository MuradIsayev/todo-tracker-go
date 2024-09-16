package main

/*

TODO: Make it more user friendly for the production step.
TODO: Think of how to improve error handling.
TODO: Reduce the repetition in the code by creating helper functions and using struct embedding more effectively.
TODO: Add more features like sorting, filtering, and searching.
TODO: Make more strucutred folder and file management for the task and project json files.
TODO: Think of storing the PROJECT_ID in the task json file to make it easier to find the tasks for a specific project.
TODO: Maybee add unit tests?
FIX: Fix the issue with the command line input not working properly when the timer is done.
*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/MuradIsayev/todo-tracker/constants"
	"github.com/MuradIsayev/todo-tracker/project"
	"github.com/MuradIsayev/todo-tracker/task"
	"github.com/olekukonko/tablewriter"
)

func startREPL(projectService *project.ProjectService, projectID string) {
	fmt.Println("Welcome to the Task Management CLI - Interactive Mode")
	reader := bufio.NewReader(os.Stdin)

	taskTable := tablewriter.NewWriter(os.Stdout)
	taskService := task.NewTaskService(projectService, projectID, taskTable)

	for {
		fmt.Print(">>> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" || input == "quit" {
			break
		}

		executeCommand(input, projectID, taskService)
	}

}

func executeCommand(input string, projectID string, taskService *task.TaskService) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case constants.ADD:
		handleAddCommand(args, taskService)
	case constants.LIST:
		handleListCommand(args, taskService)
	case constants.UPDATE:
		handleUpdateCommand(args, taskService)
	case constants.DELETE:
		handleDeleteCommand(args, taskService)
	case constants.MARK:
		handleMarkCommand(args, taskService)
	case constants.TIMER:
		handleCountdownCommand(args, projectID, taskService)
	default:
		fmt.Println("Unknown command:", command)
	}
}

func handleCountdownCommand(args []string, projectID string, taskService *task.TaskService) {
	if len(args) != 3 {
		fmt.Println("USAGE: t <task_id> --time <duration>")
		return
	}

	countdownCommand := flag.NewFlagSet("TIMER", flag.ExitOnError)
	timePtr := countdownCommand.Int("time", 1, "Specify the countdown duration in minutes")

	if err := countdownCommand.Parse(args[1:]); err != nil {
		fmt.Println("Error parsing countdown command:", err)
		return
	}

	taskID := args[0]
	controller := task.NewCountdownController()

	task, err := taskService.FindTaskById(taskID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var wg sync.WaitGroup

	// Start the countdown in a separate goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		taskService.StartCountdown(task, projectID, *timePtr, controller)
	}()

	// Display timer updates without interrupting input
	wg.Add(1)
	go func() {
		defer wg.Done()
		for displayMsg := range controller.DisplayChan {
			fmt.Print("\0337")                           // Save cursor position
			fmt.Printf("\033[1A\033[2K\r%s", displayMsg) // Clear timer line, print new time
			fmt.Print("\0338")                           // Restore cursor position
		}
	}()

	// Start a goroutine to read input commands
	reader := bufio.NewReader(os.Stdin)
	printControls()

	for {
		fmt.Print("T> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		input = strings.TrimSpace(strings.ToLower(input))

		select {
		case <-controller.DoneChan:
			return
		default:
			// Handle user commands from input
			switch input {
			case "p":
				controller.PauseChan <- true
			case "r":
				controller.ResumeChan <- true
			case "s":
				controller.StopChan <- true // Stop and update the task time
				return
			case "e":
				fmt.Println("Exiting timer mode without saving time.")
				controller.ExitChan <- true // Exit without updating the task time
				return
			default:
				fmt.Print("\0337")                                                                         // Save cursor position
				fmt.Printf("\033[2A\033[2K\rUnknown command. Use (p)ause, (r)esume, (s)top, or (e)xit.\n") // Clear and print controls
				fmt.Print("\0338")                                                                         // Restore cursor position
			}
		}
	}
}

// Function to consistently display the controls
func printControls() {
	fmt.Print("\0337")                                                                                              // Save cursor position
	fmt.Printf("\033[2A\033[2K\rControls: type (p)ause, (r)esume, (s)top, or (e)xit to control the countdown.\n\n") // Clear and print controls
	fmt.Print("\0338")                                                                                              // Restore cursor position
}

func handleAddCommand(args []string, taskService *task.TaskService) {
	if len(args) < 1 {
		fmt.Println("USAGE: add <task_name>")
		return
	}

	taskTitle := strings.Join(args, " ")
	if err := taskService.CreateTask(taskTitle); err != nil {
		fmt.Println("Error:", err)
	}
}

func handleListCommand(args []string, taskService *task.TaskService) {
	listCommand := flag.NewFlagSet(constants.LIST, flag.ExitOnError)
	listDone := listCommand.Bool("done", false, "List tasks with status DONE")
	listInProgress := listCommand.Bool("in-progress", false, "List tasks with status IN_PROGRESS")
	listTodo := listCommand.Bool("todo", false, "List tasks with status TODO")

	listCommand.Parse(args)

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

func handleUpdateCommand(args []string, taskService *task.TaskService) {
	if len(args) != 2 {
		fmt.Println("USAGE: update <task_id> \"new task name\"")
		return
	}

	if err := taskService.UpdateTaskTitle(args[0], args[1]); err != nil {
		fmt.Println("Error:", err)
	}
}

func handleDeleteCommand(args []string, taskService *task.TaskService) {
	deleteCommand := flag.NewFlagSet(constants.DELETE, flag.ExitOnError)
	deleteAll := deleteCommand.Bool("all", false, "Delete all tasks")
	deleteCommand.Parse(args)

	if *deleteAll {
		if err := taskService.DeleteAllTasks(); err != nil {
			fmt.Println("Error:", err)
		}
		return
	}

	if len(deleteCommand.Args()) != 1 {
		fmt.Println("USAGE: delete <task_id> | --all")
		return
	}

	if err := taskService.DeleteTask(deleteCommand.Args()[0]); err != nil {
		fmt.Println("Error:", err)
	}
}

func handleMarkCommand(args []string, taskService *task.TaskService) {
	if len(args) != 2 {
		fmt.Println("USAGE: mark <task_id> --done | --in-progress | --todo")
		return
	}

	taskID := args[0]
	// Create a flag set for parsing the status flags
	markCommand := flag.NewFlagSet("mark", flag.ExitOnError)

	// Define the flags
	markDone := markCommand.Bool("done", false, "Mark task as DONE")
	markInProgress := markCommand.Bool("in-progress", false, "Mark task as IN_PROGRESS")
	markTodo := markCommand.Bool("todo", false, "Mark task as TODO")

	// Parse the remaining args after extracting the task ID
	err := markCommand.Parse(args[1:])
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return
	}

	// Determine the task status based on the parsed flags
	var status task.TaskStatus
	switch {
	case *markDone:
		status = task.DONE
	case *markInProgress:
		status = task.IN_PROGRESS
	case *markTodo:
		status = task.TODO
	default:
		fmt.Println("Invalid status. Use --done, --in-progress, or --todo.")
		return
	}

	if err := taskService.UpdateTaskStatus(taskID, status); err != nil {
		fmt.Println("Error:", err)
	}
}

func main() {
	projectTable := tablewriter.NewWriter(os.Stdout)
	projectService := project.NewProjectService(constants.PROJECT_FILE_NAME, projectTable)

	if len(os.Args) < 2 {
		fmt.Println("Expected 'add', 'list', 'update', 'delete', 'mark', or 'repl' commands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case constants.REPL:
		handleREPLCommand(os.Args[2:], projectService)
	case constants.ADD:
		handleProjectAddCommand(os.Args[2:], projectService)
	case constants.LIST:
		handleProjectListCommand(os.Args[2:], projectService)
	case constants.UPDATE:
		handleProjectUpdateCommand(os.Args[2:], projectService)
	case constants.DELETE:
		handleProjectDeleteCommand(os.Args[2:], projectService)
	case constants.MARK:
		handleProjectMarkCommand(os.Args[2:], projectService)
	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}

func handleREPLCommand(args []string, projectService *project.ProjectService) {
	if len(args) != 1 {
		fmt.Println("USAGE: repl <project_id>")
		return
	}

	projectID := args[0]
	if !projectService.IsProjectExists(projectID) {
		fmt.Println("Project with ID=", projectID, "not found")
		os.Exit(1)
	}

	startREPL(projectService, projectID)
}

func handleProjectAddCommand(args []string, projectService *project.ProjectService) {
	if len(args) < 1 {
		fmt.Println("USAGE: add <project_name>")
		return
	}

	projectName := strings.Join(args, " ")
	if err := projectService.CreateProject(projectName); err != nil {
		fmt.Println("Error:", err)
	}
}

func handleProjectListCommand(args []string, projectService *project.ProjectService) {
	listCommand := flag.NewFlagSet(constants.LIST, flag.ExitOnError)
	listCompleted := listCommand.Bool("completed", false, "List projects with status COMPLETED")
	listStarted := listCommand.Bool("started", false, "List projects with status STARTED")
	listNotStarted := listCommand.Bool("not-started", false, "List projects with status NOT_STARTED")

	listCommand.Parse(args)

	var statusFilter project.ProjectStatus
	switch {
	case *listCompleted:
		statusFilter = project.COMPLETED
	case *listStarted:
		statusFilter = project.STARTED
	case *listNotStarted:
		statusFilter = project.NOT_STARTED
	default:
		statusFilter = -1
	}

	if err := projectService.ListProjects(statusFilter); err != nil {
		fmt.Println("Error:", err)
	}
}

func handleProjectUpdateCommand(args []string, projectService *project.ProjectService) {
	if len(args) != 2 {
		fmt.Println("USAGE: update <project_id> \"new project name\"")
		return
	}

	if err := projectService.UpdateProjectName(args[0], args[1]); err != nil {
		fmt.Println("Error:", err)
	}
}

func handleProjectDeleteCommand(args []string, projectService *project.ProjectService) {
	deleteCommand := flag.NewFlagSet(constants.DELETE, flag.ExitOnError)
	deleteAll := deleteCommand.Bool("all", false, "Delete all projects")
	deleteCommand.Parse(args)

	if *deleteAll {
		if err := projectService.DeleteAllProjects(); err != nil {
			fmt.Println("Error:", err)
		}
		return
	}

	if len(deleteCommand.Args()) != 1 {
		fmt.Println("USAGE: delete <project_id> | --all")
		return
	}

	if err := projectService.DeleteProject(deleteCommand.Args()[0]); err != nil {
		fmt.Println("Error:", err)
	}
}

func handleProjectMarkCommand(args []string, projectService *project.ProjectService) {
	if len(args) != 2 {
		fmt.Println("USAGE: mark <project_id> --completed | --started | --not-started")
		return
	}

	projectID := args[0]
	// Create a flag set for parsing the status flags
	markCommand := flag.NewFlagSet("mark", flag.ExitOnError)

	// Define the flags
	markProjectCompleted := markCommand.Bool("completed", false, "Mark task as COMPLETED")
	markProjectStarted := markCommand.Bool("started", false, "Mark task as STARTED")
	markProjectNotStarted := markCommand.Bool("not-started", false, "Mark task as NOT_STARTED")

	// Parse the remaining args after extracting the task ID
	err := markCommand.Parse(args[1:])
	if err != nil {
		fmt.Println("Error parsing flags:", err)
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
		fmt.Println("Invalid status. Use --completed, --started, or --not-started.")
		return
	}

	if err := projectService.UpdateProjectStatus(projectID, status); err != nil {
		fmt.Println("Error:", err)
	}
}
