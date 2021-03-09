package media

type MediaInformation struct {
	OriginalLiteral   string
	ProtocolVersion   int
	SessionOrigin     SessionOrigin
	SessionName       string
	SessionInfo       *string
	Timings           []Timing
	SessionAttributes []Attribute
	Medias            []MediaDescription
}

func (mediaInfo MediaInformation) HasSessionAttribute(key string) bool {
	for _, attr := range mediaInfo.SessionAttributes {
		if attr.Key == key {
			return true
		}
	}
	return false
}

func (mediaInfo MediaInformation) GetSessionAttribute(key string) string {
	for _, attr := range mediaInfo.SessionAttributes {
		if attr.Key == key {
			return attr.Value
		}
	}
	return ""
}
