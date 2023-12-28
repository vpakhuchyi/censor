package formatter

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/vpakhuchyi/censor/config"
	"github.com/vpakhuchyi/censor/internal/models"
)

// Formatter is used to format values.
type Formatter struct {
	// maskValue is used to mask struct fields with sensitive data.
	// The default value is stored in config.DefaultMaskValue constant.
	maskValue string
	// displayPointerSymbol is used to display '&' (pointer symbol) in the output.
	// The default value is false.
	displayPointerSymbol bool
	// displayStructName is used to display struct name in the output.
	// A struct name includes the last part of the package path.
	// The default value is false.
	displayStructName bool
	// displayMapType is used to display map type in the output.
	// The default value is false.
	displayMapType bool
	// excludePatterns contains regexp patterns that are used to identify strings that must be masked.
	excludePatterns []string
	// excludePatternsCompiled contains already compiled regexp patterns from excludePatterns.
	excludePatternsCompiled []*regexp.Regexp
}

// New returns a new instance of Formatter with default configuration.
func New() *Formatter {
	return &Formatter{
		maskValue:               config.DefaultMaskValue,
		displayPointerSymbol:    false,
		displayStructName:       false,
		displayMapType:          false,
		excludePatterns:         nil,
		excludePatternsCompiled: nil,
	}
}

// NewWithConfig returns a new instance of Formatter with given configuration.
func NewWithConfig(cfg config.Formatter) *Formatter {
	f := Formatter{
		maskValue:            cfg.MaskValue,
		displayPointerSymbol: cfg.DisplayPointerSymbol,
		displayStructName:    cfg.DisplayStructName,
		displayMapType:       cfg.DisplayMapType,
		excludePatterns:      cfg.ExcludePatterns,
	}

	if len(f.excludePatterns) != 0 {
		f.compileExcludePatterns()
	}

	return &f
}

// compileExcludePatterns compiles regexp patterns from excludePatterns.
// Note: this method may panic if regexp pattern is invalid.
func (f *Formatter) compileExcludePatterns() {
	if f.excludePatterns != nil {
		f.excludePatternsCompiled = make([]*regexp.Regexp, len(f.excludePatterns))
		for i, pattern := range f.excludePatterns {
			f.excludePatternsCompiled[i] = regexp.MustCompile(pattern)
		}
	}
}

// SetMaskValue sets a value that will be used to mask struct fields.
func (f *Formatter) SetMaskValue(mask string) {
	f.maskValue = mask
}

// DisplayPointerSymbol sets whether to display the '&' (pointer symbol) before the pointed value.
func (f *Formatter) DisplayPointerSymbol(v bool) {
	f.displayPointerSymbol = v
}

// DisplayStructName sets whether to display the name of the struct.
func (f *Formatter) DisplayStructName(v bool) {
	f.displayStructName = v
}

// DisplayMapType sets whether to display map type in the output.
func (f *Formatter) DisplayMapType(v bool) {
	f.displayMapType = v
}

// AddExcludePatterns adds regexp patterns that are used for the selection of strings that must be masked.
// Regexp patterns compilation will be triggered automatically after adding new patterns.
// Note: this method may panic if regexp pattern is invalid.
func (f *Formatter) AddExcludePatterns(patterns ...string) {
	f.excludePatterns = append(f.excludePatterns, patterns...)
	f.compileExcludePatterns()
}

//nolint:exhaustive,gocyclo
func (f *Formatter) writeValue(buf *strings.Builder, v models.Value) {
	switch v.Kind {
	case reflect.String:
		buf.WriteString(f.String(v))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf.WriteString(f.Integer(v))
	case reflect.Float32, reflect.Float64:
		buf.WriteString(f.Float(v))
	case reflect.Complex64, reflect.Complex128:
		buf.WriteString(f.Complex(v))
	case reflect.Struct:
		buf.WriteString(f.Struct(v.Value.(models.Struct)))
	case reflect.Slice, reflect.Array:
		buf.WriteString(f.Slice(v.Value.(models.Slice)))
	case reflect.Pointer:
		buf.WriteString(f.Ptr(v.Value.(models.Ptr)))
	case reflect.Map:
		buf.WriteString(f.Map(v.Value.(models.Map)))
	case reflect.Bool:
		buf.WriteString(f.Bool(v))
	case reflect.Interface:
		buf.WriteString(f.Interface(v))
	}
}

//nolint:exhaustive,gocyclo
func (f *Formatter) writeField(field models.Field, buf *strings.Builder) {
	switch field.Kind {
	case reflect.String:
		buf.WriteString(formatField(field.Name, f.String(field.Value)))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf.WriteString(formatField(field.Name, f.Integer(field.Value)))
	case reflect.Float32, reflect.Float64:
		buf.WriteString(formatField(field.Name, f.Float(field.Value)))
	case reflect.Complex64, reflect.Complex128:
		buf.WriteString(formatField(field.Name, f.Complex(field.Value)))
	case reflect.Struct:
		buf.WriteString(formatField(field.Name, f.Struct(field.Value.Value.(models.Struct))))
	case reflect.Slice, reflect.Array:
		buf.WriteString(formatField(field.Name, f.Slice(field.Value.Value.(models.Slice))))
	case reflect.Pointer:
		buf.WriteString(formatField(field.Name, f.Ptr(field.Value.Value.(models.Ptr))))
	case reflect.Bool:
		buf.WriteString(formatField(field.Name, f.Bool(field.Value)))
	case reflect.Map:
		buf.WriteString(formatField(field.Name, f.Map(field.Value.Value.(models.Map))))
	case reflect.Interface:
		buf.WriteString(formatField(field.Name, f.Interface(field.Value)))
	}
}
