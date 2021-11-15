package kdlgo

import (
	"math/big"
	"strconv"
	"strings"
)

type KDLType string

const (
	KDLBoolType      = "kdl_bool"
	KDLNumberType    = "kdl_number"
	KDLStringType    = "kdl_string"
	KDLRawStringType = "kdl_raw_string"
	KDLDocumentType  = "kdl_document"
	KDLNullType      = "kdl_null"
	KDLDefaultType   = "kdl_default"
	KDLObjectsType   = "kdl_objects"
)

type KDLNode struct {
	node         string
	declaredType string
}

type KDLValue struct {
	Bool      bool
	Number    big.Float
	String    string
	RawString string
	Document  []KDLValue
	Objects   []KDLObject

	Type         KDLType
	declaredType string
}

func (kdlValue KDLValue) RecreateKDL() (string, error) {
	switch kdlValue.Type {
	case KDLBoolType:
		return strconv.FormatBool(kdlValue.Bool), nil
	case KDLNumberType:
		num := kdlValue.Number
		f64, _ := num.Float64()
		return strconv.FormatFloat(f64, 'f', -1, 64), nil
	case KDLStringType:
		return RecreateString(kdlValue.String), nil
	case KDLRawStringType:
		return RecreateString(kdlValue.RawString), nil
	case KDLDocumentType:
		var s strings.Builder
		for i, v := range kdlValue.Document {
			str, err := v.RecreateKDL()
			if err != nil {
				return "", err
			}
			s.WriteString(str)
			if i+1 != len(kdlValue.Document) {
				s.WriteRune(' ')
			}
		}
		return s.String(), nil
	case KDLNullType:
		return "null", nil
	case KDLDefaultType:
		return "", nil
	case KDLObjectsType:
		if len(kdlValue.Objects) < 1 {
			return "", nil
		}
		var s strings.Builder
		for _, obj := range kdlValue.Objects {
			objStr, err := RecreateKDLObj(obj)
			if err != nil {
				return "", err
			}
			s.WriteString(objStr + "; ")
		}
		return "{ " + s.String() + "}", nil
	default:
		return "", invalidTypeErr()
	}
}

func RecreateString(s string) string {
	return strings.ReplaceAll(strconv.Quote(s), "/", "\\/")
}

func (kdlValue KDLValue) ToString() (string, error) {
	switch kdlValue.Type {
	case KDLStringType:
		return kdlValue.String, nil
	case KDLRawStringType:
		return kdlValue.RawString, nil
	case KDLBoolType, KDLNumberType, KDLDocumentType, KDLNullType, KDLDefaultType, KDLObjectsType:
		fallthrough
	default:
		return kdlValue.RecreateKDL()
	}
}

type KDLObject interface {
	GetKey() string
	GetValue() KDLValue
}

func RecreateKDLObj(kdlObj KDLObject) (string, error) {
	s, err := kdlObj.GetValue().RecreateKDL()
	if err != nil {
		return "", nil
	}
	if len(s) > 0 {
		s = " " + s
	}
	key := kdlObj.GetKey()
	if strings.Contains(key, " ") || strconv.Quote(key) != "\""+key+"\"" {
		key = strconv.Quote(key)
	}
	return key + s, nil
}

type KDLBool struct {
	key   string
	value KDLValue
}

func NewKDLBool(key string, value bool) KDLBool {
	return KDLBool{key: key, value: KDLValue{Bool: value, Type: KDLBoolType}}
}

func (kdlNode KDLBool) GetKey() string {
	return kdlNode.key
}

func (kdlNode KDLBool) GetValue() KDLValue {
	return kdlNode.value
}

type KDLNumber struct {
	key   string
	value KDLValue
}

func NewKDLNumber(key string, value float64) KDLNumber {
	return KDLNumber{key: key, value: KDLValue{Number: *big.NewFloat(value), Type: KDLNumberType}}
}

func (kdlNode KDLNumber) GetKey() string {
	return kdlNode.key
}

func (kdlNode KDLNumber) GetValue() KDLValue {
	return kdlNode.value
}

type KDLString struct {
	key   string
	value KDLValue
}

func NewKDLString(key string, value string) KDLString {
	value = strings.ReplaceAll(value, "\n", "\\n")
	s, _ := strconv.Unquote(`"` + value + `"`)
	return KDLString{key: key, value: KDLValue{String: s, Type: KDLStringType}}
}

func (kdlNode KDLString) GetKey() string {
	return kdlNode.key
}

func (kdlNode KDLString) GetValue() KDLValue {
	return kdlNode.value
}

type KDLRawString struct {
	key   string
	value KDLValue
}

func NewKDLRawString(key string, value string) KDLRawString {
	return KDLRawString{key: key, value: KDLValue{RawString: value, Type: KDLRawStringType}}
}

func (kdlNode KDLRawString) GetKey() string {
	return kdlNode.key
}

func (kdlNode KDLRawString) GetValue() KDLValue {
	return kdlNode.value
}

type KDLDocument struct {
	key   string
	value KDLValue
}

func NewKDLDocument(key string, value []KDLValue) KDLDocument {
	return KDLDocument{key: key, value: KDLValue{Document: value, Type: KDLDocumentType}}
}

func (kdlNode KDLDocument) GetKey() string {
	return kdlNode.key
}

func (kdlNode KDLDocument) GetValue() KDLValue {
	return kdlNode.value
}

type KDLNull struct {
	key   string
	value KDLValue
}

func NewKDLNull(key string) KDLNull {
	return KDLNull{key: key, value: KDLValue{Type: KDLNullType}}
}

func (kdlNode KDLNull) GetKey() string {
	return kdlNode.key
}

func (kdlNode KDLNull) GetValue() KDLValue {
	return kdlNode.value
}

type KDLDefault struct {
	key   string
	value KDLValue
}

func NewKDLDefault(key string) KDLDefault {
	return KDLDefault{key: key, value: KDLValue{Type: KDLDefaultType}}
}

func (kdlNode KDLDefault) GetKey() string {
	return kdlNode.key
}

func (kdlNode KDLDefault) GetValue() KDLValue {
	return kdlNode.value
}

type KDLObjects struct {
	key   string
	value KDLValue
}

func NewKDLObjects(key string, objects []KDLObject) KDLObjects {
	return KDLObjects{key: key, value: KDLValue{Objects: objects, Type: KDLObjectsType}}
}

func (kdlNode KDLObjects) GetKey() string {
	return kdlNode.key
}

func (kdlNode KDLObjects) GetValue() KDLValue {
	return kdlNode.value
}

func (kdlObjs KDLObjects) ToObjMap() KDLObjectsMap {
	ret := make(KDLObjectsMap)
	for _, obj := range kdlObjs.GetValue().Objects {
		ret[obj.GetKey()] = obj
	}
	return ret
}

func (kdlObjs KDLObjects) ToValueMap() KDLValuesMap {
	ret := make(KDLValuesMap)
	for _, obj := range kdlObjs.GetValue().Objects {
		ret[obj.GetKey()] = obj.GetValue()
	}
	return ret
}

type KDLObjectsMap map[string]KDLObject
type KDLValuesMap map[string]KDLValue
