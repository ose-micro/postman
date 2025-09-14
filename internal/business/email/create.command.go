package email

import (
	"fmt"
	"strings"

	"github.com/ose-micro/cqrs"
)

type CreateCommand struct {
	Recipient string
	Sender    string
	Data      map[string]interface{}
	Template  string
	From      string
}

// CommandName implements cqrs.Command.
func (c CreateCommand) CommandName() string {
	return "postman.template.create.command"
}

// Validate implements cqrs.Command.
func (c CreateCommand) Validate() error {
	fields := make([]string, 0)

	if c.Recipient == "" {
		fields = append(fields, "recipient is required")
	}

	if c.Template == "" {
		fields = append(fields, "template is required")
	}

	if c.From == "" {
		fields = append(fields, "from is required")
	}

	if len(fields) > 0 {
		return fmt.Errorf("%s", strings.Join(fields, ", "))
	}

	return nil
}

var _ cqrs.Command = CreateCommand{}
