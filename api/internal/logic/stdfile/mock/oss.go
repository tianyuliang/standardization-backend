// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

//go:build !mock_logic_off
// +build !mock_logic_off

package mock

import (
	"context"
	"io"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
)

// ============================================
// OSS 文件服务 Mock
//
// 替换目标: OssFileUploadDownloadUtil.*()
// ============================================

// OssUploadFile 上传文件到OSS
// MOCK: 模拟文件上传
// 替换目标: ossFileUploadDownloadUtil.uploadFile(file, ossId)
func OssUploadFile(ctx context.Context, svcCtx *svc.ServiceContext, file interface{}, ossId string) error {
	// MOCK: 直接返回成功
	// TODO: 调用 OSS 服务上传文件
	// log.Infof("上传文件到OSS: ossId=%s", ossId)
	return nil
}

// OssDownloadFile 从OSS下载文件
// MOCK: 模拟文件下载，返回空Reader
// 替换目标: ossFileUploadDownloadUtil.download(ossId, outputStream)
func OssDownloadFile(ctx context.Context, svcCtx *svc.ServiceContext, ossId string) (io.ReadCloser, string, error) {
	// MOCK: 返回空内容，实际应返回文件流和文件名
	// TODO: 调用 OSS 服务下载文件
	// return fileReader, fileName, nil
	return nil, "", nil
}

// OssDeleteFile 从OSS删除文件
// MOCK: 模拟文件删除
// 替换目标: ossFileUploadDownloadUtil.delete(ossId)
func OssDeleteFile(ctx context.Context, svcCtx *svc.ServiceContext, ossId string) error {
	// MOCK: 直接返回成功
	// TODO: 调用 OSS 服务删除文件
	// log.Infof("从OSS删除文件: ossId=%s", ossId)
	return nil
}

// OssGetFileUrl 获取文件访问URL
// MOCK: 模拟获取文件URL
// 替换目标: ossFileUploadDownloadUtil.getFileUrl(ossId)
func OssGetFileUrl(ctx context.Context, svcCtx *svc.ServiceContext, ossId string) string {
	// MOCK: 返回Mock URL
	// TODO: 调用 OSS 服务获取文件URL
	return "https://mock.oss.com/file/" + ossId
}
