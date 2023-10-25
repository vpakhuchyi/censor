package models

import "reflect"

// Value represents a value including its type and kind.
type Value struct {
	Value any
	Kind  reflect.Kind
}

// KV represents a key-value pair.
type KV struct {
	SortValue string // SortValue is used to sort the map keys to ensure consistent output.
	Key       Value
	Value     Value
}
