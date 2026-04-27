package logic

import (
	"context"

	"coupang_spider/internal/pkg/spider"
	"coupang_spider/internal/svc"
	"coupang_spider/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OkMartTrackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOkMartTrackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OkMartTrackLogic {
	return &OkMartTrackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Track 调用 OK 超商查询接口，返回解析后的轨迹。
func (l *OkMartTrackLogic) Track(req *types.TrackRequest) (*types.TrackResponse, error) {
	if req.TrackNo == "" {
		return nil, errEmptyTrackNo
	}
	raw, err := l.svcCtx.OkMart.Query(l.ctx, req.TrackNo)
	if err != nil && raw == "" {
		l.Logger.Errorf("okmart query failed: %v", err)
		return nil, err
	}
	// OK Mart 的轨迹表格 class 为 table-striped
	rows := spider.ParseHTMLTableRows(raw, "table-striped")
	if len(rows) == 0 {
		rows = spider.ParseHTMLTableRows(raw, "")
	}
	return &types.TrackResponse{
		TrackNo:    req.TrackNo,
		Carrier:    "TW_OKMART",
		StatusList: rowsToItems(rows),
		Raw:        raw,
	}, nil
}
