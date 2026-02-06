// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package stdfile

import (
	"context"
	"fmt"

	"github.com/jinguoxing/idrm-go-base/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	stdfilemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/stdfile/stdfile"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadStdFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// DownloadStdFileResult 下载结果
type DownloadStdFileResult struct {
	FileId      string
	FileName    string
	ContentType string
}

// 根据文件ID下载标准文件附件
//
// 业务流程:
//  1. 校验文件存在性
//  2. 校验附件类型（必须是FILE，不能是URL）
//  3. 返回文件下载信息
//
// 异常处理:
//   - 30201: 标准文件不存在
//   - 30230: URL类型不支持文件下载
func NewDownloadStdFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadStdFileLogic {
	return &DownloadStdFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadStdFileLogic) DownloadStdFile(id int64) (*DownloadStdFileResult, error) {
	// Step 1: 校验文件存在性
	existing, err := l.svcCtx.StdFileModel.FindOne(l.ctx, id)
	if err != nil {
		return nil, HandleError(err)
	}
	if existing == nil {
		return nil, errorx.NewWithMsg(30201, "标准文件不存在")
	}

	// Step 2: 校验附件类型（URL类型不支持文件下载）
	if existing.AttachmentType != stdfilemodel.AttachmentTypeFile {
		return nil, errorx.NewWithMsg(30230, "[URL]类型没有文件附件")
	}

	// Step 3: 返回文件下载信息
	logx.Infof("准备下载标准文件: fileId=%d, fileName=%s", id, existing.FileName)

	return &DownloadStdFileResult{
		FileId:      fmt.Sprintf("%d", existing.Id),
		FileName:    existing.FileName,
		ContentType: "application/octet-stream",
	}, nil
}
