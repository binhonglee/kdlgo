package kdlgo

import (
	"bufio"
	"math/big"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const EOF = "EOF"

func ParseFile(fullfilepath string) (KDLObjects, error) {
	var t KDLObjects
	f, err := os.Open(fullfilepath)
	if err != nil {
		return t, err
	}
	r := bufio.NewReader(f)
	return ParseReader(r)
}

func ParseReader(reader *bufio.Reader) (KDLObjects, error) {
	setReader(reader)
	return parseObjects(false, "")
}

func parseObjects(hasOpen bool, key string) (KDLObjects, error) {
	var t KDLObjects
	var objects []KDLObject
	for {
		obj, err := parseObject()
		if err == nil {
			objects = append(objects, obj)
		} else if err.Error() == EOF || err.Error() == kdlEndOfObj {
			if obj != nil {
				objects = append(objects, obj)
			}
			return NewKDLObjects(key, objects), nil
		} else {
			return t, wrapError(err)
		}
	}
}

func parseObject() (KDLObject, error) {
	for {
		r, err := peek()
		if err != nil {
			return nil, err
		}

		if r == '}' {
			discard(1)
			return nil, endOfObjErr()
		}

		if !unicode.IsSpace(r) {
			break
		}

		discard(1)
	}

	key, err := parseKey()

	if err != nil {
		if err.Error() == kdlKeyOnly {
			return NewKDLDefault(key), nil
		}
		return nil, err
	}

	var objects []KDLObject
	for {
		r, err := readRune()
		if err != nil && err.Error() != EOF {
			return nil, err
		}

		if r == '\\' {
			peek, err := peek()
			if err == nil && peek == '\n' {
				discard(1)
				continue
			}
		}

		if r == '\n' || r == ';' ||
			(err != nil && err.Error() == EOF) {
			if len(objects) == 0 {
				return NewKDLDefault(key), nil
			} else if len(objects) == 1 {
				return objects[0], nil
			} else {
				return ConvertToDocument(objects)
			}
		} else if unicode.IsSpace(r) {
			continue
		}

		skipNext := false
		peek, _ := peek()
		if r == '/' && peek == '-' {
			discard(1)
			r, err = readRune()
			if err != nil {
				return nil, err
			}
			skipNext = true
		}

		obj, err := parseValue(key, r)
		if !skipNext {
			objects = append(objects, obj)
		}
		if err != nil {
			return nil, err
		}
	}
}

func parseKey() (string, error) {
	var key strings.Builder

	for {
		r, err := readRune()
		if err != nil {
			return key.String(), err
		}

		if unicode.IsSpace(r) {
			if len(key.String()) < 1 {
				continue
			} else if r == '\n' {
				return key.String(), keyOnlyErr()
			} else {
				return key.String(), nil
			}
		}

		invalid :=
			(len(key.String()) < 1 && unicode.IsNumber(r)) ||
				unicode.IsSpace(r) || r == '=' || r == '"'
		if invalid {
			return key.String(), invalidKeyCharErr()
		}
		key.WriteRune(r)
	}
}

func parseValue(key string, r rune) (KDLObject, error) {
	if unicode.IsNumber(r) {
		return parseNumber(key, r)
	}

	switch r {
	case '"':
		return parseString(key)
	case 'n':
		return parseNull(key)
	case 't':
		fallthrough
	case 'f':
		return parseBool(key, r)
	case 'r':
		return parseRawString(key)
	case '{':
		return parseObjects(true, key)
	}

	return nil, invalidSyntaxErr()
}

func parseString(key string) (KDLString, error) {
	var kdls KDLString
	var s strings.Builder

	for {
		r, err := readRune()
		if err != nil {
			return kdls, err
		}

		if r == '"' {
			return *NewKDLString(key, s.String()), nil
		}

		s.WriteRune(r)
	}
}

func parseRawString(key string) (KDLRawString, error) {
	var kdlrs KDLRawString
	var s strings.Builder

	count := 0

	for {
		r, err := readRune()
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
		r, err := readRune()
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

				r, err := readRune()
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

func parseNumber(key string, start rune) (KDLNumber, error) {
	var kdlnum KDLNumber
	var val strings.Builder
	val.WriteRune(start)

	for {
		r, err := peek()
		if r != ';' && r != '\n' {
			discard(1)
		}
		if err != nil && err.Error() != EOF {
			return kdlnum, err
		}

		if r == ';' || unicode.IsSpace(r) || (err != nil && err.Error() == EOF) {
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

func parseNull(key string) (KDLNull, error) {
	var kdlnull KDLNull
	charset := []rune{'u', 'l', 'l'}
	i := 0

	for {
		r, err := readRune()
		if err != nil {
			return kdlnull, err
		}

		if r != charset[i] {
			return kdlnull, invalidSyntaxErr()
		}
		i++
		if i == len(charset) {
			break
		}
	}

	return *NewKDLNull(key), nil
}

func parseBool(key string, start rune) (KDLBool, error) {
	var kdlbool KDLBool
	var charset []rune
	i := 0

	if start == 't' {
		charset = []rune{'r', 'u', 'e'}
	} else if start == 'f' {
		charset = []rune{'a', 'l', 's', 'e'}
	} else {
		return kdlbool, invalidSyntaxErr()
	}

	for {
		r, err := readRune()
		if err != nil {
			return kdlbool, err
		}

		if r != charset[i] {
			return kdlbool, invalidSyntaxErr()
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
		return doc, emptyArrayErr()
	}

	key = objs[0].GetKey()
	for _, obj := range objs {
		if obj.GetKey() != key {
			return doc, differentKeysErr()
		}

		vals = append(vals, obj.GetValue())
	}
	return *NewKDLDocument(key, vals), nil
}
