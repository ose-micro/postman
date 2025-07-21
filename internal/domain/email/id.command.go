package email

import (
	"fmt"
	"strings"

	"github.com/ose-micro/cqrs"
)

type IdCommand struct {
	Id string
}

// CommandName implements cqrs.Command.
func (c IdCommand) CommandName() string {
	return UPDATED_COMMAND
}

// Validate implements cqrs.Command.
func (c IdCommand) Validate() error {
	fields := make([]string, 0)

	if c.Id == "" {
		fields = append(fields, "id is required")
	}

	if len(fields) > 0 {
		return fmt.Errorf("%s", strings.Join(fields, ", "))
	}

	return nil
}

var _ cqrs.Command = IdCommand{}
