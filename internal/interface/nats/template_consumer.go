package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/app/template"
	domain_template "github.com/moriba-cloud/ose-postman/internal/domain/template"
	"github.com/ose-micro/cqrs/bus"
)

func newTemplateConsumer(handler domain_template.Event, bus bus.Bus) {
	// Created Event
	bus.Subscribe(domain_template.CREATED_COMMAND,
		func(ctx context.Context, data any) error {
			var event template.CreatedEvent
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
	bus.Subscribe(domain_template.UPDATED_COMMAND,
		func(ctx context.Context, data any) error {
			var event template.UpdatedEvent
			raw, err := toByte(data)
			if err != nil {
				return err
			}
			
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("failed to unmarshal into UpdateEvent: %w", err)
			}
			return handler.Updated(ctx, event)
		})

	// Delete Event
	bus.Subscribe(domain_template.DELETED_COMMAND,
		func(ctx context.Context, data any) error {
			var event template.DeletedEvent
			raw, err := toByte(data)
			if err != nil {
				return err
			}
			
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("failed to unmarshal into DeletedEvent: %w", err)
			}
			return handler.Deleted(ctx, event)
		})
}
