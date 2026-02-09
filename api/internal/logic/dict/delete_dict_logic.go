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

type DeleteDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除码表
func NewDeleteDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictLogic {
	return &DeleteDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDictLogic) DeleteDict(req *types.IdReq) (resp *types.DictBaseResp, err error) {
	// Step 1: 校验码表存在性
	dict, err := l.svcCtx.DictModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, localErrorx.DictDataNotExist()
	}

	// Step 2: 删除码表（物理删除）
	if err := l.svcCtx.DictModel.Delete(l.ctx, req.Id); err != nil {
		logx.Errorf("删除码表失败: id=%d, error=%v", req.Id, err)
		return nil, localErrorx.DictInvalidParam("删除码表失败")
	}

	// Step 3: 删除关联的码值
	if err := l.svcCtx.DictEnumModel.DeleteByDictId(l.ctx, req.Id); err != nil {
		logx.Errorf("删除码值失败: dictId=%d, error=%v", req.Id, err)
		// 继续执行，不返回错误
	}

	// Step 4: 删除关联文件关系
	if err := l.svcCtx.RelationDictFileModel.DeleteByDictId(l.ctx, req.Id); err != nil {
		logx.Errorf("删除关联文件失败: dictId=%d, error=%v", req.Id, err)
		// 继续执行，不返回错误
	}

	logx.Infof("码表删除成功: id=%d, code=%d, chName=%s", req.Id, dict.Code, dict.ChName)
	return &types.DictBaseResp{Code: "0", Description: "删除成功"}, nil
}
