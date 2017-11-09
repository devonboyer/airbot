package airtable

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/net/context/ctxhttp"
)

type SortField struct {
	Field string
	Desc  bool
}

type TableHandle struct {
	client  *Client
	baseID  string
	tableID string
}

func (t *TableHandle) List() *TableListCall {
	return newTableListCall(t.client, t.baseID, t.tableID)
}

type TableListCall struct {
	client    *Client
	baseID    string
	tableID   string
	urlParams url.Values
}

func newTableListCall(client *Client, baseID, tableID string) *TableListCall {
	return &TableListCall{
		client:    client,
		baseID:    baseID,
		tableID:   tableID,
		urlParams: make(url.Values),
	}
}

func (c *TableListCall) Fields(fields []string) *TableListCall {
	for _, field := range fields {
		c.urlParams.Add("fields[]", field)
	}
	return c
}

func (c *TableListCall) FilterByFormula(formula string) *TableListCall {
	c.urlParams.Set("filterByFormula", formula)
	return c
}

func (c *TableListCall) MaxRecords(maxRecords int) *TableListCall {
	c.urlParams.Set("maxRecords", strconv.Itoa(maxRecords))
	return c
}

func (c *TableListCall) PageSize(pageSize int) *TableListCall {
	c.urlParams.Set("pageSize", strconv.Itoa(pageSize))
	return c
}

func (c *TableListCall) SortFields(sortFields []SortField) *TableListCall {
	for i, sort := range sortFields {
		c.urlParams.Add(fmt.Sprintf("sort[%d][field]", i), sort.Field)
		direction := "asc"
		if sort.Desc {
			direction = "desc"
		}
		c.urlParams.Add(fmt.Sprintf("sort[%d][direction]", i), direction)
	}
	return c
}

func (c *TableListCall) View(view string) *TableListCall {
	c.urlParams.Set("view", view)
	return c
}

func (c *TableListCall) Do(ctx context.Context, v interface{}) error {
	res, err := c.doRequest(ctx)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err := checkResponse(res); err != nil {
		return err
	}
	return json.NewDecoder(res.Body).Decode(v)
}

func (c *TableListCall) doRequest(ctx context.Context) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/%s", c.client.basePath, c.baseID, c.tableID)
	if len(c.urlParams) > 0 {
		url += "?" + c.urlParams.Encode()
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.client.apiKey))
	req.Header.Set("x-airtable-application-id", c.baseID)
	setVersionHeader(req.Header)
	if ctx == nil {
		return c.client.hc.Do(req)
	}
	return ctxhttp.Do(ctx, c.client.hc, req)
}
