// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"
	"strings"
	"time"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	rulemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/rule/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRuleStateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 停用/启用编码规则
//
// 业务流程（参考 specs/编码规则管理接口流程说明_20260204.md 第4.2节）:
//  1. 校验规则存在性
//  2. 停用时必须填写原因
//  3. 停用原因长度校验 (<= 800)
//  4. 更新状态（启用时清空停用原因）
//  5. 发送MQ消息
//
// 异常处理:
//   - 30301: 规则不存在
//   - 30302: 停用原因不能为空
//   - 30303: 停用原因长度不能超过800字符
func NewUpdateRuleStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRuleStateLogic {
	return &UpdateRuleStateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRuleStateLogic) UpdateRuleState(id int64, req *types.UpdateRuleStateReq) (resp *types.EmptyResp, err error) {
	// ====== 步骤1: 校验规则存在性 ======
	_, err = l.svcCtx.RuleModel.FindOne(l.ctx, id)
	if err != nil {
		// TODO: 返回 errorx.RuleNotExist(id) [错误码 30301]
		return nil, err
	}

	// ====== 步骤2: 停用时必须填写原因 ======
	targetState := rulemodel.GetStateInt(req.State)
	if targetState == rulemodel.StateDisable {
		if strings.TrimSpace(req.Reason) == "" {
			// TODO: 返回 errorx.RuleDisableReasonEmpty() [错误码 30302]
			return nil, err
		}
	}

	// ====== 步骤3: 停用原因长度校验 ======
	if targetState == rulemodel.StateDisable && len([]rune(req.Reason)) > 800 {
		// TODO: 返回 errorx.RuleDisableReasonTooLong() [错误码 30303]
		return nil, err
	}

	// ====== 步骤4: 更新状态 ======
	updateData := &rulemodel.Rule{
		Id:         id,
		State:      targetState,
		UpdateTime: time.Now(),
	}
	if targetState == rulemodel.StateDisable {
		updateData.DisableReason = req.Reason
	} else {
		updateData.DisableReason = ""
	}

	err = l.svcCtx.RuleModel.Update(l.ctx, updateData)
	if err != nil {
		return nil, err
	}

	// ====== 步骤5: 发送MQ消息 ======
	// TODO: 调用 SendRuleMQMessage(producer, []updateData, "update")
	// - MQ Topic: MQ_MESSAGE_SAILOR
	// - 消息格式: { "header": {}, "payload": { "type": "smart-recommendation-graph", "content": { "type": "update", "tableName": "t_rule", "entities": [...] } } }

	return &types.EmptyResp{}, nil
}
