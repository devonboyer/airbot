package messenger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type nopEventHandler struct{}

func (n nopEventHandler) HandleEvent(_ *WebhookEvent) {}

func Test_WebhookHandler(t *testing.T) {
	client := &Client{logger: nopLogger, skipVerify: true}

	tests := []struct {
		name, jsonResponse string
	}{
		{
			"message event",
			"{\"object\":\"page\",\"entry\":[{\"id\":\"233235520541490\",\"time\":1509072531432,\"messaging\":[{\"sender\":{\"id\":\"1687448547993463\"},\"recipient\":{\"id\":\"233235520541490\"},\"timestamp\":1508990647040,\"message\":{\"mid\":\"mid.$cAADUIEhRdW5liB9DAFfVtoBCMqlf\",\"seq\":63713,\"text\":\"foo\"}}]}]}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			buf.WriteString(test.jsonResponse)

			req, err := http.NewRequest("POST", "/webhook", buf)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(client.WebhookHandler(nopEventHandler{}))
			handler.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)
		})
	}
}
