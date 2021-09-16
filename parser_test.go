package kdlgo

import (
	"strconv"
	"testing"
)

func TestParseFromFile(t *testing.T) {
	objs, err := ParseFile("test.kdl")
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{
		`firstkey "first\n\ttab\nnewline\nval" r"testing""`,
		`thirdkey true null`,
		`secondkey 1.2E+01 "test" null false "testagain"`,
		`anotherkey "true" 1.23543E+05 null true`,
		`moreKeys false true`,
		`keyonly`,
		`testcomment`,
		`objects { node1 1.2E+01; node2 "string"; node3 null; }`,
		`multiline-node "random"`,
	}

	if len(objs.GetValue().Objects) != len(expected) {
		t.Fatal(
			"There should be " + strconv.Itoa(len(expected)) +
				" KDLObjects. Got " + strconv.Itoa(len(objs.GetValue().Objects)) + " instead.",
		)
	}

	for i, obj := range objs.GetValue().Objects {
		s, err := KDLObjToString(obj)
		if err != nil {
			t.Fatal(err)
			return
		}
		if s != expected[i] {
			t.Error(
				"Item number "+strconv.Itoa(i+1)+" is incorrectly parsed.\n",
				"Expected: '"+expected[i]+"' but got '"+s+"' instead",
			)
		}
	}
}
