package core

import (
	"fmt"
	"log"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/metux/go-magicdict/api"
)

type dictData struct {
	sync.RWMutex
	M map[api.Key]api.Entry
}

// Simple Dictionary, based on map of [github.com/metux/go-magicdict/api.Entry],
// implementing the [github.com/metux/go-magicdict/api.Entry] interface
type Dict struct {
	// keeping as reference instead of internal, so we can easily copy
	// this struct whithout throwing ourselves into a parallel universe ;-)
	data *dictData
}

func (d Dict) Serialize() (string, error) {
	d.data.RLock()
	text, err := yaml.Marshal(d.data.M)
	d.data.RUnlock()
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func (d Dict) initMap() {
	if d.data.M == nil {
		log.Println("=== init map")
		d.data.M = make(api.EntryMap)
	}
}

// Get an (sub-)entry by key path
func (d Dict) Get(k api.Key) (api.Entry, error) {
	if k.Empty() {
		return d, nil
	}

	if d.data.M == nil {
		return nil, api.ErrDictNotInitialized
	}

	head, tail := k.Head()

	d.data.RLock()
	sub := d.data.M[head]
	d.data.RUnlock()

	if tail.Empty() || sub == nil {
		return sub, nil
	}

	return sub.Get(tail)
}

// Return the direct sub-keys as string slice
func (d Dict) Keys() api.KeyList {
	d.initMap()

	idx := 0
	d.data.RLock()
	keys := make(api.KeyList, len(d.data.M))
	for key := range d.data.M {
		keys[idx] = api.Key(key)
		idx++
	}
	d.data.RUnlock()
	return keys
}

// Return the direct sub-elements as slice of [github.com/metux/go-magicdict/api.Entry]
func (d Dict) Elems() api.EntryList {
	d.initMap()

	idx := 0
	d.data.RLock()
	vals := make(api.EntryList, len(d.data.M))
	for _, val := range d.data.M {
		vals[idx] = val
		idx++
	}
	d.data.RUnlock()
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
	if ok, lk := head.IsListOp(); ok {
		nlist = ok
		head = lk
	}

	d.data.Lock()

	if !tail.Empty() {
		cur := d.data.M[head]
		if cur == nil {
			if nlist {
				e := EmptyList()
				d.data.M[head] = e
				d.data.Unlock()
				return e.Put(tail, v)
			} else {
				e := EmptyDict()
				d.data.M[head] = e
				d.data.Unlock()
				return e.Put(tail, v)
			}
		}
		d.data.Unlock()
		return cur.Put(tail, v)
	}

	// explicit delete
	if v == nil {
		delete(d.data.M, head)
	} else {
		d.data.M[head] = v
	}
	d.data.Unlock()
	return nil
}

// Check whether the dict is empty
func (d Dict) Empty() bool {
	d.data.RLock()
	res := len(d.data.M) == 0
	d.data.RUnlock()
	return res
}

// Does nothing, just return "". Dicts don't have a valid string representation.
func (d Dict) String() string {
	return ""
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

func (d Dict) IsList() bool {
	return false
}

func (d Dict) IsDict() bool {
	return true
}

// Dict objects aren't constant
func (d Dict) IsConst() bool {
	return false
}

func (d *Dict) UnmarshalYAML(node *yaml.Node) error {
	d.initMap()

	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("dict: not a mapping node: tag=%s val=%s at %d:%d", node.Tag, node.Value, node.Line, node.Column)
	}

	d.data.Lock()
	defer d.data.Unlock()

	idx := 0
	for idx < len(node.Content) {
		sub := node.Content[idx]

		if sub.Tag != "!!str" {
			return fmt.Errorf("dict: unhandled tag: tag=%s val=%s at %d:%d", sub.Tag, sub.Value, sub.Line, sub.Column)
		}

		name := api.Key(sub.Value)

		idx++
		if idx == len(node.Content) {
			return fmt.Errorf("dict: unexpected end")
		}

		switch sub2 := node.Content[idx]; sub2.Kind {
		case yaml.SequenceNode:
			sublist := EmptyList()
			sublist.UnmarshalYAML(sub2)
			d.data.M[name] = sublist
		case yaml.MappingNode:
			subdict := EmptyDict()
			subdict.UnmarshalYAML(sub2)
			d.data.M[name] = subdict
		case yaml.ScalarNode:
			if sub2.Tag == "!!null" {
				d.data.M[name] = nil
			} else {
				d.data.M[name] = NewScalarStr(sub2.Value)
			}
		default:
			return fmt.Errorf("dict: unhandled tag: tag=%s val=%s at %d:%d", sub2.Tag, sub2.Value, sub2.Line, sub2.Column)
		}

		idx++
	}

	return nil
}

// Implements the yaml.Marshaler interface
func (d Dict) MarshalYAML() (interface{}, error) {
	return d.data.M, nil
}

func EmptyDict() Dict {
	return Dict{data: &dictData{M: make(map[api.Key]api.Entry)}}
}
