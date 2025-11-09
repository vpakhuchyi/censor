package zerologhandler

/*
Package zerologhandler integrates github.com/vpakhuchyi/censor with github.com/rs/zerolog by providing helpers
that build and install a censor-aware marshal function. Zerolog exposes a single global hook,
zerolog.InterfaceMarshalFunc, for customizing how values passed through Any/Interface are encoded. To keep the side
effect explicit and reversible, this package lets you derive the marshal function separately from installing it.

Typical workflow:

	import (
		"os"

		"github.com/rs/zerolog"
		"github.com/vpakhuchyi/censor"
		zerologhandler "github.com/vpakhuchyi/censor/handlers/zerolog"
	)

	func main() {
		type User struct {
			Name       string `censor:"display"`
			Department string `censor:"display"`
			Email      string
			Token      string
		}

		processor := censor.New()

		// Build a marshal function that uses the custom processor (other options like WithZerolog are available).
		marshal := zerologhandler.GetMarshalFunc(
			zerologhandler.WithCensor(processor),
		)

		// Install the marshal function globally and defer restoration of the previous zerolog.InterfaceMarshalFunc.
		restore := zerologhandler.InstallMarshalFunc(marshal)
		defer restore()

		// Construct or reuse a zerolog.Logger. The helper New(...) is available for convenience.
		logger := zerologhandler.New().With().Str("component", "example").Logger()

		payload := User{
			Name:       "Hryhorii Skovoroda",
			Department: "Philosophy",
			Email:      "hryhorii.skovoroda@example.com",
			Token:      "S3cr3t!",
		}

		logger.Info().
			Any("payload", payload).
			Msg("user profile loaded")

		// Output:
		// {"level":"info","component":"example","payload":{"Name":"Hryhorii Skovoroda","Department":"Philosophy","Email":"[CENSORED]","Token":"[CENSORED]"},"message":"user profile loaded"}
	}

Important considerations:

  - Installing the marshal function affects every zerolog logger in the process until you call the restore callback.
    Always restore the previous value when the override is no longer needed (for example, at the end of main or in
    test cleanup).
  - Options such as WithCensor and WithZerolog can be supplied to New and GetMarshalFunc so that loggers and marshal
    functions share the same configuration.
  - Types other than Any/Interface are still serialized by zerolog itself; they will not pass through Censor unless
    zerolog exposes dedicated hooks for them in the future.
*/
