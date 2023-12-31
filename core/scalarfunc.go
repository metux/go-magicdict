package core

import (
	"github.com/metux/go-magicdict/api"
)

type ScalarFunc struct {
	data func() string
}

func (this ScalarFunc) Elems() api.EntryList {
	return api.EntryList{}
}

func (this ScalarFunc) Keys() api.KeyList {
	return api.KeyList{}
}

// FIXME: correct semantics ?
func (this ScalarFunc) Get(k api.Key) (api.Entry, error) {
	if k.Empty() {
		return this, nil
	}

	return nil, api.ErrSubNotSupported
}

func (this ScalarFunc) String() string {
	return this.data()
}

func (this ScalarFunc) IsConst() bool {
	return false
}

func (this ScalarFunc) Put(k api.Key, v api.Entry) error {
	return api.ErrSubNotSupported
}

func (this ScalarFunc) MayMergeDefaults() bool {
	return false
}

func (sc ScalarFunc) Empty() bool {
	return false
}

func (sc ScalarFunc) IsScalar() bool {
	return true
}

func NewScalarFunc(f func() string) api.Entry {
	return ScalarFunc{data: f}
}
