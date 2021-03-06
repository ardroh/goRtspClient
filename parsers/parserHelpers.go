package parsers

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getNumberFromLine(pattern string, line string) int {
	literal := getStringFromLine(pattern, line)
	if literal == nil {
		return 0
	}
	num, _ := strconv.Atoi(strings.TrimSpace(*literal))
	return num
}

func getStringFromLine(pattern string, line string) *string {
	r, _ := regexp.Compile(pattern)
	matches := r.FindStringSubmatch(line)
	if len(matches) < 2 {
		return nil
	}
	return &matches[1]
}

func getRtspDateParserFromLine(pattern string, line string) *time.Time {
	dateTimeLiteral := getStringFromLine(pattern, line)
	if dateTimeLiteral == nil {
		return nil
	}
	t, err := time.Parse(time.RFC1123, *dateTimeLiteral)
	if err != nil {
		return nil
	}
	return &t
}
