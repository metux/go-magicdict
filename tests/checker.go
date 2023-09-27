package tests

import (
    "runtime/debug"
    "strconv"
    "testing"
    "github.com/metux/go-magicdict/api"
    "github.com/metux/go-magicdict/core"
)

type Checker struct {
    Test * testing.T
    Root api.Entry
}

func (c Checker) AssertKeys(k api.Key, want[] string) {
    keys := fetchEntry(c.Test, c.Root, k).Keys()
    if len(keys) != len(want) {
        c.Test.Fatalf("keys size mismatch: %d should be %d (\"%s\" vs \"%s\")", len(keys), len(want), keys, want)
    }
    checkStrs(c.Test, keys, want)
}

// note that this only works if the upper layer doesn't have this entry
func (c Checker) AssertDefaultStr(k api.Key, v string) {
    putDefaultStr(c.Test, c.Root, k, v)
    c.AssertString(k, v)
}

func (c Checker) AssertMakeTree(k api.Key) {
    c.Root.Put(k.Append("@"), nil)
    c.AssertEntry(k)
}

func (c Checker) AssertListStr(k api.Key, want[] string) {
    elems := fetchEntry(c.Test, c.Root, k).Elems()
    el2 := make([]string, len(elems))

    for x,e := range elems {
        if e.IsScalar() && e.IsConst() {
            el2[x] = e.String()
            n_ent, n_err := e.Get(api.MagicAttrKey)
            if n_err != nil || n_ent.String() != strconv.Itoa(x) {
                c.Test.Fatalf("name of list entry #%d of %s broken: \"%s\" should be %d -- %s\n%s\n", x, k, n_ent.String(), x, n_err, string(debug.Stack()))
            }
        } else {
            c.Test.Fatalf("list entry #%d of %s is not scalar", x, k)
        }
    }

    checkStrs(c.Test, el2, want)
}

func (c Checker) AssertPutStr(k api.Key, v string) {
    if err := c.Root.Put(k, core.NewScalarStr(v)); err != nil {
        c.Test.Fatalf("error putting string: %s", err)
    }
    c.AssertString(k, v)
}

func (c Checker) AssertMissing(k api.Key) {
    if fetchEntry(c.Test, c.Root, k) != nil {
        c.Test.Fatalf("entry unexpected: %s", k)
    }
}

func (c Checker) AssertString(k api.Key, want string) {
    if str := c.FetchScalar(k).String(); str != want {
        c.Test.Fatalf("\"%s\" should be \"%s\" but is \"%s\"", k, want, str)
    }
    c.Test.Logf("asserted key %s is string value %s\n", k, want)
}

func (c Checker) AssertMagicName(k api.Key) {
    ent := c.AssertEntry(k)
    _,tail := k.Tail()
    magicent,_ := ent.Root.Get(api.MagicAttrKey)
    if magicent.String() != string(tail) {
        c.Test.Fatalf("wrong entry magic name of %s: %s", k, magicent.String())
    }
}

func (c Checker) FetchScalar(k api.Key) api.Entry {
    ent := c.AssertEntry(k)
    if ent.Root.IsScalar() {
        return ent.Root
    }
    c.Test.Fatalf("%s not a scalar", k)
    return nil
}

func (c Checker) AssertEntry(k api.Key) Checker {
    e := fetchEntry(c.Test, c.Root, k)
    if e == nil {
        c.Test.Fatalf("entry missing: %s", k)
    }
    return Checker { Test: c.Test, Root : e }
}

// explicitly skip nil's
func (c Checker) FetchEntry(k api.Key) api.Entry {
    v, err := c.Root.Get(k)
    if err != nil {
        c.Test.Fatalf("entry %s is not Entry: %s", k, err)
        return nil
    }
    return v
}

func (c Checker) DeleteEntry(k api.Key) {
    if err := c.Root.Put(k, nil); err != nil {
        c.Test.Fatalf("cant delete entry %s: %s", k, err)
    }
}

func (c Checker) AssertDelete(k api.Key) {
    c.DeleteEntry(k)
    c.AssertMissing(k)
}
