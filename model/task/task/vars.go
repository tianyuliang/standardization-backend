// Code scaffolded by speckit. Safe to edit.

package task

// TaskStatusEnum 任务状态枚举
const (
	TaskStatusUnhandled  int32 = 0 // 未处理
	TaskStatusProcessing int32 = 1 // 处理中
	TaskStatusCompleted  int32 = 2 // 处理完成
)

// TaskStatusText 任务状态文本
var TaskStatusText = map[int32]string{
	TaskStatusUnhandled:  "unhandled",
	TaskStatusProcessing: "processing",
	TaskStatusCompleted:  "completed",
}

// StatusToInt 状态字符串转整数
func StatusToInt(state string) int32 {
	switch state {
	case "unhandled":
		return TaskStatusUnhandled
	case "processing":
		return TaskStatusProcessing
	case "completed":
		return TaskStatusCompleted
	default:
		return TaskStatusUnhandled
	}
}

// IntToStatus 整数转状态字符串
func IntToStatus(status int32) string {
	switch status {
	case TaskStatusUnhandled:
		return "unhandled"
	case TaskStatusProcessing:
		return "processing"
	case TaskStatusCompleted:
		return "completed"
	default:
		return "unhandled"
	}
}
