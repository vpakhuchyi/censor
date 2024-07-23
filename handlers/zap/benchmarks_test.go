package zaphandler

import (
	"io"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// BenchmarkWithCensor-12    	      180955	      6057 ns/op	    4913 B/op	     146 allocs/op
// BenchmarkWithCensor-12    	      179113	      6081 ns/op	    4912 B/op	     146 allocs/op
// BenchmarkWithoutCensor-12    	 1000000	      1028 ns/op	     304 B/op	       7 allocs/op
// BenchmarkWithoutCensor-12    	 1000000	      1015 ns/op	     304 B/op	       7 allocs/op

func BenchmarkWithCensor(b *testing.B) {
	o := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return NewHandler(core)
	})

	zl := zap.New(newBenchZapCore(), o)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		zl.Info("user", zap.Any("payload", payload))
	}
}

func BenchmarkWithoutCensor(b *testing.B) {
	zl := zap.New(newBenchZapCore())

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		zl.Info("user", zap.Any("payload", payload))
	}
}

type profile struct {
	ID         int    `censor:"display"`
	Speciality string `censor:"display"`
	Grade      int
}
type user struct {
	Name      string `censor:"display"`
	Email     string
	Addresses []string
	Profile   profile `censor:"display"`
}

var payload = map[string]interface{}{
	"1": user{
		Name:      "John Doe",
		Email:     "example@gmail.com",
		Addresses: []string{"address1", "address2"},
		Profile: profile{
			ID:         123,
			Speciality: "doctor",
			Grade:      5,
		},
	},
	"2": user{
		Name:      "Ivan Sirko",
		Email:     "example2@gmail.com",
		Addresses: []string{"address1", "address2"},
		Profile: profile{
			ID:         123,
			Speciality: "doctor",
			Grade:      5,
		},
	},
}

func newBenchZapCore() zapcore.Core {
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeDuration = zapcore.NanosDurationEncoder
	ec.EncodeTime = zapcore.EpochNanosTimeEncoder
	enc := zapcore.NewJSONEncoder(ec)

	return zapcore.NewCore(
		enc,
		&Discarder{},
		zapcore.InfoLevel,
	)
}

// A Syncer is a spy for the Sync portion of zapcore.WriteSyncer.
type Syncer struct {
	err    error
	called bool
}

// SetError sets the error that the Sync method will return.
func (s *Syncer) SetError(err error) {
	s.err = err
}

// Sync records that it was called, then returns the user-supplied error (if
// any).
func (s *Syncer) Sync() error {
	s.called = true
	return s.err
}

// Called reports whether the Sync method was called.
func (s *Syncer) Called() bool {
	return s.called
}

// A Discarder sends all writes to io.Discard.
type Discarder struct{ Syncer }

// Write implements io.Writer.
func (d *Discarder) Write(b []byte) (int, error) {
	return io.Discard.Write(b)
}
