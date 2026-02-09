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

type StdCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标准创建（内部）
func NewStdCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StdCreateLogic {
	return &StdCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StdCreateLogic) StdCreate(req *types.StdRecReq) (resp *types.TaskBaseResp, err error) {
	// Step 1: 参数校验
	if req.TaskId <= 0 {
		return nil, localErrorx.TaskParamEmpty("taskId")
	}
	if req.BusinessTable == "" {
		return nil, localErrorx.TaskParamEmpty("businessTable")
	}
	if req.TableField == "" {
		return nil, localErrorx.TaskParamEmpty("tableField")
	}

	// Step 2: 执行标准创建流程（TODO: 调用推荐服务并保存结果）
	logx.Infof("执行标准创建: taskId=%d, table=%s, field=%s", req.TaskId, req.BusinessTable, req.TableField)
	// TODO: 1. 调用推荐服务获取标准推荐
	// TODO: 2. 保存推荐结果到 TaskStdCreateResultModel

	// Step 3: 返回成功
	return &types.TaskBaseResp{Code: "0", Description: "标准创建成功"}, nil
}
