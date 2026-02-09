// Code scaffolded by speckit. Safe to edit.

package pool

// PoolStatusEnum 池状态枚举
const (
	PoolStatusPending    int32 = 0 // 待处理
	PoolStatusProcessing int32 = 1 // 处理中
	PoolStatusCompleted  int32 = 2 // 已完成
	PoolStatusAdopted    int32 = 3 // 已采纳
	PoolStatusCancelled  int32 = 4 // 已撤销
)

// PoolStatusText 池状态文本
var PoolStatusText = map[int32]string{
	PoolStatusPending:    "pending",
	PoolStatusProcessing: "processing",
	PoolStatusCompleted:  "completed",
	PoolStatusAdopted:    "adopted",
	PoolStatusCancelled:  "cancelled",
}

// StatusToInt 状态字符串转整数
func PoolStatusToInt(state string) int32 {
	switch state {
	case "pending":
		return PoolStatusPending
	case "processing":
		return PoolStatusProcessing
	case "completed":
		return PoolStatusCompleted
	case "adopted":
		return PoolStatusAdopted
	case "cancelled":
		return PoolStatusCancelled
	default:
		return PoolStatusPending
	}
}

// IntToStatus 整数转状态字符串
func PoolIntToStatus(status int32) string {
	switch status {
	case PoolStatusPending:
		return "pending"
	case PoolStatusProcessing:
		return "processing"
	case PoolStatusCompleted:
		return "completed"
	case PoolStatusAdopted:
		return "adopted"
	case PoolStatusCancelled:
		return "cancelled"
	default:
		return "pending"
	}
}
