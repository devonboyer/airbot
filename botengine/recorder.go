package botengine

import "bytes"

type ResponseRecorder struct {
	Body   *bytes.Buffer
	Status Status
}

func NewRecorder() *ResponseRecorder {
	return &ResponseRecorder{
		Body:   new(bytes.Buffer),
		Status: StatusOk,
	}
}

func (rr *ResponseRecorder) Write(buf []byte) (n int, err error) {
	if rr.Body != nil {
		rr.Body.Write(buf)
	}
	return len(buf), nil
}

func (rr *ResponseRecorder) SetStatus(s Status) {
	rr.Status = s
}
