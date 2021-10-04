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

	openBracket      = '{'
	closeBracket     = '}'
	openParenthesis  = '('
	closeParenthesis = ')'
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

	skipNext, _ := kdlr.isNext([]byte{slash, dash})
	if skipNext {
		parseKey(kdlr)
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

		obj, err := parseVal(kdlr, key, r)
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
	isQuoted := false
	prev := newline

	for {
		r, err := kdlr.readRune()
		if err != nil {
			if err.Error() == eof {
				err = unexpectedEOFErr()
			}
			return key.String(), err
		}

		if (!isQuoted && unicode.IsSpace(r)) || r == newline ||
			((unicode.IsSpace(r) || r == equals) && prev == dquote) {
			if key.Len() < 1 {
				continue
			} else if r == newline {
				return checkQuotedString(key), keyOnlyErr()
			} else {
				return checkQuotedString(key), nil
			}
		}

		invalid :=
			(key.Len() < 1 && unicode.IsNumber(r)) ||
				(!isQuoted && unicode.IsSpace(r)) || r == equals
		if invalid {
			return key.String(), invalidKeyCharErr()
		}

		if key.Len() < 1 {
			isQuoted = r == dquote
		}
		if prev == backslash && r == backslash {
			prev = newline
		} else if prev == backslash && r == dquote {
			prev = newline
		} else {
			prev = r
		}
		key.WriteRune(r)
	}
}

func checkQuotedString(s strings.Builder) string {
	ss := s.String()
	unquoted, err := strconv.Unquote(ss)
	if err != nil {
		return ss
	} else {
		return unquoted
	}
}

func parseVal(kdlr *kdlReader, key string, r rune) (KDLObject, error) {
	value, err := parseValue(kdlr, key, r)
	if err == nil {
		return value, nil
	}

	node, err := parseKey(kdlr)
	if err != nil && err.Error() != KDLInvalidKeyChar {
		if err.Error() == kdlKeyOnly {
			return NewKDLDefault(node), nil
		}
		return nil, err
	}

	if kdlr.lastRead() != equals {
		return NewKDLDefault(node), nil
	}
	r, err = kdlr.peek()
	if err != nil {
		return nil, err
	}

	obj, err := parseValue(kdlr, node, r)
	if err != nil {
		return nil, err
	}

	return NewKDLObjects(key, []KDLObject{obj}), nil
}

func parseValue(kdlr *kdlReader, key string, r rune) (KDLObject, error) {
	if unicode.IsNumber(r) {
		kdlr.discard(1)
		return parseNumber(kdlr, key, r)
	}

	switch r {
	case dquote:
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
	s, err := parseQuotedString(kdlr)
	if err != nil {
		return kdls, err
	}
	return NewKDLString(key, s), nil
}

func parseQuotedString(kdlr *kdlReader) (string, error) {
	count := 2

	for {
		bytes, err := kdlr.peekX(count)
		if err != nil {
			kdlr.discard(count)
			return string(bytes[1:]), err
		}
		r := rune(bytes[len(bytes)-1])

		if r == backslash {
			bs, err := kdlr.peekX(count + 1)
			if err != nil {
				kdlr.discard(count)
				return string(bytes[1:]), err
			}
			next := bs[len(bs)-1] == byte(dquote)

			if next {
				count += 2
				continue
			}
		}

		if r == dquote {
			toRet := string(bytes[1 : len(bytes)-1])
			temp, err := kdlr.peekX(count + 1)
			if err != nil {
				if err.Error() != eof {

					return toRet, err
				}
				kdlr.discard(count)
				return toRet, nil
			}

			r = rune(temp[len(temp)-1])
			if !(unicode.IsSpace(r) || r == semicolon) {
				return toRet, invalidSyntaxErr()
			}
			kdlr.discard(count)
			return toRet, nil
		}

		count++
	}
}

func parseRawString(kdlr *kdlReader, key string) (KDLRawString, error) {
	var kdlrs KDLRawString
	count := 0
	length := 0

	for {
		length++
		bytes, err := kdlr.peekX(length)
		if err != nil {
			return kdlrs, err
		}
		r := rune(bytes[len(bytes)-1])

		if r == pound {
			count++
			continue
		}

		if r == dquote {
			break
		}

		return kdlrs, invalidSyntaxErr()
	}

	start := length
	length++
	poundCount := 0
	dqStart := false

	for {
		bytes, err := kdlr.peekX(length)
		if err != nil {
			return kdlrs, err
		}
		r := rune(bytes[len(bytes)-1])

		if r == dquote {
			dqStart = true
			length++
			continue
		}

		if dqStart && r == pound {
			poundCount++
		} else {
			poundCount = 0
			dqStart = false
		}

		if poundCount == count {
			kdlr.discard(length)
			return NewKDLRawString(key, string(bytes[start:len(bytes)-count-1])), nil
		}

		length++
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
