// Code scaffolded by goctl. Safe to edit.

package dataelement

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/dataelement"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除数据元标签
func DeleteLabelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 提取路径参数 :id
		// URL格式: /v1/dataelement/label/{id}
		pathParts := strings.Split(r.URL.Path, "/")
		id := int64(0)
		if len(pathParts) >= 5 {
			idStr := pathParts[4]
			id, _ = strconv.ParseInt(idStr, 10, 64)
		}

		l := dataelement.NewDeleteLabelLogic(r.Context(), svcCtx)
		resp, err := l.DeleteLabel(id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
