package helpers

import (
	"errors"
	"fmt"
	"os"
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
