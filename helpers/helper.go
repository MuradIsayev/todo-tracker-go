package helpers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
)

// Validates the ID and converts it to an integer
func ValidateIdAndConvertToInt(id string) (int, error) {
	var numberRegex = regexp.MustCompile(`^[0-9]+$`)

	if !numberRegex.MatchString(id) {
		return 0, errors.New("ID must only contain digits")
	}

	return strconv.Atoi(id)
}

// Removes a file by its file path
func RemoveFileByFilePath(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

// Checks if a file exists
func DoesFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

// Removes all files in a directory
func RemoveContentsOfDirectory(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// Formats the total spent time in hours, minutes, and seconds
func FormatSpendTime(totalSpentTime int) string {
	formattedSpendTime := ""

	if totalSpentTime > 0 {

		hours := totalSpentTime / 3600
		minutes := (totalSpentTime % 3600) / 60
		seconds := totalSpentTime % 60

		if hours > 0 {
			// find if it is plural or singular
			if hours > 1 {
				formattedSpendTime += fmt.Sprintf("%d hours ", hours)
			} else {
				formattedSpendTime += fmt.Sprintf("%d hour ", hours)
			}
		}

		if minutes > 0 {
			// find if it is plural or singular
			if minutes > 1 {
				formattedSpendTime += fmt.Sprintf("%d minutes ", minutes)
			} else {
				formattedSpendTime += fmt.Sprintf("%d minute ", minutes)
			}
		}

		if seconds > 0 {
			// find if it is plural or singular
			if seconds > 1 {
				formattedSpendTime += fmt.Sprintf("%d seconds ", seconds)
			} else {
				formattedSpendTime += fmt.Sprintf("%d second ", seconds)
			}
		}

	} else {
		formattedSpendTime = "0"
	}

	return formattedSpendTime
}

func BeepBeep() {
	cmd := exec.Command("afplay", "/System/Library/Sounds/Glass.aiff") // Change to a different sound if needed
	cmd.Run()
}

func DisplayHelp() {
	fmt.Println("\nWelcome to Todo Tracker CLI in Go!")
	fmt.Println("This tool helps you manage projects and tasks, set timers, and keep track of your progress.")

	fmt.Println("\n**Available Commands**")
	fmt.Println("-----------------------")
	fmt.Println("1. **Project Management (Normal Mode)**")
	fmt.Println("   - `add <project name>`  : Creates a new project with the given name.")
	fmt.Println("   - `list`                   : Lists all current projects.")
	fmt.Println("   - `delete <project ID> | --all`  : Deletes the specified project(s) and all its associated tasks.")
	fmt.Println("   - `update <project ID> <new project name>` : Renames the specified project.")
	fmt.Println("   - `mark <project ID> --done | --in-progress | --todo` : Marks the project status as done, in-progress, or to-do.")
	fmt.Println("   - `repl <project ID>`                   : Enters the REPL mode for the specified project to manage tasks.")

	fmt.Println("\n2. **Task Management (REPL Mode)**")
	fmt.Println("Contains the same commands as project management, except for the following command")
	fmt.Println("   - `t <task ID> <minutes>` : Starts a countdown timer for a specific task, and enters the Timer mode")
	fmt.Println("   (Timer will countdown from the specified minutes)")

	fmt.Println("\n3. **Timer Commands (Timer Mode)**")
	fmt.Println("   - `s`                                     : Stops the active timer and saves the focused time.")
	fmt.Println("   - `p`                                     : Pauses the active timer.")
	fmt.Println("   - `r`                                     : Resumes the paused timer.")
	fmt.Println("   - `e`                                     : Exits the timer mode and ignores the focused time.")

	fmt.Println("\n4. **General Commands**")
	fmt.Println("   - `help`                 : Shows this help message with command descriptions (Normal mode command).")
	fmt.Println("   - `exit`       : Exits the REPL mode. (REPL mode command).")

	fmt.Println("\n**Note**: For a full guide, see the README file or visit the project repository on GitHub.")
	fmt.Println("Happy tracking!")
	fmt.Println("-----------------------")
}
