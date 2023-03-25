package transmission

import (
	"context"
	"fmt"
)

type getSessionRequestArgs struct {
	SessionID string   `json:"session-id"`
	Fields    []string `json:"fields"`
}

type getSessionResponse struct {
	Result    string                 `json:"result"`
	Arguments map[string]interface{} `json:"arguments"`
}

func (t *Client) GetSession(ctx context.Context) (map[string]interface{}, error) {
	var response getSessionResponse
	req := getSessionRequestArgs{
		SessionID: t.sessionID,
	}
	if err := t.callRPC(ctx, "session-get", &req, &response); err != nil {
		return nil, err
	}
	if response.Result != "success" {
		return nil, fmt.Errorf("failed to get session stats: %s", response.Result)
	}
	return response.Arguments, nil
}
