package kdlgo

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type KDLErrorType string

const (
	KDLEmptyArray     = "Array is empty"
	KDLDifferentKey   = "All keys of KDLObject to convert to document should be the same"
	KDLInvalidKeyChar = "Invalid character for key"
	KDLInvalidSyntax  = "Invalid syntax"
	KDLInvalidType    = "Invalid KDLType"

	// These should be caught and handled internally
	kdlKeyOnly  = "Key only"
	kdlEndOfObj = "End of KDLObject"
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

func DebugCaller(level int) {
	_, file, no, ok := runtime.Caller(level)
	if ok {
		paths := strings.Split(file, "/")
		if paths[len(paths)-1] != "asm_amd64.s" {
			fmt.Println("(" + file + ":" + strconv.Itoa(no) + ")")
			DebugCaller(level + 2)
		}
	}
}
