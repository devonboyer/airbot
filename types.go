package airbot

import (
	"bytes"
	"fmt"
	"time"
)

type Show struct {
	ID          string    `json:"id"`
	Fields      Fields    `json:"fields"`
	CreatedTime time.Time `json:"createdTime"`
}

type Fields struct {
	Name      string `json:"Name"`
	DayOfWeek string `json:"Day of Week"`
	Status    string `json:"Status"`
}

type ShowList struct {
	Records []Show `json:"records"`
	Offset  string `json:"offset"`
}

func (sl *ShowList) String() string {
	buf := &bytes.Buffer{}
	for _, s := range sl.Records {
		fmt.Fprintln(buf, s.Fields.Name)
	}
	return buf.String()
}
