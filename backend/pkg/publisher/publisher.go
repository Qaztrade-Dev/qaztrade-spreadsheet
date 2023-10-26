package publisher

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type Publisher interface {
	Publish(ctx context.Context, payloads ...[]byte) error
}

type PublisherClient struct {
	publisher message.Publisher
	topic     string
}

var _ Publisher = (*PublisherClient)(nil)

func NewPublisherClient(publisher message.Publisher, topic string) *PublisherClient {
	return &PublisherClient{
		publisher: publisher,
		topic:     topic,
	}
}

func (p *PublisherClient) Publish(ctx context.Context, payloads ...[]byte) error {
	messages := make([]*message.Message, 0, len(payloads))
	for i := range payloads {
		var (
			payload = payloads[i]
			msg     = message.NewMessage(uuid.NewString(), payload)
		)

		messages = append(messages, msg)
	}

	if err := p.publisher.Publish(p.topic, messages...); err != nil {
		return err
	}
	return nil
}
