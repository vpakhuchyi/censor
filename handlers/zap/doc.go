package zaphandler

/*
Package zaphandler integrates github.com/vpakhuchyi/censor with go.uber.org/zap by wrapping an existing zapcore.Core
and sanitizing fields before they reach the underlying core. It is designed for structured logging flows where values
are passed as zap.Field arguments; message strings remain unchanged unless you handle them separately.

Typical workflow:

	import (
		"go.uber.org/zap"
		"go.uber.org/zap/zapcore"

		"github.com/vpakhuchyi/censor"
		zaphandler "github.com/vpakhuchyi/censor/handlers/zap"
	)

	func main() {
		type User struct {
			Name       string `censor:"display"`
			Department string `censor:"display"`
			Email      string
			Token      string
		}

		processor := censor.New()

		logger, err := zap.NewProduction(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zaphandler.NewHandler(core, zaphandler.WithCensor(processor))
		}))
		if err != nil {
			panic(err)
		}

		payload := User{
			Name:       "Hryhorii Skovoroda",
			Department: "Philosophy",
			Email:      "hryhorii.skovoroda@example.com",
			Token:      "S3cr3t!",
		}

		logger.Info("user",
			zap.String("name", payload.Name),
			zap.Any("payload", payload),
		)

		// Output:
		// {"level":"info","msg":"user","name":"Hryhorii Skovoroda","payload":{"Name":"Hryhorii Skovoroda","Department":"Philosophy","Email":"[CENSORED]","Token":"[CENSORED]"}}
	}

Important considerations:

  - The handler sanitizes zap.Field values (strings, reflective payloads, etc.) before delegating to the wrapped core.
    Message strings and key names are untouched, so review them separately if they may contain sensitive data.
  - WithCensor allows you to reuse a shared *censor.Processor; if omitted, NewHandler falls back to a default processor.
  - Because the handler wraps the supplied core, it inherits that core’s encoder, sampling, and level configuration.
    Apply those settings before wrapping.
  - Zap’s formatting helpers that collapse arguments into a single string (for example, Infof or Infoln) are not
    supported—the handler cannot recover individual values from those strings.

Supported zap methods:

  Unsugared:
    - Info(msg string, fields ...zap.Field)
    - Debug(msg string, fields ...zap.Field)
    - Warn(msg string, fields ...zap.Field)
    - Error(msg string, fields ...zap.Field)
    - Panic(msg string, fields ...zap.Field)
    - Fatal(msg string, fields ...zap.Field)
    - With(fields ...zap.Field).Info(msg string, fields ...zap.Field)
    - With(fields ...zap.Field).Debug(msg string, fields ...zap.Field)
    - With(fields ...zap.Field).Warn(msg string, fields ...zap.Field)
    - With(fields ...zap.Field).Error(msg string, fields ...zap.Field)
    - With(fields ...zap.Field).Panic(msg string, fields ...zap.Field)
    - With(fields ...zap.Field).Fatal(msg string, fields ...zap.Field)

  Sugared:
    - Infow(msg string, keysAndValues ...interface{})
    - Debugw(msg string, keysAndValues ...interface{})
    - Warnw(msg string, keysAndValues ...interface{})
    - Errorw(msg string, keysAndValues ...interface{})
    - Panicw(msg string, keysAndValues ...interface{})
    - Fatalw(msg string, keysAndValues ...interface{})
    - With(args ...interface{}).Infow(msg string, keysAndValues ...interface{})
    - With(args ...interface{}).Debugw(msg string, keysAndValues ...interface{})
    - With(args ...interface{}).Warnw(msg string, keysAndValues ...interface{})
    - With(args ...interface{}).Errorw(msg string, keysAndValues ...interface{})
    - With(args ...interface{}).Panicw(msg string, keysAndValues ...interface{})
    - With(args ...interface{}).Fatalw(msg string, keysAndValues ...interface{})

  Unsupported formatting helpers (non-exhaustive):
    - Infof / Debugf / Warnf / Errorf / Panicf / Fatalf
    - Infoln / Debugln / Warnln / Errorln / Panicln / Fatalln
*/
