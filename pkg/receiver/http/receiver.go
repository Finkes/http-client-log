package http

import (
	"fmt"
	"github.com/Finkes/http-client-log/pkg/receiver"
	"log"
	"net/http"
)

type Receiver struct {
	url string
}

func (r *Receiver) Name() string {
	return fmt.Sprintf("http %v", r.url)
}

func (r *Receiver) Receive(event receiver.Event) error {
	serializedEvent, err := event.Serialize()
	if err != nil {
		return err
	}
	resp, err := http.Post(r.url, "application/json", serializedEvent)
	if err != nil {
		return err
	}
	log.Println("http receiver response:" + resp.Status)
	return nil
}
