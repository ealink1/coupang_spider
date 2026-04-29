package handler

import (
	"net/http"

	"coupang_spider/internal/logic"
	"coupang_spider/internal/svc"
	"coupang_spider/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func OkmartTrackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TrackRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewOkMartTrackLogic(r.Context(), svcCtx)
		resp, err := l.Track(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
