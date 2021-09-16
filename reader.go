package kdlgo

import (
	"bufio"
)

var (
	line   int
	pos    int
	reader *bufio.Reader
)

func setReader(r *bufio.Reader) {
	line = 1
	pos = 0
	reader = r
}

func readRune() (rune, error) {
	r, _, err := reader.ReadRune()
	if r == '\n' {
		line++
		pos = 0
	} else {
		pos++
	}

	return r, err
}

func discard(count int) {
	reader.Discard(count)
}

func peekX(count int) ([]byte, error) {
	return reader.Peek(count)
}

func peek() (rune, error) {
	r, _, err := reader.ReadRune()
	if err != nil {
		return r, err
	}

	err = reader.UnreadRune()
	return r, err
}
