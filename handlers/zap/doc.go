/*
Package zaphandler provides a configurable logging handler for github.com/uber-go/zap library, allowing users
to apply censoring to log entries and fields, overriding the original values before passing them to the core.

Please pay attention to log.Infof(template, ...args).
Because it uses fmt.Sprintf, it can't be censored using all the Censor instrument.
The only thing that can be used is a comparison of all the strings against regex patterns.

Example struct with censor tags:

	type User struct {
		Name  string `censor:"display"`
		Email string
	}

	u := User{
		Name:  "John Doe",
		Email: "example@gmail.com",
	}

Usage with default options:

	handler := NewJSONHandler()
	log := slog.New(handler)
	log.Info("user", slog.Any("payload", u))

Output:
{"time":"2023-12-28T20:15:45.893115+01:00","level":"INFO","msg":"user","payload":"{Name: John Doe, Email: [CENSORED]}"}

Usage with custom options:

	censorInst := censor.New()
	opts := []Option{
		WithOut(os.Stdout),
		WithCensor(censorInst),
		WithAddSource(),
		WithReplaceAttr(func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "msg" {
				return slog.Any("msg", "replaced msg")
			}

			return a
		}),
	}

	handler := NewJSONHandler(opts...)
	log := slog.New(handler)
	log.Info("user", slog.Any("payload", u))

Output:

		{
	    	"time": "2023-12-28T20:22:44.35868+01:00",
	    	"level": "INFO",
	    	"source": {
	        	"function": "github.com/vpakhuchyi/censor/handlers/slog.TestFunc",
	        	"file": "/Users/volodymyr/Files/source/censor/handlers/slog/handler_test.go",
	        	"line": "156"
	    	},
	    	"msg": "replaced msg",
	    	"payload": "{Name: John Doe, Email: [CENSORED]}"
		}
*/
package zaphandler
