// Code scaffolded by goctl. Safe to edit.

package dataelement

import (
	"context"
	"fmt"
	"time"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/dataelement/dataelement"
	"github.com/zeromicro/go-zero/core/logx"
)

type ExportDataElementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 导出数据元
func NewExportDataElementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportDataElementLogic {
	return &ExportDataElementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportDataElementLogic) ExportDataElement(req *types.ExportDataElementReq) (resp *types.ExportVo, err error) {
	// Step 1: 构建查询选项
	opts := &dataelement.FindOptions{
		CatalogId: &req.CatalogId,
		Page:      1,
		PageSize:  10000, // 导出时使用较大的分页
		Sort:      "f_id",
		Direction: "asc",
	}

	// Step 2: 查询数据元列表
	_, totalCount, err := l.svcCtx.DataElementModel.FindByCatalogIds(l.ctx, []int64{req.CatalogId}, opts)
	if err != nil {
		logx.Errorf("导出数据元查询失败: %v", err)
		return nil, err
	}

	// Step 3: 生成Excel文件（TODO: 使用excelize库生成）
	// 临时实现：返回下载URL
	fileName := fmt.Sprintf("dataelement_%d.xlsx", time.Now().Unix())
	downloadUrl := fmt.Sprintf("/api/v1/download/%s", fileName)

	logx.Infof("导出数据元到目录 %d, 共 %d 条记录", req.CatalogId, totalCount)

	// Step 4: 创建异步导出任务
	_ = mock.SendMqMessage(l.ctx, "export", map[string]interface{}{
		"catalogId": req.CatalogId,
		"count":     totalCount,
		"fileName":  fileName,
	}, "system")

	return &types.ExportVo{
		DownloadUrl: downloadUrl,
		FileName:    fileName,
	}, nil
}
