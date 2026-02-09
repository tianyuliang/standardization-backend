// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"
	"fmt"

	localErrorx "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitDataElementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 提交选定数据元
func NewSubmitDataElementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitDataElementLogic {
	return &SubmitDataElementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitDataElementLogic) SubmitDataElement(req *types.SubmitDataElementReq) (resp *types.TaskBaseResp, err error) {
	// Step 1: 参数校验（Java源码: SubmitDeDto with id and dataElementId）
	if req.Id == "" {
		return nil, localErrorx.TaskParamEmpty("id")
	}

	// Step 2: 查询池记录是否存在（Java源码: getById(submitDeDto.getId())）
	pool, err := l.svcCtx.BusinessTablePoolModel.FindOneByBusinessTableFieldId(l.ctx, req.Id)
	if err != nil || pool == nil {
		logx.Errorf("查询池记录失败: id=%s, error=%v", req.Id, err)
		return nil, localErrorx.TaskDataNotExist()
	}

	// Step 3: 如果提供了dataElementId，验证并更新（Java源码: if StringUtils.isNotBlank(submitDeDto.getDataElementId())）
	if req.DataElementId != "" {
		// TODO: 调用mock服务验证数据元是否存在（Java源码: dataElementInfoService.getById()）
		// dataElementInfo, err := mock.GetDataElementInfo(l.ctx, l.svcCtx, dataElementId)
		// if err != nil || dataElementInfo == nil {
		//     return nil, localErrorx.TaskInvalidParam("数据元不存在")
		// }

		// TODO: 获取数据元详情并更新池记录（Java源码: dataElementInfoService.getDetailVo()）
		// dataElementDetailVo, err := mock.GetDataElementDetailVo(l.ctx, l.svcCtx, dataElementId)
		// if err != nil {
		//     return nil, localErrorx.TaskInvalidParam("获取数据元详情失败")
		// }

		// 将dataElementId转换为int64（Java源码使用Long类型）
		var dataElementId int64
		if _, err := fmt.Sscanf(req.DataElementId, "%d", &dataElementId); err != nil {
			return nil, localErrorx.TaskInvalidParam("dataElementId格式错误")
		}

		// 更新池记录的数据元ID（Java源码: entity.setDataElementId(), updateById(entity)）
		err = l.svcCtx.BusinessTablePoolModel.UpdateDataElementId(l.ctx, pool.Id, dataElementId)
		if err != nil {
			logx.Errorf("更新数据元ID失败: id=%s, dataElementId=%s, error=%v", req.Id, req.DataElementId, err)
			return nil, localErrorx.TaskInvalidParam("更新数据元ID失败")
		}

		logx.Infof("提交数据元成功: id=%s, dataElementId=%s", req.Id, req.DataElementId)
	} else {
		// Step 4: 如果没有提供dataElementId，删除关联（Java源码: businessTableStdCreatePoolMapper.deleteDeId()）
		err = l.svcCtx.BusinessTablePoolModel.DeleteDataElementId(l.ctx, pool.Id)
		if err != nil {
			logx.Errorf("删除数据元ID失败: id=%s, error=%v", req.Id, err)
			return nil, localErrorx.TaskInvalidParam("删除数据元ID失败")
		}

		logx.Infof("删除数据元关联成功: id=%s", req.Id)
	}

	return &types.TaskBaseResp{Code: "0", Description: "提交成功"}, nil
}
