package constants

const (
	ADD    string = "add"
	UPDATE string = "update"
	DELETE string = "delete"
	LIST   string = "list"
	MARK   string = "mark"
	REPL   string = "repl"
	TIMER  string = "t"
)

const (
	COLUMN_ID                    = "ID"
	COLUMN_NAME                  = "Name"
	COLUMN_STATUS                = "Status"
	COLUMN_CREATE_DATE           = "Create Date"
	COLUMN_UPDATE_DATE           = "Update Date"
	COLUMN_TOTAL_SPENT_TIME      = "Total Spent Time"
	COLUMN_TOTAL_SPENT_TIME_TASK = "Total Spent Time"
)

const DATE_FORMAT = "2006-01-02 15:04:05"

const TASK_FILE_NAME = "task.json"
const PROJECT_FILE_NAME = "output/projects.json"

// TIMER COMMANDS:
const (
	TIMER_PAUSE  string = "p"
	TIMER_RESUME string = "r"
	TIMER_STOP   string = "s"
	TIMER_EXIT   string = "e"
)
