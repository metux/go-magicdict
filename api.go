package magicdict

import (
	"fmt"

	"github.com/metux/go-magicdict/api"
	"github.com/metux/go-magicdict/core"
	"github.com/metux/go-magicdict/magic"
)

type Key = api.Key
type Entry = api.Entry
type KeyList = api.KeyList
type EntryList = api.EntryList

// loader
var (
	YamlLoad  = magic.YamlLoad
	YamlStore = core.YamlStore
)

// errors
var (
	ErrNilInterface = api.ErrNilInterface
)

// operating on on entries within entry
var (
	EntryGet    = api.GetEntry
	EntryPut    = api.SetEntry
	EntryKeys   = api.GetKeys
	EntryElems  = api.GetElems
	EntryDelete = api.Delete

	EntryStr    = api.GetStr
	EntryPutStr = api.SetStr

	EntryStrList = api.GetStrList

	EntryStrMap = api.GetStrMap

	EntryIntDef = api.GetInt
	EntryPutInt = api.SetInt

	EntryBoolDef = api.GetBool
	EntryPutBool = api.SetBool

	EntryStrListAppend = api.AppendStr
	EntryMakeList      = api.MakeList
	EntryMakeDict      = api.MakeDict
)

func RequiredEntryStr(r Entry, k Key) string {
	s := EntryStr(r, k)
	if s == "" {
		panic(fmt.Sprintf("missing entry for key %s", k))
	}
	return s
}

func EntryKey(r Entry, k Key) Key {
	return Key(EntryStr(r, k))
}

func EntryKeyList(r Entry, k Key) KeyList {
	sl := EntryStrList(r, k)
	kl := make(KeyList, 0, len(sl))
	for idx, str := range sl {
		kl[idx] = Key(str)
	}
	return kl
}

func EntryPutStrList(r Entry, k Key, v []string) error {
	if err := EntryDelete(r, k); err != nil {
		return err
	}
	for _, x := range v {
		if err := EntryStrListAppend(r, k, x); err != nil {
			return err
		}
	}
	return nil
}

func EntryPutStrMap(r Entry, k Key, v map[string]string) error {
	if err := EntryDelete(r, k); err != nil {
		return err
	}
	for k1, v := range v {
		if err := EntryPutStr(r, k.Append(Key(k1)), v); err != nil {
			return err
		}
	}
	return nil
}

func EntryMap(r Entry, k Key) map[Key]Entry {
	m := make(map[Key]Entry)
	if ent := EntryGet(r, k); ent != nil {
		for _, idx := range ent.Keys() {
			if sub, _ := ent.Get(idx); ent != nil {
				m[idx] = sub
			}
		}
	}
	return m
}

// operating on defaults within entry

func DefaultGet(r Entry, k Key) Entry {
	return EntryGet(r, k.MagicDefaults())
}

func DefaultPut(r Entry, k Key, v Entry) error {
	return EntryPut(r, k.MagicDefaults(), v)
}

func DefaultStr(r Entry, k Key) string {
	return EntryStr(r, k.MagicDefaults())
}

func DefaultKey(r Entry, k Key) Key {
	return Key(EntryStr(r, k.MagicDefaults()))
}

func DefaultPutStr(r Entry, k Key, v string) error {
	return EntryPutStr(r, k.MagicDefaults(), v)
}

func DefaultBoolDef(r Entry, k Key, dflt bool) bool {
	return EntryBoolDef(r, k.MagicDefaults(), dflt)
}

func DefaultPutBool(r Entry, k Key, v bool) error {
	return EntryPutBool(r, k.MagicDefaults(), v)
}

func DefaultIntDef(r Entry, k Key, dflt int) int {
	return EntryIntDef(r, k.MagicDefaults(), dflt)
}

func DefaultPutInt(r Entry, k Key, v int) error {
	return EntryPutInt(r, k.MagicDefaults(), v)
}

func IsScalar(e Entry) bool {
	if e == nil {
		return false
	}
	return e.IsScalar()
}

func IsList(e Entry) bool {
	if e == nil {
		return false
	}
	return e.IsList()
}

func IsDict(e Entry) bool {
	if e == nil {
		return false
	}
	return e.IsDict()
}
