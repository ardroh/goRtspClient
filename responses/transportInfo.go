package responses

import "github.com/ardroh/goRtspClient/commands"

// Transport: RTP/AVP/TCP;unicast;interleaved=0-1;ssrc=5F53FB16;mode="PLAY"

type TransportInfo struct {
	TransportType    commands.RtspTransportType
	TransmissionType commands.RtspTransmissionType
	InterleavedPair  commands.InterleavedPair
	Ssrc             string
	Mode             string
}

func (info TransportInfo) FromString(line string) {

}
