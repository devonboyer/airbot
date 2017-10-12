package main

import "time"

type Show struct {
	ID     string `json:"id"`
	Fields struct {
		Name        string   `json:"Name"`
		DayOfWeek   string   `json:"Day of Week"`
		Genres      []string `json:"Genres"`
		RunningTime string   `json:"Running Time"`
		Networks    []string `json:"Networks"`
	} `json:"fields"`
	CreatedTime time.Time `json:"createdTime"`
}

type ShowList struct {
	Records []Show `json:"records"`
	Offset  string `json:"offset"`
}
