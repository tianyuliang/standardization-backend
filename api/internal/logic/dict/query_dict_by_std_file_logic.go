// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	localErrorx "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryDictByStdFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按文件查询码表
func NewQueryDictByStdFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryDictByStdFileLogic {
	return &QueryDictByStdFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryDictByStdFileLogic) QueryDictByStdFile(req *types.QueryDictByStdFileReq) (resp *types.DictDataListResp, err error) {
	// Step 1: 参数校验 - FileId 必填
	if req.FileId <= 0 {
		return nil, localErrorx.DictParamEmpty("文件ID不能为空")
	}

	// Step 2: 按文件ID查询码表列表
	dicts, err := l.svcCtx.DictModel.FindByFileId(l.ctx, req.FileId)
	if err != nil {
		logx.Errorf("按文件ID查询码表失败: fileId=%d, error=%v", req.FileId, err)
		return &types.DictDataListResp{Data: []types.DictVo{}, TotalCount: 0}, nil
	}

	// Step 3: 转换为响应格式
	result := make([]types.DictVo, 0)
	for _, dict := range dicts {
		enums, _ := l.svcCtx.DictEnumModel.FindByDictId(l.ctx, dict.Id)
		catalogName := getCatalogName(dict.CatalogId)
		deptName, deptPathNames, _ := getDeptInfo(dict.DepartmentIds)

		result = append(result, buildDictVo(l.ctx, dict, enums, catalogName, deptName, deptPathNames, false))
	}

	logx.Infof("按文件ID查询码表成功: fileId=%d, count=%d", req.FileId, len(result))
	return &types.DictDataListResp{
		Data:       result,
		TotalCount: int64(len(result)),
	}, nil
}
