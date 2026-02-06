// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

//go:build !mock_logic_off
// +build !mock_logic_off

package mock

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
)

// ============================================
// StdFile RPC Mock
//
// 替换目标: stdFileMgrMapper.queryStdFilesByRuleId()
// ============================================

// StdFileGetById 批量获取标准文件信息
// MOCK: 模拟获取文件信息，返回空列表
// 替换目标: stdFileMgrMapper.queryStdFilesByRuleId(page, ruleId)
func StdFileGetById(ctx context.Context, svcCtx *svc.ServiceContext, fileIds []int64) []interface{} {
	// MOCK: 返回空列表
	// TODO: 调用 StdFile RPC 批量获取文件信息
	// List<StdFileMgrEntity> files = stdFileMgrMapper.queryStdFilesByRuleId(page, ruleId);
	// return files;
	return []interface{}{}
}
