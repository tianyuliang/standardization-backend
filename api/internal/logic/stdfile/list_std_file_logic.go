// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package stdfile

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	stdfilemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/stdfile/stdfile"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListStdFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标准文件-列表查询
func NewListStdFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListStdFileLogic {
	return &ListStdFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListStdFileLogic) ListStdFile(req *types.StdFileListQuery) (resp *types.StdFileDataListResp, err error) {
	// Step 1: 参数校验
	if err := ValidatePagination(req.Offset, req.Limit); err != nil {
		return nil, err
	}

	// Step 2: 验证排序字段
	if err := ValidateSortField(req.Sort); err != nil {
		return nil, err
	}

	// Step 3: 构建查询选项
	opts := &stdfilemodel.FindOptions{
		Page:      req.Offset,
		PageSize:  req.Limit,
		Sort:      req.Sort,
		Direction: req.Direction,
	}

	// 设置目录过滤条件
	if req.CatalogId > 0 {
		opts.CatalogId = &req.CatalogId
	}

	// 设置关键字搜索
	if req.Keyword != "" {
		opts.Keyword = req.Keyword
	}

	// 设置组织类型过滤
	if req.OrgType >= 0 {
		orgType := int(req.OrgType)
		opts.OrgType = &orgType
	}

	// 设置状态过滤
	if req.State != "" {
		state, err := ParseState(req.State)
		if err != nil {
			return nil, err
		}
		opts.State = &state
	}

	// 设置部门ID过滤
	if req.DepartmentId != "" {
		opts.DepartmentId = req.DepartmentId
	}

	// Step 5: 查询列表
	// 如果没有指定目录ID，使用特殊查询（不限制目录）
	var models []*stdfilemodel.StdFile
	var totalCount int64

	if req.CatalogId == 0 {
		// 不限制目录，直接查询所有
		models, totalCount, err = l.queryAllFiles(opts)
	} else {
		// 限制目录，使用 FindByCatalogIds
		models, totalCount, err = l.svcCtx.StdFileModel.FindByCatalogIds(l.ctx, opts)
	}

	if err != nil {
		return nil, HandleError(err)
	}

	// Step 6: 转换为响应对象
	data := ModelsToResp(l.ctx, l.svcCtx, models)

	resp = &types.StdFileDataListResp{
		TotalCount: totalCount,
		Data:       data,
	}

	logx.Infof("查询标准文件列表成功: offset=%d, limit=%d, total=%d",
		req.Offset, req.Limit, totalCount)
	return resp, nil
}

// queryAllFiles 查询所有文件（不限制目录）
func (l *ListStdFileLogic) queryAllFiles(opts *stdfilemodel.FindOptions) ([]*stdfilemodel.StdFile, int64, error) {
	// TODO: 实现不限制目录的查询逻辑
	// 当前暂时使用 FindByCatalogIds，传入 nil CatalogId
	return l.svcCtx.StdFileModel.FindByCatalogIds(l.ctx, opts)
}
