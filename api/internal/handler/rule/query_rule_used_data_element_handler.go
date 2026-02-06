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

// 查询引用规则的数据元
func QueryRuleUsedDataElementHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get(":id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		var req types.PageQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := rule.NewQueryRuleUsedDataElementLogic(r.Context(), svcCtx)
		resp, err := l.QueryRuleUsedDataElement(id, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
