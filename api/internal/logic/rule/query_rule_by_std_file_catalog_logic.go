// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/rule/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryRuleByStdFileCatalogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据标准文件目录查询规则
//
// 对应 Java: RuleServiceImpl.queryByStdFileCatalog() (lines 179-221)
// 业务流程:
//  1. catalog_id 为空: 返回空列表
//  2. catalog_id = -1: 返回未关联文件的规则
//  3. 校验是否为标准文件目录
//  4. 顶级目录: 返回所有规则
//  5. 获取子目录列表
//  6. 查询关联该目录下文件的所有规则
func NewQueryRuleByStdFileCatalogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRuleByStdFileCatalogLogic {
	return &QueryRuleByStdFileCatalogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRuleByStdFileCatalogLogic) QueryRuleByStdFileCatalog(req *types.QueryByStdFileCatalogReq) (resp *types.RuleListResp, err error) {
	// ====== 步骤1: catalog_id 为空: 返回空列表 ======
	// 对应 Java: if (CustomUtil.isEmpty(stdFileCatalogId)) (lines 189-191)
	if req.CatalogId == 0 {
		return &types.RuleListResp{
			Entries:    []types.RuleResp{},
			TotalCount: 0,
		}, nil
	}

	// ====== 步骤2: catalog_id = -1: 返回未关联文件的规则 ======
	// 对应 Java: if (-1 == stdFileCatalogId) (lines 199-201)
	if req.CatalogId == -1 {
		return l.findRulesNotUsedStdFile(req)
	}

	// ====== 步骤3: 校验是否为标准文件目录 ======
	// 对应 Java: if (catalog == null || !catalog.getType().equals(CatalogTypeEnum.File)) (lines 206-209)
	// MOCK: mock.CatalogIsStdFileCatalog() - 校验是否为标准文件目录
	if !mock.CatalogIsStdFileCatalog(l.ctx, l.svcCtx, req.CatalogId) {
		return &types.RuleListResp{
			Entries:    []types.RuleResp{},
			TotalCount: 0,
		}, nil
	}

	// ====== 步骤4-5: 获取子目录列表并查询规则 ======
	// 对应 Java: List<Long> catalogIds = iDeCatalogInfoService.getIDList(stdFileCatalogId) (line 217)
	//            ruleMapper.queryByStdFileCatalog(page, catalogIds, ...) (line 218)
	// MOCK: mock.CatalogGetChildIds() - 获取子目录列表
	_ = mock.CatalogGetChildIds(l.ctx, l.svcCtx, req.CatalogId)

	// TODO: 调用 Model 层的 queryByStdFileCatalog 方法
	// 当前返回空列表
	return &types.RuleListResp{
		Entries:    []types.RuleResp{},
		TotalCount: 0,
	}, nil
}

// findRulesNotUsedStdFile 查找未关联文件的规则
// 对应 Java: ruleMapper.queryDataNotUesdStdFile(page, keyword, orgType, state, ruleType) (line 200)
func (l *QueryRuleByStdFileCatalogLogic) findRulesNotUsedStdFile(req *types.QueryByStdFileCatalogReq) (*types.RuleListResp, error) {
	// TODO: 调用 Model 层的 queryDataNotUsedStdFile 方法
	return &types.RuleListResp{
		Entries:    []types.RuleResp{},
		TotalCount: 0,
	}, nil
}
