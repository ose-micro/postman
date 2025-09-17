package bus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/postman/internal/business/email"
)

func newEmailConsumer(app email.App, bus domain.Bus, log logger.Logger) {
	// Send Event
	_ = bus.Subscribe(email.SendMailEvent, "postman", func(ctx context.Context, data any) error {
		var event email.SendCommand
		raw, err := toByte(data)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(raw, &event); err != nil {
			return fmt.Errorf("failed to unmarshal into SendMailEvent: %w", err)
		}
		if _, err = app.Create(ctx, email.CreateCommand{
			Recipient: event.Recipient,
			Sender:    event.Sender,
			Data:      event.Data,
			Template:  event.Template,
			From:      event.From,
		}); err != nil {
			return err
		}

		return nil
	})
}
