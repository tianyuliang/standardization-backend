// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

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
// 业务流程:
//  1. catalog_id = -1: 返回未关联文件的规则
//  2. 校验是否为标准文件目录
//  3. 顶级目录: 返回所有规则
//  4. 获取子目录列表
//  5. 查询关联该目录下文件的所有规则
func NewQueryRuleByStdFileCatalogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRuleByStdFileCatalogLogic {
	return &QueryRuleByStdFileCatalogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRuleByStdFileCatalogLogic) QueryRuleByStdFileCatalog(req *types.QueryByStdFileCatalogReq) (resp *types.RuleListResp, err error) {
	// ====== 步骤1: catalog_id = -1: 返回未关联文件的规则 ======
	if req.CatalogId == -1 {
		return l.findRulesNotUsedStdFile(req)
	}

	// ====== 步骤2: 校验是否为标准文件目录 ======
	// TODO: 调用 Catalog RPC 校验目录类型
	// if !isStdFileCatalog(req.CatalogId) { return emptyList() }

	// ====== 步骤3: 顶级目录: 返回所有规则 ======
	// TODO: 判断是否为根目录
	// if isRootCatalog(req.CatalogId) { return l.findAllRules(req) }

	// ====== 步骤4-5: 获取子目录列表并查询规则 ======
	// TODO: 调用 Catalog RPC 获取子目录列表
	// catalogIds := []int64{req.CatalogId}

	// TODO: 查询关联该目录下文件的规则ID列表
	// 当前返回空列表
	return &types.RuleListResp{
		Entries:    []types.RuleResp{},
		TotalCount: 0,
	}, nil
}

// findRulesNotUsedStdFile 查找未关联文件的规则
func (l *QueryRuleByStdFileCatalogLogic) findRulesNotUsedStdFile(req *types.QueryByStdFileCatalogReq) (*types.RuleListResp, error) {
	// TODO: 实现查询未关联文件的规则
	return &types.RuleListResp{
		Entries:    []types.RuleResp{},
		TotalCount: 0,
	}, nil
}
