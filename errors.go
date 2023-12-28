package sfstypes

import (
	"errors"
	"fmt"
)

var (
	ErrKeyEmpty     = errors.New("key is empty")
	ErrDataNull     = errors.New("input data is null")
	ErrDecodingBool = errors.New("error decoding bool")
)

type ErrDecodingUnsupportedType struct {
	sfsType sfsDataType
	Len     int
	Cap     int
}

func (err *ErrDecodingUnsupportedType) Error() string {
	return fmt.Sprintf("can't decode type %s: len: %d cap: %d", sfsTypeToString(err.sfsType), err.Len, err.Cap)
}

type ErrInsufficientByteData struct {
	sfsType sfsDataType
	size    int
}

func (err *ErrInsufficientByteData) Error() string {
	return fmt.Sprintf("can't decode an %s, byte data is insufficient: size: %d", sfsTypeToString(err.sfsType), err.size)
}

type ErrUnsupportedType struct {
	value interface{}
}

func (err *ErrUnsupportedType) Error() string {
	return fmt.Sprintf("value %v, is of unsupported type %t", err.value, err.value)
}

type ErrWrongType struct {
	actualType sfsDataType
	wantedType sfsDataType
}

func (err *ErrWrongType) Error() string {
	return fmt.Sprintf("found %s but expected type %s", sfsTypeToString(err.actualType), sfsTypeToString(err.wantedType))
}

type ErrKeyNotFound struct {
	key string
}

func (err *ErrKeyNotFound) Error() string {
	return fmt.Sprintf("key \"%s\" not found", err.key)
}

type ErrInvalidKeySize struct {
	key    string
	length int
}

func (err *ErrInvalidKeySize) Error() string {
	return fmt.Sprintf("invalid length of key \"%s\" (%d) (key must be >0 and <256)", err.key, err.length)
}

type ErrIndexNotInRange struct {
	index int
}

func (err *ErrIndexNotInRange) Error() string {
	return fmt.Sprintf("index %d not in range", err.index)
}

type ErrReadingData struct {
	TypeToRead string
	Len        int
	Cap        int
	IoErr      error
}

func (err *ErrReadingData) Error() string {
	return fmt.Sprintf("error while reading %s: len: %d cap: %d, ioError: %s", err.TypeToRead, err.Len, err.Cap, err.IoErr)
}
