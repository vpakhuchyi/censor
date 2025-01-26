/*
Package zaphandler provides a configurable logging handler for the go.uber.org/zap library, allowing users
to apply censoring to log entries and fields, overriding the original values before passing them to the core.

Due to the diversity of the zap library usage, please pay attention to the way you use it with the censor handler.

Example of censor handler initialization:

	import (
		censorlog "github.com/vpakhuchyi/censor/handlers/zap"
		"go.uber.org/zap"
		"go.uber.org/zap/zapcore"
	)

	type User struct {
		Name  string `censor:"display"`
		Email string
	}

	u := User{
		Name:  "John Doe",
		Email: "example@gmail.com",
	}

	o := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return censorlog.NewHandler(core)
	})

	l, err := zap.NewProduction(o)
	if err != nil {
		// handle error
	}

	l.Info("user", zap.Any("payload", u))

Output:
{"level":"info",...,"msg":"user","payload":"{Name: John Doe, Email: [CENSORED]}"}

Describing the logger usage we can operate with a few keywords: "msg", "key" and "value".
For example, in a call to `l.Info("payload", zap.Any("addresses", []string{"address1", "address2"}))`:
  - "payload" is a "msg"
  - "addresses" is a "key"
  - []string{"address1", "address2"} is a "value"

By default, censor processes only "value" to minimize the overhead. The "msg" and "key" rarely contain sensitive data,
because it's usually a static string or a key name. So, it's the user's responsibility to care about these fields.
Enabling censor for "msg" and "key" would increase the overhead and may lead to performance issues.

The following options are available:
  - `WithCensor(censor *censor.Processor)` - sets the censor processor instance for the zap handler.
    If not provided, a default censor processor is used.

After the handler is initialized, it can be used as a regular zap logger. Because the censor handler is a wrapper around
go.uber.org/zap library logic, it may not be compatible with all the possible ways of the logger usage.

# Logger Usage Recommendations

While the zaphandler works seamlessly with the unsugared logger, its compatibility with the sugared logger is limited.
Supported Unsugared Methods:

	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(fields ...zap.Field).Info(msg string, fields ...zap.Field)
	With(fields ...zap.Field).Debug(msg string, fields ...zap.Field)
	With(fields ...zap.Field).Warn(msg string, fields ...zap.Field)
	With(fields ...zap.Field).Error(msg string, fields ...zap.Field)
	With(fields ...zap.Field).Panic(msg string, fields ...zap.Field)
	With(fields ...zap.Field).Fatal(msg string, fields ...zap.Field)

Supported Sugared Methods:

	Infow(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	With(args ...interface{}).Infow(msg string, keysAndValues ...interface{})
	With(args ...interface{}).Debugw(msg string, keysAndValues ...interface{})
	With(args ...interface{}).Warnw(msg string, keysAndValues ...interface{})
	With(args ...interface{}).Errorw(msg string, keysAndValues ...interface{})
	With(args ...interface{}).Panicw(msg string, keysAndValues ...interface{})
	With(args ...interface{}).Fatalw(msg string, keysAndValues ...interface{})

Other methods that are not listed here are not supported due to the way data is handled by zap before passing it
to the censor handler. Such methods concatenate arguments using fmt.Sprint-like functions and pass them to the censor
handler as single formatted strings. As a result, the censor handler cannot process individual values, and sensitive
data within these concatenated strings will not be censored.
*/
package zaphandler
