package main

type Secrets struct {
	Airtable struct {
		APIKey string
		BaseID string
	}
	Messenger struct {
		AccessToken string
		VerifyToken string
		AppSecret   string
	}
}
