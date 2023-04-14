package broker

import (
	"fmt"
	"github.com/Finkes/http-client-log/pkg/receiver"
)

type Broker struct {
	receivers []receiver.Receiver
	options   receiver.Options
}

func NewBroker(receivers []receiver.Receiver, options receiver.Options) *Broker {
	return &Broker{receivers: receivers, options: options}
}

// CaptureEvent sends all captured events to the enabled receivers
func (b *Broker) CaptureEvent(event receiver.Event) {
	for _, receiver := range b.receivers {
		if err := receiver.Receive(event); err != nil {
			fmt.Printf("failed to send message to receiver %v: %v", receiver.Name(), err)
		}
	}
}
