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
	Id             int               `json:"id"`
	Name           string            `json:"name"`
	Status         status.ItemStatus `json:"status"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
	TotalSpentTime int               `json:"totalSpentTime"`
	NbOfTotalTasks int               `json:"nbOfTotalTasks"`
}

type ProjectService struct {
	baseService *base.BaseService[Project]
	table       *tablewriter.Table
}

func NewProjectService(filePath string, table *tablewriter.Table) *ProjectService {
	table.SetHeader([]string{constants.COLUMN_ID, constants.COLUMN_NAME, constants.COLUMN_STATUS, constants.COLUMN_CREATE_DATE, constants.COLUMN_UPDATE_DATE, constants.COLUMN_TOTAL_SPENT_TIME, constants.COLUMN_TOTAL_TASKS})

	return &ProjectService{
		table: table,
		baseService: &base.BaseService[Project]{
			FilePath: filePath,
		},
	}
}

type ProjectManager interface {
	DeleteAllProjects() error
	DeleteProjectById(projectId string) error
	UpdateProjectTimer(projectId int, newDuration int) error
}

func (p *ProjectService) DeleteAllProjects() error {
	if err := p.baseService.DeleteAllItems(); err != nil {
		return err
	}

	fmt.Println("Projects deleted successfully")

	return nil
}

func (p *ProjectService) DeleteProjectById(projectId string) error {
	if err := p.baseService.DeleteItemById(projectId); err != nil {
		return err
	}

	fmt.Println("Project deleted successfully")

	return nil
}

func (p *ProjectService) UpdateProjectTimer(projectId int, newDuration int) error {
	return p.baseService.UpdateTotalSpentTime(projectId, newDuration)
}

func (s *ProjectService) UpdateProjectStatus(id string, projectStatus status.ItemStatus) error {
	if err := s.baseService.UpdateItemStatus(id, projectStatus); err != nil {
		return err
	}

	fmt.Println("Project status updated successfully")

	return nil
}

func (s *ProjectService) UpdateProjectName(id, name string) error {
	if err := s.baseService.UpdateItemName(id, name); err != nil {
		return err
	}

	fmt.Println("Project name updated successfully")

	return nil
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

	if err := s.baseService.WriteToFile(projects); err != nil {
		return err
	}

	fmt.Println("Project created successfully")

	return nil
}

func (s *ProjectService) UpdateTotalTasksOfProject(id string, nbOfTotalTasks int) error {
	projects := []Project{}
	err := s.baseService.ReadFromFile(&projects)
	if err != nil {
		return err
	}

	projectId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return err
	}

	index, project, err := s.baseService.FindItemById(projects, projectId)
	if err != nil {
		return err
	}

	project.NbOfTotalTasks = nbOfTotalTasks
	projects[index] = *project

	if err := s.baseService.WriteToFile(projects); err != nil {
		return err
	}

	return nil
}

func (s *ProjectService) FindProjectNameById(id string) string {
	projectId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return ""
	}

	projects := []Project{}
	err = s.baseService.ReadFromFile(&projects)
	if err != nil {
		return ""
	}

	_, project, err := s.baseService.FindItemById(projects, projectId)
	if err != nil {
		return ""
	}

	return project.Name
}

func (s *ProjectService) UpdateProjectSpentTime(id int, spentTime int) error {
	if err := s.baseService.UpdateTotalSpentTime(id, spentTime); err != nil {
		return err
	}

	return nil
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
			totalSpentTime := helpers.FormatSpendTime(project.TotalSpentTime)

			s.table.Append([]string{strconv.Itoa(project.Id), project.Name, project.Status.String(), createdAt, updatedAt, totalSpentTime, strconv.Itoa(project.NbOfTotalTasks)})
			if project.Status == status.TODO {
				nbOfLeftprojects++
			}
		}
	}

	s.table.SetRowLine(true)
	s.table.SetFooter([]string{"", "", "", "", "", " ", defineFooterText(nbOfLeftprojects, len(projects))})
	s.table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)
	s.table.SetFooterColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})

	s.table.Render()

	return nil
}
