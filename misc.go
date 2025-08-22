package economic

import (
	"strconv"
	"strings"
)

func ValidateDate(date string) bool { // YYYY-MM-DD
	ps := strings.Split(date, "-")
	if len(ps) != 3 {
		return false
	}

	_, err := strconv.Atoi(ps[0])
	if err != nil || len(ps[0]) != 4 {
		return false
	}

	_, err = strconv.Atoi(ps[1])
	if err != nil || len(ps[1]) < 1 || len(ps[1]) > 2 {
		return false
	}

	_, err = strconv.Atoi(ps[2])
	if err != nil || len(ps[2]) < 1 || len(ps[2]) > 2 {
		return false
	}

	return true
}
