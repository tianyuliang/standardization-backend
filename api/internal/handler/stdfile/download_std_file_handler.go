// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package stdfile

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/stdfile"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 根据文件ID下载标准文件附件
func DownloadStdFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
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

		// Step 2: 调用业务逻辑
		l := stdfile.NewDownloadStdFileLogic(r.Context(), svcCtx)
		result, err := l.DownloadStdFile(id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// Step 3: 设置响应头并流式传输文件
		// MOCK: 模拟文件下载，实际应从OSS获取文件流
		w.Header().Set("Content-Type", result.ContentType)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", url.QueryEscape(result.FileName)))

		// TODO: 从OSS获取文件流并写入响应
		// reader, _, err := mock.OssDownloadFile(r.Context(), svcCtx, result.FileId)
		// if err != nil {
		//     httpx.ErrorCtx(r.Context(), w, err)
		//     return
		// }
		// defer reader.Close()
		// io.Copy(w, reader)

		// Mock: 写入空内容
		w.WriteHeader(http.StatusOK)
	}
}
