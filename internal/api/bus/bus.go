package bus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ose-micro/core/domain"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/postman/internal/app"
	"github.com/ose-micro/postman/internal/business/email"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func InvokeConsumers(lc fx.Lifecycle, app app.Apps, log logger.Logger, bus domain.Bus) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				eventList := []string{email.SendMailEvent}
				err := bus.EnsureStream("POSTMAN", eventList...)
				if err != nil {
					log.Fatal("nats stream failed", zap.Error(err))
				}

				newEmailConsumer(app.Email, bus, log)
			}()
			return nil
		},
	})
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
