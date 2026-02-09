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

type ListDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 码表列表查询
func NewListDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDictLogic {
	return &ListDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDictLogic) ListDict(req *types.DictListQuery) (resp *types.DictDataListResp, err error) {
	// Step 1: 处理目录ID（如果提供了目录ID，需要包含子目录）
	// Step 2: 构建查询条件
	opts := &dictmodel.FindOptions{
		Page:      req.Offset,
		PageSize:  req.Limit,
		Sort:      req.Sort,
		Direction: req.Direction,
	}

	if req.CatalogId > 0 {
		opts.CatalogId = &req.CatalogId
	}

	if req.OrgType >= 0 {
		opts.OrgType = &req.OrgType
	}

	if req.State != "" {
		state := stateToInt(req.State)
		opts.State = &state
	}

	if req.Keyword != "" {
		opts.Keyword = req.Keyword
	}

	if req.DepartmentId != "" {
		opts.DepartmentId = req.DepartmentId
	}

	// Step 3: 分页查询
	dicts, totalCount, err := l.svcCtx.DictModel.FindByCatalogIds(l.ctx, opts)
	if err != nil {
		logx.Errorf("查询码表列表失败: %v", err)
		return nil, localErrorx.DictInvalidParam("查询码表列表失败")
	}

	// Step 4: 数据处理（为每个码表查询码值、部门信息）
	result := make([]types.DictVo, len(dicts))
	for i, dict := range dicts {
		enums, _ := l.svcCtx.DictEnumModel.FindByDictId(l.ctx, dict.Id)
		catalogName := getCatalogName(dict.CatalogId)
		deptName, deptPathNames, _ := getDeptInfo(dict.DepartmentIds)

		// TODO: 查询是否被引用（需要调用 dataelement 服务）
		usedFlag := false

		result[i] = buildDictVo(l.ctx, dict, enums, catalogName, deptName, deptPathNames, usedFlag)
	}

	return &types.DictDataListResp{
		Data:       result,
		TotalCount: totalCount,
	}, nil
}
