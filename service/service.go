package service

import (
	"github.com/MuradIsayev/todo-tracker/project"
	"github.com/MuradIsayev/todo-tracker/task"
)

type Manager struct {
	TaskService    task.TaskManager
	ProjectService project.ProjectManager
}

func NewManager(taskService task.TaskManager, projectService project.ProjectManager) *Manager {
	return &Manager{
		TaskService:    taskService,
		ProjectService: projectService,
	}
}

func (m *Manager) DeleteProjectAndCorrespondingTasks(projectId string) error {
	if err := m.TaskService.DeleteTasksByProjectId(projectId); err != nil {
		return err
	}

	return m.ProjectService.DeleteProjectById(projectId)
}

func (m *Manager) DeleteAllProjectWithAllTasks() error {
	if err := m.TaskService.DeleteAllTasks(); err != nil {
		return err
	}

	return m.ProjectService.DeleteAllProjects()
}

func (m *Manager) UpdateTaskAndProjectTimers(taskId, projectId int, newDuration int) error {
	if err := m.TaskService.UpdateTaskTimer(taskId, newDuration); err != nil {
		return err
	}

	return m.ProjectService.UpdateProjectTimer(projectId, newDuration)
}
