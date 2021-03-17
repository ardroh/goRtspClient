package sdp

type Sdp struct {
	OriginalLiteral   string
	ProtocolVersion   int
	SessionOrigin     SessionOrigin
	SessionName       string
	SessionInfo       *string
	Timings           []Timing
	SessionAttributes []Attribute
	Medias            []MediaDescription
}

func (sdp Sdp) HasSessionAttribute(key string) bool {
	for _, attr := range sdp.SessionAttributes {
		if attr.Key == key {
			return true
		}
	}
	return false
}

func (sdp Sdp) GetSessionAttribute(key string) string {
	for _, attr := range sdp.SessionAttributes {
		if attr.Key == key {
			return attr.Value
		}
	}
	return ""
}
