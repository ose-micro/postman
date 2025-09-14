package template

import (
	"fmt"
	"strings"

	"github.com/ose-micro/cqrs"
)

type DeleteCommand struct {
	Id string
}

// CommandName implements cqrs.Command.
func (d DeleteCommand) CommandName() string {
	return "postman.template.delete.command"
}

// Validate implements cqrs.Command.
func (d DeleteCommand) Validate() error {
	fields := make([]string, 0)

	if d.Id == "" {
		fields = append(fields, "id is required")
	}

	if len(fields) > 0 {
		return fmt.Errorf("%s", strings.Join(fields, ", "))
	}

	return nil
}

var _ cqrs.Command = DeleteCommand{}
