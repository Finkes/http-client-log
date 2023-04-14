package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Finkes/http-client-log/pkg/receiver"
	"gopkg.in/yaml.v3"
	"io"
	"time"
)

type Event struct {
	UUID            string        `json:"uuid" bson:"uuid"`
	Component       string        `json:"component" bson:"component"`
	Protocol        string        `json:"protocol" bson:"protocol"`
	ProtocolVersion string        `json:"protocolVersion" bson:"protocolVersion"`
	Request         EventRequest  `json:"request" bson:"request"`
	Response        EventResponse `json:"response" bson:"response"`
	Duration        float64       `json:"duration" bson:"duration"`
	Time            time.Time     `json:"time" bson:"time"`
}

type EventRequest struct {
	Method         string `json:"method" bson:"method"`
	Host           string `json:"host" bson:"host"`
	Path           string `json:"path" bson:"path"`
	Header         Header `json:"header" bson:"header"`
	Body           string `json:"body" bson:"body"`
	Query          string `json:"query" bson:"query"`
	BodySize       int    `json:"bodySize" bson:"bodySize"`             // contentLength is not reliable, measure body by ourselves
	BodySizePretty string `json:"bodySizePretty" bson:"bodySizePretty"` // contentLength is not reliable, measure body by ourselves

	// todo: should we track basic auth in URL? Might be missing in request header..
}

type EventResponse struct {
	Header         Header `json:"header" bson:"header"`
	Body           string `json:"body" bson:"body"`
	StatusCode     int    `json:"statusCode" bson:"statusCode"`
	BodySize       int    `json:"bodySize" bson:"bodySize"`             // contentLength is not reliable, measure body by ourselves
	BodySizePretty string `json:"bodySizePretty" bson:"bodySizePretty"` // contentLength is not reliable, measure body by ourselves
}

type Header map[string]string

func (e *Event) Url() string {
	baseUrl := fmt.Sprintf("%v://%v%v", e.Protocol, e.Request.Host, e.Request.Path)
	if e.Request.Query != "" {
		baseUrl = fmt.Sprintf("%v?%v", baseUrl, e.Request.Query)
	}
	return baseUrl
}

func (e *Event) String(options *receiver.Options) string {
	if options.Format == receiver.FormatSummary {
		return e.ShortString()
	} else if options.Format == receiver.FormatYaml {
		buffer, _ := yaml.Marshal(e)
		return string(buffer)
	} else if options.Format == receiver.FormatJSON {
		buffer, _ := json.MarshalIndent(e, "", "  ")
		return string(buffer)
	}
	return ""
}

func (e *Event) ShortString() string {
	return fmt.Sprintf("[%v] %v\t%v req:%v res:%v %v %v",
		e.Response.StatusCode,
		e.Request.Method,
		clamp(fmt.Sprintf("%v://%v%v%v", e.Protocol, e.Request.Host, e.Request.Path, e.Request.Query), 70, "left"),
		clamp(e.Request.BodySizePretty, 9, "left"),
		clamp(e.Response.BodySizePretty, 9, "left"),
		clamp(fmt.Sprint(time.Duration(e.Duration*float64(time.Second))), 20, "left"),
		e.UUID)
}

func (e *Event) Type() string {
	return "http"
}

func (e *Event) Serialize() (io.Reader, error) {
	payload, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(payload), nil
}

func clamp(text string, maxSize int, align string) string {
	if len(text) == maxSize {
		return text
	}
	if len(text) < maxSize {
		return fill(text, maxSize, align)
	}
	return cut(text, maxSize)
}

func fill(text string, maxSize int, align string) string {
	textLength := len(text)
	for i := 0; i < maxSize-textLength; i++ {
		if align == "left" {
			text += " "
		} else if align == "right" {
			text = " " + text
		}
	}
	return text
}

func cut(text string, maxSize int) string {
	return string([]rune(text)[:maxSize-3]) + "..."
}
