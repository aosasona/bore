package mimetype

import "errors"

type MimeType interface {
	mimeType() string
	String() string
}

func MimeTypeFromString(s string) (MimeType, error) {
	switch s {
	case "application/json":
		return MimeTypeApplicationJSON(s), nil
	case "text/plain":
		return MimeTypeTextPlain(s), nil
	default:
		return nil, errors.New("unknown mime type: " + s)
	}
}

type (
	MimeTypeTextPlain       string
	MimeTypeApplicationJSON string
)

// String implements MimeType.
func (m MimeTypeApplicationJSON) String() string {
	return string(m)
}

// mimeType implements MimeType.
func (m MimeTypeApplicationJSON) mimeType() string {
	return "application/json"
}

func (m MimeTypeTextPlain) mimeType() string {
	return "text/plain"
}

func (m MimeTypeTextPlain) String() string {
	return string(m)
}

var (
	_ MimeType = (*MimeTypeTextPlain)(nil)
	_ MimeType = (*MimeTypeApplicationJSON)(nil)
)
