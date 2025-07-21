package email

import (
	"fmt"
	"strings"

	"github.com/ose-micro/cqrs"
)

type SendCommand struct {
	Recipient string                 `json:"recipient"`
	Sender    string                `json:"sender"`
	Subject   string                 `json:"subject"`
	Data      map[string]interface{} `json:"data"`
	Template  string                 `json:"template"`
	From      string                 `json:"from"`
	Message   string                 `json:"message"`
}

// CommandName implements cqrs.Command.
func (c SendCommand) CommandName() string {
	return CREATED_COMMAND
}

// Validate implements cqrs.Command.
func (c SendCommand) Validate() error {
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

var _ cqrs.Command = SendCommand{}
