package core

import (
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/metux/go-magicdict/api"
)

type List struct {
	// keeping as reference instead of internal, so we can easily copy
	// this struct whithout throwing ourselves into a parallel universe ;-)
	data *api.EntryList
}

func (l List) GetIdx(idx int) (api.Entry, error) {
	if idx < len(*l.data) {
		return (*l.data)[idx], nil
	}
	return nil, nil
}

func (l List) Elems() api.EntryList {
	data := make(api.EntryList, len(*l.data))
	for x := 0; x < len(*l.data); x++ {
		data[x], _ = l.GetIdx(x)
	}
	return data
}

func (l List) Keys() api.KeyList {
	data := make(api.KeyList, len(*l.data))
	for x := 0; x < len(*l.data); x++ {
		data[x] = api.Key(strconv.Itoa(x))
	}
	return data
}

func (l List) Get(k api.Key) (api.Entry, error) {
	if k.Empty() {
		return l, nil
	}
	i, err := strconv.Atoi(string(k))
	if err != nil {
		return nil, err
	}

	return l.GetIdx(i)
}

func (l List) Put(k api.Key, v api.Entry) error {
	// append
	if k.IsAppend() {
		l.append(v)
		return nil
	}

	i, err := strconv.Atoi(string(k))
	if err != nil {
		return err
	}

	// delete
	if v == nil {
		if i >= len(*l.data) {
			return api.ErrIndexOutOfRange
		}

		dnew := make(api.EntryList, 0, len(*l.data)-1)
		for x, y := range *l.data {
			if x != i {
				dnew = append(dnew, y)
			}
		}
		*l.data = dnew
		return nil
	}

	// simple update
	if i < len(*l.data) {
		(*l.data)[i] = v
	} else {
		newdata := make(api.EntryList, len(*l.data), i)
		for idx, v := range *l.data {
			newdata[idx] = v
		}
		newdata[i] = v
		l.data = &newdata
	}

	return api.ErrSubNotSupported
}

func (l List) Empty() bool {
	return len(*l.data) == 0
}

func (l List) String() string {
	return ""
}

func (l List) MayMergeDefaults() bool {
	return false
}

func (l List) IsScalar() bool {
	return false
}

func (l List) IsList() bool {
	return true
}

func (l List) IsDict() bool {
	return false
}

func (l List) IsConst() bool {
	return false
}

func (l *List) append(val api.Entry) {
	*l.data = append(*l.data, val)
}

func (l *List) UnmarshalYAML(node *yaml.Node) error {
	l.data = new(api.EntryList)

	if node.Kind != yaml.SequenceNode {
		return fmt.Errorf("list: not a sequence node: tag=%s val=%s at %d:%d", node.Tag, node.Value, node.Line, node.Column)
	}

	for _, sub := range node.Content {
		switch sub.Kind {
		case yaml.SequenceNode:
			sublist := EmptyList()
			if err := sublist.UnmarshalYAML(sub); err != nil {
				return err
			}
			l.append(sublist)
		case yaml.MappingNode:
			subdict := EmptyDict()
			if err := subdict.UnmarshalYAML(sub); err != nil {
				return err
			}
			l.append(subdict)
		case yaml.ScalarNode:
			l.append(NewScalarStr(sub.Value))
		default:
			return fmt.Errorf("list: unhandled tag: tag=%s val=%s at %d:%d", sub.Tag, sub.Value, sub.Line, sub.Column)
		}
	}

	return nil
}

// Implements the yaml.Marshaler interface
func (l List) MarshalYAML() (interface{}, error) {
	return l.data, nil
}

func EmptyList() List {
	l := make(api.EntryList, 0)
	return List{data: &l}
}
