package option

import (
	"bytes"
	"encoding/json"
	"errors"
)

// MarshalJSON implements json.Marshaler. None is encoded as null; Some is
// encoded as the JSON representation of its value.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

// UnmarshalJSON implements json.Unmarshaler. JSON null decodes to None; any
// other valid JSON value decodes to Some.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if o == nil {
		return errors.New("option: UnmarshalJSON on nil receiver")
	}

	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*o = None[T]()
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*o = Some(value)
	return nil
}
