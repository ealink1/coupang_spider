package logic

import (
	"context"

	"coupang_spider/internal/pkg/spider"
	"coupang_spider/internal/svc"
	"coupang_spider/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HiLifeTrackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHiLifeTrackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HiLifeTrackLogic {
	return &HiLifeTrackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Track 调用莱尔富 Hi-Life 查询接口，返回解析后的轨迹。
func (l *HiLifeTrackLogic) Track(req *types.TrackRequest) (*types.TrackResponse, error) {
	if req.TrackNo == "" {
		return nil, errEmptyTrackNo
	}
	raw, err := l.svcCtx.HiLife.Query(l.ctx, req.TrackNo)
	if err != nil && raw == "" {
		l.Logger.Errorf("hilife query failed: %v", err)
		return nil, err
	}
	rows := spider.ParseHTMLTableRows(raw, "")
	return &types.TrackResponse{
		TrackNo:    req.TrackNo,
		Carrier:    "TW_HILIFE",
		StatusList: rowsToItems(rows),
		Raw:        raw,
	}, nil
}
