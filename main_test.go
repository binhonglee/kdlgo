package main

import (
	"strconv"
	"testing"
)

func TestParseFromFile(t *testing.T) {
	objs, err := ParseFile("test.kdl")
	if err != nil {
		t.Fatal(err)
	}
	if len(objs) != 5 {
		t.Fatal("There should be 5 KDLObjects. Got " + strconv.Itoa(len(objs)) + " instead.")
	}

	expected := []string{
		`firstkey {"first\n\ttab\nnewline\nval"}`,
		`thirdkey {true null}`,
		`secondkey {1.2E+01 "test" null false "testagain"}`,
		`anotherkey {"true" 1.23543E+05 null true}`,
		`moreKeys {false true}`,
	}

	for i, obj := range objs {
		s, err := KDLObjToString(obj)
		if err != nil {
			t.Fatal(err)
			return
		}
		if s != expected[i] {
			t.Error("Item number " + strconv.Itoa(i+1) + " is incorrectly parsed.")
		}
	}
}
