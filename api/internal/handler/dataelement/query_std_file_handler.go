// Code scaffolded by goctl. Safe to edit.

package dataelement

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/dataelement"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 查询数据元关联的标准文件
func QueryStdFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryStdFilePageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 提取路径参数 :id
		pathParts := strings.Split(r.URL.Path, "/")
		id := int64(0)
		if len(pathParts) >= 5 {
			idStr := pathParts[4] // /v1/dataelement/query/stdFile/{id}
			id, _ = strconv.ParseInt(idStr, 10, 64)
		}

		l := dataelement.NewQueryStdFileLogic(r.Context(), svcCtx)
		resp, err := l.QueryStdFile(id, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
