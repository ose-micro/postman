package email

import "time"

func (d *Domain) GetID() string {
	return d.id
}

func (d *Domain) GetRecipient() string {
	return d.recipient
}

func (d *Domain) GetData() map[string]interface{} {
	return d.data
}

func (d *Domain) GetSender() string {
	return d.sender
}

func (d *Domain) GetFrom() string {
	return d.from
}

func (d *Domain) GetSubject() string {
	return d.subject
}

func (d *Domain) GetTemplate() string {
	return d.template
}

func (d *Domain) GetMessage() string {
	return d.message
}

func (d *Domain) GetState() State {
	return d.state
}

func (d *Domain) GetCreatedAt() time.Time {
	return d.createdAt
}

func (d *Domain) GetUpdatedAt() time.Time {
	return d.updatedAt
}

func (d *Domain) SetState(state State) {
	d.state = state
}

func (d *Domain) MakePublic() Public {
	return Public{
		Id:        d.id,
		Recipient: d.recipient,
		Sender:    d.sender,
		Subject:   d.subject,
		Data:      d.data,
		Template:  d.template,
		From:      d.from,
		State:     d.state,
		Message:   d.message,
		CreatedAt: d.createdAt,
		UpdatedAt: d.updatedAt,
	}
}
