package api

import (
    "strings"
    "strconv"
)

// get a single string value from dict and key
func GetStr(r Entry, k Key) string {
    if r == nil {
        return ""
    }
    if e, err := r.Get(k); e == nil || err != nil {
        return ""
    } else {
        return e.String()
    }
}

// append the string value(s) of entry to given slice
// if entry is list or dict, collect the string values of its members
// (also done recursively)
func StrListAppendEntry(s[] string, e Entry) []string {
    if e == nil {
        return s
    }
    if e.IsScalar() {
        return append(s, e.String())
    }
    for _,v := range e.Elems() {
        s = StrListAppendEntry(s, v)
    }
    return s
}

// fetch entry from dict by key and return aa slice of string values
// if an entry is list or dict, collect the string values of its members
// (also done recursively)
func GetStrList(r Entry, k Key) []string {
    s := []string{}
    if r == nil {
        return s
    }
    v, err := r.Get(k);
    if v == nil || err != nil {
        return s
    }
    return StrListAppendEntry(s, v)
}

func GetBool(r Entry, k Key, dflt bool) bool {
    switch strings.ToLower(GetStr(r, k)) {
        case "1", "y", "yes", "true",  "on":  return true
        case "0", "n", "no",  "false", "off": return false
        default:                              return dflt
    }
}

func GetInt(r Entry, k Key, dflt int) int {
    if i, err := strconv.Atoi(GetStr(r, k)); err == nil {
        return i
    } else {
        return dflt
    }
}

func GetKeys(r Entry, k Key) []string {
    if r != nil {
        if ent,err := r.Get(k); err == nil || ent != nil {
            return ent.Keys()
        }
    }
    return []string{}
}

func GetElems(r Entry, k Key) []Entry {
    if r != nil {
        if ent,err := r.Get(k); err == nil || ent != nil {
            return ent.Elems()
        }
    }
    return []Entry{}
}

func SetStr(r Entry, k Key, val string) error {
    if r == nil {
        return ErrNilInterface
    }
    return r.Put(k, Scalar{Data: val})
}

//
// Append string value to a list entry.
// Automatically creates the list entry if not existing yet
func AppendStr(r Entry, k Key, val string) error {
    if r == nil {
        return ErrNilInterface
    }
    return r.Put(Key(string(k)+"[]::[]"), Scalar{Data: val})
}

func SetInt(r Entry, k Key, val int) error {
    if r == nil {
        return ErrNilInterface
    }
    return r.Put(k, Scalar{Data: strconv.Itoa(val)})
}

func SetBool(r Entry, k Key, val bool) error {
    if r == nil {
        return ErrNilInterface
    }
    return r.Put(k, Scalar{Data: strconv.FormatBool(val)})
}

// Delete an entry with given key within given root entry, by putting nil value
//
// nil-check's the root entry
func Delete(root Entry, k Key) error {
    if root == nil {
        return ErrNilInterface
    }
    return root.Put(k, nil)
}

// Create an [github.com/metux/go-magicdict/core.List] entry inside given root
// entry with given key and return it. If already existing, just return it.
// If already exists as different type, the behavior is unspecified
//
// nil-checks the root entry
func MakeList(root Entry, k Key) (Entry, error) {
    if root == nil {
        return nil, ErrNilInterface
    }
    k2 := Key(string(k)+"[]::[]")
    if err := root.Put(k2, nil); err != nil {
        return nil, err
    }
    return root.Get(k)
}

func MakeDict(root Entry, k Key) (Entry, error) {
    if root == nil {
        return nil, ErrNilInterface
    }
    if err := root.Put(k.Append("@@@"), nil); err != nil {
        return nil, err
    }
    return root.Get(k)
}

