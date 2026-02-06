// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryDataExistsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 检查数据是否存在
//
// 业务流程:
//  1. 部门ID路径处理
//  2. 检查是否存在（支持 filter_id 排除自身）
func NewQueryDataExistsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryDataExistsLogic {
	return &QueryDataExistsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryDataExistsLogic) QueryDataExists(req *types.QueryDataExistsReq) (resp *types.DataExistsResp, err error) {
	// ====== 步骤1: 部门ID路径处理 ======
	// TODO: 从 Token/部门服务获取完整路径
	deptPathIds := req.DepartmentIds

	// ====== 步骤2: 检查是否存在 ======
	_, err = l.svcCtx.RuleModel.FindDataExists(l.ctx, req.FilterId, req.Name, deptPathIds)
	// 如果找到数据，err 为 nil；如果没找到或出错，err 不为 nil
	exists := err == nil

	return &types.DataExistsResp{
		Exists: exists,
	}, nil
}
