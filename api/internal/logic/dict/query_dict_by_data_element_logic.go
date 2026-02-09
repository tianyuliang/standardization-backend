// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryDictByDataElementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询引用码表的数据元
func NewQueryDictByDataElementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryDictByDataElementLogic {
	return &QueryDictByDataElementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryDictByDataElementLogic) QueryDictByDataElement(req *types.PageQuery) (resp *types.DictDataElementDataListResp, err error) {
	// MOCK: 查询引用码表的数据元
	// TODO: 需要调用 dataelement 服务获取数据元列表
	// 当前返回空列表

	logx.Infof("查询引用码表的数据元: offset=%d, limit=%d", req.Offset, req.Limit)

	return &types.DictDataElementDataListResp{
		Data:       []interface{}{},
		TotalCount: 0,
	}, nil
}
