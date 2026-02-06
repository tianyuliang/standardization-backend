// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package stdfile

import (
	"context"
	"fmt"

	"github.com/jinguoxing/idrm-go-base/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	stdfilemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/stdfile/stdfile"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadBatchStdFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// DownloadBatchStdFileResult 批量下载结果
type DownloadBatchStdFileResult struct {
	FileName    string
	ContentType string
	FileCount   int
}

// 标准文件附件下载（批量）
//
// 业务流程:
//  1. 校验ID列表非空
//  2. 查询文件列表并校验存在性
//  3. 过滤掉URL类型和不存在的文件
//  4. 返回批量下载信息（ZIP格式）
//
// 异常处理:
//   - 30202: ID列表不能为空
//   - 30201: 标准文件不存在
//   - 30230: [URL]类型没有文件附件
func NewDownloadBatchStdFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadBatchStdFileLogic {
	return &DownloadBatchStdFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadBatchStdFileLogic) DownloadBatchStdFile(req *types.QueryByIdsReq) (*DownloadBatchStdFileResult, error) {
	// Step 1: 校验ID列表非空
	if len(req.Ids) == 0 {
		return nil, errorx.NewWithMsg(30202, "ID列表不能为空")
	}

	// Step 2: 查询文件列表
	models, err := l.svcCtx.StdFileModel.FindByIds(l.ctx, req.Ids)
	if err != nil {
		return nil, HandleError(err)
	}

	// Step 3: 校验文件存在性并过滤
	var validFiles []*stdfilemodel.StdFile
	var missingIds []int64
	idMap := make(map[int64]bool)
	for _, id := range req.Ids {
		idMap[id] = true
	}

	for _, model := range models {
		if model != nil && model.Deleted == 0 {
			idMap[model.Id] = false
			// 只处理FILE类型，URL类型不支持下载
			if model.AttachmentType == stdfilemodel.AttachmentTypeFile {
				validFiles = append(validFiles, model)
			}
		}
	}

	// 检查缺失的文件
	for id, missing := range idMap {
		if missing {
			missingIds = append(missingIds, id)
		}
	}

	if len(missingIds) > 0 {
		return nil, errorx.NewWithMsg(30201, fmt.Sprintf("标准文件不存在: ids=%v", missingIds))
	}

	if len(validFiles) == 0 {
		return nil, errorx.NewWithMsg(30230, "[URL]类型没有文件附件")
	}

	// Step 4: 生成ZIP文件名
	// Java: String.format("标准文件_%s.zip", DateFormatUtils.format(new Date(), "yyyyMMddHHmmss"))
	zipFileName := "标准文件_batch.zip"

	logx.Infof("准备批量下载标准文件: fileCount=%d", len(validFiles))

	return &DownloadBatchStdFileResult{
		FileName:    zipFileName,
		ContentType: "application/zip",
		FileCount:   len(validFiles),
	}, nil
}
