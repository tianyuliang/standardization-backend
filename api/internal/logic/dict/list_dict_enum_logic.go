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

type ListDictEnumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 码值分页查询
func NewListDictEnumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDictEnumLogic {
	return &ListDictEnumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDictEnumLogic) ListDictEnum(req *types.DictEnumListQuery) (resp *types.DictEnumDataListResp, err error) {
	// Step 1: 参数校验 - 码表ID必须大于0
	if req.DictId <= 0 {
		return nil, localErrorx.DictParamEmpty("码表ID不能为空")
	}

	// Step 2: 校验码表是否存在
	dict, err := l.svcCtx.DictModel.FindOne(l.ctx, req.DictId)
	if err != nil {
		return nil, localErrorx.DictDataNotExist()
	}
	logx.Infof("查询码值列表: dictId=%d, dictCode=%d", req.DictId, dict.Code)

	// Step 3: 分页查询码值
	enums, totalCount, err := l.svcCtx.DictEnumModel.FindPageByDictId(
		l.ctx,
		req.DictId,
		req.Keyword,
		req.Offset,
		req.Limit,
	)
	if err != nil {
		logx.Errorf("查询码值分页列表失败: dictId=%d, error=%v", req.DictId, err)
		return nil, localErrorx.DictInvalidParam("查询码值失败")
	}

	// Step 4: 数据处理（转换为响应格式）
	result := make([]types.DictEnumVo, len(enums))
	for i, enum := range enums {
		result[i] = types.DictEnumVo{
			Id:    enum.Id,
			Code:  enum.Code,
			Value: enum.Value,
		}
	}

	logx.Infof("查询码值列表成功: dictId=%d, count=%d, total=%d", req.DictId, len(enums), totalCount)
	return &types.DictEnumDataListResp{
		Data:       result,
		TotalCount: totalCount,
	}, nil
}
