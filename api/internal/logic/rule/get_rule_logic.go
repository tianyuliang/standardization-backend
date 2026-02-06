// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"
	"encoding/json"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/rule/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	rulemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/rule/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 编码规则详情查看
//
// 对应 Java: RuleServiceImpl.queryById(id) (lines 284-314)
// 业务流程:
//  1. 查询规则基本信息
//  2. 如果已删除，状态设为停用
//  3. 查询目录名称
//  4. 根据ruleType解析expression（返回regex或custom）
//  5. 查询关联文件列表
//  6. 解析部门信息（部门ID、名称、路径名称）
//  7. 查询是否被引用
//
// 特殊说明：本接口无异常抛出，记录不存在时返回null（code=0）
func NewGetRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRuleLogic {
	return &GetRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRuleLogic) GetRule(id int64) (resp *types.RuleResp, err error) {
	// ====== 步骤1: 查询规则基本信息 ======
	// 对应 Java: ruleMapper.selectById(id) (line 285)
	ruleData, err := l.svcCtx.RuleModel.FindOne(l.ctx, id)
	if err != nil || ruleData == nil {
		// 记录不存在时返回null，不抛异常
		return nil, nil
	}

	// ====== 步骤2: 如果已删除，状态设为停用 ======
	// 对应 Java: if(source.getDeleted()) target.setState(EnableDisableStatusEnum.DISABLE) (lines 291-293)
	state := ruleData.State
	if ruleData.Deleted != 0 {
		state = rulemodel.StateDisable
	}

	// ====== 步骤3: 查询目录名称 ======
	// 对应 Java: iDeCatalogInfoService.getById(target.getCatalogId()) (lines 294-297)
	// MOCK: mock.CatalogGetCatalogName() - 获取目录名称
	catalogName := mock.CatalogGetCatalogName(l.ctx, l.svcCtx, ruleData.CatalogId)

	// ====== 步骤4: 根据ruleType解析expression ======
	// 对应 Java: if (RuleTypeEnum.CUSTOM.equals(target.getRuleType())) {...} (lines 298-302)
	expression := ruleData.Expression

	// ====== 步骤5: 查询关联文件列表 ======
	// 对应 Java: relationRuleFileMapper.queryByRuleId(id) (lines 304-309)
	relations, _ := l.svcCtx.RelationRuleFileModel.FindByRuleId(l.ctx, id)
	stdFileIds := make([]int64, len(relations))
	for i, r := range relations {
		stdFileIds[i] = r.FileId
	}

	// ====== 步骤6: 解析部门信息 ======
	// 对应 Java: String deptId = StringUtil.PathSplitAfter(source.getDepartmentIds()) (line 310)
	//            Map<String, Department> deptEntityMap = TokenUtil.getMapDeptInfo(Arrays.asList(deptId)) (line 311)
	// TODO: 解析部门ID、部门名称、路径名称
	// departmentId := ""
	// departmentName := ""
	// departmentPathNames := ""

	// ====== 步骤7: 查询是否被引用 ======
	// 对应 Java: (在其他地方调用 iDataElementInfoService.ruleUsed)
	// MOCK: mock.DataElementRuleUsed() - 检查规则是否被引用
	usedFlag := mock.DataElementRuleUsed(l.ctx, l.svcCtx, id)

	// ====== 步骤8: 构建响应 ======
	// 对应 Java: CustomUtil.copyProperties(source, target) (line 290)
	return &types.RuleResp{
		Id:            ruleData.Id,
		Name:          ruleData.Name,
		CatalogId:     ruleData.CatalogId,
		CatalogName:   catalogName,
		OrgType:       ruleData.OrgType,
		Description:   ruleData.Description,
		RuleType:      intToRuleType(ruleData.RuleType),
		Version:       ruleData.Version,
		Expression:    expression,
		State:         state,
		DisableReason: ruleData.DisableReason,
		AuthorityId:   ruleData.AuthorityId,
		DepartmentIds: ruleData.DepartmentIds,
		ThirdDeptId:   ruleData.ThirdDeptId,
		StdFileIds:    stdFileIds,
		Used:          usedFlag,
		CreateTime:    timeToStr(ruleData.CreateTime),
		CreateUser:    ruleData.CreateUser,
		UpdateTime:    timeToStr(ruleData.UpdateTime),
		UpdateUser:    ruleData.UpdateUser,
	}, nil
}

// ====== 辅助函数 ======

// parseCustomFromExpression 将expression JSON解析为custom数组
func parseCustomFromExpression(expression string) ([]types.RuleCustom, error) {
	if expression == "" {
		return []types.RuleCustom{}, nil
	}
	var custom []types.RuleCustom
	err := json.Unmarshal([]byte(expression), &custom)
	if err != nil {
		return nil, err
	}
	return custom, nil
}
