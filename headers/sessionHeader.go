package headers

import (
	"errors"
	"strconv"
	"strings"
)

type SessionHeader struct {
	Id      string
	Timeout int
}

func (header SessionHeader) FromString(line string) (SessionHeader, error) {
	params := strings.Split(line, ";")
	for _, param := range params {
		if strings.HasPrefix(param, "Session:") {
			header.Id = strings.ReplaceAll(param, "Session: ", "")
		} else if strings.Contains(line, "=") {
			key, value := parseKeyValueFromParam(param)
			if key == nil || value == nil {
				return header, errors.New("Failed to parse param!")
			}
			var err error
			switch *key {
			case "timeout":
				header.Timeout, err = strconv.Atoi(*value)
				break
			}
			if err != nil {
				return header, errors.New("Failed to parse param!")
			}
		}
	}
	return header, nil
}
