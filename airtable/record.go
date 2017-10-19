package airtable

import "time"

type Record struct {
	ID          string                 `json:"id"`
	Fields      map[string]interface{} `json:"fields"`
	CreatedTime time.Time              `json:"createdTime"`
}

type RecordList struct {
	Records []Record `json:"records"`
	Offset  string   `json:"offset"`
}
