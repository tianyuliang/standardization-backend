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

type UpdateDictStateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 启用/停用码表
func NewUpdateDictStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictStateLogic {
	return &UpdateDictStateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDictStateLogic) UpdateDictState(req *types.UpdateDictStateReq, id int64) (resp *types.DictBaseResp, err error) {
	// Step 1: 校验码表存在性
	oldDict, err := l.svcCtx.DictModel.FindOne(l.ctx, id)
	if err != nil {
		return nil, localErrorx.DictDataNotExist()
	}

	// Step 2: 业务校验（停用原因必填）
	if err := checkDisableReason(req.State, req.Reason); err != nil {
		return nil, err
	}

	// Step 3: 检查状态是否变更
	newState := stateToInt(req.State)
	if oldDict.State == newState {
		return &types.DictBaseResp{Code: "0", Description: "状态未变更"}, nil
	}

	// Step 4: 更新状态
	username, _ := getUserInfo(l.ctx)
	dictData := &dictmodel.Dict{
		Id:            id,
		Code:          oldDict.Code,
		ChName:        oldDict.ChName,
		EnName:        oldDict.EnName,
		Description:   oldDict.Description,
		CatalogId:     oldDict.CatalogId,
		OrgType:       oldDict.OrgType,
		Version:       oldDict.Version,
		State:         newState,
		DisableReason: req.Reason,
		AuthorityId:   oldDict.AuthorityId,
		DepartmentIds: oldDict.DepartmentIds,
		ThirdDeptId:   oldDict.ThirdDeptId,
		CreateUser:    oldDict.CreateUser,
		UpdateUser:    username,
	}

	if err := l.svcCtx.DictModel.Update(l.ctx, dictData); err != nil {
		logx.Errorf("更新码表状态失败: id=%d, state=%s, error=%v", id, req.State, err)
		return nil, localErrorx.DictInvalidParam("更新状态失败")
	}

	logx.Infof("码表状态修改成功: id=%d, state=%s, reason=%s", id, req.State, req.Reason)
	return &types.DictBaseResp{Code: "0", Description: "状态修改成功"}, nil
}
