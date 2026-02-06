// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package stdfile

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/stdfile"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 标准文件附件下载（批量）
func DownloadBatchStdFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Step 1: 解析请求参数（JSON body）
		var req types.QueryByIdsReq
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// Step 2: 调用业务逻辑
		l := stdfile.NewDownloadBatchStdFileLogic(r.Context(), svcCtx)
		result, err := l.DownloadBatchStdFile(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// Step 3: 设置响应头并流式传输ZIP文件
		// MOCK: 模拟文件下载，实际应从OSS获取文件流并打包为ZIP
		w.Header().Set("Content-Type", result.ContentType)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", url.QueryEscape(result.FileName)))

		// TODO: 从OSS获取多个文件并打包为ZIP
		// 1. 从OSS下载所有文件
		// 2. 创建ZIP写入器
		// 3. 将文件写入ZIP
		// 4. 写入响应流

		// Mock: 写入空内容
		w.WriteHeader(http.StatusOK)
	}
}
