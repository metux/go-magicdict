package core

import (
    "strconv"
    "github.com/metux/go-magicdict/api"
)

type Scalar struct {
    data      string
}

func (sc Scalar) Size() int {
    return 0
}

func (sc Scalar) Elems() [] api.Entry {
    return []api.Entry{}
}

func (sc Scalar) Keys() [] string {
    return []string{}
}

// FIXME: correct semantics ?
func (sc Scalar) Get(k api.Key) (api.Entry, error) {
    if k.Empty() {
        return sc, nil
    }

    return nil, api.ErrSubNotSupported
}

func (sc Scalar) String() string {
    return sc.data
}

func (sc Scalar) IsConst() bool {
    return true
}

func (sc Scalar) Put(k api.Key, v api.Entry) error {
    return api.ErrSubNotSupported
}

func (sc Scalar) MayMergeDefaults() bool {
    return false
}

func (sc Scalar) Empty() bool {
    return len(sc.data) == 0
}

func (sc Scalar) IsScalar() bool {
    return true
}

func NewScalarStr(val string) api.Entry {
    return Scalar { data: val }
}

func NewScalarInt(val int) api.Entry {
    return NewScalarStr(strconv.Itoa(val))
}

func NewScalarFloat(val float64) api.Entry {
    return NewScalarStr(strconv.FormatFloat(val, 'g', 5, 64))
}
