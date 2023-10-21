package models

import "reflect"

type Value struct {
	Value any
	Kind  reflect.Kind
}

type KV struct {
	SortValue string
	Key       Value
	Value     Value
}
