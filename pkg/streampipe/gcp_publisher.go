package streampipe

type gcpPubsubImpl struct{}

// Send implements PublishProvider.
func (g *gcpPubsubImpl) Send(topic string, event Event) error {
	panic("unimplemented")
}

func NewGcpPublisher() PublishProvider {
	return &gcpPubsubImpl{}
}
