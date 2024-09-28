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
	TotalSpentTimeOfAllTasks int               `json:"totalSpentTime"`
}

type ProjectService struct {
	baseService *base.BaseService
	table       *tablewriter.Table
}

func NewProjectService(filePath string, table *tablewriter.Table) *ProjectService {
	table.SetHeader([]string{constants.COLUMN_ID, constants.COLUMN_NAME, constants.COLUMN_STATUS, constants.COLUMN_CREATE_DATE, constants.COLUMN_TOTAL_SPENT_TIME})

	return &ProjectService{
		table: table,
		baseService: &base.BaseService{
			FilePath: filePath,
		},
	}
}

func (s *ProjectService) UpdateProjectStatus(id string, projectStatus status.ItemStatus) error {
	projectId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return err
	}

	projects := []Project{}
	err = s.baseService.ReadFromFile(&projects)
	if err != nil {
		return err
	}

	index, project, err := s.findProjectById(projects, projectId)
	if err != nil {
		return err
	}

	project.Status = projectStatus
	// project.UpdatedAt = time.Now()
	projects[index] = *project

	return s.baseService.WriteToFile(projects)
}

func (s *ProjectService) UpdateProjectName(id, name string) error {
	projectId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return err
	}

	projects := []Project{}
	err = s.baseService.ReadFromFile(&projects)
	if err != nil {
		return err
	}

	index, project, err := s.findProjectById(projects, projectId)
	if err != nil {
		return err
	}

	if name != "" {
		project.Name = name
		// project.UpdatedAt = time.Now()
		projects[index] = *project
	}

	return s.baseService.WriteToFile(projects)
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

func (s *ProjectService) findProjectById(projects []Project, id int) (int, *Project, error) {
	for i, project := range projects {
		if project.Id == id {
			return i, &project, nil
		}
	}
	return -1, nil, fmt.Errorf("project with ID=%d not found", id)
}

func (s *ProjectService) DeleteProject(id string) error {
	projectId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return err
	}

	projects := []Project{}
	err = s.baseService.ReadFromFile(&projects)
	if err != nil {
		return err
	}

	index, _, err := s.findProjectById(projects, projectId)
	if err != nil {
		return err
	}

	projects = append(projects[:index], projects[index+1:]...)

	return s.baseService.WriteToFile(projects)
}

func (s *ProjectService) DeleteAllProjects() error {
	var projects []Project

	return s.baseService.WriteToFile(projects)
}

func (s *ProjectService) UpdateTotalSpentTime(id int, spentTime int) error {
	projects := []Project{}
	err := s.baseService.ReadFromFile(&projects)
	if err != nil {
		return err
	}

	index, project, err := s.findProjectById(projects, id)
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
			totalSpentTimeOfAllTasks := helpers.FormatSpendTime(project.TotalSpentTimeOfAllTasks)

			s.table.Append([]string{strconv.Itoa(project.Id), project.Name, project.Status.String(), createdAt, totalSpentTimeOfAllTasks})
			if project.Status == status.TODO {
				nbOfLeftprojects++
			}
		}
	}

	s.table.SetRowLine(true)
	s.table.SetFooter([]string{"", "", "", " ", defineFooterText(nbOfLeftprojects, len(projects))})
	s.table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)
	s.table.SetFooterColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})

	s.table.Render()

	return nil
}
