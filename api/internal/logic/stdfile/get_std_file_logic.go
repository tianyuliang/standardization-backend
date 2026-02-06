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

type GetStdFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据ID查询详情
func NewGetStdFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStdFileLogic {
	return &GetStdFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStdFileLogic) GetStdFile(req *types.StdFileDetailResp) (resp *types.StdFileDetailResp, err error) {
	// Step 1: 从响应对象中获取ID（path参数通过handler注入）
	// 注意: id字段会由handler从path参数中注入
	if req.Id <= 0 {
		return nil, errorx.NewWithMsg(30202, "文件ID不能为空")
	}

	// Step 2: 查询文件
	model, err := l.svcCtx.StdFileModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, errorx.NewWithMsg(30201, "标准文件不存在")
	}
	if model == nil {
		return nil, errorx.NewWithMsg(30201, "标准文件不存在")
	}

	// Step 3: 转换为响应对象
	resp = &types.StdFileDetailResp{}
	*resp = ModelToResp(l.ctx, l.svcCtx, model)

	logx.Infof("查询标准文件详情成功: id=%d", req.Id)
	return resp, nil
}
