// Code scaffolded by goctl. Safe to edit.

package dataelement

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/dataelement/dataelement"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListDataElementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页查询数据元
func NewListDataElementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDataElementLogic {
	return &ListDataElementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDataElementLogic) ListDataElement(req *types.DataElementListReq) (resp *types.DataElementListVo, err error) {
	// Step 1: 处理目录ID（获取当前目录及所有子目录ID列表）
	catalogIds := []int64{req.CatalogId}
	if req.CatalogId > 0 {
		childIds, err := mock.GetChildCatalogIds(l.ctx, req.CatalogId)
		if err != nil {
			logx.Errorf("获取子目录失败: %v", err)
		} else {
			catalogIds = append(catalogIds, childIds...)
		}
	}

	// Step 2: 构建查询选项
	opts := &dataelement.FindOptions{
		CatalogId:     &req.CatalogId,
		StdType:       &req.StdType,
		State:         stateToIntPtr(req.State),
		DataType:      &req.DataType,
		RelationType:  &req.RelationType,
		DepartmentIds: req.DepartmentIds,
		Keyword:       req.Keyword,
		Page:          int(req.Offset),
		PageSize:      int(req.Limit),
		Sort:          req.Sort,
		Direction:     req.Direction,
	}

	// Step 3: 分页查询
	elements, totalCount, err := l.svcCtx.DataElementModel.FindByCatalogIds(l.ctx, catalogIds, opts)
	if err != nil {
		logx.Errorf("查询数据元失败: %v", err)
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
		TotalCount: totalCount,
		Entries:    entries,
	}, nil
}

// stateToIntPtr 将状态字符串转换为int32指针
func stateToIntPtr(state string) *int32 {
	if state == "" {
		return nil
	}
	s := dataelement.StateToInt(state)
	return &s
}
