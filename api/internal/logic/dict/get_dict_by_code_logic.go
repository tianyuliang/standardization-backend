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

type GetDictByCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 码表详情（按Code）
func NewGetDictByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictByCodeLogic {
	return &GetDictByCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDictByCodeLogic) GetDictByCode(req *types.IdReq) (resp *types.DictVo, err error) {
	// Step 1: 按Code查询码表（req.Id 实际上是 code）
	dict, err := l.svcCtx.DictModel.FindByCode(l.ctx, req.Id)
	if err != nil {
		return nil, localErrorx.DictDataNotExist()
	}

	// Step 2: 查询码值明细列表
	enums, err := l.svcCtx.DictEnumModel.FindByDictId(l.ctx, dict.Id)
	if err != nil {
		logx.Errorf("查询码值失败: %v", err)
		enums = []*dictmodel.DictEnum{} // 返回空列表而不是错误
	}

	// Step 3: 查询目录信息
	catalogName := getCatalogName(dict.CatalogId)

	// Step 4: 查询部门信息
	deptName, deptPathNames, _ := getDeptInfo(dict.DepartmentIds)

	// TODO: 查询是否被引用（需要调用 dataelement 服务）
	usedFlag := false

	result := buildDictVo(l.ctx, dict, enums, catalogName, deptName, deptPathNames, usedFlag)

	logx.Infof("按编码查询码表成功: code=%d, id=%d", req.Id, dict.Id)
	return &result, nil
}
