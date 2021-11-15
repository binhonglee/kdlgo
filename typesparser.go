package kdlgo

import (
	"strconv"
	"strings"
	"unicode"
)

func checkQuotedString(s strings.Builder) string {
	ss := s.String()
	unquoted, err := strconv.Unquote(ss)
	if err != nil {
		return ss
	} else {
		return unquoted
	}
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
			return stringEscape(string(bytes[1:])), err
		}
		r := rune(bytes[len(bytes)-1])

		if r == backslash {
			bs, err := kdlr.peekX(count + 1)
			if err != nil {
				kdlr.discard(count)
				return stringEscape(string(bytes[1:])), err
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
			toRet = stringEscape(toRet)
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

func stringEscape(s string) string {
	return strings.ReplaceAll(s, "\\/", "/")
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

func parseNumber(kdlr *kdlReader, key string) (KDLNumber, error) {
	var kdlnum KDLNumber
	length := 0
	dotCount := 0

	for {
		length++
		bytes, err := kdlr.peekX(length)
		if err != nil && err.Error() != eof {
			return kdlnum, err
		}
		r := rune(bytes[len(bytes)-1])
		if r == dot {
			dotCount++
			if dotCount > 1 {
				return kdlnum, invalidNumValueErr()
			}
		}

		if r == semicolon || unicode.IsSpace(r) ||
			r == slash || (err != nil && err.Error() == eof) {
			rawStr := string(bytes[0 : len(bytes)-1])
			if err != nil && err.Error() == eof {
				rawStr = string(bytes)
			}
			str := strings.ReplaceAll(rawStr, "_", "")
			value, err := strconv.ParseFloat(str, 64)
			if err != nil {
				val, err := strconv.ParseInt(str, 0, 10)
				if err != nil {
					return kdlnum, err
				}
				value = float64(val)
			}
			kdlr.discard(length - 1)
			return NewKDLNumber(key, value), nil
		}
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
