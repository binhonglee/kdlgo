package main

import (
	"bufio"
	"errors"
	"math/big"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func ParseFile(fullfilepath string) ([]KDLObject, error) {
	f, err := os.Open(fullfilepath)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(f)
	return ParseReader(r)
}

func ParseReader(reader *bufio.Reader) ([]KDLObject, error) {
	var objects []KDLObject
	for {
		obj, err := parseObject(reader)
		if err == nil {
			objects = append(objects, obj)
		} else if err.Error() == "EOF" {
			if obj != nil {
				objects = append(objects, obj)
			}
			return objects, nil
		} else {
			return nil, err
		}
	}
}

func parseObject(reader *bufio.Reader) (KDLObject, error) {
	key, err := parseKey(reader)
	if err != nil {
		return nil, err
	}

	var objects []KDLObject
	for {
		r, _, err := reader.ReadRune()
		if err != nil && err.Error() != "EOF" {
			return nil, err
		}

		if r == '\n' || (err != nil && err.Error() == "EOF") {
			if len(objects) == 0 {
				return nil, errors.New("Missing value")
			} else if len(objects) == 1 {
				return objects[0], nil
			} else {
				return ConvertToDocument(key, objects)
			}
		} else if unicode.IsSpace(r) {
			continue
		}

		obj, err := parseValue(reader, key, r)
		objects = append(objects, obj)
		if err != nil {
			return nil, err
		}
	}
}

func parseKey(reader *bufio.Reader) (string, error) {
	var key strings.Builder

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return key.String(), err
		}

		if unicode.IsSpace(r) {
			if len(key.String()) < 1 {
				continue
			} else {
				return key.String(), nil
			}
		}

		var valid bool

		valid = unicode.IsLetter(r)
		if !valid && len(key.String()) > 1 {
			valid = unicode.IsNumber(r)
		}

		if !valid {
			return key.String(), errors.New("Invalid character for key.")
		}
		key.WriteRune(r)
	}
}

func parseValue(reader *bufio.Reader, key string, r rune) (KDLObject, error) {
	if unicode.IsNumber(r) {
		return parseNumber(reader, key, r)
	}

	switch r {
	case '"':
		return parseString(reader, key)
	case 'n':
		return parseNull(reader, key)
	case 't':
		fallthrough
	case 'f':
		return parseBool(reader, key, r)
	case 'r':
		return parseRawString(reader, key)
	}

	return nil, errors.New("Eh")
}

func parseString(reader *bufio.Reader, key string) (KDLString, error) {
	var kdls KDLString
	var s strings.Builder

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return kdls, err
		}

		if r == '"' {
			return *NewKDLString(key, s.String()), nil
		}

		s.WriteRune(r)
	}
}

func parseRawString(reader *bufio.Reader, key string) (KDLRawString, error) {
	var kdlrs KDLRawString
	var s strings.Builder

	count := 0

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return kdlrs, err
		}

		if r == '#' {
			count++
			continue
		}

		if r == '"' {
			break
		}
	}

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return kdlrs, err
		}

		for {
			if r != '"' {
				s.WriteRune(r)
				break
			}

			var temp strings.Builder
			tempCount := 0
			temp.WriteRune(r)

			for {
				if tempCount == count {
					return *NewKDLRawString(key, s.String()), nil
				}

				r, _, err := reader.ReadRune()
				if err != nil {
					return kdlrs, err
				}

				if r != '#' {
					break
				}

				tempCount++
				temp.WriteRune(r)
			}

			s.WriteString(temp.String())
		}
	}
}

func parseNumber(reader *bufio.Reader, key string, start rune) (KDLNumber, error) {
	var kdlnum KDLNumber
	var val strings.Builder
	val.WriteRune(start)

	for {
		r, _, err := reader.ReadRune()
		if err != nil && err.Error() != "EOF" {
			return kdlnum, err
		}

		if unicode.IsSpace(r) || (err != nil && err.Error() == "EOF") {
			value, err := strconv.ParseFloat(val.String(), 64)
			if err != nil {
				return kdlnum, err
			}
			kdlnum.key = key
			kdlnum.value.Number = *big.NewFloat(value)
			kdlnum.value.Type = KDLNumberType
			return kdlnum, nil
		}

		val.WriteRune(r)
	}
}

func parseNull(reader *bufio.Reader, key string) (KDLNull, error) {
	var kdlnull KDLNull
	charset := []rune{'u', 'l', 'l'}
	i := 0

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return kdlnull, err
		}

		if r != charset[i] {
			return kdlnull, errors.New("Invalid syntax")
		}
		i++
		if i == len(charset) {
			break
		}
	}

	return *NewKDLNull(key), nil
}

func parseBool(reader *bufio.Reader, key string, start rune) (KDLBool, error) {
	var kdlbool KDLBool
	var charset []rune
	i := 0

	if start == 't' {
		charset = []rune{'r', 'u', 'e'}
	} else if start == 'f' {
		charset = []rune{'a', 'l', 's', 'e'}
	} else {
		return kdlbool, errors.New("Invalid syntax")
	}

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return kdlbool, err
		}

		if r != charset[i] {
			return kdlbool, errors.New("Invalid syntax")
		}
		i++
		if i == len(charset) {
			break
		}
	}

	return *NewKDLBool(key, start == 't'), nil
}

func ConvertToDocument(objs []KDLObject) (KDLDocument, error) {
	var key string
	var vals []KDLValue
	var doc KDLDocument

	if len(objs) < 1 {
		return doc, errors.New("Empty array")
	}

	key = objs[0].GetKey()
	for _, obj := range objs {
		if obj.GetKey() != key {
			return doc, errors.New("Different key found in the array of KDLObject")
		}

		vals = append(vals, obj.GetValue())
	}
	return *NewKDLDocument(key, vals), nil
}
