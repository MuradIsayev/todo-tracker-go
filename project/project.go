package project

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MuradIsayev/todo-tracker/constants"
	"github.com/MuradIsayev/todo-tracker/helpers"
	"github.com/olekukonko/tablewriter"
)

type ProjectStatus int

const (
	NOT_STARTED ProjectStatus = iota
	STARTED
	COMPLETED
)

type Project struct {
	Id                       int           `json:"id"`
	Name                     string        `json:"name"`
	Status                   ProjectStatus `json:"status"`
	CreatedAt                time.Time     `json:"createdAt"`
	TotalSpentTimeOfAllTasks int           `json:"totalSpentTime"`
}

type ProjectService struct {
	filePath string
	table    *tablewriter.Table
}

func (projectStatus ProjectStatus) String() string {
	switch projectStatus {
	case NOT_STARTED:
		return "NOT_STARTED"
	case STARTED:
		return "STARTED"
	case COMPLETED:
		return "COMPLETED"
	default:
		return "UNKNOWN"
	}
}

func NewProjectService(filePath string, table *tablewriter.Table) *ProjectService {
	table.SetHeader([]string{constants.COLUMN_ID, constants.COLUMN_NAME, constants.COLUMN_STATUS, constants.COLUMN_CREATE_DATE, constants.COLUMN_TOTAL_SPENT_TIME})

	return &ProjectService{
		filePath: filePath,
		table:    table,
	}
}

func (s *ProjectService) getNextID(projects []Project) int {
	if len(projects) == 0 {
		return 1
	}
	return projects[len(projects)-1].Id + 1
}

func (s *ProjectService) readProjectsFromFile() ([]Project, error) {
	fileContent, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Project{}, nil
		}
		return nil, fmt.Errorf("cannot read file: %v", err)
	}

	var projects []Project
	if len(fileContent) > 0 {
		if err := json.Unmarshal(fileContent, &projects); err != nil {
			return nil, fmt.Errorf("cannot convert JSON to struct: %v", err)
		}
	}
	return projects, nil
}

func (s *ProjectService) writeProjectsToFile(projects []Project) error {
	newProjects, err := json.Marshal(projects)
	if err != nil {
		return fmt.Errorf("Cannot convert Struct to JSON: %v", err)
	}

	if err := os.WriteFile(s.filePath, newProjects, 0644); err != nil {
		return fmt.Errorf("cannot write to file: %v", err)
	}

	return nil
}

func (s *ProjectService) UpdateProjectStatus(id string, projectStatus ProjectStatus) error {
	projectId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return err
	}

	projects, err := s.readProjectsFromFile()
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

	return s.writeProjectsToFile(projects)
}

func (s *ProjectService) UpdateProjectName(id, name string) error {
	projectId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return err
	}

	projects, err := s.readProjectsFromFile()
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

	return s.writeProjectsToFile(projects)
}

func (s *ProjectService) CreateProject(name string) error {
	projects, err := s.readProjectsFromFile()
	if err != nil {
		return err
	}

	project := Project{
		Id:        s.getNextID(projects),
		Name:      name,
		Status:    NOT_STARTED,
		CreatedAt: time.Now(),
	}

	projects = append(projects, project)

	return s.writeProjectsToFile(projects)
}

func (s *ProjectService) IsProjectExists(id string) bool {
	projectId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return false
	}

	projects, err := s.readProjectsFromFile()
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

	projects, err := s.readProjectsFromFile()
	if err != nil {
		return err
	}

	index, _, err := s.findProjectById(projects, projectId)
	if err != nil {
		return err
	}

	projects = append(projects[:index], projects[index+1:]...)

	return s.writeProjectsToFile(projects)
}

func (s *ProjectService) DeleteAllProjects() error {
	var projects []Project

	return s.writeProjectsToFile(projects)
}

func (s *ProjectService) UpdateTotalSpentTime(id int, spentTime int) error {
	projects, err := s.readProjectsFromFile()
	if err != nil {
		return err
	}

	index, project, err := s.findProjectById(projects, id)
	if err != nil {
		return err
	}

	project.TotalSpentTimeOfAllTasks += spentTime

	if project.Status != COMPLETED && project.TotalSpentTimeOfAllTasks > 0 && project.Status != STARTED {
		project.Status = STARTED
	}
	projects[index] = *project

	return s.writeProjectsToFile(projects)
}

func defineFooterText(nbOfLeftProjects, nbOfTotalProjects int) string {
	if nbOfLeftProjects == 0 && nbOfTotalProjects == 0 {
		return "No projects found"
	}

	return fmt.Sprintf("Left projects: %d", nbOfLeftProjects)
}

func (s *ProjectService) ListProjects(statusFilter ProjectStatus) error {
	projects, err := s.readProjectsFromFile()
	if err != nil {
		return err
	}

	var nbOfLeftprojects int

	for _, project := range projects {
		if statusFilter == -1 || project.Status == statusFilter {
			createdAt := project.CreatedAt.Format(constants.DATE_FORMAT)
			totalSpentTimeOfAllTasks := helpers.FormatSpendTime(project.TotalSpentTimeOfAllTasks)

			s.table.Append([]string{strconv.Itoa(project.Id), project.Name, project.Status.String(), createdAt, totalSpentTimeOfAllTasks})
			if project.Status == NOT_STARTED {
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
