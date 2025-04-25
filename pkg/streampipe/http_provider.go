package streampipe

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type httpPublishProviderImpl struct {
	endpoint string
}

type HttpEvent struct {
	EventID string `json:"event_id"`
	Topic   string `json:"topic"`
	Data    string `json:"data"`
}

func NewHttpPublishProvider(endpoint string) PublishProvider {
	return &httpPublishProviderImpl{
		endpoint: endpoint,
	}
}

// Close implements PublishProvider.
func (h *httpPublishProviderImpl) Close() error {
	return nil
}

// Send implements PublishProvider.
func (h *httpPublishProviderImpl) Send(topic string, event Event) error {
	var err error

	raw, err := json.Marshal(event)
	if err != nil {
		return err
	}

	hevent := &HttpEvent{
		EventID: event.EventPath(),
		Topic:   topic,
		Data:    string(raw),
	}

	sendraw, err := json.Marshal(hevent)
	if err != nil {
		return err
	}

	res, err := http.Post(h.getUrlTopic(topic), "application/json", bytes.NewBuffer(sendraw))
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		data, _ := io.ReadAll(res.Body)
		return errors.New(string(data))
	}
	return err
}

func (h *httpPublishProviderImpl) getUrlTopic(topic string) string {
	return fmt.Sprintf("%s%s", h.endpoint, topic)
}

type httpPullProviderImpl struct{}

func NewHttpPullProvider(g *gin.RouterGroup) PullProvider {
	// g.POST()

	return &httpPullProviderImpl{}
}

// Close implements PullProvider.
func (h *httpPullProviderImpl) Close() error {
	return nil
}

// Receive implements PullProvider.
func (h *httpPullProviderImpl) Receive(handler func(event PullEvent)) error {
	panic("unimplemented")
}
