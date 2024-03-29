package transmission

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	rootURL     string
	sessionID   string
	DownloadDir string
	cli         *http.Client
	sessionInfo map[string]interface{}
}

func New(ctx context.Context, rootURL string) (*Client, error) {
	tr := &Client{
		rootURL: rootURL,
		cli: &http.Client{
			Timeout: time.Second * 10,
		},
	}
	err := tr.getSessionID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting session: %w", err)
	}
	if tr.sessionInfo, err = tr.GetSession(ctx); err != nil {
		return nil, fmt.Errorf("failed getting session info: %w", err)
	}
	if downloadDir, ok := tr.sessionInfo["download-dir"]; ok {
		tr.DownloadDir = downloadDir.(string)
	}
	return tr, nil
}

func (t *Client) getSessionID(ctx context.Context) error {
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

func (t *Client) callRPC(ctx context.Context, requestMethod string, requestArguments, response interface{}) error {
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
	if err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.rootURL+"/transmission/rpc", buffer)
	if err != nil {
		return fmt.Errorf("failed to create new http request: %w", err)
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
