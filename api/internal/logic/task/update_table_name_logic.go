// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	localErrorx "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTableNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改表名称
func NewUpdateTableNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTableNameLogic {
	return &UpdateTableNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTableNameLogic) UpdateTableName(req *types.UpdateTableNameReq) (resp *types.TaskBaseResp, err error) {
	// Step 1: 参数校验
	if req.TableId <= 0 {
		return nil, localErrorx.TaskParamEmpty("tableId")
	}
	if req.TableName == "" {
		return nil, localErrorx.TaskParamEmpty("tableName")
	}

	// Step 2: 查询业务表是否存在
	pool, err := l.svcCtx.BusinessTablePoolModel.FindOne(l.ctx, req.TableId)
	if err != nil {
		logx.Errorf("查询业务表失败: tableId=%d, error=%v", req.TableId, err)
		return nil, localErrorx.TaskNotExist()
	}
	if pool == nil {
		return nil, localErrorx.TaskDataNotExist()
	}

	// Step 3: 更新表名称
	pool.TableName = req.TableName
	err = l.svcCtx.BusinessTablePoolModel.Update(l.ctx, pool)
	if err != nil {
		logx.Errorf("更新表名称失败: tableId=%d, error=%v", req.TableId, err)
		return nil, localErrorx.TaskInvalidParam("更新表名称失败")
	}

	// Step 4: 返回成功
	logx.Infof("更新表名称成功: tableId=%d, tableName=%s", req.TableId, req.TableName)
	return &types.TaskBaseResp{Code: "0", Description: "更新成功"}, nil
}
