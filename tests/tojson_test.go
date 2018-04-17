package tests

import (
	"testing"
	"time"

	"github.com/utrack/petrichor/convertor"
)

type Types struct {
	F1 int           `json:"field_int" description:"This is an int" tags:"a, b, c, d, e" module:"module"`
	f2 string        `json:"field_string" description:"This is a string"`
	F3 time.Duration `json:"field_duration"`
	f4 float64       `json:"field_float64"`
	f5 bool          `json:"field_bool"`
	f6 []string      `json:"field_strings_enum"`
	F7 []int         `json:"field_int_enum"`
}

func TestNilStruct(t *testing.T) {
	_, err := convertor.Json.Convert(nil)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	t.Log(err)
}

func TestReflectStruct(t *testing.T) {
	types := Types{
		F1: 42,
		f2: "alibaba",
		F3: time.Hour,
		f4: 2.71,
	}
	j, err := convertor.Json.Convert(types)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(j))
}
