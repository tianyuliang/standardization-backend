// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package catalog

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/catalog"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除目录
func DeleteCatalogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Step 1: 解析路径参数 id
		// 对应 Java: @PathVariable Long id
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			// go-zero 路由参数解析，从路径中提取 :id
			// 例如: /catalog/123
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

		// Step 2: 调用业务逻辑
		l := catalog.NewDeleteCatalogLogic(r.Context(), svcCtx)
		resp, err := l.DeleteCatalog(id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
