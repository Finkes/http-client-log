package http

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Finkes/http-client-log/pkg/broker"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type InterceptedTransport struct {
	Transport http.RoundTripper
	broker    *broker.Broker
}

var pendingRequests sync.WaitGroup

// InterceptHttp Init starts collecting all http(s) communication and sends them to a pactly server
func InterceptHttp(broker *broker.Broker) error {
	interceptedTransport, err := NewInterceptedTransport(broker)
	if err != nil {
		return err
	}
	http.DefaultTransport = interceptedTransport
	return nil
}

// AwaitPendingRequests will wait for all pending http requests to finish
func AwaitPendingRequests() {
	pendingRequests.Wait()
}

func NewInterceptedTransport(broker *broker.Broker) (*InterceptedTransport, error) {
	return &InterceptedTransport{
		Transport: http.DefaultTransport,
		broker:    broker,
	}, nil
}

// RoundTrip is the core part of this module and implements http.RoundTripper.
func (t *InterceptedTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	var requestBody []byte
	var responseBody []byte
	var err error
	ctx := context.Background()
	request = request.WithContext(ctx)
	requestTime := time.Now()

	response, err := t.transport().RoundTrip(request)
	if err != nil {
		return response, err
	}

	responseTime := time.Now()

	if request.Body != nil {
		requestBodyRaw, _ := request.GetBody()
		requestBody, err = ioutil.ReadAll(requestBodyRaw)
		if err != nil {
			println(err)
		}
	}

	if response.Body != nil {
		responseBody, _ = ioutil.ReadAll(response.Body)
		err := response.Body.Close()
		if err != nil {
			println(err)
		}
		r := bytes.NewReader(responseBody)
		response.Body = io.NopCloser(r)
	}

	pendingRequests.Add(1)
	go func() {
		defer pendingRequests.Done()
		err := t.captureEvent(*request, requestBody, *response, responseBody, requestTime, responseTime)
		if err != nil {
			fmt.Printf("Failed to capture http event: %v\n", err)
		}
	}()

	return response, err
}

func (t *InterceptedTransport) captureEvent(request http.Request, requestBody []byte, response http.Response, responseBody []byte, requestTime time.Time, responseTime time.Time) error {
	requestBodyString := string(requestBody)
	responseBodyString := string(responseBody)

	httpEvent := Event{
		UUID: uuid.New().String(),
		Time: requestTime.UTC(),
		//Component:       t.PactlyComponent, todo
		Protocol:        request.URL.Scheme,
		ProtocolVersion: response.Request.Proto,
		Request: EventRequest{
			Header:         normalizeHeader(request.Header),
			Body:           requestBodyString,
			Host:           request.Host,
			Method:         request.Method,
			Path:           request.URL.Path,
			Query:          request.URL.RawQuery,
			BodySize:       len(requestBodyString),
			BodySizePretty: broker.FormatFileSize(len(requestBodyString)),
		},
		Response: EventResponse{
			Header:         normalizeHeader(response.Header),
			Body:           responseBodyString,
			BodySize:       len(responseBodyString),
			StatusCode:     response.StatusCode,
			BodySizePretty: broker.FormatFileSize(len(responseBodyString)),
		},
		Duration: responseTime.Sub(requestTime).Seconds(),
	}
	t.broker.CaptureEvent(&httpEvent)

	return nil
}

func (t *InterceptedTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

func normalizeHeader(header http.Header) Header {
	result := map[string]string{}
	for key, value := range header {
		result[key] = value[0]
	}
	return result
}
