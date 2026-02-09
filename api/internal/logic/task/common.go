// Code scaffolded by speckit. Safe to edit.

package task

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	taskmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/task/task"

	"github.com/zeromicro/go-zero/core/logx"
)

type taskCommonLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// GenerateTaskNo 生成任务编号
func GenerateTaskNo() string {
	return fmt.Sprintf("TASK%s%04d", time.Now().Format("20060102"), rand.Intn(10000))
}

// StatusToInt 状态字符串转整数
func StatusToInt(state string) int32 {
	switch state {
	case "unhandled":
		return taskmodel.TaskStatusUnhandled
	case "processing":
		return taskmodel.TaskStatusProcessing
	case "completed":
		return taskmodel.TaskStatusCompleted
	default:
		return taskmodel.TaskStatusUnhandled
	}
}

// IntToStatus 整数转状态字符串
func IntToStatus(status int32) string {
	switch status {
	case taskmodel.TaskStatusUnhandled:
		return "unhandled"
	case taskmodel.TaskStatusProcessing:
		return "processing"
	case taskmodel.TaskStatusCompleted:
		return "completed"
	default:
		return "unhandled"
	}
}

// buildTaskResp 构建任务响应
func buildTaskResp(task *taskmodel.TaskStdCreate) types.TaskResp {
	return types.TaskResp{
		Id:               task.Id,
		TaskNo:           task.TaskNo,
		Table:            task.Table,
		TableDescription: task.TableDescription,
		TableField:       task.TableField,
		Status:           task.Status,
		CreateTime:       task.CreateTime,
		CreateUser:       task.CreateUser,
		CreateUserPhone:  task.CreateUserPhone,
		Webhook:          task.Webhook,
	}
}

// buildTaskResultResp 构建任务结果响应
func buildTaskResultResp(result *taskmodel.TaskStdCreateResult) types.TaskStdCreateResult {
	return types.TaskStdCreateResult{
		Id:                    result.Id,
		TaskId:                result.TaskId,
		TableField:            result.TableField,
		TableFieldDescription: result.TableFieldDescription,
		StdRefFile:            result.StdRefFile,
		StdCode:               result.StdCode,
		RecStdCodes:           result.RecStdCodes,
		StdChName:             result.StdChName,
		StdEnName:             result.StdEnName,
	}
}
