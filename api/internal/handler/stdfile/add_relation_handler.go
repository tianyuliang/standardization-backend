// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package stdfile

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/stdfile"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 根据标准文件ID添加关联关系
func AddRelationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Step 1: 解析路径参数 :id
		idStr := r.URL.Query().Get(":id")
		if idStr == "" {
			idStr = r.PathValue("id")
		}
		if idStr == "" {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("missing path parameter: id"))
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("invalid id parameter: %w", err))
			return
		}

		// Step 2: 解析请求体（JSON）
		var req types.StdFileRelationDto
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// Step 3: 调用业务逻辑
		l := stdfile.NewAddRelationLogic(r.Context(), svcCtx)
		resp, err := l.AddRelation(id, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
