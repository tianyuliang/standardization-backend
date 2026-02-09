// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryDictByStdFileCatalogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按文件目录查询码表
func NewQueryDictByStdFileCatalogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryDictByStdFileCatalogLogic {
	return &QueryDictByStdFileCatalogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryDictByStdFileCatalogLogic) QueryDictByStdFileCatalog(req *types.QueryDictByStdFileCatalogReq) (resp *types.DictDataListResp, err error) {
	// MOCK: 按标准文件目录查询码表
	// TODO: 需要调用 stdfile 服务获取目录下的文件列表，然后查询关联的码表
	// 当前返回空列表

	logx.Infof("按标准文件目录查询码表: catalogId=%d, keyword=%s, orgType=%d", req.CatalogId, req.Keyword, req.OrgType)

	return &types.DictDataListResp{
		Data:       []types.DictVo{},
		TotalCount: 0,
	}, nil
}
