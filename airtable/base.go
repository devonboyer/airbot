package airtable

type BaseClient struct {
	client *Client
	baseID string
}

func (c *Client) Base(baseID string) *BaseClient {
	return &BaseClient{
		client: c,
		baseID: baseID,
	}
}

func (c *BaseClient) Table(name string) *TableClient {
	return &TableClient{
		client:  c.client,
		baseID:  c.baseID,
		tableID: name,
	}
}
