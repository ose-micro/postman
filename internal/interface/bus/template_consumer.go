package bus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/domain/template"
	"github.com/ose-micro/cqrs/bus"
)

func newTemplateConsumer(handler template.Event, bus bus.Bus) {
	// Created Event
	bus.Subscribe(template.CREATED_COMMAND,
		QUEUE,
		func(ctx context.Context, data any) error {
			var event template.DomainEvent
			raw, err := toByte(data)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("failed to unmarshal into DefaultEvent: %w", err)
			}
			return handler.Created(ctx, event)
		})

	// Updated Event
	bus.Subscribe(template.UPDATED_COMMAND,
		QUEUE,
		func(ctx context.Context, data any) error {
			var event template.DomainEvent
			raw, err := toByte(data)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("failed to unmarshal into DefaultEvent: %w", err)
			}
			return handler.Updated(ctx, event)
		})

	// Delete Event
	bus.Subscribe(template.DELETED_COMMAND,
		QUEUE,
		func(ctx context.Context, data any) error {
			var event template.DomainEvent
			raw, err := toByte(data)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("failed to unmarshal into DefaultEvent: %w", err)
			}
			return handler.Deleted(ctx, event)
		})
}
