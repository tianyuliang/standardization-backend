// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"fmt"
	"net/http"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/task"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 完成任务
func FinishTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract task_id from path parameter
		taskIdStr := r.PathValue("task_id")
		if taskIdStr == "" {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("missing path parameter: task_id"))
			return
		}

		// 验证taskId格式（36位UUID）
		if len(taskIdStr) != 36 {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("invalid task_id format: must be 36 characters UUID"))
			return
		}

		l := task.NewFinishTaskLogic(r.Context(), svcCtx)
		resp, err := l.FinishTask(taskIdStr)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
