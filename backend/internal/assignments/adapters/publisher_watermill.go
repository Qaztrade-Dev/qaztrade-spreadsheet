package adapters

import (
	"context"
	"strconv"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"github.com/google/uuid"
)

type PublisherWatermill struct {
	publisher message.Publisher
	topic     string
}

var _ domain.Publisher = (*PublisherWatermill)(nil)

func NewPublisherWatermill(publisher message.Publisher, topic string) *PublisherWatermill {
	return &PublisherWatermill{
		publisher: publisher,
		topic:     topic,
	}
}

func (p *PublisherWatermill) Publish(ctx context.Context, assignmentID uint64) error {
	var (
		payload = strconv.FormatUint(assignmentID, 10)
		msg     = message.NewMessage(uuid.NewString(), []byte(payload))
	)

	if err := p.publisher.Publish(p.topic, msg); err != nil {
		return err
	}
	return nil
}
