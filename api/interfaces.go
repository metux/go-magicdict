package api

type Entry interface {
	// get arbitrary key
	Get(k Key) (Entry, error)

	// Put arbitrary value
	Put(k Key, value Entry) error

	// get a list of direct subkeys
	Keys() []Key

	// get a list of values
	Elems() []Entry

	Empty() bool

	// get plain string representation (if any)
	String() string

	// when const scalars are Put'ed, their value may be taken directly
	// thus dropping the ScalarEntry object
	IsConst() bool

	// merge the entry's subs with defaults ?
	MayMergeDefaults() bool

	// Typechecks
	IsScalar() bool
}

type EntryDefaults interface {
	Entry
	SetDefaultEntry(k Key, val Entry) error
}
