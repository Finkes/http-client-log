package log

import (
	"fmt"
	"github.com/Finkes/http-client-log/pkg/receiver"
)

type Receiver struct{
	options receiver.Options
}

var defaultOptions = receiver.Options {
	Format: receiver.FormatSummary,
}

func NewReceiver(customOptions *receiver.Options) *Receiver {
	receiver:= Receiver{
		options: defaultOptions,
	}
	if customOptions != nil {
		receiver.options = *customOptions
	}
	return &receiver
}

func (r *Receiver) Name() string {
	return "log"
}

func (r *Receiver) Receive(event receiver.Event) error {
	fmt.Println(event.String(&r.options))
	return nil
}