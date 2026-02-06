// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"net/http"
	"strconv"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/rule"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 根据标准文件查询规则
func QueryRuleByStdFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileIdStr := r.URL.Query().Get("fileId")
		if fileIdStr == "" {
			httpx.ErrorCtx(r.Context(), w, nil)
			return
		}
		fileId, _ := strconv.ParseInt(fileIdStr, 10, 64)

		var req types.RuleListQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := rule.NewQueryRuleByStdFileLogic(r.Context(), svcCtx)
		resp, err := l.QueryRuleByStdFile(fileId, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
