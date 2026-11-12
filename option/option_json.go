package option

import (
	"bytes"
	"encoding/json"
	"errors"
)

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsNone() {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

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
