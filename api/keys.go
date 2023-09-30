package api

import (
    "strconv"
    "strings"
)

type Key string

func (k Key) Head() (Key, Key) {
    if i := strings.Index(string(k), "::"); i >= 0 {
        return Key(k[:i]), k[i+2:]
    }
    return k, ""
}

func (k Key) Tail() (Key, Key) {
    if i := strings.LastIndex(string(k), "::"); i >= 0 {
        return k[:i], k[i+2:]
    }
    return "", k
}

func (k Key) Empty() bool {
    return k == ""
}

func (k Key) AddPrefix(prefix string) Key {
    if prefix == "" {
        return k
    }
    return Key(prefix + "::" + string(k))
}

func (k Key) Append(suffix Key) Key {
    if k.Empty() {
        return Key(suffix)
    }

    if suffix.Empty() {
        return k
    }

    return Key(string(k) + "::" + string(suffix))
}

func (k Key) IsAppend() bool {
    return k == "[]"
}

func (k Key) AppendIdx(idx int) Key {
    return k.Append(Key(strconv.Itoa(idx)))
}
