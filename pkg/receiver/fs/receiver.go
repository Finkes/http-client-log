package fs

import (
	"fmt"
	"github.com/Finkes/http-client-log/pkg/receiver"
	"os"
	"path"
)

type Receiver struct {
	options  receiver.Options
	filePath string
}

var defaultOptions = receiver.Options{
	Format:      receiver.FormatYaml,
	LogFileName: "http-client-log.txt",
}

func NewReceiver(filePath string, customOptions *receiver.Options) *Receiver {
	receiver := Receiver{
		filePath: filePath,
		options:  defaultOptions,
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
	formattedEvent := event.String(&r.options)
	return r.appendToFile(formattedEvent)
}

func (r *Receiver) appendToFile(text string) error {
	f, err := os.OpenFile(path.Join(r.filePath, r.options.LogFileName),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("%v\n", text))
	return err
}
