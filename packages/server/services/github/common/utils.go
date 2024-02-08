package common

import (
	"regexp"
	"strconv"
	"strings"
)

// Find the last page with the Link header
// It splits the link on the comma, then uses a regular expression to find the last page number
func GetLastPage(Link string) (int, error) {
	parts := strings.Split(Link, ",")
	if len(parts) > 1 {
		re, err := regexp.Compile(`&page=(\d+)`)
		if err != nil {
			return 1, err
		}
		match := re.FindStringSubmatch(parts[1])
		if match == nil {
			return 1, nil
		}
		if len(match) > 1 {
			lastPageNumber, errr := strconv.Atoi(match[1])
			return lastPageNumber, errr
		}
	}
	return 1, nil
}
