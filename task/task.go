package task

import (
	"fmt"
	"strconv"
	"time"

	"github.com/MuradIsayev/todo-tracker/base"
	"github.com/MuradIsayev/todo-tracker/constants"
	"github.com/MuradIsayev/todo-tracker/helpers"
	"github.com/MuradIsayev/todo-tracker/project"
	"github.com/MuradIsayev/todo-tracker/status"
	"github.com/olekukonko/tablewriter"
)

type Task struct {
	Id             int               `json:"id"`
	Name           string            `json:"name"`
	Status         status.ItemStatus `json:"status"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
	TotalSpentTime int               `json:"totalSpentTime"`
	ProjectId      int               `json:"projectId"`
}

type TaskService struct {
	baseService    *base.BaseService[Task]
	table          *tablewriter.Table
	projectService *project.ProjectService
}

func NewTaskService(projectService *project.ProjectService, projectId string, table *tablewriter.Table) *TaskService {
	filePath := fmt.Sprintf("output/tasks/%s_%s", projectId, constants.TASK_FILE_NAME)

	table.SetHeader([]string{constants.COLUMN_ID, constants.COLUMN_NAME, constants.COLUMN_STATUS, constants.COLUMN_CREATE_DATE, constants.COLUMN_UPDATE_DATE, constants.COLUMN_TOTAL_SPENT_TIME_TASK})

	return &TaskService{
		table:          table,
		projectService: projectService,
		baseService: &base.BaseService[Task]{
			FilePath: filePath,
		},
	}
}

func (s *TaskService) UpdateTaskSpentTime(id int, spentTime int) error {
	tasks := []Task{}
	err := s.baseService.ReadFromFile(&tasks)
	if err != nil {
		return err
	}

	index, task, err := findTaskById(tasks, id)
	if err != nil {
		return err
	}

	s.projectService.UpdateTotalSpentTime(task.ProjectId, spentTime)
	task.TotalSpentTime += spentTime
	if task.Status != status.DONE && task.TotalSpentTime > 0 && task.Status != status.IN_PROGRESS {
		task.Status = status.IN_PROGRESS
	}

	tasks[index] = *task
	return s.baseService.WriteToFile(tasks)
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

	tasks := []Task{}
	err = s.baseService.ReadFromFile(&tasks)
	if err != nil {
		return nil, err
	}

	_, task, err := findTaskById(tasks, taskId)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) UpdateTaskStatus(id string, taskStatus status.ItemStatus) error {
	taskId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return err
	}

	tasks := []Task{}
	err = s.baseService.ReadFromFile(&tasks)
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

	return s.baseService.WriteToFile(tasks)
}

func (s *TaskService) UpdateTaskName(id, name string) error {
	taskId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return err
	}

	tasks := []Task{}
	err = s.baseService.ReadFromFile(&tasks)
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
	return s.baseService.WriteToFile(tasks)
}

func (s *TaskService) DeleteTask(id string) error {
	taskId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return err
	}

	tasks := []Task{}
	err = s.baseService.ReadFromFile(&tasks)
	if err != nil {
		return err
	}

	index, _, err := findTaskById(tasks, taskId)
	if err != nil {
		return err
	}

	tasks = append(tasks[:index], tasks[index+1:]...)

	return s.baseService.WriteToFile(tasks)
}

func (s *TaskService) DeleteAllTasks() error {
	return s.baseService.DeleteAllItems()
}

func defineTableFooterText(nbOfLeftTasks, nbOfTotalTasks int) string {
	if nbOfLeftTasks == 0 && nbOfTotalTasks == 0 {
		return "No tasks found"
	}

	return fmt.Sprintf("Left tasks: %d", nbOfLeftTasks)
}

func (s *TaskService) ListTasks(statusFilter status.ItemStatus) error {
	s.table.ClearRows()
	s.table.ClearFooter()

	tasks := []Task{}
	err := s.baseService.ReadFromFile(&tasks)
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

			if task.Status == status.TODO {
				nbOfLeftTasks++
			}
		}
	}

	s.table.SetRowLine(true)
	s.table.SetFooter([]string{"", "", "", "", " ", defineTableFooterText(nbOfLeftTasks, len(tasks))})
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
	tasks := []Task{}
	err := s.baseService.ReadFromFile(&tasks)

	if err != nil {
		return err
	}

	projectId, err := helpers.ValidateIdAndConvertToInt(projectID)
	if err != nil {
		return err
	}

	task := Task{
		Id:             s.baseService.GetNextID(tasks),
		Name:           name,
		Status:         status.TODO,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		TotalSpentTime: 0,
		ProjectId:      projectId,
	}

	tasks = append(tasks, task)

	return s.baseService.WriteToFile(tasks)
}
