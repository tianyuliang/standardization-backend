// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"
	"strconv"

	localErrorx "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	poolmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/task/pool"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryTaskStateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 任务状态查询
func NewQueryTaskStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryTaskStateLogic {
	return &QueryTaskStateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryTaskStateLogic) QueryTaskState(req *types.QueryTaskStateReq) (resp *types.BusinessTableStateListResp, err error) {
	// Step 1: 参数校验
	if len(req.BusinessTableId) == 0 {
		return nil, localErrorx.TaskParamEmpty("businessTableId")
	}

	// Step 2: 转换状态字符串为整数（Java源码: BusinessTableStdCreatePoolStateEnum.valueOf(StringUtils.upperCase(i)).getValue()）
	var states []int32
	if len(req.State) > 0 {
		for _, stateStr := range req.State {
			// Java源码使用枚举转换，Go使用PoolStatusToInt
			state := poolmodel.PoolStatusToInt(stateStr)
			states = append(states, state)
		}
	}

	// Step 3: 查询池记录（Java源码: list(queryWrapper) where businessTableId IN (?) AND state IN (?)）
	pools, err := l.svcCtx.BusinessTablePoolModel.FindByTableNamesAndStates(l.ctx, req.BusinessTableId, states)
	if err != nil {
		logx.Errorf("查询任务状态失败: businessTableId=%v, states=%v, error=%v", req.BusinessTableId, req.State, err)
		return &types.BusinessTableStateListResp{
			TotalCount: 0,
			Data:       []types.BusinessTableStateVo{},
		}, nil
	}

	// Step 4: 构建响应（Java源码: BusinessTableStateVo with businessTableFieldId and state）
	data := make([]types.BusinessTableStateVo, 0)
	for _, pool := range pools {
		// Java源码: vo.setBusinessTableFieldId(i.getBusinessTableFieldId())
		// 我们的实现使用f_id作为业务表字段ID
		businessTableFieldId := strconv.FormatInt(pool.Id, 10)

		// Java源码: vo.setState(StringUtils.lowerCase(BusinessTableStdCreatePoolStateEnum.of(i.getState()).name()))
		stateText := poolmodel.PoolIntToStatus(pool.Status)

		data = append(data, types.BusinessTableStateVo{
			BusinessTableFieldId: businessTableFieldId,
			State:                stateText,
		})
	}

	logx.Infof("查询任务状态成功: count=%d", len(data))
	return &types.BusinessTableStateListResp{
		TotalCount: int64(len(data)),
		Data:       data,
	}, nil
}
