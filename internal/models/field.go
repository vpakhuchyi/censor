package models

import (
	"github.com/vpakhuchyi/censor/internal/options"
)

// Field represents a field of a struct.
type Field struct {
	Name  string
	Value Value
	Opts  options.FieldOptions
}
