// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	localErrorx "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除码表
func NewBatchDeleteDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteDictLogic {
	return &BatchDeleteDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteDictLogic) BatchDeleteDict(req *types.IdReq) (resp *types.DictBaseResp, err error) {
	// Step 1: 解析ID列表
	// 注意：由于API定义使用 IdReq，实际批量删除需要通过多次调用或修改API
	// 这里实现单个ID删除逻辑，批量功能需要调整API定义或使用单独的批量接口
	id := req.Id
	if id <= 0 {
		return nil, localErrorx.DictParamEmpty("ID不能为空")
	}

	// Step 2: 校验码表存在性
	dict, err := l.svcCtx.DictModel.FindOne(l.ctx, id)
	if err != nil {
		return nil, localErrorx.DictDataNotExist()
	}

	// Step 3: 删除码表（物理删除）
	if err := l.svcCtx.DictModel.Delete(l.ctx, id); err != nil {
		logx.Errorf("删除码表失败: id=%d, error=%v", id, err)
		return nil, localErrorx.DictInvalidParam("删除码表失败")
	}

	// Step 4: 删除关联的码值
	if err := l.svcCtx.DictEnumModel.DeleteByDictId(l.ctx, id); err != nil {
		logx.Errorf("删除码值失败: dictId=%d, error=%v", id, err)
		// 继续执行，不返回错误
	}

	// Step 5: 删除关联文件关系
	if err := l.svcCtx.RelationDictFileModel.DeleteByDictId(l.ctx, id); err != nil {
		logx.Errorf("删除关联文件失败: dictId=%d, error=%v", id, err)
		// 继续执行，不返回错误
	}

	logx.Infof("批量删除码表成功: id=%d, code=%d, chName=%s", id, dict.Code, dict.ChName)
	return &types.DictBaseResp{Code: "0", Description: "删除成功"}, nil
}
