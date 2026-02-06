// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/rule/mock"
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
// 对应 Java: RuleServiceImpl.queryUsedDataElementByRuleId(Long id, Integer offset, Integer limit) (lines 642-651)
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
	// 对应 Java: ruleMapper.selectById(id) (lines 643-646)
	_, err = l.svcCtx.RuleModel.FindOne(l.ctx, id)
	if err != nil {
		return nil, errorx.RuleRecordNotExist()
	}

	// ====== 步骤2: 查询引用的数据元 ======
	// 对应 Java: iDataElementInfoService.queryByRuleId(ruleEntity.getId(), offset, limit) (line 647)
	// MOCK: mock.DataElementQueryByRuleId() - 查询引用该规则的数据元
	_, totalCount, err := mock.DataElementQueryByRuleId(l.ctx, l.svcCtx, id, req.Offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// 构建响应（TODO: 转换为 DataElementResp）
	return &types.DataElementListResp{
		Entries:    []types.DataElementResp{},
		TotalCount: totalCount,
	}, nil
}
