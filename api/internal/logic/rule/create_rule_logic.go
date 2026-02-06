// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"
	"time"

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
// 业务流程（参考 specs/编码规则管理接口流程说明_20260204.md 第4.1节）:
//  1. 参数校验（Handler已完成基础校验，这里处理业务相关校验）
//  2. 表达式校验
//     - REGEX类型：校验正则表达式非空且格式正确
//     - CUSTOM类型：校验custom非空、segment_length>0、value有效性
//  3. 名称唯一性校验（同一orgType下）
//  4. 目录存在性校验
//  5. 部门ID处理（从Token获取完整路径）
//  6. 数据处理（生成表达式、设置创建信息）
//  7. 保存数据库（t_rule + t_relation_rule_file）
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
	// TODO: 调用 ValidateExpression(ruleType, regex, custom)
	// - REGEX: 校验 regex 非空 + regexp.Compile()
	// - CUSTOM: 校验 custom 非空 + 遍历校验

	// ====== 步骤2: 名称唯一性校验 ======
	// TODO: 调用 CheckNameUnique(name, orgType, departmentIds)
	// - 查询数据库：SELECT * FROM t_rule WHERE f_name=? AND f_org_type=? AND f_deleted=0 AND f_department_ids=?
	// - 如果存在记录，返回 errorx.RuleNameDuplicate(name) [错误码 30310]

	// ====== 步骤3: 目录存在性校验 ======
	// TODO: 调用 Catalog RPC 校验 catalogId 是否存在
	// - 当前返回 mock 数据
	// - 如果不存在，返回 errorx.RuleCatalogNotExist(req.CatalogId) [错误码 30312]

	// ====== 步骤4: 部门ID处理 ======
	// TODO: 从 Token 解析部门完整路径
	// - 当前使用 req.DepartmentIds 原值
	// - 后续从 Token 获取部门信息

	// ====== 步骤5: 数据处理 ======
	expression := getExpression(req.RuleType, req.Regex, req.Custom)
	now := time.Now()

	ruleData := &rulemodel.Rule{
		Name:          req.Name,
		CatalogId:     req.CatalogId,
		OrgType:       req.OrgType,
		Description:   req.Description,
		RuleType:      ruleTypeToInt(req.RuleType),
		Version:       1,
		Expression:    expression,
		State:         rulemodel.StateEnable, // 默认启用状态
		DepartmentIds: req.DepartmentIds,
		CreateTime:    now,
		UpdateTime:    now,
		Deleted:       0,
	}

	// ====== 步骤6: 保存数据库 ======
	// TODO: 开启事务
	// 6.1 插入 t_rule
	id, err := l.svcCtx.RuleModel.Insert(l.ctx, ruleData)
	if err != nil {
		return nil, err
	}
	ruleData.Id = id

	// 6.2 保存关联文件 t_relation_rule_file
	if len(req.StdFileIds) > 0 {
		// TODO: 构建关联文件数据
		// TODO: 调用 RelationRuleFileModel.InsertBatch()
		// 注意：最多关联10个标准文件
		_ = len(req.StdFileIds) // TODO: Implement file relation insertion (suppress staticcheck SA9003)
	}

	// ====== 步骤7: 发送MQ消息 ======
	// TODO: 调用 SendRuleMQMessage(producer, []ruleData, "insert")
	// - MQ Topic: MQ_MESSAGE_SAILOR
	// - 消息格式: { "header": {}, "payload": { "type": "smart-recommendation-graph", "content": { "type": "insert", "tableName": "t_rule", "entities": [...] } } }

	// ====== 步骤8: 构建响应 ======
	resp = buildRuleResp(ruleData, "", false, req.StdFileIds)
	return
}
