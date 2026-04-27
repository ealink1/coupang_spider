package spider

import (
	"context"
	"net/url"
)

// HiLifeClient 台湾莱尔富 Hi-Life 单号查询。
// 接口：GET https://www.hilife.com.tw/E3-logisticSearch.aspx?no={orderNo}
type HiLifeClient struct {
	base *BaseClient
}

func NewHiLifeClient(base *BaseClient) *HiLifeClient {
	return &HiLifeClient{base: base}
}

const hiLifeURL = "https://www.hilife.com.tw/E3-logisticSearch.aspx"

// Query 执行单号查询，返回 HTML 原始响应。
func (c *HiLifeClient) Query(ctx context.Context, orderNo string) (string, error) {
	q := url.Values{}
	q.Set("no", orderNo)
	return c.base.DoGet(ctx, hiLifeURL+"?"+q.Encode(), nil)
}
