package kdlgo

import (
	"strings"
	"unicode"
)

const (
	eof = "EOF"

	asterisk   = '*'
	backslash  = '\\'
	dash       = '-'
	dot        = '.'
	dquote     = '"'
	equals     = '='
	newline    = '\n'
	pound      = '#'
	semicolon  = ';'
	slash      = '/'
	space      = ' '
	underscore = '_'

	openBracket      = '{'
	closeBracket     = '}'
	openParenthesis  = '('
	closeParenthesis = ')'
)

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
		if err != nil {
			if err.Error() == kdlEndOfObj {
				return ConvertToDocument(objects)
			}
			return nil, err
		}
		if !skipNext {
			objects = append(objects, obj)
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
				return checkQuotedString(key), keyOnlyErr()
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

func parseVal(kdlr *kdlReader, key string, r rune) (KDLObject, error) {
	value, err := parseValue(kdlr, key, r)
	if err == nil || err.Error() == KDLInvalidNumValue {
		return value, err
	}

	if err.Error() == kdlEndOfObj {
		return value, err
	}

	node, err := parseKey(kdlr)

	if err != nil && err.Error() != KDLInvalidKeyChar {
		if err.Error() == kdlKeyOnly {
			return NewKDLObjects(key, []KDLObject{NewKDLDefault(node)}), nil
		}
		return nil, err
	}

	if kdlr.lastRead() != equals {
		return NewKDLObjects(key, []KDLObject{NewKDLDefault(node)}), nil
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
	for {
		if r != space {
			break
		}
		kdlr.discard(1)

		var err error
		r, err = kdlr.peek()
		if err != nil {
			return nil, err
		}
	}

	if unicode.IsNumber(r) {
		return parseNumber(kdlr, key)
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
	case closeBracket:
		return nil, endOfObjErr()
	}

	return nil, invalidSyntaxErr()
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
