// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

//go:build !mock_logic_off
// +build !mock_logic_off

package mock

import (
	"context"
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
