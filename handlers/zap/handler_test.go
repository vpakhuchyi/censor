package zaphandler

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/vpakhuchyi/censor"
)

func TestNewHandler(t *testing.T) {
	const logFileName = "test_log"

	c, err := censor.NewWithOpts(censor.WithConfig(&censor.Config{
		Encoder: censor.EncoderConfig{
			MaskValue:            censor.DefaultMaskValue,
			DisplayPointerSymbol: false,
			DisplayStructName:    false,
			DisplayMapType:       false,
			ExcludePatterns:      []string{`#sensitive#`},
		},
		General: censor.General{
			OutputFormat: censor.OutputFormatJSON,
		},
	}))
	require.NoError(t, err)

	msg, key := "some-msg", "key"
	value := struct {
		Name  string `censor:"display"`
		Text  string `censor:"display"`
		Email string
	}{
		Name:  "Petro Perekotypole",
		Text:  `so"me text with #sensitive# data`,
		Email: "example@example.com",
	}

	// Un-sugared logger.
	t.Run("info", func(t *testing.T) {
		// GIVEN.
		core := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return NewHandler(core, WithCensor(c))
		})
		outputPath := path.Join(t.TempDir(), logFileName)
		outputFile, err := os.Create(outputPath)
		require.NoError(t, err)

		l := newTestProductionZap(t, outputPath, core)

		want := `"key":{"Name":"Petro Perekotypole","Text":"so\"me text with [CENSORED] data","Email":"[CENSORED]"}`

		// WHEN.
		l.Info(msg, zap.Any(key, value))

		// THEN.
		got := readLogs(t, outputFile)
		require.Contains(t, string(got), want)
	})

	t.Run("info with string args", func(t *testing.T) {
		// GIVEN.
		core := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return NewHandler(core, WithCensor(c))
		})
		outputPath := path.Join(t.TempDir(), logFileName)
		outputFile, err := os.Create(outputPath)
		require.NoError(t, err)

		l := newTestProductionZap(t, outputPath, core)

		want := `"Petro Perekotypole":"so\"me text with [CENSORED] data"`

		// WHEN.
		l.Info(msg, zap.String(value.Name, value.Text))

		// THEN.
		got := readLogs(t, outputFile)
		require.Contains(t, string(got), want)
	})

	t.Run("error", func(t *testing.T) {
		// GIVEN.
		core := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return NewHandler(core, WithCensor(c))
		})
		outputPath := path.Join(t.TempDir(), logFileName)
		outputFile, err := os.Create(outputPath)
		require.NoError(t, err)

		l := newTestProductionZap(t, outputPath, core)

		want := `"key":{"Name":"Petro Perekotypole","Text":"so\"me text with [CENSORED] data","Email":"[CENSORED]"}`

		// WHEN.
		l.Error(msg, zap.Any(key, value))

		// THEN.
		got := readLogs(t, outputFile)
		require.Contains(t, string(got), want)
	})

	t.Run("debug", func(t *testing.T) {
		// GIVEN.
		core := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return NewHandler(core, WithCensor(c))
		})
		outputPath := path.Join(t.TempDir(), logFileName)
		outputFile, err := os.Create(outputPath)
		require.NoError(t, err)

		l := newTestDevelopmentZap(t, outputPath, core)

		want := `{"Name":"Petro Perekotypole","Text":"so\"me text with [CENSORED] data","Email":"[CENSORED]"}`

		// WHEN.
		l.Debug(msg, zap.Any(key, value))

		// THEN.
		got := readLogs(t, outputFile)
		require.Contains(t, string(got), want)
	})

	t.Run("warn", func(t *testing.T) {
		// GIVEN.
		core := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return NewHandler(core, WithCensor(c))
		})
		outputPath := path.Join(t.TempDir(), logFileName)
		outputFile, err := os.Create(outputPath)
		require.NoError(t, err)

		l := newTestProductionZap(t, outputPath, core)

		want := `{"Name":"Petro Perekotypole","Text":"so\"me text with [CENSORED] data","Email":"[CENSORED]"}`

		// WHEN.
		l.Warn(msg, zap.Any(key, value))

		// THEN.
		got := readLogs(t, outputFile)
		require.Contains(t, string(got), want)
	})

	t.Run("panic", func(t *testing.T) {
		// GIVEN.
		core := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return NewHandler(core, WithCensor(c))
		})
		outputPath := path.Join(t.TempDir(), logFileName)
		outputFile, err := os.Create(outputPath)
		require.NoError(t, err)

		l := newTestProductionZap(t, outputPath, core)

		want := `{"Name":"Petro Perekotypole","Text":"so\"me text with [CENSORED] data","Email":"[CENSORED]"}`

		// WHEN.
		require.Panics(t, func() { l.Panic(msg, zap.Any(key, value)) })

		// THEN.
		got := readLogs(t, outputFile)
		require.Contains(t, string(got), want)
	})

	t.Run("fatal", func(t *testing.T) {
		// GIVEN.
		core := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return NewHandler(core, WithCensor(c))
		})
		outputPath := path.Join(t.TempDir(), logFileName)
		outputFile, err := os.Create(outputPath)
		require.NoError(t, err)

		l := newTestProductionZap(t, outputPath, core)

		// By default, a call to Fatal will exit the program with no possibility to validate the output.
		// To avoid this, we can use the WithFatalHook option to write the log message and then panic instead.
		// We don't need to test the Fatal method itself, so os.Exit(1) can be replaces with a panic.
		// Our goal is just to be sure that in case of such a call a censor handler works correctly.
		l = l.WithOptions(zap.WithFatalHook(zapcore.WriteThenPanic))

		want := `{"Name":"Petro Perekotypole","Text":"so\"me text with [CENSORED] data","Email":"[CENSORED]"}`

		// WHEN.
		require.Panics(t, func() { l.Fatal(msg, zap.Any(key, value)) })

		// THEN.
		got := readLogs(t, outputFile)
		require.Contains(t, string(got), want)
	})

	t.Run("with info", func(t *testing.T) {
		// GIVEN.
		core := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return NewHandler(core, WithCensor(c))
		})
		outputPath := path.Join(t.TempDir(), logFileName)
		outputFile, err := os.Create(outputPath)
		require.NoError(t, err)

		l := newTestProductionZap(t, outputPath, core)

		want := `"key":{"Name":"Petro Perekotypole","Text":"so\"me text with [CENSORED] data","Email":"[CENSORED]"},"key":{"Name":"Petro Perekotypole","Text":"so\"me text with [CENSORED] data","Email":"[CENSORED]"}`

		// WHEN.
		l.With(zap.Any(key, value)).Info(msg, zap.Any(key, value))

		// THEN.
		got := readLogs(t, outputFile)
		require.Contains(t, string(got), want)
	})

	// Sugared logger.
	t.Run("info", func(t *testing.T) {
		// GIVEN.
		core := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return NewHandler(core, WithCensor(c))
		})
		outputPath := path.Join(t.TempDir(), logFileName)
		outputFile, err := os.Create(outputPath)
		require.NoError(t, err)

		l := newTestProductionZap(t, outputPath, core)
		sl := l.Sugar()

		want := `"key":{"Name":"Petro Perekotypole","Text":"so\"me text with [CENSORED] data","Email":"[CENSORED]"}`

		// WHEN.
		//sl.Infow(msg, key, value)
		sl.With(key, value).Infow(msg, key, value)

		// THEN.
		got := readLogs(t, outputFile)
		fmt.Println(string(got))
		require.Contains(t, string(got), want)
	})

	t.Run("info with default censor", func(t *testing.T) {
		// GIVEN.
		core := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return NewHandler(core)
		})
		outputPath := path.Join(t.TempDir(), logFileName)
		outputFile, err := os.Create(outputPath)
		require.NoError(t, err)

		l := newTestProductionZap(t, outputPath, core)

		// #sensitive# is not replaced with the mask value because the default censor is used and it
		// has no configured exclude patterns.
		want := `"key":{"Name":"Petro Perekotypole","Text":"so\"me text with #sensitive# data","Email":"[CENSORED]"}`

		// WHEN.
		l.Info(msg, zap.Any(key, value))

		// THEN.
		got := readLogs(t, outputFile)
		require.Contains(t, string(got), want)
	})
}

// readLogs reads logs from the output file and returns them as a byte slice.
func readLogs(t *testing.T, f *os.File) []byte {
	fs, err := f.Stat()
	require.NoError(t, err)

	got := make([]byte, fs.Size())

	_, err = f.Read(got)
	require.NoError(t, err)

	return got
}

// newTestProductionZap creates a new Zap production logger with the output set to the given file.
func newTestProductionZap(t *testing.T, output string, core zap.Option) *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{output}
	cfg.ErrorOutputPaths = []string{output}

	l, err := cfg.Build(core)
	require.NoError(t, err)

	return l
}

// newTestDevelopmentZap creates a new Zap development logger with the output set to the given file.
// Note: please pay attention to the fact that the development logger has a different format of the output.
func newTestDevelopmentZap(t *testing.T, output string, core zap.Option) *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{output}
	cfg.ErrorOutputPaths = []string{output}

	l, err := cfg.Build(core)
	require.NoError(t, err)

	return l
}
