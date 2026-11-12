package option

import (
	"encoding/json"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		o    Option[int]
		want string
	}{
		{name: "zero value", o: Option[int]{}, want: "null"},
		{name: "none", o: None[int](), want: "null"},
		{name: "some zero", o: Some(0), want: "0"},
		{name: "some value", o: Some(42), want: "42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.o)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("json.Marshal() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMarshalJSONSomeNil(t *testing.T) {
	var value *int
	got, err := json.Marshal(Some(value))
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if string(got) != "null" {
		t.Fatalf("json.Marshal(Some(nil)) = %s, want null", got)
	}

	var decoded Option[*int]
	if err := json.Unmarshal(got, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if decoded.IsSome() {
		t.Error("JSON null should decode to None, including when marshaled from Some(nil)")
	}
}

func TestMarshalJSONPropagatesValueError(t *testing.T) {
	_, err := json.Marshal(Some(make(chan int)))
	if err == nil {
		t.Fatal("json.Marshal() error = nil, want unsupported type error")
	}
}

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

func TestUnmarshalJSON(t *testing.T) {
	var got Option[int]
	if err := json.Unmarshal([]byte("42"), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if value, ok := got.Value(); !ok || value != 42 {
		t.Errorf("json.Unmarshal() = %v, %v; want 42, true", value, ok)
	}

	if err := json.Unmarshal([]byte("null"), &got); err != nil {
		t.Fatalf("json.Unmarshal(null) error = %v", err)
	}
	if got.IsSome() {
		t.Error("json.Unmarshal(null) should produce None")
	}
}

func TestUnmarshalJSONErrorPreservesValue(t *testing.T) {
	got := Some(42)
	if err := json.Unmarshal([]byte(`"not an integer"`), &got); err == nil {
		t.Fatal("json.Unmarshal() error = nil, want type error")
	}

	if value, ok := got.Value(); !ok || value != 42 {
		t.Errorf("Option after failed unmarshal = %v, %v; want 42, true", value, ok)
	}
}

func TestUnmarshalJSONNilReceiver(t *testing.T) {
	var got *Option[int]
	if err := got.UnmarshalJSON([]byte("42")); err == nil {
		t.Fatal("UnmarshalJSON() error = nil, want nil receiver error")
	}
}
