//go:build go1.24

package option

import (
	"encoding/json"
	"testing"
)

func TestMarshalJSONOmitZero(t *testing.T) {
	type payload struct {
		Required Option[int] `json:"required"`
		Optional Option[int] `json:"optional,omitzero"`
	}

	tests := []struct {
		name string
		in   payload
		want string
	}{
		{
			name: "omits none",
			in:   payload{Required: None[int](), Optional: None[int]()},
			want: `{"required":null}`,
		},
		{
			name: "keeps some zero",
			in:   payload{Required: None[int](), Optional: Some(0)},
			want: `{"required":null,"optional":0}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.in)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("json.Marshal() = %s, want %s", got, tt.want)
			}
		})
	}
}
