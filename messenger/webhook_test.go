package messenger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_WebhookHandler(t *testing.T) {
	tests := []struct {
		name, jsonResponse string
	}{
		{
			"unmarshal message event",
			"{\"object\":\"page\",\"entry\":[{\"id\":\"233235520541490\",\"time\":1509072531432,\"messaging\":[{\"sender\":{\"id\":\"1687448547993463\"},\"recipient\":{\"id\":\"233235520541490\"},\"timestamp\":1508990647040,\"message\":{\"mid\":\"mid.$cAADUIEhRdW5liB9DAFfVtoBCMqlf\",\"seq\":63713,\"text\":\"foo\"}}]}]}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			buf.WriteString(tt.jsonResponse)

			req, err := http.NewRequest("POST", "/webhook", buf)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			eventsCh := make(chan *Event, 1)

			handler := WebhookHandler{skipVerifySignature: true, eventsCh: eventsCh}
			handler.ServeHTTP(rr, req)

			timeout := time.NewTimer(1 * time.Second)
			select {
			case <-eventsCh:
				require.Equal(t, http.StatusOK, rr.Code)
			case <-timeout.C:
				require.FailNow(t, "Timeout waiting for event")
			}
		})
	}
}
