package email

import (
	"fmt"
	"strings"

	"github.com/moriba-cloud/ose-postman/internal/common"
	"github.com/ose-micro/core/timestamp"
	"github.com/ose-micro/core/utils"
)

type emailDomain struct {
	timestamp timestamp.Timestamp
}

// Load implements IDomain.
func (e emailDomain) Existing(param Params) (*Domain, error) {
	fields := make([]string, 0)

	if param.From == "" {
		fields = append(fields, "from is required")
	}

	if param.Id == "" {
		fields = append(fields, "id is required")
	}

	if param.Message == "" {
		fields = append(fields, "message is required")
	}

	if param.State == "" {
		fields = append(fields, "state is required")
	}

	if param.Recipient == "" {
		fields = append(fields, "recipient is required")
	}

	if param.Sender == "" {
		fields = append(fields, "sender is required")
	}

	if param.Subject == "" {
		fields = append(fields, "subject is required")
	}

	if param.Template == "" {
		fields = append(fields, "template is required")
	}

	if param.CreatedAt.IsZero() {
		fields = append(fields, "createdAt is required")
	}

	if param.UpdatedAt.IsZero() {
		fields = append(fields, "updatedAt is required")
	}

	if len(fields) > 0 {
		return nil, fmt.Errorf("%s", strings.Join(fields, ", "))
	}

	return &Domain{
		timestamp: e.timestamp,
		id:        param.Id,
		recipient: param.Recipient,
		sender:    param.Sender,
		from:      param.From,
		subject:   param.Subject,
		template:  param.Template,
		data:      param.Data,
		message:   param.Message,
		state:     param.State,
		createdAt: param.CreatedAt,
		updatedAt: param.UpdatedAt,
	}, nil
}

// New implements IDomain.
func (e emailDomain) New(param Params) (*Domain, error) {
	fields := make([]string, 0)

	if param.Message == "" {
		fields = append(fields, "message is required")
	}

	if param.State == "" {
		fields = append(fields, "state is required")
	}

	if param.Recipient == "" {
		fields = append(fields, "recipient is required")
	}

	if param.Sender == "" {
		fields = append(fields, "sender is required")
	}

	if param.Subject == "" {
		fields = append(fields, "subject is required")
	}

	if param.Template == "" {
		fields = append(fields, "template is required")
	}

	if len(fields) > 0 {
		return nil, fmt.Errorf("%s", strings.Join(fields, ", "))
	}

	return &Domain{
		timestamp: e.timestamp,
		id:        utils.GenerateUUID(),
		recipient: param.Recipient,
		sender:    param.Sender,
		from:      param.From,
		subject:   param.Subject,
		template:  param.Template,
		data:      param.Data,
		message:   param.Message,
		state:     param.State,
		createdAt: e.timestamp.Now(),
		updatedAt: e.timestamp.Now(),
	}, nil
}

func NewEmailDomain(timestamp timestamp.Timestamp) common.IDomain[Domain, Params] {
	return &emailDomain{
		timestamp: timestamp,
	}
}
