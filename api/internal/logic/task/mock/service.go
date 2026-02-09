// Code scaffolded by speckit. Safe to edit.

//go:build !mock_logic_off
// +build !mock_logic_off

package mock

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
)

// ============================================
// AfService Mock
//
// 替换目标: afService.getTaskDetailDto(taskId)
// ============================================

// GetTaskDetailDto 获取任务详情
// MOCK: 模拟获取任务详情
// 替换目标: afService.getTaskDetailDto(taskId)
func GetTaskDetailDto(ctx context.Context, svcCtx *svc.ServiceContext, taskId string) (*types.TaskDetailResp, error) {
	// MOCK: 默认返回nil
	// TODO: 调用 AfService RPC 获取任务详情
	// TaskDetailDto taskDetailDto = afService.getTaskDetailDto(taskId);
	// return taskDetailDto;
	return nil, nil
}

// ============================================
// IDataElementInfoService Mock
//
// 替换目标: dataElementInfoService.getById(dataElementId)
// ============================================

// GetDataElementInfo 获取数据元信息
// MOCK: 模拟获取数据元信息
// 替换目标: dataElementInfoService.getById(dataElementId)
func GetDataElementInfo(ctx context.Context, svcCtx *svc.ServiceContext, dataElementId int64) (*types.DataElementInfo, error) {
	// MOCK: 默认返回nil
	// TODO: 调用 DataElement RPC 获取数据元信息
	// DataElementInfo dataElementInfo = dataElementInfoService.getById(dataElementId);
	// return dataElementInfo;
	return nil, nil
}

// GetDataElementDetailVo 获取数据元详情
// MOCK: 模拟获取数据元详情
// 替换目标: dataElementInfoService.getDetailVo(dataElementId)
func GetDataElementDetailVo(ctx context.Context, svcCtx *svc.ServiceContext, dataElementId int64) (*types.DataElementDetailVo, error) {
	// MOCK: 默认返回nil
	// TODO: 调用 DataElement RPC 获取数据元详情
	// DataElementDetailVo detailVo = dataElementInfoService.getDetailVo(dataElementId);
	// return detailVo;
	return nil, nil
}

// ============================================
// 推荐服务 Mock
//
// 替换目标: 调用外部推荐服务API
// ============================================

// CallStdRecService 调用标准推荐服务
// MOCK: 模拟标准推荐服务
// 替换目标: HTTP POST to 推荐服务
func CallStdRecService(ctx context.Context, svcCtx *svc.ServiceContext, req *types.StdRecReq) (*types.StdRecResp, error) {
	// MOCK: 默认返回空结果
	// TODO: 实现HTTP调用推荐服务
	// HttpClient.post(recommendationServiceUrl + "/std-rec/rec", req)
	return &types.StdRecResp{
		Code:        "0",
		Description: "成功",
		Data:        []types.StdRecItem{},
	}, nil
}

// CallRuleRecService 调用规则推荐服务
// MOCK: 模拟规则推荐服务
// 替换目标: HTTP POST to 推荐服务
func CallRuleRecService(ctx context.Context, svcCtx *svc.ServiceContext, req *types.RuleRecReq) (*types.StdRecResp, error) {
	// MOCK: 默认返回空结果
	// TODO: 实现HTTP调用规则推荐服务
	// HttpClient.post(recommendationServiceUrl + "/rule-rec/rec", req)
	return &types.StdRecResp{
		Code:        "0",
		Description: "成功",
		Data:        []types.StdRecItem{},
	}, nil
}

// ============================================
// Webhook Mock
//
// 替换目标: 发送任务完成回调
// ============================================

// SendTaskCallback 发送任务完成回调
// MOCK: 模拟发送webhook回调
// 替换目标: HTTP POST to webhook URL
func SendTaskCallback(ctx context.Context, svcCtx *svc.ServiceContext, webhook string, taskId string) error {
	// MOCK: 默认成功
	// TODO: 实现HTTP POST回调
	// HttpClient.post(webhook, callbackData)
	return nil
}
