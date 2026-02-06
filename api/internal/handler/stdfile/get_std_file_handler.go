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

// 根据ID查询详情
func GetStdFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Step 1: 解析路径参数 id
		idStr := r.PathValue("id")
		if idStr == "" {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("missing path parameter: id"))
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("invalid id parameter: %w", err))
			return
		}

		// 构建请求对象并注入ID
		req := &types.StdFileDetailResp{
			Id: id,
		}

		l := stdfile.NewGetStdFileLogic(r.Context(), svcCtx)
		resp, err := l.GetStdFile(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
