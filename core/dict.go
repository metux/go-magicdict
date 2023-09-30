package core

import (
    "github.com/metux/go-magicdict/api"
    "gopkg.in/yaml.v3"
)

// implemements the Entry interface
type Dict struct {
    data    * api.AnyMap
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

    if (writeback && sub != nil) {
        (*d.data)[string(head)] = sub
    }

    if tail.Empty() || sub == nil {
        return sub, nil
    }

    return sub.Get(tail)
}

func (d Dict) Keys() []string {
    d.initMap()

    idx := 0
    keys := make([]string, len(*d.data))
    for key := range *d.data {
        keys[idx] = key
        idx++
    }
    return keys
}

func (d Dict) Elems() []api.Entry {
    d.initMap()

    idx := 0
    vals := make([]api.Entry, len(*d.data))
    for key,val := range *d.data {
        // FIXME: handle error ?
        v,_,wb := encap(val, d)
        if wb {
            (*d.data)[key] = v
        }
        vals[idx] = v
        idx++
    }
    return vals
}

func (d Dict) Put(k api.Key, v api.Entry) error {
    if k.Empty() {
        return api.ErrKeyEmpty
    }

    d.initMap()

    head, tail := k.Head()

    if !tail.Empty() {
        cur := (*d.data)[string(head)]

        switch curVal := cur.(type) {
            case nil:
                e := make(api.AnyMap)
                (*d.data)[string(head)] = e
                return NewDict(&e).Put(tail, v)
            case api.AnyMap:
                return NewDict(&curVal).Put(tail, v)
            case api.AnyList:
                l := NewList(curVal)
                (*d.data)[string(head)] = l
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

    // check for Scalar special handling
    if v.IsScalar() && v.IsConst() {
        (*d.data)[string(head)] = v.String()
        return nil
    }

    // FIXME: policy for unboxing

    // check for YamlList special handling
//    if ymlList, ok := v.(List); ok {
//        log.Println("putting a yml list")
//        (*d.data)[string(head)] = ymlList.data
//        return nil
//    }

    // check for Dict special handling
//    if ymlDict, ok := v.(Dict); ok {
//        log.Println("putting a yml dict")
//        (*d.data)[string(head)] = ymlDict.data
//        return nil
//    }

    (*d.data)[string(head)] = v
    return nil
}

func (d Dict) Empty() bool {
    return len(*d.data) == 0
}

func (d Dict) String() string {
    return ""
}

func NewDict(val *api.AnyMap) Dict {
    if val == nil {
        m := make(api.AnyMap)
        val = &m
    }
    return Dict { data: val }
}

func (d Dict) MayMergeDefaults() bool {
    return true
}

func (d Dict) IsScalar() bool {
    return false
}

func (d Dict) IsConst() bool {
    return false
}
