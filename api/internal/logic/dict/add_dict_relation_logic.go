// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	localErrorx "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	dictmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/dict/dict"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddDictRelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 添加关联关系
func NewAddDictRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddDictRelationLogic {
	return &AddDictRelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddDictRelationLogic) AddDictRelation(req *types.AddDictRelationReq, id int64) (resp *types.DictBaseResp, err error) {
	// Step 1: 校验码表存在性
	dict, err := l.svcCtx.DictModel.FindOne(l.ctx, id)
	if err != nil {
		return nil, localErrorx.DictDataNotExist()
	}

	// Step 2: 删除旧的关联文件关系
	if err := l.svcCtx.RelationDictFileModel.DeleteByDictId(l.ctx, id); err != nil {
		logx.Errorf("删除旧关联文件失败: dictId=%d, error=%v", id, err)
		// 继续执行
	}

	// Step 3: 添加新的关联文件关系
	if len(req.StdFiles) > 0 {
		var relations []*dictmodel.RelationDictFile
		for _, fileId := range req.StdFiles {
			relations = append(relations, &dictmodel.RelationDictFile{
				Id:     generateSnowflakeCode(),
				DictId: id,
				FileId: fileId,
			})
		}
		if err := l.svcCtx.RelationDictFileModel.InsertBatch(l.ctx, relations); err != nil {
			logx.Errorf("保存关联文件失败: dictId=%d, error=%v", id, err)
			return nil, localErrorx.DictInvalidParam("保存关联文件失败")
		}
	}

	logx.Infof("添加码表关联关系成功: id=%d, code=%d, fileCount=%d", id, dict.Code, len(req.StdFiles))
	return &types.DictBaseResp{Code: "0", Description: "添加关联关系成功"}, nil
}
