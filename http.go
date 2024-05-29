package truecoach

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type JSONDate time.Time

func (j *JSONDate) UnmarshalJSON(b []byte) error {
	if b == nil || bytes.Equal(b, []byte("null")) {
		return nil
	}

	s := strings.Trim(string(b), "\"")

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	*j = JSONDate(t)

	return nil
}

type BaseResponse struct {
	Page       int `json:"page"`
	TotalPages int `json:"total_pages"`
	PerPage    int `json:"per_page"`
	TotalCount int `json:"total_count"`
}

type errorResponse struct {
	Message string `json:"error"`
}

func (e *errorResponse) String() string {
	if e.Valid() {
		return ""
	}

	return e.Message
}

func (e errorResponse) Error() string {
	return e.String()
}

func (e *errorResponse) Valid() bool {
	return e.Message == ""
}

type transport struct {
	token string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header["Authorization"] = []string{"Bearer " + t.token}
	req.Header["Role"] = []string{"Trainer"}

	return http.DefaultTransport.RoundTrip(req)
}

func (rc *Service) get(path string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, rc.origin+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}

	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed getting resource: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %w", err)
	}

	errorResp := errorResponse{}

	err = json.Unmarshal(body, &errorResp)
	if err == nil && !errorResp.Valid() {
		return nil, errorResp
	}

	return resp, nil
}

// @todo Implement write operations.
// nolint:unused
func (rc *Service) post(path string, payload interface{}) (*http.Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, rc.origin+path, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed posting resource: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %w", err)
	}

	errorResp := errorResponse{}

	err = json.Unmarshal(body, &errorResp)
	if err == nil && !errorResp.Valid() {
		return nil, fmt.Errorf("received error message: %w", errorResp)
	}

	return resp, nil
}
