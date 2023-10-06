package api

import (
	"errors"
)

var (
	ErrUnknownEntryType   = errors.New("unknown entry type")
	ErrSubNotSupported    = errors.New("sub-entries not supported on this entry type")
	ErrKeyEmpty           = errors.New("key empty")
	ErrDictNotInitialized = errors.New("dict not initialized")
	ErrIndexOutOfRange    = errors.New("index out of range")
	ErrNilInterface       = errors.New("nil interface")
)
