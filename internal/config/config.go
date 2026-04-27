package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Spider Spider
}

type Spider struct {
	HttpTimeoutSec   int    `json:",default=20"`
	DefaultUserAgent string `json:",optional"`
	MobileUserAgent  string `json:",optional"`
}
