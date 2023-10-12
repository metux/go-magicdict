package core

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/metux/go-magicdict/api"
)

// Simple Dictionary, based on map of [github.com/metux/go-magicdict/api.Entry],
// implementing the [github.com/metux/go-magicdict/api.Entry] interface
type Dict struct {
	// keeping as reference instead of internal, so we can easily copy
	// this struct whithout throwing ourselves into a parallel universe ;-)
	data *api.EntryMap
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
		*(d.data) = make(api.EntryMap)
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

	sub := (*d.data)[string(head)]

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
	for _, val := range *d.data {
		vals[idx] = val
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
		if cur == nil {
			if nlist {
				e := EmptyList()
				d.appendK(head, e)
				return e.Put(tail, v)
			} else {
				e := EmptyDict()
				d.appendK(head, e)
				return e.Put(tail, v)
			}
		}
		return cur.Put(tail, v)
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

func (d Dict) append(k string, val api.Entry) {
	(*d.data)[k] = val
}

func (d Dict) appendK(k api.Key, val api.Entry) {
	d.append(string(k), val)
}

func (d *Dict) UnmarshalYAML(node *yaml.Node) error {
	d.initMap()

	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("dict: not a mapping node: tag=%s val=%s at %d:%d", node.Tag, node.Value, node.Line, node.Column)
	}

	idx := 0
	for idx < len(node.Content) {
		sub := node.Content[idx]

		if sub.Tag != "!!str" {
			return fmt.Errorf("dict: unhandled tag: tag=%s val=%s at %d:%d", sub.Tag, sub.Value, sub.Line, sub.Column)
		}
		name := sub.Value

		idx++
		if idx == len(node.Content) {
			return fmt.Errorf("dict: unexpected end")
		}

		switch sub2 := node.Content[idx]; sub2.Kind {
		case yaml.SequenceNode:
			sublist := EmptyList()
			sublist.UnmarshalYAML(sub2)
			d.append(name, sublist)
		case yaml.MappingNode:
			subdict := EmptyDict()
			subdict.UnmarshalYAML(sub2)
			d.append(name, subdict)
		case yaml.ScalarNode:
			if sub2.Tag == "!!null" {
				d.append(name, nil)
			} else {
				d.append(name, NewScalarStr(sub2.Value))
			}
		default:
			return fmt.Errorf("dict: unhandled tag: tag=%s val=%s at %d:%d", sub2.Tag, sub2.Value, sub2.Line, sub2.Column)
		}

		idx++
	}

	return nil
}

func EmptyDict() Dict {
	m := make(api.EntryMap)
	return Dict{data: &m}
}
