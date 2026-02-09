// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dataelement

import (
	"net/http"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/dataelement"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 内部分页查询数据元
func ListDataElementInternalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageInfoWithKeyword
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := dataelement.NewListDataElementInternalLogic(r.Context(), svcCtx)
		resp, err := l.ListDataElementInternal(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
