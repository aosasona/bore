package payload

import "encoding/json"

type Payload interface {
	json.Marshaler
	Type() string
}
