// Code scaffolded by goctl. Safe to edit.

package dataelement

import (
	"net/http"
	"strings"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/dataelement"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 启用/停用数据元
func UpdateStateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateStateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 提取路径参数 :ids
		// URL格式: /v1/dataelement/state/{ids}
		pathParts := strings.Split(r.URL.Path, "/")
		ids := ""
		if len(pathParts) >= 5 {
			ids = pathParts[4] // /v1/dataelement/state/{ids}
		}

		l := dataelement.NewUpdateStateLogic(r.Context(), svcCtx)
		resp, err := l.UpdateState(ids, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
