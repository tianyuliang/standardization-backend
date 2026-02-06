// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"
	"time"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/rule/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	rulemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/rule/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改编码规则
//
// 对应 Java: RuleServiceImpl.update() (lines 474-515)
// 业务流程:
//  1. 校验规则是否存在
//  2. 名称唯一性校验（排除自身）
//  3. 目录存在性校验
//  4. 版本变更检测
//  5. 无变更直接返回原数据
//  6. 有变更则更新数据（版本号+1、更新关联文件）
//  7. 发送MQ消息
//
// 异常处理：
//   - 30301: 记录不存在
//   - 30311: 修改时名称与其他记录重复
//   - 30312: 目录不存在
func NewUpdateRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRuleLogic {
	return &UpdateRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRuleLogic) UpdateRule(id int64, req *types.UpdateRuleReq) (resp *types.RuleResp, err error) {
	// ====== 步骤1: 校验规则是否存在 ======
	// 对应 Java: ruleMapper.selectById(id) (lines 475-478)
	oldRule, err := l.svcCtx.RuleModel.FindOne(l.ctx, id)
	if err != nil || oldRule == nil {
		return nil, errorx.RuleRecordNotExist()
	}

	// ====== 步骤2: 表达式校验 ======
	// 对应 Java: checkRuleExpression() (line 479)
	if err := ValidateExpression(req.RuleType, req.Regex, req.Custom); err != nil {
		return nil, err
	}

	// ====== 步骤3: 名称唯一性校验（排除自身）=====
	// 对应 Java: checkUpdateDataExist(id, updateDto) (line 479)
	if err := CheckNameUniqueForUpdate(l.ctx, l.svcCtx, id, req.Name, req.OrgType, req.DepartmentIds); err != nil {
		return nil, err
	}

	// ====== 步骤4: 目录存在性校验 ======
	// 对应 Java: checkCatalogIdExist(updateDto.getCatalogId()) (line 480)
	// MOCK: mock.CatalogCheckExist() - 校验目录是否存在
	if !mock.CatalogCheckExist(l.ctx, l.svcCtx, req.CatalogId) {
		return nil, errorx.RuleCatalogNotExist(req.CatalogId)
	}

	// ====== 步骤5: 版本变更检测 ======
	// 对应 Java: checkVersionChange(exist, updateDto) (line 482)
	// 查询原有关联文件
	oldFiles, _ := l.svcCtx.RelationRuleFileModel.FindByRuleId(l.ctx, id)
	oldRelationFiles := make([]*RelationRuleFile, len(oldFiles))
	for i, f := range oldFiles {
		oldRelationFiles[i] = &RelationRuleFile{
			Id:     f.Id,
			RuleId: f.RuleId,
			FileId: f.FileId,
		}
	}

	needVersionIncrement := CheckVersionChange(oldRule, req, oldRelationFiles)

	if !needVersionIncrement {
		// ====== 步骤6: 无变更直接返回原数据 ======
		// 对应 Java: return Result.success(exist) (lines 483-485)
		// TODO: 查询目录名称、引用状态、关联文件
		return buildRuleResp(oldRule, "", false, req.StdFileIds), nil
	}

	// ====== 步骤7: 有变更则更新数据 ======
	// 对应 Java: exist.setXXX() + ruleMapper.updateById(exist) (lines 486-504)

	// MOCK: mock.GetDeptPathIds() - 从 Token/部门服务获取完整路径
	departmentIds, thirdDeptId := mock.GetDeptPathIds(l.ctx, req.DepartmentIds)

	// MOCK: mock.GetUserInfo() - 从 Token 获取用户信息
	_, updateUser := mock.GetUserInfo(l.ctx)

	// 构建更新数据
	newRule := &rulemodel.Rule{
		Id:            oldRule.Id,
		Name:          req.Name,
		CatalogId:     req.CatalogId,
		OrgType:       req.OrgType,
		Description:   req.Description,
		RuleType:      ruleTypeToInt(req.RuleType),
		Version:       oldRule.Version + 1, // 版本号+1
		Expression:    getExpression(req.RuleType, req.Regex, req.Custom),
		State:         oldRule.State, // 保持原有状态
		DepartmentIds: departmentIds,
		ThirdDeptId:   thirdDeptId,
		UpdateTime:    time.Now(),
		UpdateUser:    updateUser,
		// 继承原有字段
		CreateTime:    oldRule.CreateTime,
		CreateUser:    oldRule.CreateUser,
		DisableReason: oldRule.DisableReason,
		AuthorityId:   oldRule.AuthorityId,
		Deleted:       oldRule.Deleted,
	}

	// TODO: 开启事务

	// 7.1 更新 t_rule
	// 对应 Java: ruleMapper.updateById(exist) (line 504)
	err = l.svcCtx.RuleModel.Update(l.ctx, newRule)
	if err != nil {
		return nil, err
	}

	// 7.2 更新关联文件 t_relation_rule_file
	// 对应 Java: relationRuleFileMapper.deleteByRuleId(exist.getId()) (line 506)
	//            saveRelationRuleFile(exist.getId(), updateDto.getStdFiles()) (line 507)
	if err := l.svcCtx.RelationRuleFileModel.DeleteByRuleId(l.ctx, id); err != nil {
		logx.Errorf("删除旧关联文件失败: %v", err)
	}

	if len(req.StdFileIds) > 0 {
		if err := SaveRelationRuleFile(l.ctx, l.svcCtx, id, req.StdFileIds); err != nil {
			logx.Errorf("保存新关联文件失败: %v", err)
		}
	}

	// ====== 步骤8: 发送MQ消息 ======
	// 对应 Java: packageMqInfo(Arrays.asList(exist), "update") (line 511)
	//            kafkaProducerService.sendMessage(MqTopic.MQ_MESSAGE_SAILOR, mqInfo) (line 513)
	if err := SendRuleMQMessage([]*rulemodel.Rule{newRule}, "update"); err != nil {
		logx.Errorf("发送MQ消息失败: %v", err)
	}
	logx.Infof("编码规则修改成功: id=%d, name=%s, version=%d", newRule.Id, newRule.Name, newRule.Version)

	// ====== 步骤9: 构建响应 ======
	// 对应 Java: CustomUtil.copyProperties(exist, target) (line 510)
	// TODO: 查询目录名称、引用状态
	return buildRuleResp(newRule, "", false, req.StdFileIds), nil
}
