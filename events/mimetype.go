package events

type MimeType interface {
	mimeType() string
	String() string
}

type MimeTypeTextPlain string

func (m MimeTypeTextPlain) mimeType() string {
	return "text/plain"
}

func (m MimeTypeTextPlain) String() string {
	return string(m)
}

var _ MimeType = (*MimeTypeTextPlain)(nil)
