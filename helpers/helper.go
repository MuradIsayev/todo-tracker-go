package helpers

import (
	"errors"
	"fmt"
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
