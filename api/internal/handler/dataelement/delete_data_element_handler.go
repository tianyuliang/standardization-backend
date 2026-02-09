// Code scaffolded by goctl. Safe to edit.

package dataelement

import (
	"net/http"
	"strings"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/dataelement"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除数据元
func DeleteDataElementHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 提取路径参数 :ids
		// URL格式: /v1/dataelement/{ids}
		pathParts := strings.Split(r.URL.Path, "/")
		ids := ""
		if len(pathParts) >= 4 {
			ids = pathParts[3] // /v1/dataelement/{ids}
		}

		l := dataelement.NewDeleteDataElementLogic(r.Context(), svcCtx)
		resp, err := l.DeleteDataElement(ids)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
