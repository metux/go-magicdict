package api

import (
	"testing"
)

func TestKeysSplit(t *testing.T) {
	key := Key("abc::booh::cat")

	h1, t1 := key.Head()
	if h1 != "abc" || t1 != "booh::cat" {
		t.Fatalf("wrong head split: '%s' '%s'", h1, t1)
	}

	h2, t2 := key.Tail()
	if h2 != "abc::booh" || t2 != "cat" {
		t.Fatalf("wrong tail split: '%s' '%s'", h2, t2)
	}
}

func TestKeyEmpty(t *testing.T) {
	k1 := Key("abc::def::foo")
	if k1.Empty() {
		t.Fatalf("empty check 1 failed")
	}
	k2 := Key("")
	if !k2.Empty() {
		t.Fatalf("empty check 2 failed")
	}
}

func assertKeys(t *testing.T, k1 Key, k2 Key) {
	if k1 != k2 {
		t.Fatalf("key append mismatch: \"%s\" should be \"%s\"", k1, k2)
	}
}

func TestKeyAppend(t *testing.T) {
	assertKeys(t, Key("").AppendStr("wurst"), Key("wurst"))
	assertKeys(t, Key("foo").AppendStr("bar"), Key("foo::bar"))
}
