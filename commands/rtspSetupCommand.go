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
	Address         string
	Cseq            int
	Transport       RtspTransportType
	Transmission    RtspTransmissionType
	RtpPort         PortPool
	InterleavedPair InterleavedPair
}

func (setup RtspSetupCommand) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("SETUP %s RTSP/1.0\n", setup.Address))
	sb.WriteString(fmt.Sprintf("CSeq: %d\n", setup.Cseq))
	sb.WriteString(fmt.Sprintf("Transport: %s;%s;", setup.Transport, setup.Transmission))
	if setup.Transport == RtpAvp {
		sb.WriteString(fmt.Sprintf("client_port=%d-%d", setup.RtpPort.RtpPort, setup.RtpPort.RtcpPort))
	} else if setup.Transport == RtpAvpTcp {
		sb.WriteString(fmt.Sprintf("interleaved=%d-%d", setup.InterleavedPair.RangeMin, setup.InterleavedPair.RangeMax))
	}
	sb.WriteString("\n\n")
	return sb.String()
}
