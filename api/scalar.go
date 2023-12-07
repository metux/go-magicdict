package api

import (
	"fmt"
)

type Scalar struct {
	Data string
}

func (sc Scalar) Elems() EntryList {
	return EntryList{}
}

func (sc Scalar) Keys() KeyList {
	return KeyList{}
}

// FIXME: correct semantics ?
func (sc Scalar) Get(k Key) (Entry, error) {
	if k.Empty() {
		return sc, nil
	}

	return nil, fmt.Errorf("Scalar \"%s\" does not support Get(): %w", sc.Data, ErrSubNotSupported)
}

func (sc Scalar) String() string {
	return sc.Data
}

func (sc Scalar) IsConst() bool {
	return true
}

func (sc Scalar) Put(k Key, v Entry) error {
	return fmt.Errorf("Scalar \"%s\" does not support Put(): %w", sc.Data, ErrSubNotSupported)
}

func (sc Scalar) MayMergeDefaults() bool {
	return false
}

func (sc Scalar) Empty() bool {
	return len(sc.Data) == 0
}

func (sc Scalar) IsScalar() bool {
	return true
}

func (sc Scalar) IsList() bool {
	return false
}

func (sc Scalar) IsDict() bool {
	return false
}

// Implements the yaml.Marshaler interface
func (sc Scalar) MarshalYAML() (interface{}, error) {
	return sc.Data, nil
}
