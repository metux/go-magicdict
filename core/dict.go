package core

import (
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/metux/go-magicdict/api"
)

// Simple Dictionary, based on [github.com/metux/go-magicdict/api.AnyMap],
// implementing the [github.com/metux/go-magicdict/api.Entry] interface
type Dict struct {
	// keeping as reference instead of internal, so we can easily copy
	// this struct whithout throwing ourselves into a parallel universe ;-)
	data *api.AnyMap
}

func (d Dict) Serialize() (string, error) {
	text, err := yaml.Marshal(d.data)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func (d Dict) initMap() {
	if *(d.data) == nil {
		*(d.data) = make(api.AnyMap)
	}
}

// Get an (sub-)entry by key path
func (d Dict) Get(k api.Key) (api.Entry, error) {
	if k.Empty() {
		return d, nil
	}

	if d.data == nil {
		return nil, api.ErrDictNotInitialized
	}

	head, tail := k.Head()

	sub, err, writeback := encap((*d.data)[string(head)], d)
	if err != nil {
		return sub, nil
	}

	if writeback && sub != nil {
		(*d.data)[string(head)] = sub
	}

	if tail.Empty() || sub == nil {
		return sub, nil
	}

	return sub.Get(tail)
}

// Return the direct sub-keys as string slice
func (d Dict) Keys() api.KeyList {
	d.initMap()

	idx := 0
	keys := make(api.KeyList, len(*d.data))
	for key := range *d.data {
		keys[idx] = api.Key(key)
		idx++
	}
	return keys
}

// Return the direct sub-elements as slice of [github.com/metux/go-magicdict/api.Entry]
func (d Dict) Elems() api.EntryList {
	d.initMap()

	idx := 0
	vals := make(api.EntryList, len(*d.data))
	for key, val := range *d.data {
		// FIXME: handle error ?
		v, _, wb := encap(val, d)
		if wb {
			d.append(key, v)
		}
		vals[idx] = v
		idx++
	}
	return vals
}

// Put in an (sub-)entry by given key path. If the path has more than one element,
// automatically diving into (and possibly creating) sub entries.
//
// When auto-creating and some key element has an "[]" suffix, a
// [github.com/metux/go-magicdict/core.List] is created instead of
// [github.com/metux/go-magicdict/core.Dict]
//
// Put()'ing a nil value causes that entry to be deleted from the dict.
//
// Hint: if the entry is scalar and constant, directly storing it's string
// representation instead of the entry itself. But this behaviour might
// change in future.
func (d Dict) Put(k api.Key, v api.Entry) error {
	if k.Empty() {
		return api.ErrKeyEmpty
	}

	d.initMap()

	head, tail := k.Head()
	nlist := false
	if strings.HasSuffix(string(head), "[]") {
		nlist = true
		head = head[:len(head)-2]
	}

	if !tail.Empty() {
		cur := (*d.data)[string(head)]

		switch curVal := cur.(type) {
		case nil:
			if nlist {
				e := EmptyList()
				d.appendK(head, e)
				return e.Put(tail, v)
			} else {
				e := EmptyDict()
				d.appendK(head, e)
				return e.Put(tail, v)
			}
		case api.AnyMap:
			return NewDict(&curVal).Put(tail, v)
		case api.AnyList:
			l := NewList(curVal)
			d.appendK(head, l)
			return l.Put(tail, v)
		case string, int, float64:
			return api.ErrSubNotSupported
		case api.Entry:
			return curVal.Put(tail, v)
		default:
			return api.ErrUnknownEntryType
		}
	}

	// explicit delete
	if v == nil {
		delete(*d.data, string(head))
		return nil
	}

	d.appendK(head, v)
	return nil
}

// Check whether the dict is empty
func (d Dict) Empty() bool {
	return len(*d.data) == 0
}

// Does nothing, just return "". Dicts don't have a valid string representation.
func (d Dict) String() string {
	return ""
}

// Create a new dict from existing [github.com/metux/go-magicdict/api.AnyMap].
// Using a *pointer* to the AnyMap, instead of copy, thus any changes in the
// dict will be reflected in the passed AnyMap.
func NewDict(val *api.AnyMap) Dict {
	if val == nil {
		m := make(api.AnyMap)
		val = &m
	}
	return Dict{data: val}
}

// Tell [github.com/metux/go-magicdict/magic.MagicDict] that it's allowed to
// merge our keys with those of the lower default dict layer
func (d Dict) MayMergeDefaults() bool {
	return true
}

// Dict objects aren't scalar at all
func (d Dict) IsScalar() bool {
	return false
}

// Dict objects aren't constant
func (d Dict) IsConst() bool {
	return false
}

func (d Dict) append(k string, val api.Any) {
	(*d.data)[k] = val
}

func (d Dict) appendK(k api.Key, val api.Any) {
	d.append(string(k), val)
}

func EmptyDict() Dict {
	m := make(api.AnyMap)
	return Dict{data: &m}
}
