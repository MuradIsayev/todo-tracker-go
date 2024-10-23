package main

/*

TODO: Think of how to improve error handling.
TODO: Add more features like sorting, filtering, and searching.
TODO: Make it more user friendly for the production step.
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
	"github.com/MuradIsayev/todo-tracker/countdown"
	"github.com/MuradIsayev/todo-tracker/project"
	"github.com/MuradIsayev/todo-tracker/service"
	"github.com/MuradIsayev/todo-tracker/status"
	"github.com/MuradIsayev/todo-tracker/task"
	"github.com/olekukonko/tablewriter"
)

func startREPL(
	circularDependencyManager *service.Manager,
	projectId string,
	taskService *task.TaskService,
) {
	fmt.Println("Welcome to the Task Management CLI - Interactive Mode")
	reader := bufio.NewReader(os.Stdin)

	// add projectID to the task service
	taskService.AddProjectIdToTaskService(projectId)

	for {
		fmt.Print(">>> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" || input == "quit" {
			break
		}

		executeCommand(input, projectId, taskService, circularDependencyManager)
	}

}

func executeCommand(input string, projectId string, taskService *task.TaskService, circularDependencyManager *service.Manager) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case constants.ADD:
		handleAddCommand(args, projectId, taskService)
	case constants.LIST:
		handleListCommand(args, taskService)
	case constants.UPDATE:
		handleUpdateCommand(args, taskService)
	case constants.DELETE:
		handleDeleteCommand(args, taskService)
	case constants.MARK:
		handleMarkCommand(args, taskService)
	case constants.TIMER:
		handleCountdownCommand(args, taskService, circularDependencyManager)
	default:
		fmt.Println("Unknown command:", command)
	}
}

func handleCountdownCommand(args []string, taskService *task.TaskService, circularDependencyManager *service.Manager) {
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
	countdownService := countdown.NewCountdownService(taskService, circularDependencyManager)

	task, err := taskService.FindTaskById(taskID)
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup

	// Start the countdown in a separate goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		countdownService.StartCountdown(task, *timePtr)
	}()

	// Display timer updates without interrupting input
	wg.Add(1)
	go func() {
		defer wg.Done()
		for displayMsg := range countdownService.DisplayChan {
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
		case <-countdownService.DoneChan:
			return
		default:
			// Handle user commands from input
			switch input {
			case constants.TIMER_PAUSE:
				countdownService.PauseChan <- true
			case constants.TIMER_RESUME:
				countdownService.ResumeChan <- true
			case constants.TIMER_STOP:
				countdownService.StopChan <- true // Stop and update the task time
				return
			case constants.TIMER_EXIT:
				fmt.Println("Exiting timer mode without saving time.")
				countdownService.ExitChan <- true // Exit without updating the task time
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

func handleAddCommand(args []string, projectId string, taskService *task.TaskService) {
	if len(args) < 1 {
		fmt.Println("USAGE: add <task_name>")
		return
	}

	taskName := strings.Join(args, " ")
	if err := taskService.CreateTask(projectId, taskName); err != nil {
		fmt.Println("Error:", err)
	}
}

func handleListCommand(args []string, taskService *task.TaskService) {
	listCommand := flag.NewFlagSet(constants.LIST, flag.ExitOnError)
	listDone := listCommand.Bool("done", false, "List tasks with status DONE")
	listInProgress := listCommand.Bool("in-progress", false, "List tasks with status IN_PROGRESS")
	listTodo := listCommand.Bool("todo", false, "List tasks with status TODO")

	listCommand.Parse(args)

	var statusFilter status.ItemStatus
	switch {
	case *listDone:
		statusFilter = status.DONE
	case *listInProgress:
		statusFilter = status.IN_PROGRESS
	case *listTodo:
		statusFilter = status.TODO
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

	if err := taskService.UpdateTaskName(args[0], args[1]); err != nil {
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
	var statusFilter status.ItemStatus
	switch {
	case *markDone:
		statusFilter = status.DONE
	case *markInProgress:
		statusFilter = status.IN_PROGRESS
	case *markTodo:
		statusFilter = status.TODO
	default:
		fmt.Println("Invalid status. Use --done, --in-progress, or --todo.")
		return
	}

	if err := taskService.UpdateTaskStatus(taskID, statusFilter); err != nil {
		fmt.Println("Error:", err)
	}
}

func main() {
	projectTable := tablewriter.NewWriter(os.Stdout)
	projectService := project.NewProjectService(constants.PROJECT_FILE_NAME, projectTable)

	taskTable := tablewriter.NewWriter(os.Stdout)
	taskService := task.NewTaskService(projectService, taskTable)

	circularDependencyManager := service.NewManager(taskService, projectService)

	if len(os.Args) < 2 {
		fmt.Println("Expected 'add', 'list', 'update', 'delete', 'mark', or 'repl' commands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case constants.REPL:
		handleREPLCommand(os.Args[2:], projectService, taskService, circularDependencyManager)
	case constants.ADD:
		handleProjectAddCommand(os.Args[2:], projectService)
	case constants.LIST:
		handleProjectListCommand(os.Args[2:], projectService)
	case constants.UPDATE:
		handleProjectUpdateCommand(os.Args[2:], projectService)
	case constants.DELETE:
		handleProjectDeleteCommand(os.Args[2:], circularDependencyManager)
	case constants.MARK:
		handleProjectMarkCommand(os.Args[2:], projectService)
	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}

func handleREPLCommand(args []string, projectService *project.ProjectService, taskService *task.TaskService, circularDependencyManager *service.Manager) {
	if len(args) != 1 {
		fmt.Println("USAGE: repl <project_id>")
		return
	}

	projectID := args[0]
	if !projectService.IsProjectExists(projectID) {
		fmt.Println("Project with ID=", projectID, "not found")
		os.Exit(1)
	}

	startREPL(
		// projectService,
		circularDependencyManager,
		projectID,
		taskService,
	)
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

	listDone := listCommand.Bool("done", false, "List projects with status DONE")
	listInProgress := listCommand.Bool("in-progress", false, "List projects with status IN_PROGRESS")
	listTodo := listCommand.Bool("todo", false, "List projects with status TODO")

	listCommand.Parse(args)

	var statusFilter status.ItemStatus
	switch {
	case *listDone:
		statusFilter = status.DONE
	case *listInProgress:
		statusFilter = status.IN_PROGRESS
	case *listTodo:
		statusFilter = status.TODO
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

func handleProjectDeleteCommand(args []string, circularDependencyManager *service.Manager) {
	deleteCommand := flag.NewFlagSet(constants.DELETE, flag.ExitOnError)
	deleteAll := deleteCommand.Bool("all", false, "Delete all projects")
	deleteCommand.Parse(args)

	if *deleteAll {
		fmt.Println("Are you sure you want to delete all projects? This action will also delete all tasks. (y/n)")
		reader := bufio.NewReader(os.Stdin)

		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)

		if response != "y" {
			fmt.Println("Command cancelled.")
			return
		}

		if err := circularDependencyManager.DeleteAllProjectWithAllTasks(); err != nil {
			fmt.Println("Error:", err)
		}

		return
	}

	if len(deleteCommand.Args()) != 1 {
		fmt.Println("USAGE: delete <project_id> | --all")
		return
	}

	if err := circularDependencyManager.DeleteProjectAndCorrespondingTasks(deleteCommand.Args()[0]); err != nil {
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
	markProjectDone := markCommand.Bool("done", false, "Mark project as DONE")
	markProjectInProgress := markCommand.Bool("in-progress", false, "Mark project as IN_PROGRESS")
	markProjectTodo := markCommand.Bool("todo", false, "Mark project as TODO")

	// Parse the remaining args after extracting the task ID
	err := markCommand.Parse(args[1:])
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return
	}

	var statusFilter status.ItemStatus
	switch {
	case *markProjectDone:
		statusFilter = status.DONE
	case *markProjectInProgress:
		statusFilter = status.IN_PROGRESS
	case *markProjectTodo:
		statusFilter = status.TODO
	default:
		fmt.Println("Invalid status. Use --completed, --started, or --not-started.")
		return
	}

	if err := projectService.UpdateProjectStatus(projectID, statusFilter); err != nil {
		fmt.Println("Error:", err)
	}
}
