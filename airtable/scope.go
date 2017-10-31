package airtable

type TableScopedClient struct {
	*Client
	BaseID  string
	TableID string
}

func (c *Client) WithTableScope(baseID, tableID string) *TableScopedClient {
	return &TableScopedClient{
		Client:  c,
		BaseID:  baseID,
		TableID: tableID,
	}
}

func (c *TableScopedClient) List() *TableListCall {
	return newTableListCall(c.Client, c.BaseID, c.TableID)
}
