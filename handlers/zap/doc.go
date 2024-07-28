/*
Package zaphandler provides a configurable logging handler for go.uber.org/zap library, allowing users
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
For example, in a call to l.Info("payload", zap.Any("addresses", []string{"address1", "address2"})):
  - "payload" is a "msg"
  - "addresses" is a "key"
  - []string{"address1", "address2"} is a "value"

By default, censor processes only "value" to minimize the overhead. The "msg" and "key" rarely contain sensitive data.
However, there are predefined configuration options that can be passed to NewHandler() function to customize the
censor handler behavior.

The following options are available:
1. WithCensor(censor *censor.Processor) - sets the censor processor instance for the zap handler.
If not provided, a default censor processor is used.
2. WithMessagesFormat() - enables the censoring of a log "msg" values if such are present.
3. WithKeysFormat() - enables the censoring of log "key" values.

After the handler is initialized, it can be used as a regular zap logger. Because the censor handler is a wrapper around
go.uber.org/zap library logic, it may not be compatible with all the possible ways of the logger usage.

That's why it's recommended to use it with the following constructions:

  - l.Info(msg string, fields ...zap.Field)

  - l.With(fields ...zap.Field).Info(msg string, fields ...zap.Field)

    In both cases, Info could be replaced with Debug, Warn, Error, Panic, Fatal.

With sugared logger, the following constructions are supported:

  - l.Info(args ...interface{})

  - l.Infof(template string, args ...interface{})

  - l.Infow(msg string, keysAndValues ...interface{})

  - l.Infoln(args ...interface{})

  - l.With(args ...interface{}).Info(args ...interface{})

    In all cases, Info could be replaced with Debug, Warn, Error, Panic, Fatal.

Methods ending in "f", "ln" and log.Print-style (l.Info) in a sugared logger can't use all the
censor handler features. Due to the nature of the zap sugared logger, censor accepts formatted strings and has no
possibility to use its parsing. However, a features like regexp matching are still available.
*/
package zaphandler
