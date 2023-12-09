package magic

import (
	"github.com/metux/go-magicdict/api"
	"github.com/metux/go-magicdict/core"
)

type magicScalar struct {
	Root *MagicDict
	Path api.Key
	Data string
}

func (this magicScalar) Get(k api.Key) (api.Entry, error) {
	k = xlateKey(k)

	head, tail := k.Head()
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
		v2, e2 := this.Root.Get(p1)
		if tail.Empty() {
			return v2, e2
		} else {
			return v2.Get(tail)
		}
	}
	return nil, api.ErrSubNotSupported
}

func (this magicScalar) IsScalar() bool {
	return true
}

func (this magicScalar) IsList() bool {
	return false
}

func (this magicScalar) IsDict() bool {
	return false
}

func (this magicScalar) IsConst() bool {
	return true
}

func (this magicScalar) String() string {
	return this.Data
}

func (this magicScalar) Elems() []api.Entry {
	return []api.Entry{}
}

func (this magicScalar) Keys() []api.Key {
	return []api.Key{}
}

func (this magicScalar) Put(k api.Key, v api.Entry) error {
	return api.ErrSubNotSupported
}

func (this magicScalar) Empty() bool {
	return true
}

// maybe we'll wanna have different modes here
func (this magicScalar) MayMergeDefaults() bool {
	return false
}
