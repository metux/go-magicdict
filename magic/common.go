package magic

import (
	"github.com/metux/go-magicdict/api"
)

type Key = api.Key

const (
	KeyParent2 = Key("@@PARENT::@@PARENT")
	KeyParent3 = Key("@@PARENT::@@PARENT::@@PARENT")
	KeyParent4 = Key("@@PARENT::@@PARENT::@@PARENT::@@PARENT")
	KeyParent5 = Key("@@PARENT::@@PARENT::@@PARENT::@@PARENT::@@PARENT")
	KeyParent6 = Key("@@PARENT::@@PARENT::@@PARENT::@@PARENT::@@PARENT::@@PARENT")
)

func xlateKey(k api.Key) api.Key {

	if k.Empty() {
		return k
	}

	head, tail := k.Head()

	switch head {

	case api.MagicAttrShortKey:
		return xlateKey(api.MagicAttrKey.Append(tail))

	case api.MagicAttrShortParent:
		return xlateKey(api.MagicAttrParent.Append(tail))

	case api.MagicAttrShortParent2:
		return xlateKey(KeyParent2.Append(tail))

	case api.MagicAttrShortParent3:
		return xlateKey(KeyParent3.Append(tail))

	case api.MagicAttrShortParent4:
		return xlateKey(KeyParent4.Append(tail))

	case api.MagicAttrShortParent5:
		return xlateKey(KeyParent5.Append(tail))

	}

	return k
}
