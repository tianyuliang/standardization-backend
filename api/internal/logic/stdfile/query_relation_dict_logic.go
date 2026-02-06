// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package stdfile

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/stdfile/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryRelationDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询关联的码表
//
// 业务流程:
//  1. 校验文件存在性
//  2. 调用 Dict 服务查询关联码表
//
// 异常处理:
//   - 30201: 标准文件不存在
func NewQueryRelationDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRelationDictLogic {
	return &QueryRelationDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRelationDictLogic) QueryRelationDict(id int64, req *types.StdFileRelationQuery) (resp *types.StdFileDataListResp, err error) {
	// Step 1: 校验文件存在性
	existing, err := l.svcCtx.StdFileModel.FindOne(l.ctx, id)
	if err != nil {
		return nil, HandleError(err)
	}
	if existing == nil {
		return &types.StdFileDataListResp{
			TotalCount: 0,
			Data:       []types.StdFileDetailResp{},
		}, nil
	}

	// Step 2: 调用 Dict 服务查询关联码表 (Mock)
	_, err = mock.DictQueryPageByFileId(l.ctx, l.svcCtx, id, req.Offset, req.Limit)
	if err != nil {
		return nil, err
	}

	logx.Infof("查询标准文件关联码表成功: fileId=%d", id)

	// TODO: 将返回的 Dict 列表转换为 StdFileDetailResp
	return &types.StdFileDataListResp{
		TotalCount: 0,
		Data:       []types.StdFileDetailResp{},
	}, nil
}
