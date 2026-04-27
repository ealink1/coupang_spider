package logic

import (
	"context"

	"coupang_spider/internal/pkg/spider"
	"coupang_spider/internal/svc"
	"coupang_spider/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type KerryTrackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewKerryTrackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KerryTrackLogic {
	return &KerryTrackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Track 调用嘉里大荣查询接口，返回解析后的轨迹。
func (l *KerryTrackLogic) Track(req *types.TrackRequest) (*types.TrackResponse, error) {
	if req.TrackNo == "" {
		return nil, errEmptyTrackNo
	}
	raw, err := l.svcCtx.Kerry.Query(l.ctx, req.TrackNo)
	if err != nil && raw == "" {
		l.Logger.Errorf("kerry query failed: %v", err)
		return nil, err
	}
	// Kerry 的轨迹表格 class 为 table_style
	rows := spider.ParseHTMLTableRows(raw, "table_style")
	if len(rows) == 0 {
		rows = spider.ParseHTMLTableRows(raw, "")
	}
	return &types.TrackResponse{
		TrackNo:    req.TrackNo,
		Carrier:    "TW_KERRY",
		StatusList: rowsToItems(rows),
		Raw:        raw,
	}, nil
}
