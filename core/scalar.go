package core

import (
	"strconv"

	"github.com/metux/go-magicdict/api"
)

type Scalar = api.Scalar

func NewScalarStr(val string) api.Entry {
	return Scalar{Data: val}
}

func NewScalarInt(val int) api.Entry {
	return NewScalarStr(strconv.Itoa(val))
}

func NewScalarFloat(val float64) api.Entry {
	return NewScalarStr(strconv.FormatFloat(val, 'g', 5, 64))
}

// Create a new (free standing) `Scalar` entry object with given boolean value,
// returned as `api.Entry` interface
// Internally, converted to string representation
func NewScalarBool(val bool) api.Entry {
	return NewScalarStr(strconv.FormatBool(val))
}
