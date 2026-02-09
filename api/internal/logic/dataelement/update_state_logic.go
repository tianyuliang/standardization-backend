// Code scaffolded by goctl. Safe to edit.

package dataelement

import (
	"context"
	"strconv"
	"strings"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/dataelement/dataelement"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 启用/停用数据元
func NewUpdateStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStateLogic {
	return &UpdateStateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateStateLogic) UpdateState(idsStr string, req *types.UpdateStateReq) (resp *types.EmptyResp, err error) {
	// Step 1: 校验状态值
	if req.State != "enable" && req.State != "disable" {
		return nil, errorx.InvalidParameter("state", "状态值必须为enable或disable")
	}

	// Step 2: 解析ID列表
	if idsStr == "" {
		return nil, errorx.ParameterEmpty("ids")
	}

	idStrs := strings.Split(idsStr, ",")
	ids := make([]int64, 0, len(idStrs))
	for _, idStr := range idStrs {
		id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if err != nil || id <= 0 {
			return nil, errorx.InvalidParameter("ids", "ID格式错误")
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, errorx.ParameterEmpty("ids")
	}

	// Step 3: 校验停用原因
	if req.State == "disable" && req.Reason == "" {
		return nil, errorx.InvalidParameter("reason", "停用时必须填写原因")
	}

	// Step 4: 转换状态值
	state := dataelement.StateToInt(req.State)

	// Step 5: 更新状态
	updateUser := "system" // TODO: 从token获取
	err = l.svcCtx.DataElementModel.UpdateState(l.ctx, ids, state, req.Reason, updateUser)
	if err != nil {
		logx.Errorf("更新数据元状态失败: %v", err)
		return nil, err
	}

	// Step 6: 发送MQ消息
	for _, id := range ids {
		_ = mock.SendMqMessage(l.ctx, "updateState", map[string]interface{}{
			"id":     id,
			"state":  req.State,
			"reason": req.Reason,
		}, updateUser)
	}

	return &types.EmptyResp{}, nil
}
