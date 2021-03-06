package kdlgo

import (
	"errors"
	"strconv"
)

type KDLErrorType string

const (
	KDLEmptyArray      = "Array is empty"
	KDLDifferentKey    = "All keys of KDLObject to convert to document should be the same"
	KDLInvalidKeyChar  = "Invalid character for key"
	KDLInvalidNumValue = "Invalid numeric value"
	KDLInvalidSyntax   = "Invalid syntax"
	KDLInvalidType     = "Invalid KDLType"
	KDLUnexpectedEOF   = "Unexpected end of file"

	// These should be caught and handled internally
	kdlKeyOnly     = "Internal only: Key only"
	kdlEndOfObj    = "Internal only: End of KDLObject"
	kdlNothingLeft = "Internal only: Nothing else left to parse"
)

func wrapError(kdlr *kdlReader, err error) error {
	return errors.New(
		err.Error() + "\nOn line " + strconv.Itoa(kdlr.line) +
			" column " + strconv.Itoa(kdlr.pos),
	)
}

func differentKeysErr() error {
	return errors.New(KDLDifferentKey)
}

func emptyArrayErr() error {
	return errors.New(KDLEmptyArray)
}

func invalidKeyCharErr() error {
	return errors.New(KDLInvalidKeyChar)
}

func invalidNumValueErr() error {
	return errors.New(KDLInvalidNumValue)
}

func invalidSyntaxErr() error {
	return errors.New(KDLInvalidSyntax)
}

func invalidTypeErr() error {
	return errors.New(KDLInvalidType)
}

func keyOnlyErr() error {
	return errors.New(kdlKeyOnly)
}

func endOfObjErr() error {
	return errors.New(kdlEndOfObj)
}

func nothingLeftErr() error {
	return errors.New(kdlNothingLeft)
}

func unexpectedEOFErr() error {
	return errors.New(KDLUnexpectedEOF)
}
