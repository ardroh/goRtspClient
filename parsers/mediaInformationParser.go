package parsers

import "github.com/ardroh/goRtspClient/media"

type MediaInformationParser struct {
}

func (parser MediaInformationParser) FromString(literal string) (*media.MediaInformation, error) {
	mediaInformation := &media.MediaInformation{OriginalLiteral: literal}
	return mediaInformation, nil
}
