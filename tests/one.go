package tests

import (
	"testing"

	"github.com/metux/go-magicdict/api"
	"github.com/metux/go-magicdict/core"
)

func RunTestOne(t *testing.T, r api.Entry) {

	c := Checker{Test: t, Root: r}

	c.AssertKeys("", []string{"schinken", "escapetest", "xxx", "foo", "bar", "yyy", "hello", "ref", "butter", "zzz", "x123"})
	c.AssertKeys("bar", []string{"x", "y", "tree"})
	c.AssertKeys("bar::tree", []string{"leaf", "leaf2", "list", "knollo", "knollo2", "myname"})

	c.AssertEntry("bar::tree::leaf")

	// FIXME: yaml parser ignores empty dicts, so no assert here !
	c.FetchEntry("bar::tree::leaf2")

	c.AssertString("butter::0", "b1")
	c.AssertString("bar::tree::list::0", "a")
	c.AssertString("bar::x", "x-one")
	c.AssertString("bar::y", "y-one")
	c.AssertString("bar::tree::leaf", "1.2")
	c.AssertString("hello", "huhu")
	c.AssertString("ref::x", "x-one")
	c.AssertString("foo::0", "one")
	c.AssertString("bar::tree::knollo", "foobar-one")
	c.AssertString("bar::tree::knollo2", "foobar-one-for-me")

	c.AssertMagicName("bar::tree")

	c.AssertEntry("bar").AssertEntry("tree")

	c.AssertListStr("foo", []string{"one", "two", "three"})
	c.AssertListStr("bar::tree::list", []string{"a", "b"})

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
	c.AssertKeys("butter", []string{"0", "1"})
	c.AssertListStr("butter", []string{"b1", "b3"})

	c.AssertDelete("butter")

	// test putting a function
	c.Root.Put(api.Key("bar::testfunc"), core.NewScalarFunc(func() string { return "HELLO" }))
	c.AssertString("bar::testfunc", "HELLO")

	// test lazy proxy
	c.Root.Put(api.Key("proxytest"), core.NewSymlink(c.Root, api.Key("bar::tree"), false))
	c.AssertListStr("proxytest::list", []string{"a", "b"})

	// checking @@KEY magic key
	c.AssertEntry("zzz")
	c.AssertEntry("bar::tree::@@KEY")
	c.AssertString("zzz", "zzzxx")

	c.AssertString("x123::one::two", "name-is-two")
	c.AssertString("x123::one::three", "parent-is-one")

	// checking defaults set
	api.SetDefaultStr(c.Root, api.Key("bar::tree::leaf23"), "18")
	c.AssertString("bar::tree::leaf23", "18")
	c.AssertEntry("bar::tree").AssertString("leaf23", "18")

	// checking defaults set
	api.SetDefaultStr(c.Root, api.Key("one::two::three::five::six"), "23")
	c.AssertString("one::two::three::five::six", "23")
	c.AssertEntry("one::two::three::five").AssertString("six", "23")

	api.SetStr(c.Root, api.Key("one::two::three::five::six1"), "24")
	c.AssertString("one::two::three::five::six1", "24")
	c.AssertEntry("one::two::three::five").AssertString("six1", "24")

	// check escaping
	c.AssertString("escapetest", "foo-${foo}")

	// check missing vars
	c.AssertString("schinken", "")

	// check set/get
	api.SetStr(c.Root, "knollo::hundert::foo", "hello")
	c.AssertString("knollo::hundert::foo", "hello")

	// test makedict
	newdict, err := api.MakeDict(c.Root, api.Key("foox::bar::xxx"))
	if err != nil {
		t.Fatalf("MakeDict() failed: %s", err)
	}

	api.SetStr(newdict, api.Key("wurst"), "salami")
	c.AssertString("foox::bar::xxx::wurst", "salami")

	RunTestLiteral(c)
	RunTestListCreateAppend(c)
}
