// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package stdfile

import (
	"context"

	"github.com/jinguoxing/idrm-go-base/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/stdfile/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据标准文件ID添加关联关系
//
// 业务流程:
//  1. 校验文件存在性
//  2. 事务性添加关联关系（DE、Dict、Rule）
//
// 异常处理:
//   - 30201: 标准文件不存在
func NewAddRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRelationLogic {
	return &AddRelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddRelationLogic) AddRelation(id int64, req *types.StdFileRelationDto) (resp *types.BaseResp, err error) {
	// Step 1: 校验文件存在性
	existing, err := l.svcCtx.StdFileModel.FindOne(l.ctx, id)
	if err != nil {
		return nil, HandleError(err)
	}
	if existing == nil || existing.Deleted != 0 {
		return nil, errorx.NewWithMsg(30201, "标准文件不存在")
	}

	// Step 2: 事务性添加关联关系 (Mock)
	// Java: transactionTemplate.execute(status -> {...})
	// 调用三个服务添加关联关系
	if len(req.DeIds) > 0 {
		if err := mock.DataElementAddRelation(l.ctx, l.svcCtx, id, req.DeIds); err != nil {
			return nil, err
		}
	}

	if len(req.DictIds) > 0 {
		if err := mock.DictAddRelation(l.ctx, l.svcCtx, id, req.DictIds); err != nil {
			return nil, err
		}
	}

	if len(req.RuleIds) > 0 {
		if err := mock.RuleAddRelation(l.ctx, l.svcCtx, id, req.RuleIds); err != nil {
			return nil, err
		}
	}

	logx.Infof("添加标准文件关联关系成功: fileId=%d, deIds=%v, dictIds=%v, ruleIds=%v",
		id, req.DeIds, req.DictIds, req.RuleIds)

	return &types.BaseResp{
		Code:        "0",
		Description: "操作成功",
	}, nil
}
