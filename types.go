package airbot

import "time"

type Show struct {
	ID     string `json:"id"`
	Fields struct {
		Name           string   `json:"Name"`
		Notes          string   `json:"Notes"`
		DayOfWeek      string   `json:"Day of Week"`
		Genres         []string `json:"Genres"`
		RunningTime    string   `json:"Running Time"`
		Status         string   `json:"Status"`
		Networks       []string `json:"Networks"`
		PersonalRating string   `json:"Personal Rating"`
	} `json:"fields"`
	CreatedTime time.Time `json:"createdTime"`
}

type ShowList struct {
	Records []Show `json:"records"`
	Offset  string `json:"offset"`
}
