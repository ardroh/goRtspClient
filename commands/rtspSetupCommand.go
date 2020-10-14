package commands

import (
	"fmt"
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

type PortPool struct {
	RtpPort  int
	RtcpPort int
}

type InterleavedPair struct {
	RangeMin int
	RangeMax int
}

type RtspSetupCommand struct {
	Transport       RtspTransportType
	Transmission    RtspTransmissionType
	RtpPort         PortPool
	InterleavedPair InterleavedPair
}

func (cmd RtspSetupCommand) GetCommandType() RtspCommandType {
	return RtspSetup
}

func (cmd RtspSetupCommand) GetParamsString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Transport: %s;%s;", cmd.Transport, cmd.Transmission))
	if cmd.Transport == RtpAvp {
		sb.WriteString(fmt.Sprintf("client_port=%d-%d", cmd.RtpPort.RtpPort, cmd.RtpPort.RtcpPort))
	} else if cmd.Transport == RtpAvpTcp {
		sb.WriteString(fmt.Sprintf("interleaved=%d-%d", cmd.InterleavedPair.RangeMin, cmd.InterleavedPair.RangeMax))
	}
	sb.WriteString("\n")
	return sb.String()
}
