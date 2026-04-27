package spider

import (
	"context"
	"net/url"
)

// SevenClient 台湾 7-ELEVEN 交货便单号查询。
// 接口：POST https://myship.7-11.com.tw/Order/QueryShipmentInfo
// Referer 必须携带，否则会被 403。
type SevenClient struct {
	base *BaseClient
}

func NewSevenClient(base *BaseClient) *SevenClient {
	return &SevenClient{base: base}
}

const (
	sevenQueryURL = "https://myship.7-11.com.tw/Order/QueryShipmentInfo"
	sevenReferer  = "https://myship.7-11.com.tw/order/query_shipment"
)

// Query 查询单号，返回原始响应 HTML。解析逻辑由调用方根据页面内容做，
// 这里保留 Raw 给 logic 做容错 + 后续对结构的适配。
func (c *SevenClient) Query(ctx context.Context, orderNo string) (string, error) {
	form := url.Values{}
	form.Set("order_no", orderNo)
	headers := map[string]string{
		"Referer": sevenReferer,
	}
	return c.base.DoPostFormMobile(ctx, sevenQueryURL, form, headers)
}
