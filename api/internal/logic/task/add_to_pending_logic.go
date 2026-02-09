// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	localErrorx "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	poolmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/task/pool"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddToPendingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 添加至待新建
func NewAddToPendingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddToPendingLogic {
	return &AddToPendingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddToPendingLogic) AddToPending(req *types.AddToPendingReq) (resp *types.TaskBaseResp, err error) {
	// Step 1: 参数校验
	if req.BusinessTable == "" {
		return nil, localErrorx.TaskParamEmpty("businessTable")
	}

	// Step 2: 保存到待新建表
	pool := &poolmodel.BusinessTablePool{
		TableName:        req.BusinessTable,
		TableDescription: req.BusinessTableDescription,
		TableField:       req.TableField,
		Status:           poolmodel.PoolStatusPending,
		CreateUserPhone:  req.CreateUserPhone,
		Deleted:          0,
	}

	id, err := l.svcCtx.BusinessTablePoolModel.Insert(l.ctx, pool)
	if err != nil {
		logx.Errorf("添加至待新建失败: %v", err)
		return nil, localErrorx.TaskInvalidParam("添加至待新建失败")
	}

	// Step 3: 返回成功
	logx.Infof("添加至待新建成功: id=%d, tableName=%s", id, req.BusinessTable)
	return &types.TaskBaseResp{Code: "0", Description: "添加成功"}, nil
}
