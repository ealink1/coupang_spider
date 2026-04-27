package logic

import (
	"context"

	"coupang_spider/internal/pkg/spider"
	"coupang_spider/internal/svc"
	"coupang_spider/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SevenTrackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSevenTrackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SevenTrackLogic {
	return &SevenTrackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Track 调用 7-ELEVEN 查询接口，返回解析后的轨迹。
func (l *SevenTrackLogic) Track(req *types.TrackRequest) (*types.TrackResponse, error) {
	if req.TrackNo == "" {
		return nil, errEmptyTrackNo
	}
	raw, err := l.svcCtx.Seven.Query(l.ctx, req.TrackNo)
	if err != nil && raw == "" {
		l.Logger.Errorf("7-11 query failed: %v", err)
		return nil, err
	}
	// 7-11 页面返回的 HTML 中包含轨迹表格，尝试解析。
	rows := spider.ParseHTMLTableRows(raw, "")
	return &types.TrackResponse{
		TrackNo:    req.TrackNo,
		Carrier:    "TW_7ELEVEN",
		StatusList: rowsToItems(rows),
		Raw:        raw,
	}, nil
}
