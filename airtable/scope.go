package airtable

type BaseScopedClient struct {
	client *Client
	BaseID string
}

func (c *Client) WithBaseScope(baseID string) *BaseScopedClient {
	return &BaseScopedClient{
		client: c,
		BaseID: baseID,
	}
}

func (c *BaseScopedClient) Table(name string) *TableHandle {
	return &TableHandle{
		client:  c.client,
		baseID:  c.BaseID,
		tableID: name,
	}
}

func (c *BaseScopedClient) WithTableScope(tableID string) *TableScopedClient {
	return c.client.WithTableScope(c.BaseID, tableID)
}

type TableScopedClient struct {
	client  *Client
	BaseID  string
	TableID string
}

func (c *Client) WithTableScope(baseID, tableID string) *TableScopedClient {
	return &TableScopedClient{
		client:  c,
		BaseID:  baseID,
		TableID: tableID,
	}
}

func (c *TableScopedClient) List() *TableListCall {
	return newTableListCall(c.client, c.BaseID, c.TableID)
}
