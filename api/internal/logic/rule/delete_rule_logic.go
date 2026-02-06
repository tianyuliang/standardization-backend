// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"
	"strconv"
	"strings"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 编码规则删除&批量删除
//
// 业务流程（参考 specs/编码规则管理接口流程说明_20260204.md 第4.5节）:
//  1. 解析ID列表（多个ID用英文逗号分隔）
//  2. 校验ids不为空
//  3. 物理删除t_rule记录
//  4. 同步删除t_relation_rule_file关联记录
//  5. 发送MQ消息（操作类型：delete）
//
// 异常处理：
//   - 30335: ids 不能为空
func NewDeleteRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRuleLogic {
	return &DeleteRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRuleLogic) DeleteRule(idsStr string) (resp *types.EmptyResp, err error) {
	// ====== 步骤1: 解析ID列表 ======
	ids := parseIds(idsStr)

	// ====== 步骤2: 校验ids不为空 ======
	if len(ids) == 0 {
		return nil, errorx.RuleIdsEmpty() // [错误码 30335]
	}

	// ====== 步骤3: 物理删除t_rule记录 ======
	// TODO: 开启事务
	err = l.svcCtx.RuleModel.DeleteByIds(l.ctx, ids)
	if err != nil {
		return nil, err
	}

	// ====== 步骤4: 同步删除t_relation_rule_file关联记录 ======
	err = l.svcCtx.RelationRuleFileModel.DeleteByRuleIds(l.ctx, ids)
	if err != nil {
		return nil, err
	}

	// ====== 步骤5: 发送MQ消息 ======
	// TODO: 调用 SendRuleMQMessage(producer, rules, "delete")
	// - 需要先查询被删除的规则用于MQ消息

	return &types.EmptyResp{}, nil
}

// ====== 辅助函数 ======

// parseIds 解析ID列表（多个ID用英文逗号分隔）
func parseIds(idsStr string) []int64 {
	if idsStr == "" {
		return []int64{}
	}

	parts := strings.Split(idsStr, ",")
	ids := make([]int64, 0, len(parts))
	for _, part := range parts {
		id, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
		if err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}
