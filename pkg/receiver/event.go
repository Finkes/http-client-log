package receiver

import (
	"io"
)

type Event interface {
	Serialize() (io.Reader, error)
	String(options *Options) string
	Type() string
}
