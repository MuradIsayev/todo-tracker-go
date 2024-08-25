package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

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

type TaskService struct {
	filePath string
}

func NewTaskService(filePath string) *TaskService {
	return &TaskService{filePath: filePath}
}

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

func getNextID(tasks []Task) int {
	if len(tasks) == 0 {
		return 1
	}
	return tasks[len(tasks)-1].Id + 1
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
		fmt.Println("Title is empty", task.UpdatedAt)
		task.Title = title
		task.UpdatedAt = time.Now()
		fmt.Println("2 Title is empty", task.UpdatedAt)
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

func (s *TaskService) ListTasks(statusFilter TaskStatus) error {
	tasks, err := s.readTasksFromFile()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if statusFilter == -1 || task.Status == statusFilter {
			createdAt := task.CreatedAt.Format("15:04:05 (Mon) - 02/01/2006")
			updatedAt := task.UpdatedAt.Format("15:04:05 (Mon) - 02/01/2006")
			fmt.Printf("ID: %d | Title: %s | Status: %s | CreatedAt: %s | UpdatedAt: %s\n", task.Id, task.Title, task.Status, createdAt, updatedAt)
		}
	}

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
