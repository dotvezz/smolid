package smolid

import (
	"encoding"
	"encoding/json"
	"fmt"
)

var _ json.Marshaler = ID{}
var _ json.Unmarshaler = (*ID)(nil)
var _ encoding.TextMarshaler = ID{}
var _ encoding.TextUnmarshaler = (*ID)(nil)

// MarshalJSON implements the json.Marshaler interface for ID.
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for ID.
func (id *ID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return fmt.Errorf("empty ID string")
	}

	newID, err := FromString(s)
	if err != nil {
		return err
	}
	*id = newID
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for ID.
func (id ID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for ID.
func (id *ID) UnmarshalText(text []byte) error {
	newID, err := FromString(string(text))
	if err != nil {
		return err
	}
	*id = newID
	return nil
}
