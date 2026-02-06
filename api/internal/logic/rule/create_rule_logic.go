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

type CreateRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增编码规则
//
// 对应 Java: RuleServiceImpl.create() (lines 319-348)
// 业务流程:
//  1. 表达式校验 (checkRuleExpression)
//  2. 名称唯一性校验 (checkAddDataExist)
//  3. 目录存在性校验 (checkCatalogIdExist)
//  4. 部门ID处理 (从Token获取完整路径)
//  5. 数据处理 (生成表达式、设置创建信息)
//  6. 保存数据库 (t_rule)
//  7. 保存关联文件 (t_relation_rule_file)
//  8. 发送MQ消息
//
// 异常处理：
//   - 30310: 规则名称已存在
//   - 30312: 目录不存在
//   - 30320: 正则表达式为空
//   - 30321: 正则表达式非法
//   - 30330: 自定义配置为空
//   - 30331: segment_length <= 0
//   - 30332: 码表不存在
//   - 30333: 日期格式不支持
//   - 30334: value为空
func NewCreateRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRuleLogic {
	return &CreateRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRuleLogic) CreateRule(req *types.CreateRuleReq) (resp *types.RuleResp, err error) {
	// ====== 步骤1: 表达式校验 ======
	// 对应 Java: RuleServiceImpl.checkRuleExpression() (lines 381-438)
	if err := ValidateExpression(req.RuleType, req.Regex, req.Custom); err != nil {
		return nil, err
	}

	// ====== 步骤2: 名称唯一性校验 ======
	// 对应 Java: RuleServiceImpl.checkAddDataExist() (lines 676-683)
	if err := CheckNameUnique(l.ctx, l.svcCtx, req.Name, req.OrgType, req.DepartmentIds); err != nil {
		return nil, err
	}

	// ====== 步骤3: 目录存在性校验 ======
	// 对应 Java: RuleServiceImpl.checkCatalogIdExist() (lines 699-706)
	// MOCK: mock.CatalogCheckExist() - 校验目录是否存在
	if !mock.CatalogCheckExist(l.ctx, l.svcCtx, req.CatalogId) {
		return nil, errorx.RuleCatalogNotExist(req.CatalogId)
	}

	// ====== 步骤4: 部门ID处理 ======
	// 对应 Java: TokenUtil.getDeptPathIds() (line 334)
	// MOCK: mock.GetDeptPathIds() - 从 Token/部门服务获取完整路径
	departmentIds, thirdDeptId := mock.GetDeptPathIds(l.ctx, req.DepartmentIds)

	// ====== 步骤5: 数据处理 ======
	// 对应 Java: RuleServiceImpl.getExpression() (lines 449-455)
	expression := getExpression(req.RuleType, req.Regex, req.Custom)
	now := time.Now()

	// MOCK: mock.GetUserInfo() - 从 Token 获取用户信息
	authorityId, createUser := mock.GetUserInfo(l.ctx)
	updateUser := createUser

	ruleData := &rulemodel.Rule{
		Name:          req.Name,
		CatalogId:     req.CatalogId,
		OrgType:       req.OrgType,
		Description:   req.Description,
		RuleType:      ruleTypeToInt(req.RuleType),
		Version:       1,
		Expression:    expression,
		State:         rulemodel.StateEnable, // 默认启用状态
		AuthorityId:   authorityId,
		DepartmentIds: departmentIds,
		ThirdDeptId:   thirdDeptId,
		CreateTime:    now,
		UpdateUser:    updateUser,
		UpdateTime:    now,
		Deleted:       0,
	}
	ruleData.CreateUser = createUser

	// ====== 步骤6: 保存数据库 ======
	// 对应 Java: ruleMapper.insert(insert) (line 338)
	// TODO: 开启事务
	id, err := l.svcCtx.RuleModel.Insert(l.ctx, ruleData)
	if err != nil {
		return nil, err
	}
	ruleData.Id = id

	// ====== 步骤7: 保存关联文件 ======
	// 对应 Java: RuleServiceImpl.saveRelationRuleFile() (lines 457-469)
	if len(req.StdFileIds) > 0 {
		// 注意：最多关联10个标准文件
		if len(req.StdFileIds) > 10 {
			logx.Infof("标准文件数量超过10个: %d", len(req.StdFileIds))
		}
		if err := SaveRelationRuleFile(l.ctx, l.svcCtx, ruleData.Id, req.StdFileIds); err != nil {
			logx.Errorf("保存关联文件失败: %v", err)
			// TODO: 是否需要回滚事务
		}
	}

	// ====== 步骤8: 发送MQ消息 ======
	// 对应 Java: RuleServiceImpl.packageMqInfo() + kafkaProducerService.sendMessage() (lines 344-346)
	if err := SendRuleMQMessage([]*rulemodel.Rule{ruleData}, "insert"); err != nil {
		logx.Errorf("发送MQ消息失败: %v", err)
	}
	logx.Infof("编码规则新增成功: id=%d, name=%s", ruleData.Id, ruleData.Name)

	// ====== 步骤9: 构建响应 ======
	// 对应 Java: CustomUtil.copyProperties(insert, target) (line 342)
	// TODO: 查询目录名称、引用状态
	resp = buildRuleResp(ruleData, "", false, req.StdFileIds)
	return
}
