package hashtable

import "errors"

var (
	ErrKeyInvalid   = errors.New("key must contain at least two letters")
	ErrKeyAlphabet  = errors.New("key must use a single supported alphabet")
	ErrDuplicateKey = errors.New("key already exists")
	ErrTableFull    = errors.New("hash table is full")
	ErrNotFound     = errors.New("key not found")
	ErrTableSize    = errors.New("table size is below minimum")
)
