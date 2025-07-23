package bus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/domain/email"
	"github.com/ose-micro/cqrs/bus"
)

func newEmailConsumer(handler email.Event, app email.App, bus bus.Bus) {
	// Created Event
	bus.Subscribe(email.CREATED_COMMAND,
		QUEUE,
		func(ctx context.Context, data any) error {
			var event email.DomainEvent
			raw, err := toByte(data)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("failed to unmarshal into CreatedEvent: %w", err)
			}
			return handler.Created(ctx, event)
		})

	// Updated Event
	bus.Subscribe(email.UPDATED_COMMAND,
		QUEUE,
		func(ctx context.Context, data any) error {
			var event email.DomainEvent
			raw, err := toByte(data)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("failed to unmarshal into UpdatedEvent: %w", err)
			}
			return handler.Updated(ctx, event)
		})

	// Send Event
	bus.Subscribe(email.SEND_MAIL_EVENT,
		QUEUE,
		func(ctx context.Context, data any) error {
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
