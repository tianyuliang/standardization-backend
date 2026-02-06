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

type UpdateStdFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据ID修改标准文件
func NewUpdateStdFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStdFileLogic {
	return &UpdateStdFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateStdFileLogic) UpdateStdFile(req *types.UpdateStdFileReq) (resp *types.StdFileDetailResp, err error) {
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

	// Step 3: 校验文件是否存在
	existing, err := l.svcCtx.StdFileModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, errorx.NewWithMsg(30201, "标准文件不存在")
	}
	if existing == nil {
		return nil, errorx.NewWithMsg(30201, "标准文件不存在")
	}

	// Step 4: 检查标准编号是否重复（排除自身）
	if req.Number != "" && req.Number != existing.Number {
		duplicate, err := l.svcCtx.StdFileModel.FindByNumber(l.ctx, req.Number)
		if err != nil {
			return nil, HandleError(err)
		}
		if len(duplicate) > 0 {
			return nil, errorx.NewWithMsg(30210, "标准编号已存在")
		}
	}

	// Step 5: 检查名称+组织类型是否重复（排除自身）
	if req.Name != existing.Name || int(req.OrgType) != existing.OrgType {
		duplicate, err := l.svcCtx.StdFileModel.FindByNameAndOrgType(l.ctx, req.Name, int(req.OrgType))
		if err != nil {
			return nil, HandleError(err)
		}
		// 过滤掉自身
		for _, item := range duplicate {
			if item.Id != req.Id {
				return nil, errorx.NewWithMsg(30204, "标准文件名称+组织类型已存在")
			}
		}
	}

	// Step 6: 校验目录是否存在
	if err := ValidateCatalogId(req.CatalogId); err != nil {
		return nil, err
	}

	// Step 7: 解析日期
	actDate, err := ParseActDate(req.ActDate)
	if err != nil {
		return nil, err
	}

	publishDate, err := ParseActDate(req.PublishDate)
	if err != nil {
		return nil, err
	}

	// Step 8: 解析状态
	state := existing.State // 默认保持原状态
	if req.State != "" {
		state, err = ParseState(req.State)
		if err != nil {
			return nil, err
		}
	}

	// Step 9: 验证文件扩展名（如果是文件类型附件）
	if attachmentType == stdfilemodel.AttachmentTypeFile && req.AttachmentUrl != "" {
		if err := ValidateFileExtension(req.AttachmentUrl); err != nil {
			return nil, err
		}
	}

	// Step 10: 检查版本变更（version递增逻辑）
	newVersion := existing.Version + 1
	// 如果状态从启用变为停用，或者状态从停用变为启用，version需要递增
	// 其他情况如果内容有变化也需要递增

	// Step 11: 构建更新数据模型
	updateModel := &stdfilemodel.StdFile{
		Id:             req.Id,
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
		DepartmentIds:  req.DepartmentIds,
		ThirdDeptId:    existing.ThirdDeptId,
		Version:        newVersion,
		DisableReason:  existing.DisableReason,
		FileName:       existing.FileName,
		DisableDate:    existing.DisableDate,
		CreateTime:     existing.CreateTime,
		CreateUser:     existing.CreateUser,
		UpdateUser:     existing.UpdateUser,
		Deleted:        existing.Deleted,
	}

	// Step 12: 更新数据库
	if err := l.svcCtx.StdFileModel.Update(l.ctx, updateModel); err != nil {
		return nil, HandleError(err)
	}

	// Step 13: 查询更新后的完整数据
	updated, err := l.svcCtx.StdFileModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, HandleError(err)
	}

	// Step 14: 转换为响应对象
	resp = &types.StdFileDetailResp{}
	*resp = ModelToResp(l.ctx, l.svcCtx, updated)

	logx.Infof("修改标准文件成功: id=%d, number=%s, name=%s, version=%d->%d",
		req.Id, existing.Number, req.Name, existing.Version, newVersion)
	return resp, nil
}
