package countdown

import (
	"fmt"
	"time"

	"github.com/MuradIsayev/todo-tracker/service"
	"github.com/MuradIsayev/todo-tracker/task"
)

type CountdownService struct {
	PauseChan                 chan bool
	ResumeChan                chan bool
	StopChan                  chan bool
	DoneChan                  chan bool
	DisplayChan               chan string
	ExitChan                  chan bool
	taskService               *task.TaskService
	circularDependencyManager *service.Manager
}

func NewCountdownService(taskService *task.TaskService, circularDependencyManager *service.Manager) *CountdownService {
	return &CountdownService{
		PauseChan:                 make(chan bool),
		ResumeChan:                make(chan bool),
		StopChan:                  make(chan bool),
		DoneChan:                  make(chan bool),
		DisplayChan:               make(chan string),
		ExitChan:                  make(chan bool),
		taskService:               taskService,
		circularDependencyManager: circularDependencyManager,
	}
}

func (cs *CountdownService) StartCountdown(task *task.Task, countdownMinutes int) {
	remainingSeconds := countdownMinutes * 60

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	paused := false

	for remainingSeconds > 0 {
		select {
		case <-cs.StopChan:
			cs.DisplayChan <- fmt.Sprintf("Countdown stopped early for the task --> \"%s\".", task.Name)
			cs.circularDependencyManager.UpdateTaskAndProjectTimers(task.Id, task.ProjectId, countdownMinutes*60-remainingSeconds) // Save the elapsed time
			// cs.taskService.UpdateTaskSpentTime(task.Id, countdownMinutes*60-remainingSeconds)                                      // Save the elapsed time
			close(cs.DoneChan) // Signal that the countdown has ended
			return
		case <-cs.ExitChan:
			cs.DisplayChan <- fmt.Sprintf("Countdown session for task --> \"%s\" ignored.", task.Name)
			close(cs.DoneChan) // Exit without saving any time
			return
		case <-cs.PauseChan:
			paused = true
			cs.DisplayChan <- "Countdown paused. Type (r)esume to continue."
		case <-cs.ResumeChan:
			if paused {
				paused = false
				cs.DisplayChan <- "Countdown resumed."
			}
		case <-ticker.C:
			if !paused {
				remainingSeconds--
				cs.DisplayChan <- fmt.Sprintf("Task --> \"%s\": %d:%02d", task.Name, remainingSeconds/60, remainingSeconds%60)
			}
		}
	}

	cs.DisplayChan <- fmt.Sprintf("Countdown complete for the task --> \"%s\". Now press (e) to exit", task.Name)
	cs.circularDependencyManager.UpdateTaskAndProjectTimers(task.Id, task.ProjectId, countdownMinutes*60-remainingSeconds) // Save the elapsed time
	close(cs.DoneChan)                                                                                                     // Signal that the countdown has ended
}
