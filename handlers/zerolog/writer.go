package zerologhandler

import "io"

type writer struct {
	out     io.Writer
	handler Handler
}

func (w writer) Write(p []byte) (n int, err error) {
	r := w.handler.censor.Format(string(p))

	return w.out.Write([]byte(r))
}
