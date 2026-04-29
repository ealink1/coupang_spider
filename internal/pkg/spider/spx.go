package spider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

// SpxClient 台湾虾皮 SPX 单号查询。
// SPX 官网详情页是前端 SPA；前端实际调用 /api/v2/fleet_order/tracking/search。
type SpxClient struct {
	base *BaseClient
}

func NewSpxClient(base *BaseClient) *SpxClient {
	return &SpxClient{base: base}
}

const (
	spxTrackingURL  = "https://spx.tw/api/v2/fleet_order/tracking/search"
	spxDetailURL    = "https://spx.tw/detail/"
	spxTrackingSalt = "MGViZmZmZTYzZDJhNDgxY2Y1N2ZlN2Q1ZWJkYzlmZDY="
)

// SpxTrackingEvent 虾皮返回事件节点（按主要字段建模，其他字段可延后补齐）
type SpxTrackingEvent struct {
	Description string `json:"description"`
	Message     string `json:"message,omitempty"`
	Timestamp   int64  `json:"timestamp"`
	Location    string `json:"location,omitempty"`
	Status      string `json:"status,omitempty"`
}

// SpxTrackingResp 虾皮返回结构（宽松建模，保留 Raw 便于后续扩展）
type SpxTrackingResp struct {
	RetCode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		TrackingList []SpxTrackingEvent `json:"tracking_list"`
	} `json:"data"`
}

// QueryRaw 执行 SPX 前端同款接口查询，返回 JSON 原始字符串。
func (c *SpxClient) QueryRaw(ctx context.Context, trackNo string) (string, error) {
	q := url.Values{}
	q.Set("sls_tracking_number", buildSpxSignedTrackingNo(trackNo, time.Now().Unix()))
	headers := map[string]string{
		"Accept":  "application/json, text/plain, */*",
		"Referer": spxDetailURL + url.PathEscape(trackNo),
	}
	return c.base.DoGet(ctx, spxTrackingURL+"?"+q.Encode(), headers)
}

// Query 查询并把 JSON 反序列化为结构体。解析失败时仍返回 Raw 字符串。
func (c *SpxClient) Query(ctx context.Context, trackNo string) (*SpxTrackingResp, string, error) {
	raw, err := c.QueryRaw(ctx, trackNo)
	if err != nil {
		return nil, raw, err
	}
	out := &SpxTrackingResp{}
	if jerr := json.Unmarshal([]byte(raw), out); jerr != nil {
		return nil, raw, jerr
	}
	return out, raw, nil
}

func buildSpxSignedTrackingNo(trackNo string, timestamp int64) string {
	ts := strconv.FormatInt(timestamp, 10)
	sum := sha256.Sum256([]byte(trackNo + ts + spxTrackingSalt))
	return trackNo + "|" + ts + hex.EncodeToString(sum[:])
}
