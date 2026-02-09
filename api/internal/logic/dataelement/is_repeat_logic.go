// Code scaffolded by goctl. Safe to edit.

package dataelement

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type IsRepeatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 检查名称重复
func NewIsRepeatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsRepeatLogic {
	return &IsRepeatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IsRepeatLogic) IsRepeat(req *types.IsRepeatReq) (resp *types.IsRepeatVo, err error) {
	// Step 1: 参数校验
	if req.NameCn == "" {
		return nil, errorx.ParameterEmpty("nameCn")
	}

	// Step 2: 检查中文名称是否存在
	// 如果没有指定StdType，则查询所有类型
	if req.StdType == 0 {
		// 查询所有标准类型中是否存在同名
		_, err := l.svcCtx.DataElementModel.FindDataExists(l.ctx, req.NameCn, 0, 0, "")
		if err != nil {
			logx.Errorf("检查名称重复失败: %v", err)
			return nil, err
		}

		// TODO: 更精确的重复检查逻辑
		return &types.IsRepeatVo{
			Result: false,
		}, nil
	}

	// Step 3: 按指定StdType检查是否存在
	exists, err := l.svcCtx.DataElementModel.CheckNameCnExists(l.ctx, req.NameCn, req.StdType, 0, "")
	if err != nil {
		logx.Errorf("检查名称重复失败: %v", err)
		return nil, err
	}

	return &types.IsRepeatVo{
		Result: exists,
	}, nil
}
