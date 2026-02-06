// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package stdfile

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

// 查询数据是否存在
//
// 业务流程:
//  1. 根据查询条件检查数据是否存在
//  2. 返回检查结果
//
// 说明: 用于前端表单重复性校验
func NewQueryDataExistsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryDataExistsLogic {
	return &QueryDataExistsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryDataExistsLogic) QueryDataExists(req *types.QueryDataExistsReq) (resp *types.BaseResp, err error) {
	// Step 1: 根据查询条件检查数据是否存在
	existing, err := l.svcCtx.StdFileModel.FindDataExists(l.ctx,
		req.FilterId,
		req.Number,
		int(req.OrgType),
		req.Name,
		req.DepartmentIds)

	// Step 2: 返回检查结果
	if err != nil {
		// 发生错误时，返回不存在（允许继续操作）
		return &types.BaseResp{
			Code:        "0",
			Description: "不存在",
		}, nil
	}

	if existing != nil {
		return &types.BaseResp{
			Code:        "1",
			Description: "已存在",
		}, nil
	}

	return &types.BaseResp{
		Code:        "0",
		Description: "不存在",
	}, nil
}
