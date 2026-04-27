package logic

import (
	"context"
	"encoding/json"

	"coupang_spider/internal/pkg/spider"
	"coupang_spider/internal/svc"
	"coupang_spider/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FamilyMartTrackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// familyMartEnvelope 对应全家 ASP.NET WebMethod 返回格式。
type familyMartEnvelope struct {
	D string `json:"d"`
}

// familyMartDetailResponse 对应全家货件明细响应。
type familyMartDetailResponse struct {
	ErrorCode    string                 `json:"ErrorCode"`
	ErrorMessage string                 `json:"ErrorMessage"`
	List         []familyMartDetailItem `json:"List"`
}

// familyMartDetailItem 对应全家单笔轨迹资料。
type familyMartDetailItem struct {
	OrderDate       string `json:"ORDER_DATE"`
	OrderTime       string `json:"ORDER_TIME"`
	StatusD         string `json:"STATUS_D"`
	ReceiveStore    string `json:"RCV_STORE_NAME"`
	ReceiveStoreAdr string `json:"RCV_STORE_ADDRESS"`
}

func NewFamilyMartTrackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FamilyMartTrackLogic {
	return &FamilyMartTrackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Track 调用全家 FamilyMart 查询接口，返回解析后的轨迹。
func (l *FamilyMartTrackLogic) Track(req *types.TrackRequest) (*types.TrackResponse, error) {
	if req.TrackNo == "" {
		return nil, errEmptyTrackNo
	}
	raw, err := l.svcCtx.FamilyMart.Query(l.ctx, req.TrackNo)
	if err != nil && raw == "" {
		l.Logger.Errorf("familymart query failed: %v", err)
		return nil, err
	}
	// 新版全家接口返回 JSON；解析失败时保留旧 HTML 表格解析兜底。
	items := parseFamilyMartJSONItems(raw)
	if len(items) == 0 {
		rows := spider.ParseHTMLTableRows(raw, "gvDetail")
		if len(rows) == 0 {
			rows = spider.ParseHTMLTableRows(raw, "")
		}
		items = rowsToItems(rows)
	}
	return &types.TrackResponse{
		TrackNo:    req.TrackNo,
		Carrier:    "TW_FAMILYMART",
		StatusList: items,
		Raw:        raw,
	}, nil
}

// parseFamilyMartJSONItems 将全家新版明细 JSON 转成统一轨迹条目。
func parseFamilyMartJSONItems(raw string) []types.TrackItem {
	// 外层 d 是字符串形式 JSON，需要解两次。
	var envelope familyMartEnvelope
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil || envelope.D == "" {
		return nil
	}
	var detail familyMartDetailResponse
	if err := json.Unmarshal([]byte(envelope.D), &detail); err != nil || detail.ErrorCode != "000" {
		return nil
	}
	items := make([]types.TrackItem, 0, len(detail.List))
	for _, row := range detail.List {
		// 日期和时间拆字段返回，这里拼成现有 TrackItem.Time。
		item := types.TrackItem{
			Time:     row.OrderDate + " " + row.OrderTime,
			Status:   row.StatusD,
			Location: row.ReceiveStore,
		}
		if item.Location == "" {
			item.Location = row.ReceiveStoreAdr
		}
		if item.Time == " " && item.Status == "" {
			continue
		}
		items = append(items, item)
	}
	return items
}
