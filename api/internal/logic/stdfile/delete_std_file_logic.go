// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package stdfile

import (
	"context"

	"github.com/jinguoxing/idrm-go-base/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteStdFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标准文件管理-批量删除
//
// 业务流程:
//  1. 参数校验（Ids非空）
//  2. 批量软删除
//
// 异常处理:
//   - 30202: ID列表不能为空
func NewDeleteStdFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteStdFileLogic {
	return &DeleteStdFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteStdFileLogic) DeleteStdFile(ids []int64) (resp *types.BaseResp, err error) {
	// Step 1: 参数校验
	if len(ids) == 0 {
		return nil, errorx.NewWithMsg(30202, "ID列表不能为空")
	}

	// Step 2: 批量软删除
	if err := l.svcCtx.StdFileModel.DeleteByIds(l.ctx, ids); err != nil {
		return nil, HandleError(err)
	}

	logx.Infof("标准文件批量删除成功: count=%d", len(ids))

	return &types.BaseResp{
		Code:        "0",
		Description: "删除成功",
	}, nil
}
