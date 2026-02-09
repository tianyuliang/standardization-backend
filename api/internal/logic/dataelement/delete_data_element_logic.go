// Code scaffolded by goctl. Safe to edit.

package dataelement

import (
	"context"
	"fmt"
	"strings"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/dataelement/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDataElementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除数据元
func NewDeleteDataElementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDataElementLogic {
	return &DeleteDataElementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDataElementLogic) DeleteDataElement(idsStr string) (resp *types.EmptyResp, err error) {
	// Step 1: 校验ID列表
	if idsStr == "" {
		return nil, errorx.ParameterEmpty("ids")
	}

	// Step 2: 解析ID列表
	idStrs := strings.Split(idsStr, ",")
	ids := make([]int64, 0, len(idStrs))
	for _, idStr := range idStrs {
		id := int64(0)
		_, err := fmt.Sscanf(idStr, "%d", &id)
		if err != nil || id <= 0 {
			return nil, errorx.InvalidParameter("ids", "ID格式错误")
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, errorx.ParameterEmpty("ids")
	}

	// Step 3: 删除关联文件关系
	err = l.svcCtx.RelationDeFileModel.DeleteByDeIds(l.ctx, ids)
	if err != nil {
		logx.Errorf("删除关联文件关系失败: %v", err)
	}

	// Step 4: 物理删除数据元
	err = l.svcCtx.DataElementModel.DeleteByIds(l.ctx, ids)
	if err != nil {
		logx.Errorf("删除数据元失败: %v", err)
		return nil, err
	}

	// Step 5: 发送MQ消息
	for _, id := range ids {
		_ = mock.SendMqMessage(l.ctx, "delete", map[string]interface{}{"id": id}, "system")
	}

	return &types.EmptyResp{}, nil
}
