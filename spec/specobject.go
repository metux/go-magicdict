package spec

import (
    "github.com/metux/go-magicdict/api"
    "github.com/metux/go-magicdict/core"
    "github.com/metux/go-magicdict/macro"
    "github.com/metux/go-magicdict/utils"
)

// This serves as root as well as intermediate nodes
// the root node has Root==nil and Path=="" and defaults!=nil / while the
// others have defaults==nil
// Splittng it up would just save one PTR per instance, but add extra
// complexity and allocations
type SpecObject struct {
    Root * SpecObject
    Data api.Entry

    // this is always from the root (need to prepend Path)
    Defaults api.Entry

    Path api.Key
}

// Box an entry into SpecObject wrapper, so that defaults still work
// trivial scalars and nil aren't boxed
//
// FIXME: right now just boxing Dict and List
// we'll later box everything that's not inside a Boxed
//
// FIXME: add variable substitution
func (this SpecObject) box(k api.Key, v api.Entry) (api.Entry, error) {

    r := this.Root
    if r == nil {
        r = &this
    }

    switch v.(type) {
        case core.Dict, core.List:
            // FIXME: decide when to box ?
            // const strings or those w/ string interface should not be boxed ?
            sp := SpecObject {
                Root:     r,
                Path:     this.Path.Append(k),
                Data:     v,
                Defaults: this.Defaults,
            }
            return sp.Init(), nil
        default:
            return macro.ProcessVars(v, r)
    }
}

func (this SpecObject) Get(k api.Key) (api.Entry, error) {

    if k.Empty() {
        return this, nil
    }

    head, tail := k.Head()

    if !tail.Empty() {
        parent, err := this.Get(head)
        if err != nil {
            return nil, err
        }
        if parent == nil {
            return nil, nil
        }
        return parent.Get(tail)
    }

    ent, err := this.Data.Get(k)
    if err != nil {
        return nil, err
    }
    if ent != nil {
        return this.box(k, ent)
    }

    ent, err = this.Defaults.Get(this.Path.Append(k))
    if err != nil {
        return nil, err
    }
    if ent != nil {
        return this.box(k, ent)
    }

    return nil, nil
}

func (this SpecObject) IsScalar() bool {
    return false
}

func (this SpecObject) IsConst() bool {
    return false
}

func (this SpecObject) String() string {
    return ""
}

// NOTE: this needs pointer receiver
func (this * SpecObject) SetDefaultEntry(k api.Key, val api.Entry) error {
    // FIXME: maybe nil should mean clear all ?
    if val == nil {
        return api.ErrNilInterface
    }
    return this.Defaults.Put(k, val)
}

// FIXME: add AddDefaults() from list of key+value

// this prevents unwanted merging (eg. on lists)
func (this SpecObject) mergeDef() api.Entry {
    if this.MayMergeDefaults() {
        if d,_ := this.Defaults.Get(this.Path); d != nil && !d.Empty() {
            return d
        }
    }
    return nil
}

func (this SpecObject) Elems() [] api.Entry {
    if d := this.mergeDef(); d != nil {
        elems := make([]api.Entry, 0)
        for _,k := range this.Keys() {
            e,_ := this.Get(api.Key(k))
            elems = append(elems, e)
        }
        return elems
    }
    return this.Data.Elems()
}

func (this SpecObject) Keys() [] string {
    if d := this.mergeDef(); d != nil {
        return utils.UnionSlice(this.Data.Keys(), d.Keys())
    }
    return this.Data.Keys()
}

func (this SpecObject) Put(k api.Key, v api.Entry) error {
    return this.Data.Put(k, v)
}

func (this SpecObject) Empty() bool {
    return this.Data.Empty() && this.Defaults.Empty()
}

// maybe we'll wanna have different modes here
func (this SpecObject) MayMergeDefaults() bool {
    return this.Data.MayMergeDefaults()
}

func (this * SpecObject) Init() * SpecObject {
    if this.Data == nil {
        this.Data = core.NewDict(nil)
    }
    if this.Defaults == nil {
        this.Defaults = core.NewDict(nil)
    }
    return this
}

func (this * SpecObject) InitData(data api.Entry, defaults api.Entry) * SpecObject {
    this.Data = data
    this.Defaults = defaults
    return this.Init()
}

//
// only create it via constructor, since some fields *MUST* be initialized
//
func NewSpecFromDict(d api.Entry, dflt api.Entry) * SpecObject {
    if dflt == nil {
        dflt = core.NewDict(nil)
    }

    sp := SpecObject {
        Data: d,
        Defaults: dflt,
    }

    return sp.Init()
}
