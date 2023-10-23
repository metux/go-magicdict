package api

import (
	"strconv"
	"strings"
)

const (
	KeyDelimiter  = "::"
	KeyAppendList = "[]::[]"
	RefStart      = "${"
	RefEnd        = "}"
)

// Wraps a key or key path (string) with additional operations
// Key pathes are similar to file system pathes, but using "::" as delimiter
//
// Special cases:
//
//   - Keys with "@@" prefix have special meaning: processed by the leading
//     entry itself -- see Key constants with "MagicAttr" prefix:
//   - Put(('ing a Key "[]" within a [github.com/metux/go-magicdict/core.List]
//     appends the element to the existing list.
//   - When Put()'ing a key path with non-existing leading elements, these are auto-created
//     as [github.com/metux/go-magicdict/core.Dict], if yet missing
//   - If some key path element has "[]" suffix (eg. "foo[]", instead of "foo::[]")
//     then a [github.com/metux/go-magicdict/core.List] is created instead
type Key string

// Split off the first element (head) and also return the remaining part (tail):
// -> (head, tail)
//
// Empty keys return both results empty ("")
// If just one element (single key), returning it as head and "" as tail
func (k Key) Head() (Key, Key) {
	if i := strings.Index(string(k), KeyDelimiter); i >= 0 {
		return Key(k[:i]), k[i+2:]
	}
	return k, ""
}

// Split off the last element (tail) and also return the remaining part (head)
// -> (head, tail)
// Empty keys return both results empty ("")
// If just one element (single key), returning it as tail and "" as head
func (k Key) Tail() (Key, Key) {
	if i := strings.LastIndex(string(k), KeyDelimiter); i >= 0 {
		return k[:i], k[i+2:]
	}
	return "", k
}

// Check whether key is empty ("")
func (k Key) Empty() bool {
	return k == ""
}

// Create a new key path with the given prefix
// Like: <prefix> + "::" + <key>
func (k Key) AddPrefix(prefix string) Key {
	if prefix == "" {
		return k
	}
	return Key(prefix + KeyDelimiter + string(k))
}

// Create a new key path by appending key and suffix
// Like: <key> + "::" + <suffix>
func (k Key) AppendStr(suffix string) Key {
	if k.Empty() {
		return Key(suffix)
	}

	if suffix == "" {
		return k
	}

	return Key(string(k) + KeyDelimiter + suffix)
}

// Create a new key path by appending key and suffix
// Like: <key> + "::" + <suffix>
func (k Key) Append(suffix Key) Key {
	return k.AppendStr(string(suffix))
}

// Check whether it's the special append key
func (k Key) IsAppend() bool {
	return k == MagicAttrAppend
}

// Similar to Append(), but using an it. Used for List element addressing
func (k Key) AppendIdx(idx int) Key {
	return k.AppendStr(strconv.Itoa(idx))
}

func (k Key) MagicDefaults() Key {
	return MagicAttrDefaults.Append(k)
}

func (k Key) MagicLiteralPost() Key {
	return k.Append(MagicAttrLiteral)
}

func (k Key) MagicLiteralPre() Key {
	return MagicAttrLiteral.Append(k)
}

func (k Key) S() string {
	return string(k)
}

// converts the key to a variable reference to the key
func (k Key) R() string {
	return "${" + string(k) + "}"
}

// stringer interface
func (k Key) String() string {
	return string(k)
}

// convert into a key for creating new list
func (k Key) MkAppendList() Key {
	return Key(string(k) + KeyAppendList)
}
