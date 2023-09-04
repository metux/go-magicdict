package tests

import (
    "testing"
    "github.com/metux/go-magicdict/api"
    "github.com/metux/go-magicdict/core"
)

func RunTestOne(t * testing.T, r api.Entry) {

    c := Checker { Test: t, Root: r }

    c.AssertKeys("",          []string { "xxx", "foo", "bar", "yyy", "hello", "ref", "butter" })
    c.AssertKeys("bar",       []string { "x", "y", "tree" })
    c.AssertKeys("bar::tree", []string { "leaf", "leaf2", "list", "knollo", "knollo2" })

    c.AssertEntry("bar::tree::leaf")

    // FIXME: yaml parser ignores empty dicts, so no assert here !
    c.FetchEntry("bar::tree::leaf2")

    c.AssertString("butter::0",          "b1")
    c.AssertString("bar::tree::list::0", "a")
    c.AssertString("bar::x",             "x-one")
    c.AssertString("bar::y",             "y-one")
    c.AssertString("bar::tree::leaf",    "1.2")
    c.AssertString("hello",              "huhu")
    c.AssertString("ref::x",             "x-one")
    c.AssertString("foo::0",             "one")
    c.AssertString("bar::tree::knollo",  "foobar-one")
    c.AssertString("bar::tree::knollo2",  "foobar-one-for-me")

    c.AssertEntry("bar::tree")

    c.AssertEntry("bar").AssertEntry("tree")

    c.AssertListStr("foo", []string { "one", "two", "three" })
    c.AssertListStr("bar::tree::list", []string { "a", "b" })

    c.AssertMakeTree("peter::knollo::xy")
    c.AssertMakeTree("abc")
    c.AssertMakeTree("abc::def::123")
    c.AssertMakeTree("knollo::wurst::peter")

    c.AssertPutStr("ncc170d::deck1::door", "one")

    c.AssertDefaultStr("ncc170d::deck2::door", "two")

    c.AssertString("butter::1::0", "abc1")

    // test deleting list element
    c.FetchEntry("butter::1")
    c.DeleteEntry("butter::1")
    c.AssertKeys("butter",       []string { "0", "1" })
    c.AssertListStr("butter", []string { "b1", "b3" })

    c.AssertDelete("butter")

    // test putting a function
    c.Root.Put(api.Key("bar::testfunc"), core.NewScalarFunc(func() string { return "HELLO" }))
    c.AssertString("bar::testfunc", "HELLO")
}
