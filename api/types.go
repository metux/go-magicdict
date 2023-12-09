package api

import (
	"sync"
)

type EntryList = []Entry
type KeyList = []Key
type EntryMap = map[Key]Entry

// magic attributes (anything w/ prefix "@@") are handled by the currently
// parsed node itself and have very special meanings. they can be used
// to eg. construct relative references.
//
// NOTE: currently only implemented MagicDict's (and their boxed values)
// core types like Dict or List ignore them (they don't record the data at all)
const (
	// name of the node's key within its parent (string)
	MagicAttrKey      = Key("@@KEY")
	MagicAttrShortKey = Key("@@/")

	// reference to the node's parent entry itself (entry)
	MagicAttrParent       = Key("@@PARENT")
	MagicAttrShortParent  = Key("@@^")
	MagicAttrShortParent2 = Key("@@^2")
	MagicAttrShortParent3 = Key("@@^3")
	MagicAttrShortParent4 = Key("@@^4")
	MagicAttrShortParent5 = Key("@@^5")

	// reference to the parent's id
	MagicAttrParentKey = Key("@@PARENT::@@KEY")
	// name of the node's full path inside MagicDict object (string)
	MagicAttrPath = Key("@@PATH")
	// magic key for appending a list element
	MagicAttrAppend = Key("[]")
	// magic key prefix for accessing defaults
	MagicAttrDefaults = Key("@@DEFAULTS")
	// magic key for disabling variable/macro substitution
	MagicAttrLiteral = Key("@@LITERAL")
)

type EntMap = struct {
	sync.RWMutex
	M map[Key]Entry
}

func NewEntMap() *EntMap {
	return &EntMap{M: make(map[Key]Entry)}
}
