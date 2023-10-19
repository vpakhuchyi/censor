package models

import "reflect"

type Value struct {
	Value any
	Kind  reflect.Kind
}
