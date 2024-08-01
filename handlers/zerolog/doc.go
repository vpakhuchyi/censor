/*
Package zerologhandler provides a configurable logging handler for github.com/rs/zerolog library, allowing users
to apply censoring to log entries and fields, overriding the original values before passing them to the core.

Due to the diversity of the zerolog library usage, please pay attention to the way you use it with the censor handler.

Example of censor handler initialization:

	import (
		"os"

		"github.com/rs/zerolog"
		"github.com/vpakhuchyi/censor"
		zerologhandler "github.com/vpakhuchyi/censor/handlers/zerolog"
	)

	type User struct {
		Name  string `censor:"display"`
		Email string
	}

	u := User{
		Name:  "John Doe",
		Email: "example@gmail.com",
	}

	c := censor.New()
	handler := zerologhandler.New(zerologhandler.WithCensor(c))

	zerolog.InterfaceMarshalFunc = handler.NewInterfaceMarshal
	l := zerolog.New(handler.writer(os.Stderr))

	l.Info().Any("user", u).Send()

Output:
{"level":"info","user":{Name: John Doe, Email: [CENSORED]}}
*/
package zerologhandler
