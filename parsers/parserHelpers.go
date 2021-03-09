package parsers

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getNumberFromLine(line string, pattern string) int {
	literal := getStringFromLine(line, pattern)
	if literal == nil {
		return 0
	}
	num, _ := strconv.Atoi(strings.TrimSpace(*literal))
	return num
}

func getStringFromLine(line string, pattern string) *string {
	r, _ := regexp.Compile(pattern)
	matches := r.FindStringSubmatch(line)
	if len(matches) < 2 {
		return nil
	}
	return &matches[1]
}

func getRtspDateParserFromLine(line string, pattern string) *time.Time {
	dateTimeLiteral := getStringFromLine(line, pattern)
	if dateTimeLiteral == nil {
		return nil
	}
	t, err := time.Parse(time.RFC1123, *dateTimeLiteral)
	if err != nil {
		return nil
	}
	return &t
}
