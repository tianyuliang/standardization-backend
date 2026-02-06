// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package stdfile

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/stdfile"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 标准文件管理-批量删除
func DeleteStdFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Step 1: 解析路径参数 :ids (逗号分隔的ID列表)
		idsStr := r.URL.Query().Get(":ids")
		if idsStr == "" {
			idsStr = r.PathValue("ids")
		}
		if idsStr == "" {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("missing path parameter: ids"))
			return
		}

		// 解析ID列表 (支持逗号分隔)
		var ids []int64
		for _, idStr := range strings.Split(idsStr, ",") {
			id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, fmt.Errorf("invalid id parameter: %w", err))
				return
			}
			ids = append(ids, id)
		}

		// Step 2: 调用业务逻辑
		l := stdfile.NewDeleteStdFileLogic(r.Context(), svcCtx)
		resp, err := l.DeleteStdFile(ids)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
