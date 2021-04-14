package yaml2json

import "testing"
import "encoding/json"
import "strings"

import "reflect"

func roundtrip(in string, options Options, target interface{}) error {
	result, err := Process([]byte(in), options)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(strings.NewReader(result))

	err = dec.Decode(target)
	return err
}

func TestScalarString(t *testing.T) {
	in := "a simple string"

	var result string

	err := roundtrip(in, Options{CompactJSON: true, MultipleYAML: false}, &result)

	if err != nil {
		t.Error(err)
	}

	if strings.Compare(result, in) != 0 {
		t.Errorf("Got '%s' but wanted '%s'", result, in)
	}
}

func TestScalarBoolTrue(t *testing.T) {
	in := "true"

	var result bool

	err := roundtrip(in, Options{CompactJSON: true, MultipleYAML: false}, &result)

	if err != nil {
		t.Error(err)
	}

	if !result {
		t.Errorf("Got '%t' but wanted '%s'", result, in)
	}
}

func TestScalarBoolFalse(t *testing.T) {
	in := "false"

	var result bool

	err := roundtrip(in, Options{CompactJSON: true, MultipleYAML: false}, &result)

	if err != nil {
		t.Error(err)
	}

	if result {
		t.Errorf("Got '%t' but wanted '%s'", result, in)
	}
}

func TestScalarFloatIntegerZero(t *testing.T) {
	in := "0"

	var result float64

	err := roundtrip(in, Options{CompactJSON: true, MultipleYAML: false}, &result)

	if err != nil {
		t.Error(err)
	}

	if result != 0 {
		t.Errorf("Got '%G' but wanted '%s'", result, in)
	}
}

func TestScalarFloatInteger(t *testing.T) {
	in := "172537253"

	var result float64

	err := roundtrip(in, Options{CompactJSON: true, MultipleYAML: false}, &result)

	if err != nil {
		t.Error(err)
	}

	if result != 172537253 {
		t.Errorf("Got '%G' but wanted '%s'", result, in)
	}
}

func TestScalarFloat(t *testing.T) {
	in := "172537253.12345"

	var result float64

	err := roundtrip(in, Options{CompactJSON: true, MultipleYAML: false}, &result)

	if err != nil {
		t.Error(err)
	}

	if result != 172537253.12345 {
		t.Errorf("Got '%G' but wanted '%s'", result, in)
	}
}

func TestScalarFloatExponent(t *testing.T) {
	in := "172537253.12345e5"

	var result float64

	err := roundtrip(in, Options{CompactJSON: true, MultipleYAML: false}, &result)

	if err != nil {
		t.Error(err)
	}

	if result != 172537253.12345e5 {
		t.Errorf("Got '%G' but wanted '%s'", result, in)
	}
}

type TestSimpleObjectT struct {
	Value1 bool    `json:"value1"`
	Value2 float64 `json:"value2"`
	Value3 string  `json:"value3"`
	Value4 string  `json:"VALue4"` // intentionally mixed casing
}

func TestSimpleObject(t *testing.T) {
	in := `value1: true
value2: 1.35e91
value3: |- # |- removes trailing new line characters but retains embedded newlines
  :>< a test
  value
VALue4: "{}"`

	want := TestSimpleObjectT{
		Value1: true,
		Value2: 1.35e91,
		Value3: ":>< a test\nvalue",
		Value4: "{}",
	}

	var result TestSimpleObjectT

	err := roundtrip(in, Options{CompactJSON: true, MultipleYAML: false}, &result)

	if err != nil {
		t.Error(err)
	}

	if result != want {
		t.Errorf("Got '%+v' but wanted '%+v'", result, want)
	}
}

type TestComplexObjectChildT struct {
	Value1 string `json:"value1"`
	Value2 string `json:"value2"`
}

type TestComplexObjectT struct {
	Value1   bool                      `json:"value1"`
	Value2   float64                   `json:"value2"`
	Value3   string                    `json:"value3"`
	Value4   string                    `json:"VALue4"` // intentionally mixed casing
	Children []TestComplexObjectChildT `json:"Children"`
}

func TestComplexObject(t *testing.T) {
	in := `value1: true
value2: 1.35e91
value3: |- # |- removes trailing new line characters but retains embedded newlines
  :>< a test
  value
VALue4: "{}"
children:
  - value1: child1
    value2: data1
  - value1: child2
    value2: data2
  - value1: child3
    value2: data3`

	want := TestComplexObjectT{
		Value1: true,
		Value2: 1.35e91,
		Value3: ":>< a test\nvalue",
		Value4: "{}",
		Children: []TestComplexObjectChildT{
			TestComplexObjectChildT{Value1: "child1", Value2: "data1"},
			TestComplexObjectChildT{Value1: "child2", Value2: "data2"},
			TestComplexObjectChildT{Value1: "child3", Value2: "data3"},
		},
	}

	var result TestComplexObjectT

	err := roundtrip(in, Options{CompactJSON: true, MultipleYAML: false}, &result)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("Got '%+v' but wanted '%+v'", result, want)
	}
}

func TestSimpleObjectMulti(t *testing.T) {
	in := `value1: true
value2: 1.35e91
value3: |- # |- removes trailing new line characters but retains embedded newlines
  :>< a test
  value
VALue4: "{}"
---
value1: true
value2: 1.35e91
value3: second one
VALue4: ""
---
value1: false
value2: 0
value3: third one
VALue4: "---"
`

	want := []TestSimpleObjectT{
		TestSimpleObjectT{true, 1.35e91, ":>< a test\nvalue", "{}"},
		TestSimpleObjectT{true, 1.35e91, "second one", ""},
		TestSimpleObjectT{false, 0, "third one", "---"},
	}

	result := make([]TestSimpleObjectT, 0)

	err := roundtrip(in, Options{CompactJSON: true, MultipleYAML: true}, &result)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("Got '%+v' but wanted '%+v'", result, want)
	}
}
