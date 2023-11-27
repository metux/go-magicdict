package magic

import (
	"log"
	"fmt"

	"github.com/metux/go-magicdict/api"
	"github.com/metux/go-magicdict/core"
	"github.com/metux/go-magicdict/macro"
	"github.com/metux/go-magicdict/parser"
	"github.com/metux/go-magicdict/utils"
)

// This serves as root as well as intermediate nodes
// the root node has Root==nil and Path=="" and defaults!=nil / while the
// others have defaults==nil
// Splittng it up would just save one PTR per instance, but add extra
// complexity and allocations
//
// Only create it via NewMagicFromDict() or don't forget to call Init()
type MagicDict struct {
	Root *MagicDict
	Data api.Entry

	// this is always from the root (need to prepend Path)
	Defaults api.Entry

	// disable variable/macro substitution
	Literal bool

	Path api.Key
}

// Box an entry into MagicDict wrapper, so that defaults still work
// trivial scalars and nil aren't boxed
//
// FIXME: right now just boxing Dict and List
// we'll later box everything that's not inside a Boxed
//
// FIXME: add variable substitution
func (this MagicDict) box(k api.Key, v api.Entry) (api.Entry, error) {

	r := this.Root
	if r == nil {
		r = &this
	}

	switch v.(type) {
	case core.Dict, core.List:
		// FIXME: decide when to box ?
		// const strings or those w/ string interface should not be boxed ?
		sp := MagicDict{
			Root:     r,
			Literal:  this.Literal,
			Path:     this.Path.Append(k),
			Data:     v,
			Defaults: this.Defaults,
		}
		sp.Init()
		return sp, nil

	case core.Scalar:
		v = magicScalar{
			Root: r,
			Path: this.Path.Append(k),
			Data: v.String(),
		}
	}

	if this.Literal {
		return v, nil
	}

	return macro.ProcessVars(v, r)
}

func (this MagicDict) getSpecial(k api.Key, defkey api.Key) (api.Entry, error) {
	log.Println("=========> special loop begin: ", defkey)
	whead, wtail := defkey.Head()
	h := api.Key("")
	for !wtail.Empty() {
		log.Printf("whead \"%s\" wtail \"%s\"\n", whead, wtail)
		ent1, err1 := this.Defaults.Get(whead)
		if ent1 == nil {
			log.Println("xx1: ent==nil", err1)
		} else if ent1.IsScalar() {
			if refname, ok := parser.ParseSimpleRefExpr(ent1.String()); ok {
				log.Println("now should try to fetch ref", refname, "and", wtail)
				referredEntry, referredErr := this.Root.Get(api.Key(refname))
				if referredErr != nil {
					log.Println("referredErr", referredErr)
				}
				if referredEntry == nil {
					log.Println("referredEntry is NIL")
				} else {
					log.Println("referredEntr is OK, calling for", wtail)
					subEnt, subErr := referredEntry.Get(wtail)
					if subErr != nil {
						log.Println("subErr", subErr)
					}
					if subEnt == nil {
						log.Println("subEnt is NIL")
					} else {
						log.Println("subEnt is OK", subEnt.String())
					}
					log.Println("=========> leaving special loop")
					return subEnt, subErr
				}
			} else {
				log.Println("failed parsing refname")
				return nil, fmt.Errorf("failed parsing refname \"%s\" at key %s tail %s", ent1.String(), whead, wtail)
			}
		} else {
			log.Println("not scalar")
		}
		h, wtail = wtail.Head()
		whead = whead.Append(h)
	}
	log.Println("=========> special loop end")

	return nil, nil
}

func (this MagicDict) getdef(walkKey api.Key, defdict api.Entry) (api.Entry, error) {
	head, tail := walkKey.Tail()
	ent, err := defdict.Get(head)

	// FIXME: check tail.Empty()
	if err != nil {
		log.Println("default get ", head, "failed")
		return ent, nil
	}
	if ent == nil {
		log.Println("default get ", head, "NIL")
		return ent, nil
	}
	if ent.IsScalar() {
		log.Println("default path is scalar: ", ent.String())

		if refname, ok := parser.ParseSimpleRefExpr(ent.String()); ok {
			log.Println("now should try to fetch ref", refname, "and", tail)
			referredEntry, referredErr := this.Root.Get(api.Key(refname))
			if referredErr != nil {
				log.Println("referredErr", referredErr)
				return referredEntry, referredErr
			}
			if referredEntry == nil {
				log.Println("referredEntry is NIL")
				return referredEntry, referredErr
			}
			// FIXME: relative keys probably wont work yet
			ent = referredEntry
		}
	}

	if tail.Empty() {
		log.Println("reached final")
		return ent, err
	} else {
		log.Println("got tail", tail)
		return ent.Get(tail)
	}
}

