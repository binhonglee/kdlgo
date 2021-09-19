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
		`firstkey "first\n\ttab\nnewline\"\nval" "testing\""`,
		`numbers 543 234 85720394`,
		`thirdkey true null`,
		`secondkey 12 "test" null false "testagain"`,
		`anotherkey "true" 123.543 null true`,
		`moreKeys false true`,
		`keyonly`,
		`testcomment`,
		`objects { node1 12; node2 "string"; node3 null; }`,
		`multiline-node "random"`,
		`title "Some title"`,
		`"quoted node" "quoted value"`,
		`"quoted node for numbers" 21 43 465 "string"`,
		`smile "😁"`,
		`foo123~!@#$%^&*.:'|/?+ "weeee"`,
		`test "value"`,
	}

	if len(objs.GetValue().Objects) != len(expected) {
		t.Fatal(
			"There should be " + strconv.Itoa(len(expected)) +
				" KDLObjects. Got " + strconv.Itoa(len(objs.GetValue().Objects)) + " instead.",
		)
	}

	for i, obj := range objs.GetValue().Objects {
		s, err := RecreateKDLObj(obj)
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
