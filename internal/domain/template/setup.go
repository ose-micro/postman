package template

import (
	"fmt"
	"strings"

	"github.com/moriba-cloud/ose-postman/internal/common"
	"github.com/ose-micro/core/timestamp"
	"github.com/ose-micro/core/utils"
)

type roleDomain struct {
	timestamp timestamp.Timestamp
}

// Load implements IDomain.
func (r roleDomain) Existing(param Params) (*Domain, error) {
	fields := make([]string, 0)

	if param.Id == "" {
		fields = append(fields, "id is required")
	}

	if param.Subject == "" {
		fields = append(fields, "subject is required")
	}

	if param.Content == "" {
		fields = append(fields, "content is required")
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
		timestamp:    r.timestamp,
		id:           param.Id,
		subject:      param.Subject,
		content:      param.Content,
		placeholders: param.Placeholders,
		createdAt:    param.CreatedAt,
		updatedAt:    param.UpdatedAt,
	}, nil
}

// New implements IDomain.
func (r roleDomain) New(param Params) (*Domain, error) {
	fields := make([]string, 0)

	if param.Subject == "" {
		fields = append(fields, "subject is required")
	}

	if param.Content == "" {
		fields = append(fields, "content is required")
	}

	if len(fields) > 0 {
		return nil, fmt.Errorf("%s", strings.Join(fields, ", "))
	}

	return &Domain{
		timestamp:    r.timestamp,
		id:           utils.GenerateUUID(),
		subject:      param.Subject,
		content:      param.Content,
		placeholders: param.Placeholders,
		createdAt:    r.timestamp.Now(),
		updatedAt:    r.timestamp.Now(),
	}, nil
}

func NewTemplateDomain(timestamp timestamp.Timestamp) common.IDomain[Domain, Params] {
	return &roleDomain{
		timestamp: timestamp,
	}
}
