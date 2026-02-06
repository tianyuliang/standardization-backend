// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/rule/mock"
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
// 对应 Java: (查询关联文件逻辑)
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
		return nil, errorx.RuleRecordNotExist()
	}

	// ====== 步骤2: 查询关联的文件ID列表 ======
	relations, err := l.svcCtx.RelationRuleFileModel.FindByRuleId(l.ctx, id)
	if err != nil {
		return nil, err
	}

	// 提取文件ID
	fileIds := make([]int64, 0, len(relations))
	for _, r := range relations {
		fileIds = append(fileIds, r.FileId)
	}

	// ====== 步骤3: 查询文件详情 ======
	// MOCK: mock.StdFileGetById() - 批量获取文件信息
	_ = mock.StdFileGetById(l.ctx, l.svcCtx, fileIds)

	// 构建响应（TODO: 转换为 StdFileResp）
	return &types.StdFileListResp{
		Entries:    []types.StdFileResp{},
		TotalCount: int64(len(fileIds)),
	}, nil
}
