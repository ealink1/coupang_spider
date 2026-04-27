// Package spider 提供对各家物流/便利店单号查询接口的爬虫 client。
package spider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// BaseClient 封装通用 HTTP 调用，供各家爬虫 client 组合使用。
type BaseClient struct {
	HttpClient       *http.Client
	DefaultUserAgent string
	MobileUserAgent  string
}

func NewBaseClient(httpClient *http.Client, defaultUA, mobileUA string) *BaseClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if defaultUA == "" {
		defaultUA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	}
	if mobileUA == "" {
		mobileUA = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148"
	}
	return &BaseClient{
		HttpClient:       httpClient,
		DefaultUserAgent: defaultUA,
		MobileUserAgent:  mobileUA,
	}
}

// DoGet 发起 GET 请求，返回响应体字符串。
func (b *BaseClient) DoGet(ctx context.Context, rawURL string, headers map[string]string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return "", err
	}
	b.applyHeaders(req, headers, false)
	return b.do(req)
}

// DoPostForm 以 application/x-www-form-urlencoded 方式 POST。
func (b *BaseClient) DoPostForm(ctx context.Context, rawURL string, form url.Values, headers map[string]string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	b.applyHeaders(req, headers, false)
	return b.do(req)
}

// DoPostFormMobile 使用移动端 UA 发起表单 POST（部分站点针对移动端放宽校验）。
func (b *BaseClient) DoPostFormMobile(ctx context.Context, rawURL string, form url.Values, headers map[string]string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	b.applyHeaders(req, headers, true)
	return b.do(req)
}

// DoPostJSON 以 application/json 方式 POST 字符串 JSON。
func (b *BaseClient) DoPostJSON(ctx context.Context, rawURL string, body string, headers map[string]string) (string, error) {
	// 用字符串 body 保持调用方对 JSON 字段名和 null 值的完全控制。
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, strings.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	b.applyHeaders(req, headers, false)
	return b.do(req)
}

func (b *BaseClient) applyHeaders(req *http.Request, headers map[string]string, mobile bool) {
	if req.Header.Get("User-Agent") == "" {
		if mobile {
			req.Header.Set("User-Agent", b.MobileUserAgent)
		} else {
			req.Header.Set("User-Agent", b.DefaultUserAgent)
		}
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func (b *BaseClient) do(req *http.Request) (string, error) {
	resp, err := b.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return string(body), fmt.Errorf("http status %d", resp.StatusCode)
	}
	return string(body), nil
}
