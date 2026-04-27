package types
// Code scaffolded by hand (compatible with goctl). Safe to edit.

// TrackRequest 通用爬虫追踪请求
type TrackRequest struct {
	TrackNo string `json:"trackNo"`
}

// TrackItem 轨迹条目
type TrackItem struct {
	Time     string `json:"time"`
	Status   string `json:"status"`
	Location string `json:"location,omitempty"`
}

// TrackResponse 通用爬虫返回
type TrackResponse struct {
	TrackNo    string      `json:"trackNo"`
	Carrier    string      `json:"carrier"`
	StatusList []TrackItem `json:"statusList"`
	Raw        string      `json:"raw,omitempty"`
}
