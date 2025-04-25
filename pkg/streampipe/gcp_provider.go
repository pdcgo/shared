package streampipe

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/pdcgo/shared/pkg/secret"
)

type gcpPubsubImpl struct {
	ctx    context.Context
	client *pubsub.Client
	topics map[string]*pubsub.Topic
}

func NewGcpPublishProvider(
	ctx context.Context,
) PublishProvider {

	projectID := secret.GetProjectID()
	if projectID == "" {
		panic("project id gcp not set")
	}
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}

	return &gcpPubsubImpl{
		ctx:    ctx,
		client: client,
		topics: map[string]*pubsub.Topic{},
	}
}

// Close implements PublishProvider.
func (g *gcpPubsubImpl) Close() error {
	return g.client.Close()
}

// Send implements PublishProvider.
func (g *gcpPubsubImpl) Send(topic string, event Event) error {
	var err error
	if g.topics[topic] == nil {
		g.topics[topic] = g.client.Topic(topic)
	}
	t := g.topics[topic]

	rawevent, err := json.Marshal(event)
	if err != nil {
		return err
	}
	res := t.Publish(g.ctx, &pubsub.Message{
		Attributes: map[string]string{
			"event_path": event.EventPath(),
		},
		Data: rawevent,
	})

	_, err = res.Get(g.ctx)
	if err != nil {
		return err
	}

	return nil
}

type gcpPullProviderImpl struct {
	ctx    context.Context
	client *pubsub.Client
	subID  string
}

func NewGcpPullProvider(
	ctx context.Context,
	subID string,
) PullProvider {
	projectID := secret.GetProjectID()
	if projectID == "" {
		panic("project id gcp not set")
	}

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}

	return &gcpPullProviderImpl{
		ctx:    ctx,
		client: client,
		subID:  subID,
	}
}

type pubEvent struct {
	msg *pubsub.Message
}

// Ack implements PullEvent.
func (p *pubEvent) Ack() {
	p.msg.Ack()
}

// Decode implements PullEvent.
func (p *pubEvent) Decode(v any) error {
	return json.Unmarshal(p.msg.Data, v)
}

// EventPath implements PullEvent.
func (p *pubEvent) EventPath() string {
	return p.msg.Attributes["event_path"]
}

// Receive implements PullProvider.
func (g *gcpPullProviderImpl) Receive(handler func(event PullEvent)) error {

	sub := g.client.Subscription(g.subID)
	err := sub.Receive(g.ctx, func(_ context.Context, msg *pubsub.Message) {
		event := pubEvent{
			msg: msg,
		}

		handler(&event)
	})
	return err
}

// Close implements PullProvider.
func (g *gcpPullProviderImpl) Close() error {
	return g.client.Close()
}
