package core

import (
	"github.com/metux/go-magicdict/api"
)

func EmptyListOrDict(listordict bool) api.Entry {
	if listordict {
		return EmptyList()
	} else {
		return EmptyDict()
	}
}
