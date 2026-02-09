// Code scaffolded by goctl. Safe to edit.

package dataelement

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type QueryStdFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询数据元关联的标准文件
func NewQueryStdFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryStdFileLogic {
	return &QueryStdFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryStdFileLogic) QueryStdFile(id int64, req *types.QueryStdFilePageReq) (resp *types.StdFileListVo, err error) {
	// Step 1: 校验数据元ID
	if id <= 0 {
		return nil, errorx.InvalidParameter("id", "数据元ID必须大于0")
	}

	// Step 2: 查询关联的文件ID列表
	fileIds, err := l.svcCtx.RelationDeFileModel.FindFileIdsByDeId(l.ctx, id)
	if err != nil {
		logx.Errorf("查询关联文件失败: %v", err)
		return nil, err
	}

	// Step 3: 获取文件详细信息（从StdFile服务）
	fileNames, _ := mock.GetStdFileByIds(l.ctx, fileIds)

	// Step 4: 构建响应并处理分页
	totalCount := int64(len(fileIds))
	start := int(req.Offset - 1)
	end := start + int(req.Limit)

	if start < 0 {
		start = 0
	}
	if end > len(fileIds) {
		end = len(fileIds)
	}

	if start >= len(fileIds) {
		return &types.StdFileListVo{
			TotalCount: totalCount,
			Entries:    []types.StdFileVo{},
		}, nil
	}

	// Step 5: 构建分页后的条目
	entries := make([]types.StdFileVo, 0, end-start)
	for i := start; i < end; i++ {
		fileId := fileIds[i]
		fileName := ""
		if name, ok := fileNames[fileId]; ok {
			fileName = name
		}
		entries = append(entries, types.StdFileVo{
			Id:       fileId,
			FileName: fileName,
		})
	}

	return &types.StdFileListVo{
		TotalCount: totalCount,
		Entries:    entries,
	}, nil
}
