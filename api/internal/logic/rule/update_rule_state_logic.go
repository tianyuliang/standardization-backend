// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"
	"strings"
	"time"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/rule/mock"
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
// 对应 Java: RuleServiceImpl.updateState() (lines 753-785)
// 业务流程:
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
	// 对应 Java: ruleMapper.selectById(id) (lines 754-757)
	exist, err := l.svcCtx.RuleModel.FindOne(l.ctx, id)
	if err != nil || exist == nil {
		return nil, errorx.RuleDataNotExist()
	}

	// ====== 步骤2: 停用时必须填写原因 ======
	// 对应 Java: if (EnableDisableStatusEnum.DISABLE.equals(state)) (lines 759-764)
	targetState := rulemodel.GetStateInt(req.State)
	if targetState == rulemodel.StateDisable {
		if strings.TrimSpace(req.Reason) == "" {
			return nil, errorx.RuleDisableReasonEmpty()
		}
	}

	// ====== 步骤3: 停用原因长度校验 ======
	// 对应 Java: if (disableReason.length() > 800) (lines 765-770)
	if targetState == rulemodel.StateDisable && len([]rune(req.Reason)) > 800 {
		return nil, errorx.RuleDisableReasonTooLong()
	}

	// ====== 步骤4: 更新状态 ======
	// 对应 Java: ruleMapper.updateState(id, state, disableReason) (line 778)
	// MOCK: mock.GetUserInfo() - 从 Token 获取用户信息
	_, updateUser := mock.GetUserInfo(l.ctx)

	var disableReason string
	if targetState == rulemodel.StateDisable {
		disableReason = req.Reason
	}

	err = l.svcCtx.RuleModel.UpdateState(l.ctx, id, targetState, disableReason)
	if err != nil {
		return nil, err
	}

	// ====== 步骤5: 发送MQ消息 ======
	// 对应 Java: packageMqInfo(Arrays.asList(exist), "update") (line 781)
	//            kafkaProducerService.sendMessage(MqTopic.MQ_MESSAGE_SAILOR, mqInfo) (line 783)
	// 构建更新后的规则用于MQ消息
	updatedRule := &rulemodel.Rule{
		Id:            exist.Id,
		Name:          exist.Name,
		CatalogId:     exist.CatalogId,
		OrgType:       exist.OrgType,
		Description:   exist.Description,
		RuleType:      exist.RuleType,
		Version:       exist.Version,
		Expression:    exist.Expression,
		State:         targetState,
		DisableReason: disableReason,
		AuthorityId:   exist.AuthorityId,
		DepartmentIds: exist.DepartmentIds,
		ThirdDeptId:   exist.ThirdDeptId,
		UpdateTime:    time.Now(),
		UpdateUser:    updateUser,
	}

	if err := SendRuleMQMessage([]*rulemodel.Rule{updatedRule}, "update"); err != nil {
		logx.Errorf("发送MQ消息失败: %v", err)
	}
	logx.Infof("编码规则状态更新成功: id=%d, state=%d", id, targetState)

	return &types.EmptyResp{}, nil
}
