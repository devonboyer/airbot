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

type TableClient struct {
	client  *Client
	baseID  string
	tableID string
}

func (c *TableClient) List(ctx context.Context, opts *ListOptions, v interface{}) error {
	res, err := c.doRequest(ctx, opts.urlParams())
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err := checkResponse(res); err != nil {
		return err
	}
	return json.NewDecoder(res.Body).Decode(v)
}

func (c *TableClient) doRequest(ctx context.Context, urlParams url.Values) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/%s", c.client.basePath, c.baseID, c.tableID)
	if len(urlParams) > 0 {
		url += "?" + urlParams.Encode()
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

type ListOptions struct {
	Fields          []string
	SortFields      []SortField
	FilterByFormula string
	MaxRecords      int
	PageSize        int
	View            string
}

type SortField struct {
	Field string
	Desc  bool
}

func (o *ListOptions) urlParams() url.Values {
	urlParams := make(url.Values)
	for _, field := range o.Fields {
		urlParams.Add("fields[]", field)
	}
	for i, sort := range o.SortFields {
		urlParams.Add(fmt.Sprintf("sort[%d][field]", i), sort.Field)
		direction := "asc"
		if sort.Desc {
			direction = "desc"
		}
		urlParams.Add(fmt.Sprintf("sort[%d][direction]", i), direction)
	}
	if o.FilterByFormula != "" {
		urlParams.Set("filterByFormula", o.FilterByFormula)
	}
	if o.MaxRecords != 0 {
		urlParams.Set("maxRecords", strconv.Itoa(o.MaxRecords))
	}
	if o.PageSize != 0 {
		urlParams.Set("pageSize", strconv.Itoa(o.PageSize))
	}
	if o.View != "" {
		urlParams.Set("view", o.View)
	}
	return urlParams
}
