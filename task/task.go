package task

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MuradIsayev/todo-tracker/constants"
	"github.com/MuradIsayev/todo-tracker/helpers"
	"github.com/MuradIsayev/todo-tracker/project"
	"github.com/olekukonko/tablewriter"
)

type TaskStatus int

const (
	TODO TaskStatus = iota
	IN_PROGRESS
	DONE
)

type Task struct {
	Id             int        `json:"id"`
	Name           string     `json:"name"`
	Status         TaskStatus `json:"status"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	TotalSpentTime int        `json:"totalSpentTime"`
	ProjectId      int        `json:"projectId"`
}

type TaskService struct {
	filePath       string
	table          *tablewriter.Table
	projectService *project.ProjectService
}

func NewTaskService(projectService *project.ProjectService, projectId string, table *tablewriter.Table) *TaskService {
	table.SetHeader([]string{constants.COLUMN_ID, constants.COLUMN_NAME, constants.COLUMN_STATUS, constants.COLUMN_CREATE_DATE, constants.COLUMN_UPDATE_DATE, constants.COLUMN_TOTAL_SPENT_TIME_TASK})

	filePath := fmt.Sprintf("output/tasks/%s_%s", projectId, constants.TASK_FILE_NAME)

	return &TaskService{
		filePath:       filePath,
		table:          table,
		projectService: projectService,
	}
}

func (s *TaskService) UpdateTaskSpentTime(id int, spentTime int) error {
	tasks, err := s.readTasksFromFile()
	if err != nil {
		return err
	}

	index, task, err := findTaskById(tasks, id)
	if err != nil {
		return err
	}

	// also update the project's total spent time
	s.projectService.UpdateTotalSpentTime(task.ProjectId, spentTime)
	task.TotalSpentTime += spentTime
	if task.Status != DONE && task.TotalSpentTime > 0 && task.Status != IN_PROGRESS {
		task.Status = IN_PROGRESS
	}
	// task.UpdatedAt = time.Now()
	tasks[index] = *task

	return s.writeTasksToFile(tasks)
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

func findTaskById(tasks []Task, id int) (int, *Task, error) {
	for i, task := range tasks {
		if task.Id == id {
			return i, &task, nil
		}
	}
	return -1, nil, fmt.Errorf("task with ID=%d not found", id)
}

func (s *TaskService) FindTaskById(id string) (*Task, error) {
	taskId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return nil, err
	}
	tasks, err := s.readTasksFromFile()
	if err != nil {
		return nil, err
	}

	_, task, err := findTaskById(tasks, taskId)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) UpdateTaskStatus(id string, taskStatus TaskStatus) error {
	taskId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return err
	}

	tasks, err := s.readTasksFromFile()
	if err != nil {
		return err
	}

	index, task, err := findTaskById(tasks, taskId)
	if err != nil {
		return err
	}

	task.Status = taskStatus
	task.UpdatedAt = time.Now()
	tasks[index] = *task

	return s.writeTasksToFile(tasks)
}
func (s *TaskService) UpdateTaskName(id, name string) error {
	taskId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return err
	}

	tasks, err := s.readTasksFromFile()
	if err != nil {
		return err
	}

	index, task, err := findTaskById(tasks, taskId)
	if err != nil {
		return err
	}

	if name != "" {
		task.Name = name
		task.UpdatedAt = time.Now()
		tasks[index] = *task
	}

	return s.writeTasksToFile(tasks)
}

func (s *TaskService) DeleteTask(id string) error {
	taskId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return err
	}

	tasks, err := s.readTasksFromFile()
	if err != nil {
		return err
	}

	index, _, err := findTaskById(tasks, taskId)
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
	s.table.ClearRows()
	s.table.ClearFooter()

	tasks, err := s.readTasksFromFile()
	if err != nil {
		return err
	}

	var nbOfLeftTasks int

	for _, task := range tasks {
		if statusFilter == -1 || task.Status == statusFilter {
			formatSpendTime := helpers.FormatSpendTime(task.TotalSpentTime)
			createdAt := task.CreatedAt.Format(constants.DATE_FORMAT)
			updatedAt := task.UpdatedAt.Format(constants.DATE_FORMAT)

			s.table.Append([]string{strconv.Itoa(task.Id), task.Name, task.Status.String(), createdAt, updatedAt, formatSpendTime})

			if task.Status == TODO {
				nbOfLeftTasks++
			}
		}
	}

	s.table.SetRowLine(true)
	s.table.SetFooter([]string{"", "", "", "", " ", defineFooterText(nbOfLeftTasks, len(tasks))})
	s.table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)
	s.table.SetFooterColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})

	s.table.Render()

	return nil
}

func (s *TaskService) CreateTask(projectID string, name string) error {
	tasks, err := s.readTasksFromFile()
	if err != nil {
		return err
	}

	projectId, err := helpers.ValidateIdAndConvertToInt(projectID)
	if err != nil {
		return err
	}

	task := Task{
		Id:             s.getNextID(tasks),
		Name:           name,
		Status:         TODO,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		TotalSpentTime: 0,
		ProjectId:      projectId,
	}

	tasks = append(tasks, task)

	return s.writeTasksToFile(tasks)
}
