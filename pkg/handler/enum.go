package handler

// Format
type Format string

const (
	FormatBase64    Format = "base64"
	FormatPlainText Format = "plain"
)

func (f Format) String() string {
	return string(f)
}

// Source
type Source string

const (
	SourceBore   Source = "bore"
	SourceSystem Source = "system"
)

func (s Source) String() string {
	return string(s)
}
