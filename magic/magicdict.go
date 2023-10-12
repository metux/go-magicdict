package magic

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
type MagicDict struct {
	Root *MagicDict
	Data api.Entry

	// this is always from the root (need to prepend Path)
	Defaults api.Entry

	Path api.Key
}

// Box an entry into MagicDict wrapper, so that defaults still work
// trivial scalars and nil aren't boxed
//
// FIXME: right now just boxing Dict and List
// we'll later box everything that's not inside a Boxed
//
// FIXME: add variable substitution
func (this MagicDict) box(k api.Key, v api.Entry) (api.Entry, error) {

	r := this.Root
	if r == nil {
		r = &this
	}

	switch v.(type) {
	case core.Dict, core.List:
		// FIXME: decide when to box ?
		// const strings or those w/ string interface should not be boxed ?
		sp := MagicDict{
			Root:     r,
			Path:     this.Path.Append(k),
			Data:     v,
			Defaults: this.Defaults,
		}
		sp.Init()
		return sp, nil

	case core.Scalar:
		v = magicScalar{
			Root: r,
			Path: this.Path.Append(k),
			Data: v.String(),
		}
	}

	return macro.ProcessVars(v, r)
}

func (this MagicDict) Get(k api.Key) (api.Entry, error) {

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

	switch head {
	case api.MagicAttrPath:
		return core.NewScalarStr(string(this.Path)), nil

	case api.MagicAttrKey:
		_, p1 := this.Path.Tail()
		return core.NewScalarStr(string(p1)), nil

	case api.MagicAttrParent:
		p1, _ := this.Path.Tail()
		if this.Root == nil {
			return nil, nil
		}
		return this.Root.Get(p1)

	case api.MagicAttrDefaults:
		return this.Defaults, nil
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

func (this MagicDict) IsScalar() bool {
	return this.Data.IsScalar()
}

func (this MagicDict) IsConst() bool {
	return this.Data.IsConst()
}

func (this MagicDict) String() string {
	return this.Data.String()
}

// NOTE: this needs pointer receiver
func (this *MagicDict) SetDefaultEntry(k api.Key, val api.Entry) error {
	// FIXME: maybe nil should mean clear all ?
	if val == nil {
		return api.ErrNilInterface
	}
	return this.Defaults.Put(k, val)
}

// FIXME: add AddDefaults() from list of key+value

// this prevents unwanted merging (eg. on lists)
func (this MagicDict) mergeDef() api.Entry {
	if this.MayMergeDefaults() {
		if d, _ := this.Defaults.Get(this.Path); d != nil && !d.Empty() {
			return d
		}
	}
	return nil
}

func (this MagicDict) Elems() []api.Entry {
	elems := make([]api.Entry, 0)
	for _, k := range this.Keys() {
		e, _ := this.Get(k)
		elems = append(elems, e)
	}
	return elems
}

func (this MagicDict) Keys() []api.Key {
	if d := this.mergeDef(); d != nil {
		return utils.UnionSlice(this.Data.Keys(), d.Keys())
	}
	return this.Data.Keys()
}

func (this MagicDict) Put(k api.Key, v api.Entry) error {
	if head, tail := k.Head(); head == api.MagicAttrDefaults {
		return this.Defaults.Put(tail, v)
	}
	return this.Data.Put(k, v)
}

func (this MagicDict) Empty() bool {
	return this.Data.Empty() && this.Defaults.Empty()
}

// maybe we'll wanna have different modes here
func (this MagicDict) MayMergeDefaults() bool {
	return this.Data.MayMergeDefaults()
}

func (this *MagicDict) Init() {
	if this.Data == nil {
		this.Data = core.EmptyDict()
	}
	if this.Defaults == nil {
		this.Defaults = core.EmptyDict()
	}
}

// only create it via constructor, since some fields *MUST* be initialized
func NewMagicFromDict(d api.Entry, dflt api.Entry) *MagicDict {
	sp := MagicDict{
		Data:     d,
		Defaults: dflt,
	}
	sp.Init()
	return &sp
}
