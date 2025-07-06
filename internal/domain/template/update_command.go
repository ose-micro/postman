package template

import (
	"fmt"
	"strings"

	"github.com/ose-micro/cqrs"
)


type UpdateCommand struct {
	Id           string
	Content      string
	Subject      string
	Placeholders []string
}

// CommandName implements cqrs.Command.
func (u UpdateCommand) CommandName() string {
	return UPDATED_COMMAND
}

// Validate implements cqrs.Command.
func (u UpdateCommand) Validate() error {
	fields := make([]string, 0)

	if u.Id == "" {
		fields = append(fields, "id is required")
	}

	if u.Content == "" {
		fields = append(fields, "content is required")
	}

	if u.Subject == "" {
		fields = append(fields, "subject is required")
	}

	if len(fields) > 0 {
		return fmt.Errorf("%s", strings.Join(fields, ", "))
	}

	return nil
}

var _ cqrs.Command = UpdateCommand{}
