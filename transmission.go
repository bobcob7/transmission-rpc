package transmission

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Transmission struct {
	rootURL     string
	sessionID   string
	downloadDir string
	cli         *http.Client
}

func New(ctx context.Context, rootURL string) (*Transmission, error) {
	tr := &Transmission{
		rootURL: rootURL,
		cli: &http.Client{
			Timeout: time.Second * 10,
		},
	}
	if err := tr.getSessionID(ctx); err != nil {
		return nil, fmt.Errorf("failed getting session: %w", err)
	}
	return tr, nil
}

func (t *Transmission) getSessionID(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.rootURL+"/transmission/rpc", nil)
	if err != nil {
		return err
	}
	resp, err := t.cli.Do(req)
	if err != nil {
		return err
	}
	sessionID := resp.Header.Get("X-Transmission-Session-Id")
	if sessionID == "" {
		return fmt.Errorf("missing header :%#v", resp.Header)
	}
	t.sessionID = sessionID
	return nil
}

type genericRequest struct {
	Method    string      `json:"method"`
	Arguments interface{} `json:"arguments"`
	Tag       string      `json:"tag"`
}

func (t *Transmission) callRPC(ctx context.Context, requestMethod string, requestArguments, response interface{}) error {
	var resp *http.Response
	if t.sessionID == "" {
		if err := t.getSessionID(ctx); err != nil {
			return fmt.Errorf("error getting session ID: %w", err)
		}
	}
	buffer := new(bytes.Buffer)
	request := genericRequest{
		Method:    requestMethod,
		Arguments: requestArguments,
	}
	err := json.NewEncoder(buffer).Encode(request)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.rootURL+"/transmission/rpc", buffer)
	if err != nil {
		return err
	}
	for i := 0; i < 2; i++ {
		req.Header.Add("X-Transmission-Session-Id", t.sessionID)
		resp, err = t.cli.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != 409 {
			break
		}
		if err := t.getSessionID(ctx); err != nil {
			return err
		}
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return err
	}
	return nil
}
