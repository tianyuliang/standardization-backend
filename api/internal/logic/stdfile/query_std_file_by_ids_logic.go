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

type QueryStdFileByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据id列表查询
//
// 业务流程:
//  1. 参数校验（Ids非空）
//  2. 根据ID列表查询
//  3. 转换为响应对象
//
// 异常处理:
//   - 30202: ID列表不能为空
func NewQueryStdFileByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryStdFileByIdsLogic {
	return &QueryStdFileByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryStdFileByIdsLogic) QueryStdFileByIds(req *types.QueryByIdsReq) (resp *types.StdFileDataListResp, err error) {
	// Step 1: 参数校验
	if len(req.Ids) == 0 {
		return nil, errorx.NewWithMsg(30202, "ID列表不能为空")
	}

	// Step 2: 根据ID列表查询
	models, err := l.svcCtx.StdFileModel.FindByIds(l.ctx, req.Ids)
	if err != nil {
		return nil, HandleError(err)
	}

	// Step 3: 转换为响应对象
	data := ModelsToResp(l.ctx, l.svcCtx, models)

	logx.Infof("根据ID列表查询标准文件成功: count=%d", len(data))

	return &types.StdFileDataListResp{
		TotalCount: int64(len(data)),
		Data:       data,
	}, nil
}
