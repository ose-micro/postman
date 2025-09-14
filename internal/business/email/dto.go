package email

func (d *Domain) Recipient() string {
	return d.recipient
}

func (d *Domain) Data() map[string]interface{} {
	return d.data
}

func (d *Domain) Sender() string {
	return d.sender
}

func (d *Domain) From() string {
	return d.from
}

func (d *Domain) Subject() string {
	return d.subject
}

func (d *Domain) Template() string {
	return d.template
}

func (d *Domain) Message() string {
	return d.message
}

func (d *Domain) State() State {
	return d.state
}

func (d *Domain) SetState(state State) {
	d.state = state
}

func (d *Domain) Public() Public {
	return Public{
		Id:        d.ID(),
		Recipient: d.recipient,
		Sender:    d.sender,
		Subject:   d.subject,
		Data:      d.data,
		Template:  d.template,
		From:      d.from,
		State:     d.state,
		Message:   d.message,
		CreatedAt: d.CreatedAt(),
		UpdatedAt: d.UpdatedAt(),
	}
}
