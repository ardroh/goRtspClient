package goRtspClient

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
	rtpPort  int
	rtcpPort int
}

type InterleavedPair struct {
	rangeMin int
	rangeMax int
}

type RtspSetupCommand struct {
	address         string
	cseq            int
	transport       RtspTransportType
	transmission    RtspTransmissionType
	rtpPort         PortPool
	interleavedPair InterleavedPair
}

func (setup RtspSetupCommand) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("SETUP %s RTSP/1.0\n", setup.address))
	sb.WriteString(fmt.Sprintf("CSeq: %d\n", setup.cseq))
	sb.WriteString(fmt.Sprintf("Transport: %s;%s;", setup.transport, setup.transmission))
	if setup.transport == RtpAvp {
		sb.WriteString(fmt.Sprintf("client_port=%d-%d", setup.rtpPort.rtpPort, setup.rtpPort.rtcpPort))
	} else if setup.transport == RtpAvpTcp {
		sb.WriteString(fmt.Sprintf("interleaved=%d-%d", setup.interleavedPair.rangeMin, setup.interleavedPair.rangeMax))
	}
	sb.WriteString("\n\n")
	return sb.String()
}
