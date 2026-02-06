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

type QueryRelationRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询关联编码规则
//
// 业务流程:
//  1. 校验文件存在性
//  2. 调用 Rule 服务查询关联编码规则
//
// 异常处理:
//   - 30201: 标准文件不存在
func NewQueryRelationRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRelationRuleLogic {
	return &QueryRelationRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRelationRuleLogic) QueryRelationRule(id int64, req *types.StdFileRelationQuery) (resp *types.StdFileDataListResp, err error) {
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

	// Step 2: 调用 Rule 服务查询关联编码规则 (Mock)
	_, err = mock.RuleQueryPageByFileId(l.ctx, l.svcCtx, id, req.Offset, req.Limit)
	if err != nil {
		return nil, err
	}

	logx.Infof("查询标准文件关联编码规则成功: fileId=%d", id)

	// TODO: 将返回的 Rule 列表转换为 StdFileDetailResp
	return &types.StdFileDataListResp{
		TotalCount: 0,
		Data:       []types.StdFileDetailResp{},
	}, nil
}
