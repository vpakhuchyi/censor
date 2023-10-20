package models

import (
	"reflect"

	"github.com/vpakhuchyi/sanitiser/internal/options"
)

type Field struct {
	Name  string
	Tag   string
	Value Value
	Opts  options.FieldOptions
	Kind  reflect.Kind
}
