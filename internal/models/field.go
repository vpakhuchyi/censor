package models

import (
	"reflect"

	"github.com/vpakhuchyi/censor/internal/options"
)

// Field represents a field of a struct.
type Field struct {
	Name  string
	Tag   string
	Value Value
	Opts  options.FieldOptions
	Kind  reflect.Kind
}
