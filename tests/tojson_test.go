package tests

import (
	"time"
	"testing"

	"github.com/utrack/petrichor/convertor"
)

type Types struct {
	F1 int           `json:field_int,description:"This is an int",tags:["tag1", "Tag2"],module:"module"`
	f2 string        `json:field_string,description:"This is a string"`
	F3 time.Duration `json:field_duration`
	f4 float64       `json:field_float64`
}

func TestNilStruct(t *testing.T) {
	j, err := convertor.Json.C(t, nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(j)
}

func TestReflectStruct(t *testing.T) {
	types := Types{
		F1: 42,
		f2: "alibaba",
		F3: time.Hour,
		f4: 2.71,
	}
	j, err := convertor.Json.C(t, types)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(j))
}