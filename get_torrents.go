package transmission

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

var torrentFields []string

func getJSONTags(v reflect.Type) []string {
	fields := make([]string, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Type.Kind() == reflect.Struct && f.Anonymous {
			fields = append(fields, getJSONTags(f.Type)...)
		}
		jsonTag, ok := f.Tag.Lookup("json")
		if !ok {
			continue
		}
		jsonTagElements := strings.SplitN(jsonTag, ",", 1)
		if jsonTagElements[0] == "" {
			continue
		}
		fields = append(fields, jsonTagElements[0])
	}
	return fields
}

func init() {
	v := reflect.ValueOf(Torrent{})
	torrentFields = getJSONTags(v.Type())
}

type TorrentFile struct {
	BytesCompleted int    `json:"bytesCompleted"`
	Length         int    `json:"length"`
	Name           string `json:"name"`
}

type TorrentFileStats struct {
	BytesCompleted int  `json:"bytesCompleted"`
	Wanted         bool `json:"wanted"`
	Priority       int  `json:"priority"`
}

type TorrentPeer struct {
	Address            string  `json:"address"`
	ClientName         string  `json:"clientName"`
	ClientIsChoked     bool    `json:"clientIsChoked"`
	ClientIsInterested bool    `json:"clientIsInterested"`
	FlagStr            string  `json:"flagStr"`
	IsDownloadingFrom  bool    `json:"isDownloadingFrom"`
	IsEncrypted        bool    `json:"isEncrypted"`
	IsIncoming         bool    `json:"isIncoming"`
	IsUploadingTo      bool    `json:"isUploadingTo"`
	IsUTP              bool    `json:"isUTP"`
	PeerIsChoked       bool    `json:"peerIsChoked"`
	PeerIsInterested   bool    `json:"peerIsInterested"`
	Port               int     `json:"port"`
	Progress           float64 `json:"progress"`
	RateToClient       int     `json:"rateToClient"`
	RateToPeer         int     `json:"rateToPeer"`
}

type TorrentFromPeer struct {
	// TR_PEER_FROM_INCOMING = 0, /* connections made to the listening port */
	// TR_PEER_FROM_LPD, /* peers found by local announcements */
	// TR_PEER_FROM_TRACKER, /* peers found from a tracker */
	// TR_PEER_FROM_DHT, /* peers found from the DHT */
	// TR_PEER_FROM_PEX, /* peers found from PEX */
	// TR_PEER_FROM_RESUME, /* peers found in the .resume file */
	// TR_PEER_FROM_LTEP, /* peer address provided in an LTEP handshake */
	FromCache    bool `json:"fromCache"`
	FromDHT      bool `json:"fromDht"`
	FromIncoming bool `json:"fromIncoming"`
	FromLPD      bool `json:"fromLPD"`
	FromLTEP     bool `json:"fromLTEP"`
	FromPEX      bool `json:"fromPex"`
	FromTracker  bool `json:"fromTracker"`
}

