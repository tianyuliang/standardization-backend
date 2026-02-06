// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryRuleUsedDataElementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询引用规则的数据元
//
// 业务流程:
//  1. 校验规则存在
//  2. 查询引用的数据元
//
// 异常处理:
//   - 30301: 规则不存在
func NewQueryRuleUsedDataElementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRuleUsedDataElementLogic {
	return &QueryRuleUsedDataElementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRuleUsedDataElementLogic) QueryRuleUsedDataElement(id int64, req *types.PageQuery) (resp *types.DataElementListResp, err error) {
	// ====== 步骤1: 校验规则存在 ======
	_, err = l.svcCtx.RuleModel.FindOne(l.ctx, id)
	if err != nil {
		// TODO: 返回 errorx.RuleNotExist(id) [错误码 30301]
		return nil, err
	}

	// ====== 步骤2: 查询引用的数据元 ======
	// TODO: 调用 DataElement RPC 查询引用该规则的数据元
	// 当前返回空列表
	return &types.DataElementListResp{
		Entries:    []types.DataElementResp{},
		TotalCount: 0,
	}, nil
}
