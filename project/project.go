package project

import (
	"fmt"
	"strconv"
	"time"

	"github.com/MuradIsayev/todo-tracker/base"
	"github.com/MuradIsayev/todo-tracker/constants"
	"github.com/MuradIsayev/todo-tracker/helpers"
	"github.com/MuradIsayev/todo-tracker/status"
	"github.com/olekukonko/tablewriter"
)

type Project struct {
	Id                       int               `json:"id"`
	Name                     string            `json:"name"`
	Status                   status.ItemStatus `json:"status"`
	CreatedAt                time.Time         `json:"createdAt"`
	UpdatedAt                time.Time         `json:"updatedAt"`
	TotalSpentTimeOfAllTasks int               `json:"totalSpentTime"`
}

type ProjectService struct {
	baseService *base.BaseService[Project]
	table       *tablewriter.Table
}

func NewProjectService(filePath string, table *tablewriter.Table) *ProjectService {
	table.SetHeader([]string{constants.COLUMN_ID, constants.COLUMN_NAME, constants.COLUMN_STATUS, constants.COLUMN_CREATE_DATE, constants.COLUMN_UPDATE_DATE, constants.COLUMN_TOTAL_SPENT_TIME})

	return &ProjectService{
		table: table,
		baseService: &base.BaseService[Project]{
			FilePath: filePath,
		},
	}
}

func (s *ProjectService) UpdateProjectStatus(id string, projectStatus status.ItemStatus) error {
	return s.baseService.UpdateItemStatus(id, projectStatus)
}

func (s *ProjectService) UpdateProjectName(id, name string) error {
	return s.baseService.UpdateItemName(id, name)
}

func (s *ProjectService) DeleteProject(id string) error {
	return s.baseService.DeleteItem(id)
}

func (s *ProjectService) DeleteAllProjects() error {
	return s.baseService.DeleteAllItems()
}

func (s *ProjectService) CreateProject(name string) error {
	projects := []Project{}
	err := s.baseService.ReadFromFile(&projects)
	if err != nil {
		return err
	}

	project := Project{
		Id:        s.baseService.GetNextID(projects),
		Name:      name,
		Status:    status.TODO,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	projects = append(projects, project)

	return s.baseService.WriteToFile(projects)
}

func (s *ProjectService) IsProjectExists(id string) bool {
	projectId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return false
	}

	projects := []Project{}
	err = s.baseService.ReadFromFile(&projects)
	if err != nil {
		return false
	}

	for _, project := range projects {
		if project.Id == projectId {
			return true
		}
	}

	return false
}

func (s *ProjectService) UpdateTotalSpentTime(id int, spentTime int) error {
	projects := []Project{}
	err := s.baseService.ReadFromFile(&projects)
	if err != nil {
		return err
	}

	index, project, err := s.baseService.FindItemById(projects, id)
	if err != nil {
		return err
	}

	project.TotalSpentTimeOfAllTasks += spentTime

	if project.Status != status.DONE && project.TotalSpentTimeOfAllTasks > 0 && project.Status != status.IN_PROGRESS {
		project.Status = status.IN_PROGRESS
	}
	projects[index] = *project

	return s.baseService.WriteToFile(projects)
}

func defineFooterText(nbOfLeftProjects, nbOfTotalProjects int) string {
	if nbOfLeftProjects == 0 && nbOfTotalProjects == 0 {
		return "No projects found"
	}

	return fmt.Sprintf("Left projects: %d", nbOfLeftProjects)
}

func (s *ProjectService) ListProjects(statusFilter status.ItemStatus) error {
	projects := []Project{}
	err := s.baseService.ReadFromFile(&projects)
	if err != nil {
		return err
	}

	var nbOfLeftprojects int

	for _, project := range projects {
		if statusFilter == -1 || project.Status == statusFilter {
			createdAt := project.CreatedAt.Format(constants.DATE_FORMAT)
			updatedAt := project.UpdatedAt.Format(constants.DATE_FORMAT)
			totalSpentTimeOfAllTasks := helpers.FormatSpendTime(project.TotalSpentTimeOfAllTasks)

			s.table.Append([]string{strconv.Itoa(project.Id), project.Name, project.Status.String(), createdAt, updatedAt, totalSpentTimeOfAllTasks})
			if project.Status == status.TODO {
				nbOfLeftprojects++
			}
		}
	}

	s.table.SetRowLine(true)
	s.table.SetFooter([]string{"", "", "", "", " ", defineFooterText(nbOfLeftprojects, len(projects))})
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
