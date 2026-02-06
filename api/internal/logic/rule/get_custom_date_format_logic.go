// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	rulemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/rule/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCustomDateFormatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取自定义日期格式列表
//
// 业务流程:
//
//	返回预定义的日期格式列表
func NewGetCustomDateFormatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCustomDateFormatLogic {
	return &GetCustomDateFormatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCustomDateFormatLogic) GetCustomDateFormat() (resp *types.CustomDateFormatResp, err error) {
	return &types.CustomDateFormatResp{
		Data: rulemodel.CustomDateFormat,
	}, nil
}
