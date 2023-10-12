package api

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
	MagicAttrKey = Key("@@KEY")
	// reference to the node's parent entry itself (entry)
	MagicAttrParent = Key("@@PARENT")
	// name of the node's full path inside MagicDict object (string)
	MagicAttrPath = Key("@@PATH")
	// magic key for appending a list element
	MagicAttrAppend = Key("[]")
)
