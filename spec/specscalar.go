package spec

import (
    "github.com/metux/go-magicdict/api"
    "github.com/metux/go-magicdict/core"
)

type SpecScalar struct {
    Root * SpecObject
    Path api.Key
    Data string
}

func (this SpecScalar) Get(k api.Key) (api.Entry, error) {
    head, tail := k.Head()
    switch head {
        case api.MagicAttrPath:
            return core.NewScalarStr(string(this.Path)), nil

        case api.MagicAttrKey:
            _,p1 := this.Path.Tail()
            return core.NewScalarStr(string(p1)), nil

        case api.MagicAttrParent:
            p1,_ := this.Path.Tail()
            if this.Root == nil {
                return nil, nil
            }
            v2,e2 := this.Root.Get(p1)
            if tail.Empty() {
                return v2, e2
            } else {
                return v2.Get(tail)
            }
    }
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
