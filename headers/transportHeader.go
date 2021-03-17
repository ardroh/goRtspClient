package headers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type RtspTransportType string

const (
	RtpAvp    RtspTransportType = "RTP/AVP"
	RtpAvpTcp                   = "RTP/AVP/TCP"
)

type RtspTransmissionType string

const (
	Unicast   RtspTransmissionType = "unicast"
	Multicast                      = "multicast"
)

type InterleavedPair struct {
	RangeMin int
	RangeMax int
}

type RtspTransportMode string

const (
	NoneTransportMode   RtspTransportMode = ""
	PlayTransportMode                     = "PLAY"
	RecordTransportMode                   = "RECORD"
)

type PortPool struct {
	RtpPort  int
	RtcpPort int
}

type TransportHeader struct {
	TransportType    RtspTransportType
	TransmissionType RtspTransmissionType
	InterleavedPair  *InterleavedPair
	Ssrc             string
	Mode             RtspTransportMode
	ClientPort       *PortPool
}

// Transport: RTP/AVP/TCP;unicast;interleaved=0-1;ssrc=5F53FB16;mode="PLAY"
func (info TransportHeader) FromString(line string) (TransportHeader, error) {
	params := strings.Split(line, ";")
	for _, param := range params {
		if strings.HasPrefix(param, "Transport: ") {
			transportTypeLiteral := strings.ReplaceAll(param, "Transport: ", "")
			info.TransportType = RtspTransportType(transportTypeLiteral)
		} else if param == "unicast" || param == "multicast" {
			info.TransmissionType = RtspTransmissionType(param)
		} else if strings.Contains(param, "=") {
			if !info.parseParam(param) {
				return info, errors.New("Failed to parse param!")
			}
		}
	}
	return info, nil
}

func (info *TransportHeader) parseParam(param string) bool {
	key, value := parseKeyValueFromParam(param)
	if key == nil || value == nil {
		return false
	}
	switch *key {
	case "ssrc":
		info.Ssrc = *value
		break
	case "interleaved":
		interleavedPair := parseInterleaved(*value)
		if interleavedPair == nil {
			return false
		}
		info.InterleavedPair = interleavedPair
		break
	case "mode":
		info.Mode = RtspTransportMode(strings.ReplaceAll(*value, "\"", ""))
		break
	}
	return true
}

func parseInterleaved(value string) *InterleavedPair {
	splitMinMax := strings.Split(value, "-")
	if len(splitMinMax) != 2 {
		return nil
	}
	rangeMin, err := strconv.Atoi(splitMinMax[0])
	rangeMax, err := strconv.Atoi(splitMinMax[1])
	interleavedPair := &InterleavedPair{
		RangeMin: rangeMin,
		RangeMax: rangeMax,
	}
	if err != nil {
		return nil
	}
	return interleavedPair
}

func parseKeyValueFromParam(line string) (*string, *string) {
	line = strings.TrimSpace(line)
	splitKeyValue := strings.Split(line, "=")
	if len(splitKeyValue) != 2 {
		return nil, nil
	}
	return &splitKeyValue[0], &splitKeyValue[1]
}

func (info TransportHeader) ToString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Transport: %s;%s", info.TransportType, info.TransmissionType))
	if info.TransportType == RtpAvp && info.ClientPort != nil {
		sb.WriteString(fmt.Sprintf(";client_port=%d-%d", info.ClientPort.RtpPort, info.ClientPort.RtcpPort))
	} else if info.TransportType == RtpAvpTcp && info.InterleavedPair != nil {
		sb.WriteString(fmt.Sprintf(";interleaved=%d-%d", info.InterleavedPair.RangeMin, info.InterleavedPair.RangeMax))
	}
	if len(info.Ssrc) > 0 {
		sb.WriteString(fmt.Sprintf(";ssrc=%s", info.Ssrc))
	}
	if info.Mode != NoneTransportMode {
		sb.WriteString(fmt.Sprintf(";mode=\"%s\"", info.Mode))
	}
	return sb.String()
}
