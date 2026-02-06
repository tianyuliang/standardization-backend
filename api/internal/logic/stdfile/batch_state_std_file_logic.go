// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package stdfile

import (
	"context"
	"strings"

	"github.com/jinguoxing/idrm-go-base/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/stdfile/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	stdfilemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/stdfile/stdfile"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchStateStdFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量启用/停用
//
// 业务流程:
//  1. 校验请求参数（ID列表、状态）
//  2. 停用时必须填写原因
//  3. 停用原因长度校验 (<= 800)
//  4. 批量更新状态
//
// 异常处理:
//   - 30202: ID列表不能为空
//   - 30220: 停用原因不能为空
//   - 30221: 停用原因长度不能超过800字符
func NewBatchStateStdFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchStateStdFileLogic {
	return &BatchStateStdFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchStateStdFileLogic) BatchStateStdFile(req *types.BatchStateReq) (resp *types.BaseResp, err error) {
	// Step 1: 校验请求参数
	if len(req.Ids) == 0 {
		return nil, errorx.NewWithMsg(30202, "ID列表不能为空")
	}

	// Step 2: 解析目标状态
	targetState, err := ParseState(req.State)
	if err != nil {
		return nil, err
	}

	// Step 3: 停用时必须填写原因
	if targetState == stdfilemodel.StateDisable {
		if strings.TrimSpace(req.Reason) == "" {
			return nil, errorx.NewWithMsg(30220, "停用原因不能为空")
		}
	}

	// Step 4: 停用原因长度校验
	if targetState == stdfilemodel.StateDisable && len([]rune(req.Reason)) > 800 {
		return nil, errorx.NewWithMsg(30221, "停用原因长度不能超过800字符")
	}

	// Step 5: 批量更新状态
	var disableReason string
	if targetState == stdfilemodel.StateDisable {
		disableReason = req.Reason
	}

	// 获取当前用户信息（Mock）
	_, _ = mock.GetUserInfo(l.ctx)

	if err := l.svcCtx.StdFileModel.BatchUpdateState(l.ctx, req.Ids, targetState, disableReason); err != nil {
		return nil, HandleError(err)
	}

	logx.Infof("批量更新标准文件状态成功: count=%d, state=%s", len(req.Ids), req.State)

	return &types.BaseResp{
		Code:        "0",
		Description: "操作成功",
	}, nil
}
