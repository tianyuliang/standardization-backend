// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryStdFilesByRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询规则关联的标准文件
//
// 业务流程:
//  1. 校验规则存在
//  2. 查询关联的文件ID列表
//  3. 查询文件详情
func NewQueryStdFilesByRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryStdFilesByRuleLogic {
	return &QueryStdFilesByRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryStdFilesByRuleLogic) QueryStdFilesByRule(id int64, req *types.PageQuery) (resp *types.StdFileListResp, err error) {
	// ====== 步骤1: 校验规则存在 ======
	_, err = l.svcCtx.RuleModel.FindOne(l.ctx, id)
	if err != nil {
		// TODO: 返回 errorx.RuleNotExist(id) [错误码 30301]
		return nil, err
	}

	// ====== 步骤2: 查询关联的文件ID列表 ======
	relations, err := l.svcCtx.RelationRuleFileModel.FindByRuleId(l.ctx, id)
	if err != nil {
		return nil, err
	}

	// ====== 步骤3: 查询文件详情 ======
	// TODO: 调用 StdFile RPC 批量获取文件信息
	// 当前返回空列表
	fileIds := make([]int64, 0, len(relations))
	for _, r := range relations {
		fileIds = append(fileIds, r.FileId)
	}

	return &types.StdFileListResp{
		Entries:    []types.StdFileResp{},
		TotalCount: int64(len(fileIds)),
	}, nil
}
