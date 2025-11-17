package zaphandler

import (
	"go.uber.org/zap/zapcore"

	"github.com/vpakhuchyi/censor"
)

type handler struct {
	zapcore.Core
	censor *censor.Processor
}

// NewHandler returns a new zap logs handler (core) along with a censor processor.
// Options can be provided to configure the Censor processor. If no options are provided,
// a default configuration is used. See the Option documentation for more details.
// By default, the censoring of log fields only is enabled.
func NewHandler(core zapcore.Core, opts ...Option) zapcore.Core {
	cc := handler{Core: core}

	for _, o := range opts {
		o(&cc)
	}

	if cc.censor == nil {
		cc.censor = censor.New()
	}

	if cc.censor.OutputFormat() != censor.OutputFormatJSON {
		panic("zaphandler: censor processor must use json output format")
	}

	return &cc
}

// Write applies censoring to the log entry and fields, overriding the original values.
// Future processing of the log entry and fields will use the given zap core.
func (h handler) Write(e zapcore.Entry, fields []zapcore.Field) error {
	for i := range fields {
		h.censorField(&fields[i])
	}

	return h.Core.Write(e, fields)
}

// Check adds this handler to the CheckedEntry (if the entry should be logged) and returns the result.
func (h handler) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if h.Enabled(e.Level) {
		return ce.AddCore(e, h)
	}

	return ce
}

// With applies censoring to the log fields, overriding the original values before passing them to the core.
func (h handler) With(fields []zapcore.Field) zapcore.Core {
	for i := range fields {
		h.censorField(&fields[i])
	}

	return &handler{
		Core: h.Core.With(fields),
		// Censor instance is shared between the handler instances to avoid additional allocations.
		censor: h.censor,
	}
}

//nolint:exhaustive
func (h handler) censorField(f *zapcore.Field) {
	switch f.Type {
	case zapcore.StringType:
		if f.String != "" {
			f.String = string(h.censor.String(f.String))
		}
	case zapcore.ReflectType:
		if f.Interface != nil {
			f.Interface = &rawJSONValue{data: h.censor.Any(f.Interface)}
		}
	default:
	}
}

type rawJSONValue struct {
	data []byte
}

func (r *rawJSONValue) MarshalJSON() ([]byte, error) {
	return r.data, nil
}
