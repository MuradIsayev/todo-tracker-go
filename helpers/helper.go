package helpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func ValidateIdAndConvertToInt(id string) (int, error) {
	var numberRegex = regexp.MustCompile(`^[0-9]+$`)

	if !numberRegex.MatchString(id) {
		return 0, errors.New("ID must only contain digits")
	}

	return strconv.Atoi(id)
}

func RemoveFileByFilePath(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

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
