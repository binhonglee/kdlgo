package kdlgo

import (
	"bufio"
	"bytes"
)

type kdlReader struct {
	line   int
	pos    int
	reader *bufio.Reader
}

func newKDLReader(r *bufio.Reader) *kdlReader {
	return &kdlReader{line: 1, pos: 0, reader: r}
}

func (kdlr *kdlReader) readRune() (rune, error) {
	r, _, err := kdlr.reader.ReadRune()
	if r == '\n' {
		kdlr.line++
		kdlr.pos = 0
	} else {
		kdlr.pos++
	}

	return r, err
}

func (kdlr *kdlReader) discardLine() error {
	_, err := kdlr.reader.ReadString('\n')
	if err != nil {
		return err
	}

	err = kdlr.reader.UnreadByte()
	return err
}

func (kdlr *kdlReader) discard(count int) {
	s, _ := kdlr.peekX(count)
	for _, b := range s {
		var nl byte = '\n'
		if b == nl {
			kdlr.line++
			kdlr.pos = 0
		} else {
			kdlr.pos++
		}
	}
	kdlr.reader.Discard(count)
}

func (kdlr *kdlReader) peekX(count int) ([]byte, error) {
	return kdlr.reader.Peek(count)
}

func (kdlr *kdlReader) peek() (rune, error) {
	r, _, err := kdlr.reader.ReadRune()
	if err != nil {
		return r, err
	}

	err = kdlr.reader.UnreadRune()
	return r, err
}

func (kdlr *kdlReader) unreadRune() error {
	err := kdlr.reader.UnreadRune()
	if err != nil {
		return err
	}

	peek, _ := kdlr.reader.Peek(1)
	var b byte = '\n'
	if peek[0] == b {
		kdlr.line--
	} else {
		kdlr.pos--
	}

	return nil
}

func (kdlr *kdlReader) isNext(charset []byte) (bool, error) {
	peek, err := kdlr.peekX(len(charset))
	if err != nil {
		return false, err
	}

	if bytes.Compare(peek, charset) == 0 {
		kdlr.discard(len(charset))
		return true, nil
	}

	return false, nil
}
