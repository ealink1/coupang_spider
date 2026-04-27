package spider

import (
	"context"
	"net/url"
)

// KerryClient 台湾嘉里大荣货态查询。
// 接口：GET https://www.kerrytj.com/ZH/search/search_track.aspx?track_no={trackNo}
type KerryClient struct {
	base *BaseClient
}

func NewKerryClient(base *BaseClient) *KerryClient {
	return &KerryClient{base: base}
}

const kerryURL = "https://www.kerrytj.com/ZH/search/search_track.aspx"

// Query 执行单号查询，返回 HTML 原始响应。
func (c *KerryClient) Query(ctx context.Context, trackNo string) (string, error) {
	q := url.Values{}
	q.Set("track_no", trackNo)
	return c.base.DoGet(ctx, kerryURL+"?"+q.Encode(), nil)
}
