package kdlgo

import (
	"bufio"
	"os"
	"strings"
)

func ParseFile(fullfilepath string) (KDLObjects, error) {
	var t KDLObjects
	f, err := os.Open(fullfilepath)
	if err != nil {
		return t, err
	}
	r := bufio.NewReader(f)
	return ParseReader(r)
}

func ParseString(toParse string) (KDLObjects, error) {
	return ParseReader(bufio.NewReader(strings.NewReader(toParse)))
}

func ParseReader(reader *bufio.Reader) (KDLObjects, error) {
	r := newKDLReader(reader)
	return parseObjects(r, false, "")
}

func ConvertToDocument(objs []KDLObject) (KDLDocument, error) {
	var key string
	var vals []KDLValue
	var doc KDLDocument

	if len(objs) < 1 {
		return doc, emptyArrayErr()
	}

	key = objs[0].GetKey()
	for _, obj := range objs {
		if obj.GetKey() != key {
			return doc, differentKeysErr()
		}

		vals = append(vals, obj.GetValue())
	}
	return NewKDLDocument(key, vals), nil
}
