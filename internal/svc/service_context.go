package svc

import (
	"net/http"
	"time"

	"coupang_spider/internal/config"
	"coupang_spider/internal/pkg/spider"
)

type ServiceContext struct {
	Config     config.Config
	HttpClient *http.Client

	Seven      *spider.SevenClient
	FamilyMart *spider.FamilyMartClient
	HiLife     *spider.HiLifeClient
	OkMart     *spider.OkMartClient
	Kerry      *spider.KerryClient
	Spx        *spider.SpxClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	timeout := time.Duration(c.Spider.HttpTimeoutSec) * time.Second
	if timeout <= 0 {
		timeout = 20 * time.Second
	}

	httpClient := &http.Client{Timeout: timeout}

	base := spider.NewBaseClient(httpClient, c.Spider.DefaultUserAgent, c.Spider.MobileUserAgent)

	return &ServiceContext{
		Config:     c,
		HttpClient: httpClient,
		Seven:      spider.NewSevenClient(base),
		FamilyMart: spider.NewFamilyMartClient(base),
		HiLife:     spider.NewHiLifeClient(base),
		OkMart:     spider.NewOkMartClient(base),
		Kerry:      spider.NewKerryClient(base),
		Spx:        spider.NewSpxClient(base),
	}
}
