package airbot

type Secrets struct {
	Airtable struct {
		APIKey string `json:"api_key"`
		BaseID string `json:"base_id"`
	} `json:"airtable"`
	Messenger struct {
		AccessToken string `json:"access_token"`
		VerifyToken string `json:"verify_token"`
		AppSecret   string `json:"app_secret"`
	} `json:"messenger"`
}
