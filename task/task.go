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

type TaskManager interface {
	DeleteTasksByProjectId(projectId string) error
	DeleteAllTasks() error
	UpdateTaskTimer(taskId int, newDuration int) error
}

func (s *TaskService) DeleteAllTasks() error {
	return helpers.RemoveContentsOfDirectory("output/tasks")
}

func (t *TaskService) DeleteTasksByProjectId(projectId string) error {
	filePathOfTask := fmt.Sprintf("output/tasks/%s_%s", projectId, constants.TASK_FILE_NAME)

	return helpers.RemoveFileByFilePath(filePathOfTask)
}

func (t *TaskService) UpdateTaskTimer(taskId int, newDuration int) error {
	return t.baseService.UpdateTotalSpentTime(taskId, newDuration)
}

type TaskService struct {
	baseService    *base.BaseService[Task]
	table          *tablewriter.Table
	projectService *project.ProjectService
}

func NewTaskService(projectService *project.ProjectService, table *tablewriter.Table) *TaskService {
	table.SetHeader([]string{constants.COLUMN_ID, constants.COLUMN_NAME, constants.COLUMN_STATUS, constants.COLUMN_CREATE_DATE, constants.COLUMN_UPDATE_DATE, constants.COLUMN_TOTAL_SPENT_TIME_TASK})

	return &TaskService{
		table:          table,
		projectService: projectService,
	}
}

func (s *TaskService) AddProjectIdToTaskService(projectId string) *TaskService {
	filePath := fmt.Sprintf("output/tasks/%s_%s", projectId, constants.TASK_FILE_NAME)

	s.baseService = &base.BaseService[Task]{
		FilePath: filePath,
	}

	return s

}

func (s *TaskService) UpdateTaskSpentTime(id int, spentTime int) error {
	// find task by id
	task, err := s.FindTaskById(strconv.Itoa(id))
	if err != nil {
		return err
	}

	s.projectService.UpdateProjectSpentTime(task.ProjectId, spentTime)

	return s.baseService.UpdateTotalSpentTime(id, spentTime)
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

	_, task, err := s.baseService.FindItemById(tasks, taskId)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) UpdateTaskStatus(id string, taskStatus status.ItemStatus) error {
	return s.baseService.UpdateItemStatus(id, taskStatus)
}

func (s *TaskService) UpdateTaskName(id, name string) error {
	return s.baseService.UpdateItemName(id, name)
}

func (s *TaskService) DeleteTask(id string) error {
	return s.baseService.DeleteItem(id)
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
