package spider

import (
	"context"
	"net/url"
)

// OkMartClient 台湾 OK 超商单号查询。
// 接口：POST https://webapi.okmart.com.tw/CargoQuery/Index
// 表单参数：CargoNo=单号
type OkMartClient struct {
	base *BaseClient
}

func NewOkMartClient(base *BaseClient) *OkMartClient {
	return &OkMartClient{base: base}
}

const okMartURL = "https://webapi.okmart.com.tw/CargoQuery/Index"

// Query 执行单号查询，返回 HTML 原始响应。
func (c *OkMartClient) Query(ctx context.Context, cargoNo string) (string, error) {
	form := url.Values{}
	form.Set("CargoNo", cargoNo)
	return c.base.DoPostForm(ctx, okMartURL, form, nil)
}
