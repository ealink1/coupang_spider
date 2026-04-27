package logic

import (
	"errors"
	"strings"

	"coupang_spider/internal/types"
)

// errEmptyTrackNo 统一的空单号错误。
var errEmptyTrackNo = errors.New("trackNo is empty")

// rowsToItems 把二维表格转为 TrackItem 列表，默认按 [时间, 状态, 位置] 读取。
// 空行、仅包含表头关键字的行会被过滤。
func rowsToItems(rows [][]string) []types.TrackItem {
	items := make([]types.TrackItem, 0, len(rows))
	for _, row := range rows {
		if len(row) == 0 {
			continue
		}
		first := row[0]
		// 跳过纯表头行
		if isHeaderToken(first) {
			continue
		}
		item := types.TrackItem{Time: first}
		if len(row) >= 2 {
			item.Status = row[1]
		}
		if len(row) >= 3 {
			item.Location = row[2]
		}
		if item.Time == "" && item.Status == "" {
			continue
		}
		items = append(items, item)
	}
	return items
}

func isHeaderToken(v string) bool {
	v = strings.TrimSpace(v)
	switch v {
	case "", "時間", "时间", "日期", "日期時間", "日期时间", "Time", "time":
		return true
	}
	return false
}
