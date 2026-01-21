// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package catalog

import (
	"net/http"

	"github.com/tianyuliang/standardization-backend/api/internal/logic/catalog"
	"github.com/tianyuliang/standardization-backend/api/internal/svc"
	"github.com/tianyuliang/standardization-backend/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 按名称检索目录
func QueryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := catalog.NewQueryLogic(r.Context(), svcCtx)
		resp, err := l.Query(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
