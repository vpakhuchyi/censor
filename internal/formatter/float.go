package formatter

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/vpakhuchyi/censor/internal/models"
)

// Float formats a value as a float.
// The value is formatted with up to 7 significant figures for float32 and up to 15 significant figures for float64.
// Note: this method panics if the provided value is not a float.
func (f *Formatter) Float(v models.Value) string {
	if v.Kind != reflect.Float32 && v.Kind != reflect.Float64 {
		panic("provided value is not a float")
	}

	// Significant figures are the meaningful digits in a number. They indicate
	// the precision of a measurement. For example, in the number 123.4558, there
	// are seven significant figures, and rounding to four significant figures
	// would result in 123.5.
	// More details about significant figures: https://en.wikipedia.org/wiki/Significant_figures.
	var maxSignificantFigures string
	if v.Kind == reflect.Float32 {
		maxSignificantFigures = strconv.Itoa(f.float32MaxSignificantFigures)
	} else {
		maxSignificantFigures = strconv.Itoa(f.float64MaxSignificantFigures)
	}

	return fmt.Sprintf("%."+maxSignificantFigures+"g", v.Value)
}
