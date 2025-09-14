package template

import (
	"fmt"
	"strings"

	"github.com/ose-micro/cqrs"
)

type CreateCommand struct {
	Content      string
	Subject      string
	Placeholders []string
}

// CommandName implements cqrs.Command.
func (c CreateCommand) CommandName() string {
	return "postman.template.create.command"
}

// Validate implements cqrs.Command.
func (c CreateCommand) Validate() error {
	fields := make([]string, 0)

	if c.Content == "" {
		fields = append(fields, "content is required")
	}

	if c.Subject == "" {
		fields = append(fields, "subject is required")
	}

	if len(fields) > 0 {
		return fmt.Errorf("%s", strings.Join(fields, ", "))
	}

	return nil
}

var _ cqrs.Command = CreateCommand{}
