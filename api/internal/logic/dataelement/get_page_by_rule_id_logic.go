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

type GetPageByRuleIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按规则ID分页查询数据元
func NewGetPageByRuleIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPageByRuleIdLogic {
	return &GetPageByRuleIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPageByRuleIdLogic) GetPageByRuleId(req *types.GetPageByRuleIdReq) (resp *types.DataElementListVo, err error) {
	// Step 1: 校验规则ID
	if req.RuleId <= 0 {
		return nil, errorx.InvalidParameter("ruleId", "规则ID必须大于0")
	}

	// Step 2: 构建查询选项
	opts := &dataelement.FindOptions{
		Page:      int(req.Offset),
		PageSize:  int(req.Limit),
		Sort:      "f_id",
		Direction: "desc",
	}

	// Step 3: 按规则ID分页查询
	elements, totalCount, err := l.svcCtx.DataElementModel.FindByRuleId(l.ctx, req.RuleId, opts)
	if err != nil {
		logx.Errorf("按规则ID分页查询数据元失败: %v", err)
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
