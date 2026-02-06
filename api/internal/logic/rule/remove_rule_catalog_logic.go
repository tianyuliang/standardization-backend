// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/rule/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveRuleCatalogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 编码规则目录移动
//
// 对应 Java: RuleServiceImpl.removeCatalog(List<Long> ids, Long catalogId) (lines 634-639)
// 业务流程:
//  1. 参数校验（Ids非空、CatalogId有效）
//  2. 目录存在性校验
//  3. 批量更新 catalog_id、版本号+1、记录更新用户
//
// 异常处理:
//   - 30312: 目录不存在
func NewRemoveRuleCatalogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveRuleCatalogLogic {
	return &RemoveRuleCatalogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveRuleCatalogLogic) RemoveRuleCatalog(req *types.RemoveCatalogReq) (resp *types.EmptyResp, err error) {
	// ====== 步骤1: 参数校验 ======
	if len(req.Ids) == 0 {
		return nil, errorx.RuleIdsEmpty()
	}

	// ====== 步骤2: 目录存在性校验 ======
	// 对应 Java: checkCatalogIdExist(catalogId) (line 635)
	// MOCK: mock.CatalogCheckExist() - 校验目录是否存在
	if !mock.CatalogCheckExist(l.ctx, l.svcCtx, req.CatalogId) {
		return nil, errorx.RuleCatalogNotExist(req.CatalogId)
	}

	// ====== 步骤3: 批量更新 catalog_id、版本号+1、记录更新用户 ======
	// 对应 Java: ruleMapper.removeCatalog(ids, catalogId, userInfo.getUserName()) (line 637)
	// MOCK: mock.GetUserInfo() - 从 Token 获取用户信息
	_, updateUser := mock.GetUserInfo(l.ctx)

	err = l.svcCtx.RuleModel.RemoveCatalog(l.ctx, req.Ids, req.CatalogId, updateUser)
	if err != nil {
		return nil, err
	}

	logx.Infof("编码规则目录移动成功: ids=%v, catalogId=%d", req.Ids, req.CatalogId)
	return &types.EmptyResp{}, nil
}
