package parsers

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/ardroh/goRtspClient/sdp"
)

type SdpParser struct {
}

func (parser SdpParser) FromString(literal string) (*sdp.Sdp, error) {
	outputSdp := &sdp.Sdp{OriginalLiteral: literal}
	lines := strings.Split(literal, "\n")
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if strings.HasPrefix(line, "v=") {
			outputSdp.ProtocolVersion = getNumberFromLine(line, "v=(.*)")
		} else if strings.HasPrefix(line, "o=") {
			sessionOrigin := parseSessionOrigin(line)
			if sessionOrigin == nil {
				return nil, errors.New("Can't parse session origin!")
			}
			outputSdp.SessionOrigin = *sessionOrigin
		} else if strings.HasPrefix(line, "s=") {
			sessionName := getStringFromLine(line, "s=(.*)")
			if sessionName == nil {
				return nil, errors.New("Can't parse session name!")
			}
			outputSdp.SessionName = *sessionName
		} else if strings.HasPrefix(line, "i=") {
			outputSdp.SessionInfo = getStringFromLine(line, "i=(.*)")
		} else if strings.HasPrefix(line, "t=") {
			timing := parseTiming(line)
			if timing == nil {
				return nil, errors.New("Can't parse timings!")
			}
			outputSdp.Timings = append(outputSdp.Timings, *timing)
		} else if strings.HasPrefix(line, "a=") {
			attribute := parseAttribute(line)
			if attribute == nil {
				return nil, errors.New("Can't parse timings!")
			}
			outputSdp.SessionAttributes = append(outputSdp.SessionAttributes, *attribute)
		} else if strings.HasPrefix(line, "m=") {
			break
		}
	}
	var mediaTagsIdxPairs []struct {
		startIdx int
		endIdx   int
	}
	for idx, line := range lines {
		line := strings.TrimSpace(line)
		if !strings.HasPrefix(line, "m=") {
			continue
		}
		if len(mediaTagsIdxPairs) == 0 {
			mediaTagsIdxPairs = append(mediaTagsIdxPairs, struct {
				startIdx int
				endIdx   int
			}{startIdx: idx})
		} else {
			mediaTagsIdxPairs[len(mediaTagsIdxPairs)-1].endIdx = idx
			mediaTagsIdxPairs = append(mediaTagsIdxPairs, struct {
				startIdx int
				endIdx   int
			}{startIdx: idx})
		}
	}
	mediaTagsIdxPairs[len(mediaTagsIdxPairs)-1].endIdx = len(lines)
	for _, pair := range mediaTagsIdxPairs {
		description := parseMediaDescription(lines[pair.startIdx:pair.endIdx])
		if description == nil {
			return nil, errors.New("Can't parse media description!")
		}
		outputSdp.Medias = append(outputSdp.Medias, *description)
	}

	return outputSdp, nil
}

func parseSessionOrigin(line string) *sdp.SessionOrigin {
	r, _ := regexp.Compile("o=(.*?) (.*?) (.*?) (.*?) (.*?) (.*)")
	matches := r.FindStringSubmatch(line)
	if len(matches) < 7 {
		return nil
	}
	sessionVersion, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil
	}
	return &sdp.SessionOrigin{
		Username:       matches[1],
		SessionId:      matches[2],
		SessionVersion: sessionVersion,
		NetType:        matches[4],
		AddrType:       matches[5],
		UnicastAddress: matches[6],
	}
}

func parseTiming(line string) *sdp.Timing {
	r, _ := regexp.Compile("t=(.*?) (.*)")
	matches := r.FindStringSubmatch(line)
	if len(matches) < 3 {
		return nil
	}
	startTime, err := strconv.ParseInt(matches[1], 10, 64)
	endTime, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return nil
	}
	return &sdp.Timing{StartTime: startTime, EndTime: endTime}
}

func parseAttribute(line string) *sdp.Attribute {
	r, _ := regexp.Compile("a=(.*?):(.*)")
	matches := r.FindStringSubmatch(line)
	if len(matches) < 3 {
		r, _ = regexp.Compile("a=(.*)")
		matches = r.FindStringSubmatch(line)
		if len(matches) != 2 {
			return nil
		}
	}
	if len(matches) == 3 {
		return &sdp.Attribute{
			Key:   matches[1],
			Value: matches[2],
		}
	} else {
		return &sdp.Attribute{
			Key: matches[1],
		}
	}
}

func parseMediaDescription(lines []string) *sdp.MediaDescription {
	mediaDescr := &sdp.MediaDescription{}
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if strings.HasPrefix(line, "m=") {
			r, _ := regexp.Compile("m=(.*?) (.*?) (.*?) (.*)")
			matches := r.FindStringSubmatch(line)
			if len(matches) < 5 {
				return nil
			}
			mediaDescr.Media = matches[1]
			if strings.Contains(matches[2], "/") {
				splitPortInfo := strings.Split(matches[2], "/")
				if len(splitPortInfo) != 2 {
					return nil
				}
				initPort, err := strconv.Atoi(splitPortInfo[0])
				numOfPorts, err := strconv.Atoi(splitPortInfo[1])
				if err != nil {
					return nil
				}
				mediaDescr.Port.Port = initPort
				mediaDescr.Port.NumOfPorts = &numOfPorts
			} else {
				initPort, err := strconv.Atoi(matches[2])
				if err != nil {
					return nil
				}
				mediaDescr.Port.Port = initPort
			}
			mediaDescr.Proto = matches[3]
			fmt, err := strconv.Atoi(matches[4])
			if err != nil {
				return nil
			}
			mediaDescr.Fmt = fmt
		} else if strings.HasPrefix(line, "b=") {
			mediaDescr.Bandwidth = parseBandwidth(line)
		} else if strings.HasPrefix(line, "a=") {
			attr := parseAttribute(line)
			if attr == nil {
				return nil
			}
			mediaDescr.Attributes = append(mediaDescr.Attributes, *attr)
		}
	}
	return mediaDescr
}

func parseBandwidth(line string) *sdp.Bandwidth {
	r, _ := regexp.Compile("b=(.*?):(.*)")
	matches := r.FindStringSubmatch(line)
	if len(matches) < 3 {
		return nil
	}
	bandwidthValue, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return &sdp.Bandwidth{BwType: matches[1], Bandwidth: bandwidthValue}
}
