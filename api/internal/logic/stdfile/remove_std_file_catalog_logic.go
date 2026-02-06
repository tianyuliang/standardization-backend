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

type RemoveStdFileCatalogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 移动到指定目录-目录移动
//
// 业务流程:
//  1. 参数校验（Ids非空、CatalogId有效）
//  2. 目录存在性校验
//  3. 批量更新 catalog_id、记录更新用户
//
// 异常处理:
//   - 30202: ID列表不能为空
//   - 30204: 目录不存在
func NewRemoveStdFileCatalogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveStdFileCatalogLogic {
	return &RemoveStdFileCatalogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveStdFileCatalogLogic) RemoveStdFileCatalog(req *types.RemoveCatalogReq) (resp *types.BaseResp, err error) {
	// Step 1: 参数校验
	if len(req.Ids) == 0 {
		return nil, errorx.NewWithMsg(30202, "ID列表不能为空")
	}

	// Step 2: 目录存在性校验
	if !mock.CatalogCheckExist(l.ctx, l.svcCtx, req.CatalogId) {
		return nil, errorx.NewWithMsg(30204, "目录不存在")
	}

	// Step 3: 批量更新 catalog_id
	_, updateUser := mock.GetUserInfo(l.ctx)

	err = l.svcCtx.StdFileModel.RemoveCatalog(l.ctx, req.Ids, req.CatalogId, updateUser)
	if err != nil {
		return nil, HandleError(err)
	}

	logx.Infof("标准文件目录移动成功: ids=%v, catalogId=%d", req.Ids, req.CatalogId)

	return &types.BaseResp{
		Code:        "0",
		Description: "操作成功",
	}, nil
}
