package api

type Scalar struct {
	Data string
}

func (sc Scalar) Elems() []Entry {
	return []Entry{}
}

func (sc Scalar) Keys() []Key {
	return []Key{}
}

// FIXME: correct semantics ?
func (sc Scalar) Get(k Key) (Entry, error) {
	if k.Empty() {
		return sc, nil
	}

	return nil, ErrSubNotSupported
}

func (sc Scalar) String() string {
	return sc.Data
}

func (sc Scalar) IsConst() bool {
	return true
}

func (sc Scalar) Put(k Key, v Entry) error {
	return ErrSubNotSupported
}

func (sc Scalar) MayMergeDefaults() bool {
	return false
}

func (sc Scalar) Empty() bool {
	return len(sc.Data) == 0
}

func (sc Scalar) IsScalar() bool {
	return true
}
