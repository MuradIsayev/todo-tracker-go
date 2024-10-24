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

func NewTaskService(projectService *project.ProjectService, table *tablewriter.Table) *TaskService {
	table.SetHeader([]string{constants.COLUMN_ID, constants.COLUMN_NAME, constants.COLUMN_STATUS, constants.COLUMN_CREATE_DATE, constants.COLUMN_UPDATE_DATE, constants.COLUMN_TOTAL_SPENT_TIME})

	return &TaskService{
		table:          table,
		projectService: projectService,
	}
}

type TaskManager interface {
	DeleteAllTasks(projectId string, shouldAlterTasksCounter bool) error
	DeleteTasksByProjectId(projectId string) error
	UpdateTaskTimer(taskId int, newDuration int) error
}

func (s *TaskService) DeleteAllTasks(projectId string, shouldAlterTasksCounter bool) error {
	if err := helpers.RemoveContentsOfDirectory("output/tasks"); err != nil {
		return err
	}

	if shouldAlterTasksCounter {
		if err := s.projectService.UpdateTotalTasksOfProject(projectId, 0); err != nil {
			return err
		}
	}

	fmt.Println("Tasks deleted successfully")

	return nil
}

func (t *TaskService) DeleteTasksByProjectId(projectId string) error {
	filePathOfTask := fmt.Sprintf("output/tasks/%s_%s", projectId, constants.TASK_FILE_NAME)

	if err := helpers.RemoveFileByFilePath(filePathOfTask); err != nil {
		return err
	}

	fmt.Println("Tasks deleted successfully")

	return nil
}

func (t *TaskService) UpdateTaskTimer(taskId int, newDuration int) error {
	return t.baseService.UpdateTotalSpentTime(taskId, newDuration)
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

	// update project total spent time
	if err := s.projectService.UpdateProjectTimer(task.ProjectId, spentTime); err != nil {
		return err
	}

	// update task total spent time
	if err := s.baseService.UpdateTotalSpentTime(id, spentTime); err != nil {
		return err
	}

	return nil
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
	if err := s.baseService.UpdateItemStatus(id, taskStatus); err != nil {
		return err
	}

	fmt.Println("Task status updated successfully")

	return nil
}

func (s *TaskService) UpdateTaskName(id, name string) error {
	if err := s.baseService.UpdateItemName(id, name); err != nil {
		return err
	}

	fmt.Println("Task name updated successfully")

	return nil
}

func (s *TaskService) DeleteTask(id string, projectId string) error {
	if err := s.baseService.DeleteItemById(id); err != nil {
		return err
	}

	tasks := []Task{}
	err := s.baseService.ReadFromFile(&tasks)
	if err != nil {
		return err
	}

	if err := s.projectService.UpdateTotalTasksOfProject(projectId, len(tasks)); err != nil {
		return err
	}

	fmt.Println("Task deleted successfully")

	return nil
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

	if err := s.baseService.WriteToFile(tasks); err != nil {
		return err
	}

	if err := s.projectService.UpdateTotalTasksOfProject(projectID, len(tasks)); err != nil {
		return err
	}

	fmt.Println("Task created successfully")

	return nil
}
