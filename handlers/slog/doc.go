package sloghandler

/*
Package sloghandler integrates github.com/vpakhuchyi/censor with log/slog by providing drop-in handlers that mask
sensitive values before they are rendered. The handler wraps slog’s JSON/Text handlers, forwarding every record after
running it through a censor. This is useful when you want to keep using slog’s native APIs while ensuring that only the
fields explicitly tagged for display remain readable.

Typical workflow:

	import (
		"log/slog"
		"os"

		"github.com/vpakhuchyi/censor"
		sloghandler "github.com/vpakhuchyi/censor/handlers/slog"
	)

	func main() {
		type User struct {
			Name       string `censor:"display"`
			Department string `censor:"display"`
			Email      string
			Token      string
		}

		payload := User{
			Name:       "Hryhorii Skovoroda",
			Department: "Philosophy",
			Email:      "h.skovoroda@example.com",
			Token:      "S3cr3t!",
		}

		handler := sloghandler.NewJSONHandler()
		logger := slog.New(handler)

		logger.Info("profile loaded", slog.Any("payload", payload))

		// Output:
		// {"time":"2024-06-02T15:04:05Z","level":"INFO","msg":"profile loaded","payload":{"Name":"Hryhorii Skovoroda","Department":"Philosophy","Email":"[CENSORED]","Token":"[CENSORED]"}}
	}

Important considerations:

  - Only attribute values are masked. The log message itself (the string passed to logger.Info) and attribute keys remain
    untouched, so avoid placing sensitive data there.
  - The handler defaults to JSON output writing to os.Stdout. Use WithOut to point it at another io.Writer, or
    WithAddSource to include caller information.
  - Providing WithCensor lets you reuse an existing *censor.Processor, keeping configuration consistent across services.
  - WithReplaceAttr is applied after censoring. If you replace the value with a string, ensure you do not reintroduce
    sensitive information.

Supported options:

  - WithCensor(*censor.Processor) — reuse a prepared processor instead of the default.
  - WithOut(io.Writer) — change the target writer.
  - WithAddSource() — include source metadata (file/line/function) in each record.
  - WithReplaceAttr(func([]string, slog.Attr) slog.Attr) — post-process attributes (called after censoring).

By composing these options, you can keep using slog as usual while guaranteeing that sensitive payload data is
automatically redacted.
*/
