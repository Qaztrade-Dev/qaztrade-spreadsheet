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

func (p *PublisherWatermill) Publish(ctx context.Context, assignmentIDs ...uint64) error {
	messages := make([]*message.Message, 0, len(assignmentIDs))
	for i := range assignmentIDs {
		var (
			assignmentID = assignmentIDs[i]
			payload      = strconv.FormatUint(assignmentID, 10)
			msg          = message.NewMessage(uuid.NewString(), []byte(payload))
		)

		messages = append(messages, msg)
	}

	if err := p.publisher.Publish(p.topic, messages...); err != nil {
		return err
	}
	return nil
}
