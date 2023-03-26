package transmission

import (
	"context"
	"fmt"
	"path"
)

type addTransmissionResponse struct {
	Result    string                      `json:"result"`
	Arguments addTransmissionResponseArgs `json:"arguments"`
	Tag       string                      `json:"tag"`
}

type addTransmissionRequestArgs struct {
	Paused      string `json:"paused"`
	DownloadDir string `json:"download-dir"`
	Filename    string `json:"filename"`
}

type addTransmissionResponseArgs struct {
	TorrentAdded struct {
		HashString string `json:"hashString"`
		ID         int    `json:"id"`
		Name       string `json:"name"`
	} `json:"torrent-added"`
	TorrentDuplicate struct {
		HashString string `json:"hashString"`
		ID         int    `json:"id"`
		Name       string `json:"name"`
	} `json:"torrent-duplicate"`
	Tag string `json:"tag"`
}

type AddMagnetLinkOption func(*addTransmissionRequestArgs)

func DownloadDirOption(absolutePath string) AddMagnetLinkOption {
	return func(req *addTransmissionRequestArgs) {
		req.DownloadDir = absolutePath
	}
}

func DownloadSubDirOption(subDir string) AddMagnetLinkOption {
	return func(req *addTransmissionRequestArgs) {
		req.DownloadDir = path.Join(req.DownloadDir, subDir)
	}
}

func (t *Client) AddMagnetLink(ctx context.Context, link string, opts ...AddMagnetLinkOption) (int, error) {
	var response addTransmissionResponse
	req := addTransmissionRequestArgs{
		Filename:    link,
		DownloadDir: t.DownloadDir,
	}
	for _, opt := range opts {
		opt(&req)
	}
	if err := t.callRPC(ctx, "torrent-add", &req, &response); err != nil {
		return 0, err
	}
	if response.Result == "success" {
		return response.Arguments.TorrentAdded.ID, nil
	}
	if response.Result == "duplicate torrent" {
		return response.Arguments.TorrentDuplicate.ID, nil
	}
	return 0, fmt.Errorf("failed to add magnet link: %s", response.Result)
}
