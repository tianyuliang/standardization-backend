// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

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
// 业务流程（参考 specs/编码规则管理接口流程说明_20260204.md 第4.3节）:
//  1. 参数校验（Ids非空、CatalogId有效）
//  2. 目录存在性校验
//  3. 批量更新 catalog_id
//  4. 版本号 +1
//  5. 记录更新用户
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
		// TODO: 返回 errorx.RuleIdsEmpty()
		return nil, err
	}

	// ====== 步骤2: 目录存在性校验 ======
	// TODO: 调用 CheckCatalogIdExist(req.CatalogId)
	// - 当前返回 mock 数据
	// - 如果不存在，返回 errorx.RuleCatalogNotExist(req.CatalogId) [错误码 30312]

	// ====== 步骤3: 批量更新 catalog_id、版本号+1、记录更新用户 ======
	updateUser := "" // TODO: 从 Token 获取
	err = l.svcCtx.RuleModel.RemoveCatalog(l.ctx, req.Ids, req.CatalogId, updateUser)
	if err != nil {
		return nil, err
	}

	// ====== 步骤4: 发送MQ消息 ======
	// TODO: 调用 SendRuleMQMessage(producer, rules, "update")
	// - MQ Topic: MQ_MESSAGE_SAILOR

	return &types.EmptyResp{}, nil
}
