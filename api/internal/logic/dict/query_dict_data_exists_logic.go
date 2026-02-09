// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	localErrorx "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	dictmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/dict/dict"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryDictDataExistsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询数据是否存在
func NewQueryDictDataExistsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryDictDataExistsLogic {
	return &QueryDictDataExistsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryDictDataExistsLogic) QueryDictDataExists(req *types.QueryDictDataExistsReq) (resp *types.DataExistsResp, err error) {
	// Step 1: 参数校验 - OrgType 必填
	if req.OrgType <= 0 {
		return nil, localErrorx.DictParamEmpty("组织类型不能为空")
	}

	// Step 2: 构建查询条件
	// 使用 Keyword 进行名称搜索（支持中文名称或英文名称）
	keyword := ""
	if req.ChName != "" {
		keyword = req.ChName
	} else if req.EnName != "" {
		keyword = req.EnName
	}

	orgType := req.OrgType
	opts := &dictmodel.FindOptions{
		OrgType:  &orgType,
		Keyword:  keyword,
		Page:     1,
		PageSize: 1,
	}

	// Step 3: 查询是否存在
	_, totalCount, err := l.svcCtx.DictModel.FindByCatalogIds(l.ctx, opts)
	if err != nil {
		logx.Errorf("查询码表存在性失败: %v", err)
		return &types.DataExistsResp{Exists: false}, nil
	}

	exists := totalCount > 0
	logx.Infof("查询码表存在性: orgType=%d, keyword=%s, exists=%t", req.OrgType, keyword, exists)
	return &types.DataExistsResp{Exists: exists}, nil
}
