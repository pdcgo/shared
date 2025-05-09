package streampipe

import "encoding/json"

type mockPubsubImpl struct{}

// Close implements PublishProvider.
func (m *mockPubsubImpl) Close() error {
	return nil
}

// Send implements PublishProvider.
func (m *mockPubsubImpl) Send(topic string, event Event) error {
	return nil
}

func NewMockPublishProvider() PublishProvider {
	return &mockPubsubImpl{}
}

type mockPullEvent struct {
	event Event
}

// Ack implements PullEvent.
func (m *mockPullEvent) Ack() {
	return
}

// Decode implements PullEvent.
func (m *mockPullEvent) Decode(v any) error {
	raw, err := json.Marshal(m.event)
	if err != nil {
		return err
	}

	return json.Unmarshal(raw, v)
}

// EventPath implements PullEvent.
func (m *mockPullEvent) EventPath() string {
	return m.event.EventPath()
}

func NewMockPullEvent(event Event) PullEvent {
	return &mockPullEvent{
		event: event,
	}
}
