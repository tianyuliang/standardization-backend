// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

//go:build !mock_logic_off
// +build !mock_logic_off

package mock

import (
	"context"
	"strings"
)

// ============================================
// Token 服务 Mock
//
// 替换目标: TokenUtil.getUser(), TokenUtil.getDeptPathIds()
// ============================================

// GetUserInfo 从 Token 获取用户信息
// MOCK: 模拟从 Token 解析用户信息
// 替换目标: TokenUtil.getUser() 或从中间件获取用户信息
func GetUserInfo(ctx context.Context) (userId, userName string) {
	// MOCK: 返回默认系统用户
	// TODO: 从 Context 获取真实用户信息
	// return userInfo.UserId, userInfo.UserName
	return "system-user-id", "system"
}

// GetDeptPathIds 从 Token 获取部门完整路径
// MOCK: 模拟从 Token 解析部门路径
// 替换目标: TokenUtil.getDeptPathIds(departmentIds)
func GetDeptPathIds(ctx context.Context, departmentIds string) (pathIds, thirdDeptId string) {
	// MOCK: 返回原值
	// TODO: 从 Token/部门服务获取完整部门路径
	// Department department = TokenUtil.getDeptPathIds(departmentIds);
	// return department.getPathId(), department.getThirdDeptId()
	return departmentIds, ""
}

// Department 部门信息结构
type Department struct {
	Id       string
	Name     string
	PathId   string
	PathName string
}

// GetMapDeptInfo 批量获取部门信息
// MOCK: 模拟获取部门名称，返回Mock数据
// 替换目标: TokenUtil.getMapDeptInfo(deptIds)
func GetMapDeptInfo(ctx context.Context, deptIds []string) map[string]*Department {
	// MOCK: 返回空map
	// TODO: 调用部门服务获取部门信息
	// Map<String, Department> deptEntityMap = TokenUtil.getMapDeptInfo(deptIds);
	result := make(map[string]*Department)
	for _, deptId := range deptIds {
		if deptId != "" {
			result[deptId] = &Department{
				Id:       deptId,
				Name:     "Mock部门-" + deptId,
				PathId:   deptId,
				PathName: "Mock部门路径-" + deptId,
			}
		}
	}
	return result
}

// PathSplitAfter 提取路径中的最后一个ID
// 替换目标: StringUtil.PathSplitAfter(pathIds)
func PathSplitAfter(pathIds string) string {
	if pathIds == "" {
		return ""
	}
	parts := strings.Split(pathIds, "/")
	return parts[len(parts)-1]
}
