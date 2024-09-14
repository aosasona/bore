package handler

import (
	"encoding/base64"
	"fmt"
)

// DecodeToFormat decodes the content to the specified (and supported) format
func (h *Handler) Decode(content []byte, format Format) ([]byte, error) {
	switch format {
	case FormatBase64:
		destination := make([]byte, base64.StdEncoding.DecodedLen(len(content)))
		if _, err := base64.StdEncoding.Decode(destination, content); err != nil {
			return nil, err
		}
		return destination, nil
	case FormatPlainText:
		return content, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

func ValidateFormat(format Format) bool {
	return format == FormatBase64 || format == FormatPlainText
}
