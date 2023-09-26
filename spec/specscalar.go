package spec

import (
    "github.com/metux/go-magicdict/api"
)

type SpecScalar struct {
    Root * SpecObject
    Path api.Key
    Data string
}

func (this SpecScalar) Get(k api.Key) (api.Entry, error) {
    return nil, api.ErrSubNotSupported
}

func (this SpecScalar) IsScalar() bool {
    return true
}

func (this SpecScalar) IsConst() bool {
    return true
}

func (this SpecScalar) String() string {
    return this.Data
}

func (this SpecScalar) Elems() [] api.Entry {
    return []api.Entry{}
}

func (this SpecScalar) Keys() [] string {
    return []string{}
}

func (this SpecScalar) Put(k api.Key, v api.Entry) error {
    return api.ErrSubNotSupported
}

func (this SpecScalar) Empty() bool {
    return true
}

// maybe we'll wanna have different modes here
func (this SpecScalar) MayMergeDefaults() bool {
    return false
}
