package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/app/email"
	domain_email "github.com/moriba-cloud/ose-postman/internal/domain/email"
	"github.com/ose-micro/cqrs/bus"
)

func newEmailConsumer(handler domain_email.Event, bus bus.Bus) {
	// Created Event
	bus.Subscribe(domain_email.CREATED_COMMAND,
		func(ctx context.Context, data any) error {
			var event email.CreatedEvent
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
	bus.Subscribe(domain_email.UPDATED_COMMAND,
		func(ctx context.Context, data any) error {
			var event email.UpdatedEvent
			raw, err := toByte(data)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("failed to unmarshal into UpdatedEvent: %w", err)
			}
			return handler.Updated(ctx, event)
		})

	// Delete Event
	bus.Subscribe(domain_email.DELETED_COMMAND,
		func(ctx context.Context, data any) error {
			var event email.DeletedEvent
			raw, err := toByte(data)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("failed to unmarshal into DeletedEvent: %w", err)
			}
			return handler.Deleted(ctx, event)
		})

	// Send Event
	bus.Subscribe(domain_email.SEND_MAIL_EVENT,
		func(ctx context.Context, data any) error {
			var event email.SendMailEvent
			raw, err := toByte(data)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("failed to unmarshal into SendMailEvent: %w", err)
			}
			return handler.SendMail(ctx, event)
		})
}
