package tests

func RunTestLiteral(c Checker) {
	c.Test.Logf("testing MagicAttrLiteral")

	c.AssertString("@@LITERAL::bar::tree::myname", "${@@KEY}")
	c.AssertString("bar::@@LITERAL::tree::myname", "${@@KEY}")
	c.AssertString("bar::tree::@@LITERAL::myname", "${@@KEY}")

	c.AssertString("@@LITERAL::hello", "${xxx}")
	c.AssertString("@@LITERAL::ref", "${bar}")
	c.AssertString("@@LITERAL::zzz", "${@@KEY}xx")

	c.AssertString("@@LITERAL::x123::one::two", "name-is-${@@KEY}")
	c.AssertString("x123::@@LITERAL::one::two", "name-is-${@@KEY}")
	c.AssertString("x123::one::@@LITERAL::two", "name-is-${@@KEY}")

	c.AssertString("@@LITERAL::x123::one::three", "parent-is-${@@PARENT::@@KEY}")
	c.AssertString("x123::@@LITERAL::one::three", "parent-is-${@@PARENT::@@KEY}")
	c.AssertString("x123::one::@@LITERAL::three", "parent-is-${@@PARENT::@@KEY}")

	c.AssertString("@@LITERAL::schinken", "${brot::wurst}")

	c.Test.Logf("testing MagicAttrLiteral done")
}
