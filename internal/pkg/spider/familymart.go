package spider

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
)

// FamilyMartClient 台湾全家 FamilyMart 单号查询。
// 新版官网入口为 fmec.famiport.com.tw，实际 iframe 调用 ecfme.fme.com.tw 的 JSON 接口。
type FamilyMartClient struct {
	base *BaseClient
}

// NewFamilyMartClient 创建全家查询客户端。
func NewFamilyMartClient(base *BaseClient) *FamilyMartClient {
	return &FamilyMartClient{base: base}
}

const familyMartLegacyURL = "https://ecfmeve.famiport.com.tw/fms/Query/TK_Query.aspx"
const familyMartInquiryURL = "https://ecfme.fme.com.tw/FMEDCFPWebV2_II/list.aspx/InquiryOrders"
const familyMartDetailURL = "https://ecfme.fme.com.tw/FMEDCFPWebV2_II/list.aspx/GetOrderDetail"

// familyMartEnvelope 对应 ASP.NET WebMethod 返回格式，业务 JSON 包在 d 字段中。
type familyMartEnvelope struct {
	D string `json:"d"`
}

// familyMartInquiryResponse 对应全家货件列表查询结果。
type familyMartInquiryResponse struct {
	ErrorCode    string                  `json:"ErrorCode"`
	ErrorMessage string                  `json:"ErrorMessage"`
	List         []familyMartInquiryItem `json:"List"`
}

// familyMartInquiryItem 对应全家单笔货件列表资料。
type familyMartInquiryItem struct {
	ECOrderNo    string  `json:"EC_ORDER_NO"`
	CNT          float64 `json:"CNT"`
	OrderNo      string  `json:"ORDER_NO"`
	OrderMessage string  `json:"ORDERMESSAGE"`
}

// Query 执行全家新版 JSON 查询，返回 GetOrderDetail 的原始 JSON。
func (c *FamilyMartClient) Query(ctx context.Context, trackNo string) (string, error) {
	// 先调用列表接口，拿到唯一的 ORDER_NO 后才能继续查明细。
	inquiryBody, err := json.Marshal(map[string]string{"ListEC_ORDER_NO": trackNo})
	if err != nil {
		return "", err
	}
	inquiryRaw, err := c.base.DoPostJSON(ctx, familyMartInquiryURL, string(inquiryBody), familyMartHeaders())
	if err != nil && inquiryRaw == "" {
		return "", err
	}
	inquiry, err := parseFamilyMartInquiry(inquiryRaw)
	if err != nil {
		return inquiryRaw, err
	}
	if inquiry.ErrorCode != "000" {
		return inquiryRaw, errors.New("familymart: " + inquiry.ErrorMessage)
	}
	if len(inquiry.List) == 0 || inquiry.List[0].CNT == 0 {
		return inquiryRaw, nil
	}
	if inquiry.List[0].CNT > 1 {
		return inquiryRaw, errors.New("familymart: multiple orders require receiver name")
	}

	// 再调用明细接口，RCV_USER_NAME 为 nil 会编码成 null，与官网单笔订单查询行为一致。
	orderNo := defaultIfEmpty(inquiry.List[0].OrderNo, trackNo)
	detailBody, err := json.Marshal(map[string]any{
		"EC_ORDER_NO":   trackNo,
		"RCV_USER_NAME": nil,
		"ORDER_NO":      orderNo,
	})
	if err != nil {
		return "", err
	}
	return c.base.DoPostJSON(ctx, familyMartDetailURL, string(detailBody), familyMartHeaders())
}

// queryLegacy 执行旧版 WebForms 查询流程，保留作排查旧站回滚时使用。
func (c *FamilyMartClient) queryLegacy(ctx context.Context, trackNo string) (string, error) {
	// Step 1: 获取初始页面，提取 hidden input
	pageHTML, err := c.base.DoGet(ctx, familyMartLegacyURL, map[string]string{
		"User-Agent": c.base.MobileUserAgent,
	})
	if err != nil && pageHTML == "" {
		return "", err
	}
	hidden := extractHiddenInputs(pageHTML)
	if hidden["__VIEWSTATE"] == "" || hidden["__EVENTVALIDATION"] == "" {
		return pageHTML, errors.New("familymart: failed to extract __VIEWSTATE / __EVENTVALIDATION")
	}

	// Step 2: 带 token POST 查询
	form := url.Values{}
	form.Set("__VIEWSTATE", hidden["__VIEWSTATE"])
	form.Set("__VIEWSTATEGENERATOR", defaultIfEmpty(hidden["__VIEWSTATEGENERATOR"], "87B579A2"))
	form.Set("__EVENTVALIDATION", hidden["__EVENTVALIDATION"])
	form.Set("txtInvNum", trackNo)
	form.Set("btnQuery", "查詢")

	return c.base.DoPostFormMobile(ctx, familyMartLegacyURL, form, nil)
}

// parseFamilyMartInquiry 解析 ASP.NET WebMethod 包裹的列表查询响应。
func parseFamilyMartInquiry(raw string) (*familyMartInquiryResponse, error) {
	// 外层 d 是字符串形式 JSON，需要解两次。
	var envelope familyMartEnvelope
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		return nil, err
	}
	var resp familyMartInquiryResponse
	if err := json.Unmarshal([]byte(envelope.D), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// familyMartHeaders 返回全家新版 JSON 接口需要的基础请求头。
func familyMartHeaders() map[string]string {
	return map[string]string{
		"Accept":  "application/json, text/javascript, */*; q=0.01",
		"Origin":  "https://ecfme.fme.com.tw",
		"Referer": "https://ecfme.fme.com.tw/FMEDCFPWebV2_II/list.aspx",
	}
}

// defaultIfEmpty 在值为空时返回默认值。
func defaultIfEmpty(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
