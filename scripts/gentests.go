package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/binhonglee/kdlgo"
)

func main() {
	files, err := ioutil.ReadDir("../tests/kdls/")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		name := file.Name()
		if !strings.HasSuffix(name, ".kdl") {
			continue
		}

		objs, _ := kdlgo.ParseFile("../tests/kdls/" + name)

		f, _ := os.Create("../tests/testers/" + strings.ReplaceAll(name, ".kdl", "_test.go"))
		defer f.Close()

		f.WriteString(`
package testers

import (
	"strconv"
	"testing"

	"github.com/binhonglee/kdlgo"
)

func Test` + strings.ToUpper(strings.Join(strings.Split(strings.TrimRight(name, ".kdl"), "_"), "")) + `(t *testing.T) {
`)
		f.WriteString("	objs, err := kdlgo.ParseFile(\"../kdls/" + name + "\")")
		f.WriteString(`
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{
`)
		for _, obj := range objs.GetValue().Objects {
			s, _ := kdlgo.RecreateKDLObj(obj)
			f.WriteString("		`" + s + "`,")
		}

		f.WriteString(`
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
`)
	}
}
