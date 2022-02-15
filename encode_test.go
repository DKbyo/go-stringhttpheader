package stringhttpheader

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestHeader_types(t *testing.T) {
	str := "string"
	strPtr := &str
	timeVal := time.Date(2000, 1, 1, 12, 34, 56, 0, time.UTC)

	tests := []struct {
		in   interface{}
		want []string
	}{
		{
			// basic primitives
			struct {
				A string `header:"Z"`
				B int
				C uint
				D float32
				E bool
			}{},
			[]string{
				"Z: ",
				"B: 0",
				"C: 0",
				"D: 0",
				"E: false",
			},
		},
		{
			// pointers
			struct {
				A *string
				B *int
				C **string
				D *time.Time
			}{
				A: strPtr,
				C: &strPtr,
				D: &timeVal,
			},
			[]string{
				fmt.Sprintf("A: %s", str),
				fmt.Sprintf("B: %s", ""),
				fmt.Sprintf("C: %s", str),
				fmt.Sprintf("D: %s", "Sat, 01 Jan 2000 12:34:56 GMT"),
			},
		},
		{
			// slices and arrays
			struct {
				A []string
				B []*string
				C [2]string
				D []bool `header:",int"`
			}{
				A: []string{"a", "b"},
				B: []*string{&str, &str},
				C: [2]string{"a", "b"},
				D: []bool{true, false},
			},
			[]string{
				fmt.Sprintf("A: %s", "a"),
				fmt.Sprintf("A: %s", "b"),
				fmt.Sprintf("B: %s", str),
				fmt.Sprintf("B: %s", str),
				fmt.Sprintf("C: %s", "a"),
				fmt.Sprintf("C: %s", "b"),
				fmt.Sprintf("D: %s", "0"),
				fmt.Sprintf("D: %s", "1"),
			},
		},
		{
			// other types
			struct {
				A time.Time
				B time.Time `header:",unix"`
				C bool      `header:",int"`
				D bool      `header:",int"`
				E http.Header
			}{
				A: time.Date(2000, 1, 1, 12, 34, 56, 0, time.UTC),
				B: time.Date(2000, 1, 1, 12, 34, 56, 0, time.UTC),
				C: true,
				D: false,
				E: http.Header{
					"F": []string{"f1"},
					"G": []string{"gg"},
				},
			},
			[]string{
				fmt.Sprintf("A: %s", "Sat, 01 Jan 2000 12:34:56 GMT"),
				fmt.Sprintf("B: %s", "946730096"),
				fmt.Sprintf("C: %s", "1"),
				fmt.Sprintf("D: %s", "0"),
				fmt.Sprintf("F: %s", "f1"),
				fmt.Sprintf("G: %s", "gg"),
			},
		},
		{
			nil,
			[]string{},
		},
		{
			&struct {
				A string
			}{"test"},
			[]string{
				fmt.Sprintf("A: %s", "test"),
			},
		},
	}

	for i, tt := range tests {
		v, err := Header(tt.in)
		if err != nil {
			t.Errorf("%d. Header(%+v) returned error: %v", i, tt.in, err)
		}

		if !reflect.DeepEqual(tt.want, v) {
			t.Errorf("%d. Header(%+v) returned %#v, want %#v", i, tt.in, v, tt.want)
		}
	}
}

func TestHeader_omitEmpty(t *testing.T) {
	str := ""
	s := struct {
		a string
		A string
		B string    `header:",omitempty"`
		C string    `header:"-"`
		D string    `header:"omitempty"` // actually named omitempty, not an option
		E *string   `header:",omitempty"`
		F bool      `header:",omitempty"`
		G int       `header:",omitempty"`
		H uint      `header:",omitempty"`
		I float32   `header:",omitempty"`
		J time.Time `header:",omitempty"`
		K struct{}  `header:",omitempty"`
	}{E: &str}

	v, err := Header(s)
	if err != nil {
		t.Errorf("Header(%#v) returned error: %v", s, err)
	}

	want := []string{
		"A: ",
		"E: ", // E is included because the pointer is not empty, even though the string being pointed to is
		"omitempty: ",
	}
	if !reflect.DeepEqual(want, v) {
		t.Errorf("Header(%#v) returned %v, want %v", s, v, want)
	}
}

type A struct {
	B
}

type B struct {
	C string
}

type D struct {
	B
	C string
}

type e struct {
	B
	C string
}

type F struct {
	e
}

func TestHeader_embeddedStructs(t *testing.T) {
	tests := []struct {
		in   interface{}
		want []string
	}{
		{
			A{B{C: "foo"}},
			[]string{"C: foo"},
		},
		{
			D{B: B{C: "bar"}, C: "foo"},
			[]string{"C: bar", "C: foo"},
		},
		{
			F{e{B: B{C: "bar"}, C: "foo"}}, // With unexported embed
			[]string{"C: bar", "C: foo"},
		},
	}

	for i, tt := range tests {
		v, err := Header(tt.in)
		if err != nil {
			t.Errorf("%d. Header(%+v) returned error: %v", i, tt.in, err)
		}

		if !reflect.DeepEqual(tt.want, v) {
			t.Errorf("%d. Header(%+v) returned %v, want %v", i, tt.in, v, tt.want)
		}
	}
}

func TestHeader_invalidInput(t *testing.T) {
	_, err := Header("")
	if err == nil {
		t.Errorf("expected Header() to return an error on invalid input")
	}
}

type EncodedArgs []string

func (m EncodedArgs) EncodeHeader(key string, v []string) ([]string, error) {
	newV := v
	for i, arg := range m {
		newV = append(newV, fmt.Sprintf("%s.%d: %s", key, i, arg))
	}
	return newV, nil
}

func TestHeader_Marshaler(t *testing.T) {
	s := struct {
		Args EncodedArgs `header:"Arg"`
	}{[]string{"a", "b", "c"}}
	v, err := Header(s)
	if err != nil {
		t.Errorf("Header(%+v) returned error: %v", s, err)
	}

	want := []string{
		"Arg.0: a",
		"Arg.1: b",
		"Arg.2: c",
	}
	if !reflect.DeepEqual(want, v) {
		t.Errorf("Header(%+v) returned %v, want %v", s, v, want)
	}
}

func TestHeader_MarshalerWithNilPointer(t *testing.T) {
	s := struct {
		Args *EncodedArgs `header:"Arg"`
	}{}
	v, err := Header(s)
	if err != nil {
		t.Errorf("Header(%+v) returned error: %v", s, err)
	}

	want := []string{}
	if !reflect.DeepEqual(want, v) {
		t.Errorf("Header(%+v) returned %v, want %v", s, v, want)
	}
}

func TestTagParsing(t *testing.T) {
	name, opts := parseTag("field,foobar,foo")
	if name != "field" {
		t.Fatalf("name = %+v, want field", name)
	}
	for _, tt := range []struct {
		opt  string
		want bool
	}{
		{"foobar", true},
		{"foo", true},
		{"bar", false},
		{"field", false},
	} {
		if opts.Contains(tt.opt) != tt.want {
			t.Errorf("Contains(%+v) = %v", tt.opt, !tt.want)
		}
	}
}
