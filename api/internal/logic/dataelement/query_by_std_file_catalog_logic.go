// Code scaffolded by goctl. Safe to edit.

package dataelement

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/dataelement/dataelement"
	"github.com/zeromicro/go-zero/core/logx"
)

type QueryByStdFileCatalogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按标准文件目录查询数据元
func NewQueryByStdFileCatalogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryByStdFileCatalogLogic {
	return &QueryByStdFileCatalogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryByStdFileCatalogLogic) QueryByStdFileCatalog(req *types.DataElementQueryByFileCatalogReq) (resp *types.DataElementListVo, err error) {
	// Step 1: 构建查询选项
	opts := &dataelement.FindOptions{
		FileCatalogId: &req.CatalogId,
		Page:          int(req.Offset),
		PageSize:      int(req.Limit),
		Sort:          req.Sort,
		Direction:     req.Direction,
	}

	// Step 2: 按文件目录查询数据元
	elements, totalCount, err := l.svcCtx.DataElementModel.FindByFileCatalog(l.ctx, opts)
	if err != nil {
		logx.Errorf("按文件目录查询数据元失败: %v", err)
		return nil, err
	}

	// Step 3: 构建响应
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
		TotalCount: totalCount,
		Entries:    entries,
	}, nil
}
