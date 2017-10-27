package airtable

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TableListCall(t *testing.T) {
	var tests = []struct {
		desc string
		call *TableListCall
		qs   string
	}{
		{
			"no options",
			newTableListCall(),
			"",
		},
		{
			"many list options",
			newTableListCall().
				Fields([]string{"foo", "bar"}).
				FilterByFormula("foo").
				MaxRecords(100).
				PageSize(10).
				SortFields([]SortField{SortField{"foo", true}, SortField{"bar", false}}).
				View("foo"),
			"fields%5B%5D=foo&fields%5B%5D=bar&filterByFormula=foo&maxRecords=100&pageSize=10&sort%5B0%5D%5Bdirection%5D=desc&sort%5B0%5D%5Bfield%5D=foo&sort%5B1%5D%5Bdirection%5D=asc&sort%5B1%5D%5Bfield%5D=bar&view=foo",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			require.Equal(t, test.qs, test.call.urlParams.Encode())
		})
	}
}

func newTableListCall() *TableListCall {
	return &TableListCall{
		client:    &Client{},
		urlParams: make(url.Values),
	}
}
