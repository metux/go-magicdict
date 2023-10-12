package core

import (
	"strconv"

	"github.com/metux/go-magicdict/api"
)

type List struct {
	data *api.AnyList
}

func (l List) GetIdx(idx int) (api.Entry, error) {
	v, e, wb := encap((*l.data)[idx], l)
	if wb {
		(*l.data)[idx] = v
	}
	return v, e
}

func (l List) Elems() []api.Entry {
	data := make([]api.Entry, len(*l.data))
	for x := 0; x < len(*l.data); x++ {
		data[x], _ = l.GetIdx(x)
	}
	return data
}

func (l List) Keys() []api.Key {
	data := make([]api.Key, len(*l.data))
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
		*(l.data) = append(*(l.data), v)
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

		dnew := make(api.AnyList, 0, len(*l.data)-1)
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
		newdata := make(api.AnyList, len(*l.data), i)
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

func NewList(val api.AnyList) api.Entry {
	return List{data: &val}
}

func (l List) MayMergeDefaults() bool {
	return false
}

func (l List) IsScalar() bool {
	return false
}

func (l List) IsConst() bool {
	return false
}
