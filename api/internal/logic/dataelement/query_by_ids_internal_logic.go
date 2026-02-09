// Code scaffolded by goctl. Safe to edit.

package dataelement

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/dataelement/dataelement"
	"github.com/zeromicro/go-zero/core/logx"
)

type QueryByIdsInternalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 内部批量查询数据元
func NewQueryByIdsInternalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryByIdsInternalLogic {
	return &QueryByIdsInternalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryByIdsInternalLogic) QueryByIdsInternal(req *types.QueryByIdsReq) (resp *types.DataElementListVo, err error) {
	// Step 1: 校验ID列表
	if len(req.Ids) == 0 {
		return nil, errorx.ParameterEmpty("ids")
	}

	// Step 2: 限制批量查询数量
	const maxBatchSize = 1000
	if len(req.Ids) > maxBatchSize {
		return nil, errorx.InvalidParameter("ids", "批量查询数量不能超过1000")
	}

	// Step 3: 批量查询数据元
	elements, err := l.svcCtx.DataElementModel.FindByIds(l.ctx, req.Ids)
	if err != nil {
		logx.Errorf("内部批量查询数据元失败: %v", err)
		return nil, err
	}

	// Step 4: 构建响应
	entries := make([]types.DataElementVo, len(elements))
	for i, de := range elements {
		// 转换指针类型
		dataLen := int32(0)
		dataPrec := int32(0)
		ruleId := int64(0)
		labelId := int64(0)
		if de.DataLength != nil {
			dataLen = int32(*de.DataLength)
		}
		if de.DataPrecision != nil {
			dataPrec = int32(*de.DataPrecision)
		}
		if de.RuleId != nil {
			ruleId = *de.RuleId
		}
		if de.LabelId != nil {
			labelId = *de.LabelId
		}

		entries[i] = types.DataElementVo{
			Id:            de.Id,
			Code:          de.Code,
			NameEn:        de.NameEn,
			NameCn:        de.NameCn,
			Synonym:       de.Synonym,
			StdType:       de.StdType,
			DataType:      de.DataType,
			DataLength:    dataLen,
			DataPrecision: dataPrec,
			DictCode:      "", // TODO: 转换
			RuleId:        ruleId,
			RelationType:  de.RelationType,
			CatalogId:     de.CatalogId,
			LabelId:       labelId,
			Description:   de.Description,
			Version:       int32(de.Version),
			State:         dataelement.IntToState(de.State),
			DepartmentIds: de.DepartmentIds,
			ThirdDeptId:   de.ThirdDeptId,
			DisableReason: de.DisableReason,
			CreateUser:    de.CreateUser,
			CreateTime:    de.CreateTime,
			UpdateTime:    de.UpdateTime,
		}
	}

	return &types.DataElementListVo{
		TotalCount: int64(len(entries)),
		Entries:    entries,
	}, nil
}
