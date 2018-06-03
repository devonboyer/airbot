package airtable

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListOptions(t *testing.T) {
	var tests = []struct {
		desc string
		opts *ListOptions
		qs   string
	}{
		{
			"no options",
			&ListOptions{},
			"",
		},
		{
			"many list options",
			&ListOptions{
				Fields:          []string{"foo", "bar"},
				FilterByFormula: "foo",
				MaxRecords:      100,
				PageSize:        10,
				SortFields:      []SortField{SortField{"foo", true}, SortField{"bar", false}},
				View:            "bar",
			},
			"fields%5B%5D=foo&fields%5B%5D=bar&filterByFormula=foo&maxRecords=100&pageSize=10&sort%5B0%5D%5Bdirection%5D=desc&sort%5B0%5D%5Bfield%5D=foo&sort%5B1%5D%5Bdirection%5D=asc&sort%5B1%5D%5Bfield%5D=bar&view=bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			require.Equal(t, tt.qs, tt.opts.urlParams().Encode())
		})
	}
}
