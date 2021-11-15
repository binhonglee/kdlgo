
package testers

import (
	"strconv"
	"testing"

	"github.com/binhonglee/kdlgo"
)

func TestRAWSTRINGARG(t *testing.T) {
	objs, err := kdlgo.ParseFile("../kdls/raw_string_arg.kdl")
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{
		`node_1 "" { "g\\n\""; } { node_2; } "\"arg\\n\"and stuff"`,		`node_3 "#\"arg\\n\"#and stuff"`,
	}

	if len(objs.GetValue().Objects) != len(expected) {
		t.Fatal(
			"There should be " + strconv.Itoa(len(expected)) +
				" KDLObjects. Got " + strconv.Itoa(len(objs.GetValue().Objects)) + " instead.",
		)
	}

	for i, obj := range objs.GetValue().Objects {
		s, err := kdlgo.RecreateKDLObj(obj)
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
