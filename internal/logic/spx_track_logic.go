package logic

import (
	"context"
	"time"

	"coupang_spider/internal/svc"
	"coupang_spider/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SpxTrackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSpxTrackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SpxTrackLogic {
	return &SpxTrackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Track 调用虾皮 SPX 查询接口（JSON），返回解析后的轨迹。
func (l *SpxTrackLogic) Track(req *types.TrackRequest) (*types.TrackResponse, error) {
	if req.TrackNo == "" {
		return nil, errEmptyTrackNo
	}
	resp, raw, err := l.svcCtx.Spx.Query(l.ctx, req.TrackNo)
	if err != nil {
		l.Logger.Errorf("spx query failed: %v, raw=%s", err, truncate(raw, 500))
		if resp == nil {
			return &types.TrackResponse{
				TrackNo: req.TrackNo,
				Carrier: "TW_SPX",
				Raw:     raw,
			}, nil
		}
	}

	items := make([]types.TrackItem, 0)
	if resp != nil {
		for _, ev := range resp.Data.TrackingList {
			t := ""
			if ev.Timestamp > 0 {
				t = time.Unix(ev.Timestamp, 0).Format("2006-01-02 15:04:05")
			}
			items = append(items, types.TrackItem{
				Time:     t,
				Status:   firstNonEmpty(ev.Description, ev.Message, ev.Status),
				Location: ev.Location,
			})
		}
	}
	return &types.TrackResponse{
		TrackNo:    req.TrackNo,
		Carrier:    "TW_SPX",
		StatusList: items,
		Raw:        raw,
	}, nil
}

func firstNonEmpty(vs ...string) string {
	for _, v := range vs {
		if v != "" {
			return v
		}
	}
	return ""
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
