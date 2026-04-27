package spider

import (
	"context"
	"encoding/json"
	"net/url"
)

// SpxClient 台湾虾皮 SPX 单号查询。
// 接口：GET https://spx.tw/api/v1/tracking?order_sn={trackNo}
// 需要携带 Referer: https://spx.tw/
type SpxClient struct {
	base *BaseClient
}

func NewSpxClient(base *BaseClient) *SpxClient {
	return &SpxClient{base: base}
}

const (
	spxURL     = "https://spx.tw/api/v1/tracking"
	spxReferer = "https://spx.tw/"
)

// SpxTrackingEvent 虾皮返回事件节点（按主要字段建模，其他字段可延后补齐）
type SpxTrackingEvent struct {
	Description string `json:"description"`
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

// QueryRaw 执行单号查询，返回 JSON 原始字符串。
func (c *SpxClient) QueryRaw(ctx context.Context, trackNo string) (string, error) {
	q := url.Values{}
	q.Set("order_sn", trackNo)
	headers := map[string]string{
		"Referer": spxReferer,
		"Accept":  "application/json",
	}
	return c.base.DoGet(ctx, spxURL+"?"+q.Encode(), headers)
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
