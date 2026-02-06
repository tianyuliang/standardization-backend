// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

//go:build !mock_logic_off
// +build !mock_logic_off

package mock

import (
	"context"
)

// ============================================
// Dict RPC Mock
//
// 替换目标: dictService.queryById()
// ============================================

// DictCheckExist 校验码表是否存在
// MOCK: 模拟码表校验，默认返回存在
// 替换目标: dictService.queryById(dictId)
func DictCheckExist(ctx context.Context, dictId int64) bool {
	// MOCK: 默认码表存在
	// TODO: 调用 Dict RPC 校验码表是否存在
	// DictVo data = dictService.queryById(dictId);
	// return data != nil
	return true
}
