
package testers

import (
	"strconv"
	"testing"

	"github.com/binhonglee/kdlgo"
)

func TestNULLPROPKEY(t *testing.T) {
	objs, err := kdlgo.ParseFile("../kdls/null_prop_key.kdl")
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{
		`node null {  1; }`,
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
