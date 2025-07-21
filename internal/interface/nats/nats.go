package nats

import (
	"encoding/json"
	"fmt"

	"github.com/moriba-cloud/ose-postman/internal/app"
	"github.com/moriba-cloud/ose-postman/internal/events"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/cqrs/bus"
)

func InvokeConsumers(events events.Events, app app.Apps, log logger.Logger, bus bus.Bus) {
	newTemplateConsumer(events.Template, bus)
	newEmailConsumer(events.Email, app.Email, bus)
}

func toByte(data interface{}) ([]byte, error) {
    mapData, ok := data.(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("invalid message format")
    }

    // Marshal it to JSON
    raw, err := json.Marshal(mapData)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal map: %w", err)
    }

	return raw, nil
}