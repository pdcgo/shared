package streampipe

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