// https://github.com/transmission/transmission/blob/6ca0ce683a5aaa8991b16ab7d93722b8861f626b/libtransmission/transmission.h#L1508
type TorrentStatistics struct {
	/** The last time we uploaded or downloaded piece data on this torrent. */
	ActivityDate uint64 `json:"activityDate"`
	/** When the torrent was first added. */
	AddedDate uint64 `json:"addedDate"`
	/** Byte count of all the corrupt data you've ever downloaded for
	  this torrent. If you're on a poisoned torrent, this number can
	  grow very large. */
	CorruptEver uint64 `json:"corruptEver"`
	/** Byte count of all the piece data we want and don't have yet,
	  but that a connected peer does have. [0...leftUntilDone] */
	DesiredAvailable uint64 `json:"desiredAvailable"`
	/** When the torrent finished downloading. */
	DoneDate uint64 `json:"doneDate"`
	/** Byte count of all the non-corrupt data you've ever downloaded
	  for this torrent. If you deleted the files and downloaded a second
	  time, this will be 2*totalSize.. */
	DownloadedEver uint64 `json:"downloadedEver"`
	/** The last time during this session that a rarely-changing field
	  changed -- e.g. any tr_torrent_metainfo field (trackers, filenames, name)
	  or download directory. RPC clients can monitor this to know when
	  to reload fields that rarely change. */
	EditDate uint64 `json:"editDate"`
	/** If downloading, estimated number of seconds left until the torrent is done.
	  If seeding, estimated number of seconds left until seed ratio is reached. */
	ETA int `json:"eta"`
	/** If seeding, number of seconds left until the idle time limit is reached. */
	ETAIdle int `json:"etaIdle"`
	/** Byte count of all the partial piece data we have for this torrent.
	  As pieces become complete, this value may decrease as portions of it
	  are moved to `corrupt' or `haveValid'. */
	HaveUnchecked uint64 `json:"haveUnchecked"`
	/** Byte count of all the checksum-verified data we have for this torrent.
	 */
	HaveValid uint64 `json:"haveValid"`
	/** Number of seconds since the last activity (or since started).
	  -1 if activity is not seeding or downloading. */
	IdleSecs int `json:"idleSecs"`
	/** A torrent is considered finished if it has met its seed ratio.
	  As a result, only paused torrents can be finished. */
	IsFinished bool `json:"isFinished"`
	/** True if the torrent is running, but has been idle for long enough
	  to be considered stalled.  @see tr_sessionGetQueueStalledMinutes() */
	IsStalled int `json:"IsStalled"`
	/** Byte count of how much data is left to be downloaded until we've got
	  all the pieces that we want. [0...tr_stat.sizeWhenDone] */
	LeftUntilDone uint64 `json:"leftUntilDone"`
	/** How much of the metadata the torrent has.
	  For torrents added from a torrent this will always be 1.
	  For magnet links, this number will from from 0 to 1 as the metadata is downloaded.
	  Range is [0..1] */
	MetadataPercentComplete float32 `json:"metadataPercentComplete"`
	/** Number of peers that we're connected to */
	PeersConnected int `json:"peersConnected"`
	/** How many peers we found out about from the tracker, or from pex,
	  or from incoming connections, or from our resume file. */
	PeersFrom interface{} `json:"peersFrom"`
	/** Number of peers that we're sending data to */
	PeersGettingFromUs int `json:"peersGettingFromUs"`
	/** Number of peers that are sending data to us. */
	PeersSendingToUs int `json:"peersSendingToUs"`
	/** How much has been downloaded of the entire torrent.
	  Range is [0..1] */
	PercentComplete float32 `json:"percentComplete"`
	/** How much has been downloaded of the files the user wants. This differs
	  from percentComplete if the user wants only some of the torrent's files.
	  Range is [0..1]
	  @see tr_stat.leftUntilDone */
	PercentDone float32 `json:"percentDone"`
	/** Speed all piece being received for this torrent.
	  This ONLY counts piece data. */
	PieceDownloadSpeed float32 `json:"pieceDownloadSpeed"`
	/** Speed all piece being sent for this torrent.
	  This ONLY counts piece data. */
	PieceUploadSpeed float32 `json:"pieceUploadSpeed"`
	/** This torrent's queue position.
	  All torrents have a queue position, even if it's not queued. */
	QueuePosition int `json:"queuePosition"`
	/** Total uploaded bytes / sizeWhenDone.
	  NB: In Transmission 3.00 and earlier, this was total upload / download,
	  which caused edge cases when total download was less than sizeWhenDone. */
	Ratio float32 `json:"ratio"`
	/** When tr_stat.activity is TR_STATUS_CHECK or TR_STATUS_CHECK_WAIT,
	  this is the percentage of how much of the files has been
	  verified. When it gets to 1, the verify process is done.
	  Range is [0..1]
	  @see tr_stat.activity */
	RecheckProgress float32 `json:"recheckProgress"`
	/** Cumulative seconds the torrent's ever spent downloading */
	SecondsDownloading int `json:"secondsDownloading"`
	/** Cumulative seconds the torrent's ever spent seeding */
	SecondsSeeding int `json:"secondsSeeding"`
	/** How much has been uploaded to satisfy the seed ratio.
	  This is 1 if the ratio is reached or the torrent is set to seed forever.
	  Range is [0..1] */
	SeedRatioPercentDone float32 `json:"seedRatioPercentDone"`
	/** Byte count of all the piece data we'll have downloaded when we're done,
	  whether or not we have it yet. This may be less than tr_torrentTotalSize()
	  if only some of the torrent's files are wanted.
	  [0...tr_torrentTotalSize()] */
	SizeWhenDone uint64 `json:"sizeWhenDone"`
	/** When the torrent was last started. */
	StartDate   uint64 `json:"startDate"`
	TrackerList int    `json:"trackerList"` // one per line
	TotalSize   int    `json:"totalSize"`
	/** Byte count of all data you've ever uploaded for this torrent. */
	UploadedEver uint64 `json:"uploadedEver"`
	// wanted     int     `json:"wanted"`
	WebWeeds []string `json:"webseeds"`
	/** Number of webseeds that are sending data to us. */
	WebWeedsSendingToUs int `json:"webseedsSendingToUs"`
}

