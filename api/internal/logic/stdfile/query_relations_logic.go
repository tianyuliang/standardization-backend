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

type QueryRelationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标准文件关联关系查询
//
// 业务流程:
//  1. 校验文件存在性
//  2. 查询所有关联关系（DE、Dict、Rule）
//  3. 返回关联ID列表
//
// 异常处理:
//   - 30201: 标准文件不存在
func NewQueryRelationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRelationsLogic {
	return &QueryRelationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRelationsLogic) QueryRelations(id int64) (resp *types.StdFileRelationResp, err error) {
	// Step 1: 校验文件存在性
	existing, err := l.svcCtx.StdFileModel.FindOne(l.ctx, id)
	if err != nil {
		return nil, HandleError(err)
	}
	if existing == nil || existing.Deleted != 0 {
		return nil, errorx.NewWithMsg(30201, "标准文件不存在")
	}

	// Step 2: 查询所有关联关系 (Mock)
	// Java: dataElementInfoService.queryByFileId(id)
	//       dictService.queryByFileId(id)
	//       ruleService.queryByFileId(id)
	deIds := mock.DataElementQueryByFileId(l.ctx, l.svcCtx, id)
	dictIds := mock.DictQueryByFileId(l.ctx, l.svcCtx, id)
	ruleIds := mock.RuleQueryByFileId(l.ctx, l.svcCtx, id)

	logx.Infof("查询标准文件关联关系成功: fileId=%d, deCount=%d, dictCount=%d, ruleCount=%d",
		id, len(deIds), len(dictIds), len(ruleIds))

	return &types.StdFileRelationResp{
		DeIds:   deIds,
		DictIds: dictIds,
		RuleIds: ruleIds,
	}, nil
}
