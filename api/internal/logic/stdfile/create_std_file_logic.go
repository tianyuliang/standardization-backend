// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package stdfile

import (
	"context"

	"github.com/jinguoxing/idrm-go-base/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	stdfilemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/stdfile/stdfile"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateStdFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增标准文件
func NewCreateStdFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateStdFileLogic {
	return &CreateStdFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateStdFileLogic) CreateStdFile(req *types.CreateStdFileReq) (resp *types.StdFileDetailResp, err error) {
	// Step 1: 参数校验
	if err := ValidateRequiredString(req.Name, "标准文件名称"); err != nil {
		return nil, err
	}

	if err := ValidateOrgType(int(req.OrgType)); err != nil {
		return nil, err
	}

	// Step 2: 解析附件类型
	attachmentType, err := ParseAttachmentType(req.AttachmentType)
	if err != nil {
		return nil, err
	}

	// Step 3: 校验目录是否存在
	if req.CatalogId <= 0 {
		req.CatalogId = 44 // 默认目录
	}

	// Step 4: 检查标准编号是否重复
	if req.Number != "" {
		existing, err := l.svcCtx.StdFileModel.FindByNumber(l.ctx, req.Number)
		if err != nil {
			return nil, HandleError(err)
		}
		if len(existing) > 0 {
			return nil, errorx.NewWithMsg(30210, "标准编号已存在")
		}
	}

	// Step 5: 检查名称+组织类型是否重复
	existing, err := l.svcCtx.StdFileModel.FindByNameAndOrgType(l.ctx, req.Name, int(req.OrgType))
	if err != nil {
		return nil, HandleError(err)
	}
	if len(existing) > 0 {
		return nil, errorx.NewWithMsg(30204, "标准文件名称+组织类型已存在")
	}

	// Step 6: 解析日期
	actDate, err := ParseActDate(req.ActDate)
	if err != nil {
		return nil, err
	}

	publishDate, err := ParseActDate(req.PublishDate)
	if err != nil {
		return nil, err
	}

	// Step 7: 解析状态
	state := stdfilemodel.StateEnable // 默认启用
	if req.State != "" {
		state, err = ParseState(req.State)
		if err != nil {
			return nil, err
		}
	}

	// Step 8: 验证文件扩展名（如果是文件类型附件）
	if attachmentType == stdfilemodel.AttachmentTypeFile && req.AttachmentUrl != "" {
		if err := ValidateFileExtension(req.AttachmentUrl); err != nil {
			return nil, err
		}
	}

	// Step 9: 构建数据模型
	model := &stdfilemodel.StdFile{
		Number:         req.Number,
		Name:           req.Name,
		CatalogId:      req.CatalogId,
		OrgType:        int(req.OrgType),
		ActDate:        actDate,
		PublishDate:    publishDate,
		Description:    req.Description,
		AttachmentType: attachmentType,
		AttachmentUrl:  req.AttachmentUrl,
		State:          state,
		AuthorityId:    req.DepartmentIds,
		Version:        1,
	}

	// Step 10: 插入数据库
	id, err := l.svcCtx.StdFileModel.Insert(l.ctx, model)
	if err != nil {
		return nil, HandleError(err)
	}

	// Step 11: 查询插入后的完整数据
	created, err := l.svcCtx.StdFileModel.FindOne(l.ctx, id)
	if err != nil {
		return nil, HandleError(err)
	}

	// Step 12: 转换为响应对象
	resp = &types.StdFileDetailResp{}
	*resp = ModelToResp(l.ctx, l.svcCtx, created)

	logx.Infof("创建标准文件成功: id=%d, number=%s, name=%s", id, req.Number, req.Name)
	return resp, nil
}
