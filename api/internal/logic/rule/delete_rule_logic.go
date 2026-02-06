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
	rulemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/rule/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 编码规则删除&批量删除
//
// 对应 Java: RuleServiceImpl.deleteBatch(ids) (lines 609-630)
// 业务流程:
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
	// 对应 Java: String[] idArray = ids.split(",") (line 610)
	ids := parseIds(idsStr)

	// ====== 步骤2: 校验ids不为空 ======
	// 对应 Java: if (CustomUtil.isEmpty(idList)) (implicitly checked)
	if len(ids) == 0 {
		return nil, errorx.RuleIdsEmpty()
	}

	// ====== 步骤3: 查询待删除的规则（用于MQ消息） ======
	// 对应 Java: 构建删除实体列表 (lines 619-625)
	rules, _ := l.svcCtx.RuleModel.FindByIds(l.ctx, ids)
	mqRules := make([]*rulemodel.Rule, len(rules))
	for i, id := range ids {
		mqRules[i] = &rulemodel.Rule{Id: id}
	}

	// TODO: 开启事务

	// ====== 步骤4: 物理删除t_rule记录 ======
	// 对应 Java: ruleMapper.deleteByIds(idList) (line 618)
	err = l.svcCtx.RuleModel.DeleteByIds(l.ctx, ids)
	if err != nil {
		return nil, err
	}

	// ====== 步骤5: 同步删除t_relation_rule_file关联记录 ======
	// 对应 Java: (implicitly handled by foreign key or separate call)
	err = l.svcCtx.RelationRuleFileModel.DeleteByRuleIds(l.ctx, ids)
	if err != nil {
		logx.Errorf("删除关联文件失败: %v", err)
	}

	// ====== 步骤6: 发送MQ消息 ======
	// 对应 Java: packageMqInfo(list, "delete") (line 626)
	//            kafkaProducerService.sendMessage(MqTopic.MQ_MESSAGE_SAILOR, mqInfo) (line 628)
	if err := SendRuleMQMessage(mqRules, "delete"); err != nil {
		logx.Errorf("发送MQ消息失败: %v", err)
	}
	logx.Infof("编码规则删除成功: ids=%v", ids)

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
		if err == nil && id > 0 {
			ids = append(ids, id)
		}
	}
	return ids
}
