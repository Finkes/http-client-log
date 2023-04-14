package http_client_log

import (
	"github.com/Finkes/http-client-log/pkg/broker"
	"github.com/Finkes/http-client-log/pkg/interceptor/http"
	"github.com/Finkes/http-client-log/pkg/receiver"
	"github.com/Finkes/http-client-log/pkg/receiver/fs"
	"github.com/Finkes/http-client-log/pkg/receiver/log"
)

type Config struct {
	receivers []receiver.Receiver
	options   receiver.Options
}

type Option func(options *Config)

// WithFsTarget writes all events to the filesystem (single file vs multiple files)
func WithFsTarget() {}

// WithRemoteTarget sends all captured events to an HTTP server using JSON
func WithRemoteTarget() {}

// WithLogReceiver logs all events to the log
func WithLogReceiver(options *receiver.Options) Option {
	return func(config *Config) {
		config.receivers = append(config.receivers, log.NewReceiver(options))
	}
}

func WithFileReceiver(filePath string, options *receiver.Options) Option {
	return func(config *Config) {
		config.receivers = append(config.receivers, fs.NewReceiver(filePath, options))
	}
}

func WithFormat(format string) Option {
	return func(config *Config) {
		config.options.Format = format
	}
}

// WithMockServerEnabled allows to supply mocked responses in realtime from a remote client (e.g. as part of an e2e test)
// using GRPC or websocket or plain http?
func WithMockServerEnabled() {}

func defaultConfig() *Config {
	return &Config{receivers: []receiver.Receiver{}}
}

func configWithOptions(options []Option) *Config {
	config := defaultConfig()

	for _, option := range options {
		option(config)
	}

	if len(config.receivers) == 0 {
		WithLogReceiver(nil)(config)
	}

	return config
}

func Init(options ...Option) error {
	config := configWithOptions(options)
	broker := broker.NewBroker(config.receivers, config.options)

	return http.InterceptHttp(broker)
}

// Cleanup should be used in case the programm should wait for any ongoing i/o events.
// It's recommended to use `defer Cleanup()`
func Cleanup() {
	http.AwaitPendingRequests()
}
