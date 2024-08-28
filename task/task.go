package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/MuradIsayev/todo-tracker/constants"
	"github.com/olekukonko/tablewriter"
)

type TaskStatus int

const (
	TODO TaskStatus = iota
	IN_PROGRESS
	DONE
)

type Task struct {
	Id        int        `json:"id"`
	Title     string     `json:"title"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type TaskService struct {
	filePath string
	table    *tablewriter.Table
}

func NewTaskService(id string, table *tablewriter.Table) *TaskService {
	filePath := fmt.Sprintf("%s_%s", id, constants.TASK_FILE_NAME)
	table.SetHeader([]string{"ID", "Title", "Status", "Create Date", "Update Date"})

	return &TaskService{
		filePath: filePath,
		table:    table,
	}
}

func (taskStatus TaskStatus) String() string {
	switch taskStatus {
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

func (s *TaskService) readTasksFromFile() ([]Task, error) {
	fileContent, err := os.ReadFile(s.filePath)
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

func (s *TaskService) writeTasksToFile(tasks []Task) error {
	newTasks, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("Cannot convert Struct to JSON: %v", err)
	}

	if err := os.WriteFile(s.filePath, newTasks, 0644); err != nil {
		return fmt.Errorf("cannot write to file: %v", err)
	}
	return nil
}

func (s *TaskService) getNextID(tasks []Task) int {
	if len(tasks) == 0 {
		return 1
	}
	return tasks[len(tasks)-1].Id + 1
}

func (s *TaskService) validateID(id string) (int, error) {
	var numberRegex = regexp.MustCompile(`^[0-9]+$`)

	if !numberRegex.MatchString(id) {
		return 0, errors.New("ID must only contain digits")
	}

	return strconv.Atoi(id)
}

func (s *TaskService) findTaskById(tasks []Task, id int) (int, *Task, error) {
	for i, task := range tasks {
		if task.Id == id {
			return i, &task, nil
		}
	}
	return -1, nil, fmt.Errorf("task with ID=%d not found", id)
}

func (s *TaskService) UpdateTaskStatus(id string, taskStatus TaskStatus) error {
	taskId, err := s.validateID(id)
	if err != nil {
		return err
	}

	tasks, err := s.readTasksFromFile()
	if err != nil {
		return err
	}

	index, task, err := s.findTaskById(tasks, taskId)
	if err != nil {
		return err
	}

	task.Status = taskStatus
	task.UpdatedAt = time.Now()
	tasks[index] = *task

	return s.writeTasksToFile(tasks)
}
func (s *TaskService) UpdateTaskTitle(id, title string) error {
	taskId, err := s.validateID(id)
	if err != nil {
		return err
	}

	tasks, err := s.readTasksFromFile()
	if err != nil {
		return err
	}

	index, task, err := s.findTaskById(tasks, taskId)
	if err != nil {
		return err
	}

	if title != "" {
		task.Title = title
		task.UpdatedAt = time.Now()
		tasks[index] = *task
	}

	return s.writeTasksToFile(tasks)
}

func (s *TaskService) DeleteTask(id string) error {
	taskId, err := s.validateID(id)
	if err != nil {
		return err
	}

	tasks, err := s.readTasksFromFile()
	if err != nil {
		return err
	}

	index, _, err := s.findTaskById(tasks, taskId)
	if err != nil {
		return err
	}

	tasks = append(tasks[:index], tasks[index+1:]...)

	return s.writeTasksToFile(tasks)
}

func (s *TaskService) DeleteAllTasks() error {
	var tasks []Task

	return s.writeTasksToFile(tasks)
}

func defineFooterText(nbOfLeftTasks, nbOfTotalTasks int) string {
	if nbOfLeftTasks == 0 && nbOfTotalTasks == 0 {
		return "No tasks found"
	}

	return fmt.Sprintf("Left tasks: %d", nbOfLeftTasks)
}

func (s *TaskService) ListTasks(statusFilter TaskStatus) error {
	tasks, err := s.readTasksFromFile()
	if err != nil {
		return err
	}

	var nbOfLeftTasks int

	for _, task := range tasks {
		if statusFilter == -1 || task.Status == statusFilter {
			createdAt := task.CreatedAt.Format(constants.DATE_FORMAT)
			updatedAt := task.UpdatedAt.Format(constants.DATE_FORMAT)
			s.table.Append([]string{strconv.Itoa(task.Id), task.Title, task.Status.String(), createdAt, updatedAt})

			if task.Status == TODO {
				nbOfLeftTasks++
			}
		}
	}

	s.table.SetRowLine(true)
	s.table.SetFooter([]string{"", "", "", " ", defineFooterText(nbOfLeftTasks, len(tasks))})
	s.table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)
	s.table.SetFooterColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})
	// s.table.SetCaption(true, "Tasks List")

	s.table.Render()

	return nil
}

func (s *TaskService) CreateTask(title string) error {
	tasks, err := s.readTasksFromFile()
	if err != nil {
		return err
	}

	task := Task{
		Id:        s.getNextID(tasks),
		Title:     title,
		Status:    TODO,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tasks = append(tasks, task)

	return s.writeTasksToFile(tasks)
}
