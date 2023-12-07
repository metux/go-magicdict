package core

import (
	"log"

	"github.com/metux/go-magicdict/api"
)

// Mimics the behaviour of symbolic links in file systems:
// Linking to an path inside some parent entry
// Each operation first fetches the target entry and then calls it
// With caching enabled, result of first (successful) fetch is reused.
//
// NOTE: in most cases, using variable substitution in MagicDict
// is the better alternative
type Symlink struct {
	Parent   api.Entry
	Path     api.Key
	Caching  bool
	Cacheval api.Entry
}

func (this Symlink) fetch() (api.Entry, error) {
	if !this.Caching {
		return this.Parent.Get(this.Path)
	}
	if this.Cacheval == nil {
		cv, err := this.Parent.Get(this.Path)
		if err != nil {
			return nil, err
		}
		this.Cacheval = cv
	}
	return this.Cacheval, nil
}

func (this Symlink) Elems() api.EntryList {
	if orig, err := this.fetch(); err == nil {
		return orig.Elems()
	} else {
		log.Printf("symlink fetch error: %v", err)
		return api.EntryList{}
	}
}

func (this Symlink) Keys() api.KeyList {
	if orig, err := this.fetch(); err == nil {
		return orig.Keys()
	} else {
		log.Printf("symlink fetch error: %v", err)
		return api.KeyList{}
	}
}

func (this Symlink) Get(k api.Key) (api.Entry, error) {
	if orig, err := this.fetch(); err == nil {
		return orig.Get(k)
	} else {
		log.Printf("Symlink fetch error: %v", err)
		return nil, err
	}
}

func (this Symlink) String() string {
	if orig, err := this.fetch(); err == nil {
		return orig.String()
	} else {
		log.Printf("Symlink fetch error: %v", err)
		return ""
	}
}

func (this Symlink) IsConst() bool {
	if orig, err := this.fetch(); err == nil {
		return orig.IsConst()
	} else {
		log.Printf("Symlink fetch error: %v", err)
		return true
	}
}

func (this Symlink) Put(k api.Key, v api.Entry) error {
	if orig, err := this.fetch(); err == nil {
		return orig.Put(k, v)
	} else {
		log.Printf("Symlink fetch error: %v", err)
		return err
	}
}

func (this Symlink) MayMergeDefaults() bool {
	if orig, err := this.fetch(); err == nil {
		return orig.MayMergeDefaults()
	} else {
		log.Printf("Symlink fetch error: %v", err)
		return false
	}
}

func (this Symlink) Empty() bool {
	if orig, err := this.fetch(); err == nil {
		return orig.Empty()
	} else {
		log.Printf("Symlink fetch error: %v", err)
		return true
	}
}

func (this Symlink) IsScalar() bool {
	if orig, err := this.fetch(); err == nil {
		return orig.IsScalar()
	} else {
		log.Printf("Symlink fetch error: %v", err)
		return true
	}
}

func (this Symlink) IsList() bool {
	if orig, err := this.fetch(); err == nil {
		return orig.IsList()
	} else {
		log.Printf("Symlink fetch error: %v", err)
		return true
	}
}

func (this Symlink) IsDict() bool {
	if orig, err := this.fetch(); err == nil {
		return orig.IsDict()
	} else {
		log.Printf("Symlink fetch error: %v", err)
		return true
	}
}

// Implementing yml.Marshaler interface
// Returning the entry that the link is pointing to, thus will be mashaled
// like as the referenced entry would be directly here.
func (this Symlink) MarshalYAML() (interface{}, error) {
	// Not having a MarshalYAML() causes infinite loop in yaml encoder,
	// so we're doing lookup and return the proxied entry.
	// OTOH, we could also send out the entry referred by Path
	// The best would be telling yaml encoder to emit a reference
	return this.Parent.Get(this.Path)
}

func NewSymlink(parent api.Entry, path api.Key, caching bool) api.Entry {
	return Symlink{Parent: parent, Path: path, Caching: caching}
}