// https://github.com/transmission/transmission/blob/main/docs/rpc-spec.md#33-torrent-accessor-torrent-get
type Torrent struct {
	TorrentStatistics
	Availability        int                `json:"availability"` // Should be an array of ints, but is actually just a single int
	BandwidthPriority   int                `json:"bandwidthPriority"`
	Comment             string             `json:"comment"`
	Creator             string             `json:"creator"`
	DateCreated         int                `json:"dateCreated"`
	DownloadDir         string             `json:"downloadDir"`
	DownloadLimit       int                `json:"downloadLimit"`
	DownloadLimited     bool               `json:"downloadLimited"`
	Error               int                `json:"error"`
	ErrorString         string             `json:"errorString"`
	FileCount           int                `json:"file-count"`
	Files               []TorrentFile      `json:"files"`
	FileStats           []TorrentFileStats `json:"fileStats"`
	Group               int                `json:"group"`
	HashString          string             `json:"hashString"`
	HonorsSessionLimits bool               `json:"honorsSessionLimits"`
	ID                  int                `json:"id"`
	IsPrivate           bool               `json:"isPrivate"`
	Labels              []string           `json:"labels"`
	MagnetLink          string             `json:"magnetLink"`
	ManualAnnounceTime  int                `json:"manualAnnounceTime"`
	MaxConnectedPeers   int                `json:"maxConnectedPeers"`
	Name                string             `json:"name"`
	PeerLimit           int                `json:"peer-limit"`
	Peers               []TorrentPeer      `json:"peers"`
	Pieces              string             `json:"pieces"` // Base64 bitfield
	PieceCount          int                `json:"pieceCount"`
	PieceSize           int                `json:"pieceSize"`
	Priorities          []int              `json:"priorities"`
	PrimaryMIMEType     int                `json:"primary-mime-type"`
	RateDownload        int                `json:"rateDownload"`
	RateUpload          int                `json:"rateUpload"`
	SeedIdleLimit       int                `json:"seedIdleLimit"`
	SeedIdleMode        int                `json:"seedIdleMode"`
	SeedRatioLimit      float64            `json:"seedRatioLimit"`
	SeedRatioMode       int                `json:"seedRatioMode"`
	// 0	Torrent is stopped
	// 1	Torrent is queued to verify local data
	// 2	Torrent is verifying local data
	// 3	Torrent is queued to download
	// 4	Torrent is downloading
	// 5	Torrent is queued to seed
	// 6	Torrent is seeding
	Status        int     `json:"status"`
	TorrentFile   string  `json:"torrentFile"`
	UploadLimit   int     `json:"uploadLimit"`
	UploadLimited bool    `json:"uploadLimited"`
	UploadRatio   float64 `json:"uploadRatio"`
}

type listTorrentsRequestArgs struct {
	IDs    []int    `json:"ids,omitempty"`
	Fields []string `json:"fields"`
}

type listTorrentsResponse struct {
	Result    string                   `json:"result"`
	Arguments listTorrentsResponseArgs `json:"arguments"`
	Tag       string                   `json:"tag"`
}

type listTorrentsResponseArgs struct {
	Torrents []Torrent `json:"torrents"`
}

func (t *Transmission) GetTorrents(ctx context.Context, ids ...int) ([]Torrent, error) {
	var response listTorrentsResponse
	req := listTorrentsRequestArgs{
		Fields: torrentFields,
	}
	if ids != nil {
		req.IDs = ids
	}
	if err := t.callRPC(ctx, "torrent-get", &req, &response); err != nil {
		return nil, err
	}
	if response.Result != "success" {
		return nil, fmt.Errorf("failed to get session stats: %s", response.Result)
	}
	return response.Arguments.Torrents, nil
}
