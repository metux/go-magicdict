package tests

import (
	"github.com/metux/go-magicdict/api"
)

func RunTestListCreateAppend(c Checker) {
	c.Test.Logf("testing list create/append")

	api.SetStr(c.Root, "foolist::l1[]::[]", "one")
	c.AssertListStr("foolist::l1", []string{"one"})
	api.SetStr(c.Root, "foolist::l1[]::[]", "two")
	c.AssertListStr("foolist::l1", []string{"one", "two"})

	c.Test.Logf("testing list create/append: done")
}
