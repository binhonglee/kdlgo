package kdlgo

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const (
	eof = "EOF"

	asterisk  = '*'
	backslash = '\\'
	dash      = '-'
	dquote    = '"'
	equals    = '='
	newline   = '\n'
	pound     = '#'
	semicolon = ';'
	slash     = '/'

	openBracket  = '{'
	closeBracket = '}'
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

func ParseReader(reader *bufio.Reader) (KDLObjects, error) {
	r := newKDLReader(reader)
	return parseObjects(r, false, "")
}

func parseObjects(kdlr *kdlReader, hasOpen bool, key string) (KDLObjects, error) {
	var t KDLObjects
	var objects []KDLObject
	for {
		obj, err := parseObject(kdlr)
		if err == nil {
			if obj != nil {
				objects = append(objects, obj)
			}
		} else if err.Error() == eof || err.Error() == kdlEndOfObj {
			if obj != nil {
				objects = append(objects, obj)
			}
			return NewKDLObjects(key, objects), nil
		} else {
			return t, wrapError(kdlr, err)
		}
	}
}

func parseObject(kdlr *kdlReader) (KDLObject, error) {
	for {
		err := blockComment(kdlr)
		if err != nil {
			return nil, err
		}

		r, err := kdlr.peek()
		if err != nil {
			return nil, err
		}

		if r == closeBracket {
			kdlr.discard(1)
			return nil, endOfObjErr()
		}

		skipLine, err := lineComment(kdlr)
		if err != nil {
			if err.Error() == eof && skipLine {
				return nil, nil
			}
			return nil, err
		}

		if skipLine {
			continue
		}

		if unicode.IsSpace(r) {
			kdlr.discard(1)
			continue
		}

		break
	}

	key, err := parseKey(kdlr)

	if err != nil {
		if err.Error() == kdlKeyOnly {
			return NewKDLDefault(key), nil
		}
		return nil, err
	}

	var objects []KDLObject
	for {
		err = blockComment(kdlr)
		if err != nil && err.Error() != eof {
			return nil, err
		}

		r, err := kdlr.readRune()
		if err != nil && err.Error() != eof {
			return nil, err
		}

		if r == backslash {
			peek, err := kdlr.peek()
			if err == nil && peek == newline {
				kdlr.discard(1)
				continue
			}
		}

		if r == newline || r == semicolon ||
			(err != nil && err.Error() == eof) {
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

		kdlr.unreadRune()
		skipNext, _ := kdlr.isNext([]byte{slash, dash})
		if skipNext {
			r, err = kdlr.peek()
			if err != nil {
				if err.Error() == eof {
					return ConvertToDocument(objects)
				}
				return nil, err
			}
		}

		skipLine, err := lineComment(kdlr)
		if err != nil {
			if err.Error() == eof && skipLine {
				return ConvertToDocument(objects)
			}
			return nil, err
		}

		if skipLine {
			continue
		}

		obj, err := parseValue(kdlr, key, r)
		if !skipNext {
			objects = append(objects, obj)
		}
		if err != nil {
			return nil, err
		}
	}
}

func parseKey(kdlr *kdlReader) (string, error) {
	var key strings.Builder

	for {
		r, err := kdlr.readRune()
		if err != nil {
			return key.String(), err
		}

		if unicode.IsSpace(r) {
			if len(key.String()) < 1 {
				continue
			} else if r == newline {
				return key.String(), keyOnlyErr()
			} else {
				return key.String(), nil
			}
		}

		invalid :=
			(len(key.String()) < 1 && unicode.IsNumber(r)) ||
				unicode.IsSpace(r) || r == equals || r == dquote
		if invalid {
			return key.String(), invalidKeyCharErr()
		}
		key.WriteRune(r)
	}
}

func parseValue(kdlr *kdlReader, key string, r rune) (KDLObject, error) {
	if unicode.IsNumber(r) {
		kdlr.discard(1)
		return parseNumber(kdlr, key, r)
	}

	switch r {
	case dquote:
		kdlr.discard(1)
		return parseString(kdlr, key)
	case 'n':
		return parseNull(kdlr, key)
	case 't':
		fallthrough
	case 'f':
		return parseBool(kdlr, key, r)
	case 'r':
		kdlr.discard(1)
		return parseRawString(kdlr, key)
	case openBracket:
		kdlr.discard(1)
		return parseObjects(kdlr, true, key)
	}

	return nil, invalidSyntaxErr()
}

func parseString(kdlr *kdlReader, key string) (KDLString, error) {
	var kdls KDLString
	var s strings.Builder

	for {
		r, err := kdlr.readRune()
		if err != nil {
			return kdls, err
		}

		if r == backslash {
			var b byte = '"'
			next, err := kdlr.isNext([]byte{b})
			if err != nil {
				return kdls, err
			}

			if next {
				s.WriteRune(r)
				s.WriteByte(b)
				continue
			}
		}

		if r == dquote {
			return NewKDLString(key, s.String()), nil
		}

		s.WriteRune(r)
	}
}

func parseRawString(kdlr *kdlReader, key string) (KDLRawString, error) {
	var kdlrs KDLRawString
	var s strings.Builder

	count := 0

	for {
		r, err := kdlr.readRune()
		if err != nil {
			return kdlrs, err
		}

		if r == pound {
			count++
			continue
		}

		if r == dquote {
			break
		}
	}

	for {
		r, err := kdlr.readRune()
		if err != nil {
			return kdlrs, err
		}

		for {
			if r != dquote {
				s.WriteRune(r)
				break
			}

			var temp strings.Builder
			tempCount := 0
			temp.WriteRune(r)

			for {
				if tempCount == count {
					return NewKDLRawString(key, s.String()), nil
				}

				r, err := kdlr.readRune()
				if err != nil {
					return kdlrs, err
				}

				if r != pound {
					break
				}

				tempCount++
				temp.WriteRune(r)
			}

			s.WriteString(temp.String())
		}
	}
}

func parseNumber(kdlr *kdlReader, key string, start rune) (KDLNumber, error) {
	var kdlnum KDLNumber
	var val strings.Builder
	val.WriteRune(start)

	for {
		r, err := kdlr.peek()
		if err != nil && err.Error() != eof {
			return kdlnum, err
		}
		if r != semicolon && r != newline && r != slash {
			kdlr.discard(1)
		}

		if r == semicolon || unicode.IsSpace(r) ||
			r == slash || (err != nil && err.Error() == eof) {
			value, err := strconv.ParseFloat(val.String(), 64)
			if err != nil {
				return kdlnum, err
			}
			return NewKDLNumber(key, value), nil
		}

		val.WriteRune(r)
	}
}

func parseNull(kdlr *kdlReader, key string) (KDLNull, error) {
	var kdlnull KDLNull
	charset := []byte{'n', 'u', 'l', 'l'}
	next, err := kdlr.isNext(charset)
	if err != nil {
		return kdlnull, err
	}

	if next {
		return NewKDLNull(key), nil
	}

	return kdlnull, invalidSyntaxErr()
}

func parseBool(kdlr *kdlReader, key string, start rune) (KDLBool, error) {
	var kdlbool KDLBool
	var charset []byte

	if start == 't' {
		charset = []byte{'t', 'r', 'u', 'e'}
	} else if start == 'f' {
		charset = []byte{'f', 'a', 'l', 's', 'e'}
	} else {
		return kdlbool, invalidSyntaxErr()
	}

	next, err := kdlr.isNext(charset)
	if err != nil {
		return kdlbool, err
	}

	if next {
		return NewKDLBool(key, start == 't'), nil
	}
	return kdlbool, invalidSyntaxErr()
}

func lineComment(kdlr *kdlReader) (bool, error) {
	skipLine, _ := kdlr.isNext([]byte{slash, slash})
	if skipLine {
		err := kdlr.discardLine()
		if err != nil && err.Error() != eof {
			return false, err
		}
		return true, err
	}
	return false, nil
}

func blockComment(kdlr *kdlReader) error {
	count := 0
	open := []byte{slash, asterisk}
	close := []byte{asterisk, slash}

	for {
		isBlock, err := kdlr.isNext(open)
		if err != nil {
			return err
		}

		if isBlock {
			count++
		}

		break
	}

	for {
		if count == 0 {
			return nil
		}

		isOpen, err := kdlr.isNext(open)
		if err != nil {
			return err
		}

		if isOpen {
			count++
			continue
		}

		isClose, err := kdlr.isNext(close)
		if err != nil {
			return err
		}

		if isClose {
			count--
			continue
		}

		kdlr.discard(1)
	}
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
