package transmission

import (
	"context"
	"fmt"
)

type sessionStatsResponse struct {
	Result    string  `json:"result"`
	Arguments Session `json:"arguments"`
}

// https://github.com/transmission/transmission/blob/main/docs/rpc-spec.md#42-session-statistics
type Session struct {
	ActiveTorrentCount int               `json:"activeTorrentCount"`
	PausedTorrentCount int               `json:"pausedTorrentCount"`
	TorrentCount       int               `json:"torrentCount"`
	UploadSpeed        int               `json:"downloadSpeed"`
	DownloadSpeed      int               `json:"uploadSpeed"`
	CumulativeStats    SessionStatistics `json:"cumulativeStatistics"`
	CurrentStats       SessionStatistics `json:"currentStatistics"`
}

// https://github.dev/transmission/transmission/blob/6ca0ce683a5aaa8991b16ab7d93722b8861f626b/libtransmission/transmission.h#L423-L431
type SessionStatistics struct {
	Ratio           float32 `json:"ratio"`
	UploadedBytes   int     `json:"uploadedBytes"`
	DownloadedBytes int     `json:"downloadedBytes"`
	FilesAdded      int     `json:"filesAdded"`
	SessionCount    int     `json:"sessionCount"`
	SecondsActive   int     `json:"secondsActive"`
}

func (t *Client) GetSessionStats(ctx context.Context) (*Session, error) {
	var response sessionStatsResponse
	if err := t.callRPC(ctx, "session-stats", nil, &response); err != nil {
		return nil, err
	}
	if response.Result != "success" {
		return nil, fmt.Errorf("failed to get session stats: %s", response.Result)
	}
	return &response.Arguments, nil
}