func (this MagicDict) GetDefault(k api.Key) (api.Entry, error) {
	defkey := this.Path.Append(k)
//	ent, err = this.Defaults.Get(defkey)
	return this.getdef(defkey, this.Defaults)
}

func (this MagicDict) Get(k api.Key) (api.Entry, error) {

	log.Printf("MagicDict::Get() path=%s key=%s\n", this.Path, k)

	if k.Empty() {
		log.Println("final")
		return this, nil
	}

	head, tail := k.Head()

	if !tail.Empty() {
		log.Println("tail not empty", tail, "looking for head", head)
		parent, err := this.Get(head)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			log.Println("parent nil, bailing out")
			return nil, nil
		}
		log.Println("asking parent for ", tail)
		return parent.Get(tail)
	}

	switch head {
	case api.MagicAttrLiteral:
		return MagicDict{
			Root:     this.Root,
			Literal:  true,
			Path:     this.Path,
			Data:     this.Data,
			Defaults: this.Defaults,
		}, nil

	case api.MagicAttrPath:
		return core.NewScalarStr(string(this.Path)), nil

	case api.MagicAttrKey, api.MagicAttrShortKey:
		_, p1 := this.Path.Tail()
		return core.NewScalarStr(string(p1)), nil

	case api.MagicAttrParent, api.MagicAttrShortParent:
		p1, _ := this.Path.Tail()
		if this.Root == nil {
			return nil, nil
		}
		return this.Root.Get(p1)

	case api.MagicAttrDefaults:
		return api.MakeDict(this.Defaults, this.Path)
	}

	ent, err := this.Data.Get(k)
	if err != nil {
		log.Println("MagicDict: Get() error: ", err, "on key ", k)
		return nil, err
	}
	if ent != nil {
		return this.box(k, ent)
	}

	log.Println("MagicDict: trying default: ", this.Path.Append(k))

//	defkey := this.Path.Append(k)
//	ent, err = this.Defaults.Get(defkey)
	ent, err = this.GetDefault(k)
	if err != nil {
		log.Println("defaults get error", err)
//		return this.getSpecial(k, defkey)
		return ent, err
	}

	if ent != nil {
		return this.box(k, ent)
	}

	return nil, nil
}

func (this MagicDict) IsScalar() bool {
	return this.Data.IsScalar()
}

func (this MagicDict) IsConst() bool {
	return this.Data.IsConst()
}

func (this MagicDict) String() string {
	return this.Data.String()
}

// this prevents unwanted merging (eg. on lists)
func (this MagicDict) mergeDef() api.Entry {
	if this.MayMergeDefaults() {
		if d, _ := this.Defaults.Get(this.Path); d != nil && !d.Empty() {
			return d
		}
	}
	return nil
}

func (this MagicDict) Elems() []api.Entry {
	elems := make([]api.Entry, 0)
	for _, k := range this.Keys() {
		e, _ := this.Get(k)
		elems = append(elems, e)
	}
	return elems
}

func (this MagicDict) Keys() []api.Key {
	if d := this.mergeDef(); d != nil {
		return utils.UnionSlice(this.Data.Keys(), d.Keys())
	}
	return this.Data.Keys()
}

func (this MagicDict) Put(k api.Key, v api.Entry) error {
	head, tail := k.Head()

	if head == api.MagicAttrDefaults {
		// FIXME: this cannot walk across references yet
		return this.Defaults.Put(this.Path.Append(tail), v)
	}

	// this is a bit ugly, since we have to walk across references
	if !tail.Empty() {
		nlist, head := head.IsListOp()
		cur, _ := this.Get(head)
		if cur == nil {
			if nlist {
				cur = core.EmptyList()
			} else {
				cur = core.EmptyDict()
			}
			if err := this.Data.Put(head, cur); err != nil {
				return err
			}
			cur, _ = this.Get(head)
			if cur == nil {
				panic("cur 2nd time nil")
			}
		}
		return cur.Put(tail, v)
	}
	return this.Data.Put(k, v)
}

func (this MagicDict) Empty() bool {
	return this.Data.Empty() && this.Defaults.Empty()
}

// maybe we'll wanna have different modes here
func (this MagicDict) MayMergeDefaults() bool {
	return this.Data.MayMergeDefaults()
}

func (this *MagicDict) Init() {
	if this.Data == nil {
		this.Data = core.EmptyDict()
	}
	if this.Defaults == nil {
		this.Defaults = core.EmptyDict()
	}
}

// only create it via constructor, since some fields *MUST* be initialized
func NewMagicFromDict(d api.Entry, dflt api.Entry) MagicDict {
	sp := MagicDict{
		Data:     d,
		Defaults: dflt,
	}
	sp.Init()
	return sp
}
