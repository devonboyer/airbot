package airtable

type BaseHandle struct {
	client *Client
	baseID string
}

func (c *Client) Base(baseID string) *BaseHandle {
	return &BaseHandle{
		client: c,
		baseID: baseID,
	}
}

func (b *BaseHandle) Table(name string) *TableHandle {
	return &TableHandle{
		client: b.client,
		baseID: b.baseID,
		name:   name,
	}
}
