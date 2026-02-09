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

type UpdateDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改码表
func NewUpdateDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictLogic {
	return &UpdateDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDictLogic) UpdateDict(req *types.UpdateDictReq, id int64) (resp *types.DictVo, err error) {
	// Step 1: 校验码表存在性
	oldDict, err := l.svcCtx.DictModel.FindOne(l.ctx, id)
	if err != nil {
		return nil, localErrorx.DictDataNotExist()
	}

	// Step 2: 业务校验（名称唯一性、码值唯一性）
	// 修改时排除当前记录的中文名称校验
	if oldDict.ChName != req.ChName {
		if err := checkChNameUniqueExcludeSelf(l.svcCtx.DictModel, id, req.ChName, req.OrgType); err != nil {
			return nil, err
		}
	}

	// 修改时排除当前记录的英文名称校验
	if oldDict.EnName != req.EnName {
		if err := checkEnNameUniqueExcludeSelf(l.svcCtx.DictModel, id, req.EnName, req.OrgType); err != nil {
			return nil, err
		}
	}

	// 校验码值唯一性
	if err := checkEnumCodesUnique(req.Enums); err != nil {
		return nil, err
	}

	// 获取用户信息
	username, _ := getUserInfo(l.ctx)

	// Step 3: 检查是否需要更新版本
	needVersionUpdate := checkVersionChange(oldDict, *req, req.Enums, req.StdFiles)

	// 更新数据
	dictData := &dictmodel.Dict{
		Id:            id,
		Code:          oldDict.Code,
		ChName:        req.ChName,
		EnName:        req.EnName,
		Description:   req.Description,
		CatalogId:     req.CatalogId,
		OrgType:       req.OrgType,
		Version:       oldDict.Version,
		State:         oldDict.State,
		DepartmentIds: req.DepartmentIds,
		CreateUser:    oldDict.CreateUser,
		UpdateUser:    username,
	}

	if needVersionUpdate {
		dictData.Version = oldDict.Version + 1
	}

	// Step 4: 有变更则更新（更新码表、版本号+1、重新保存码值、更新关联文件）
	if err := l.svcCtx.DictModel.Update(l.ctx, dictData); err != nil {
		logx.Errorf("更新码表失败: %v", err)
		return nil, localErrorx.DictInvalidParam("更新码表失败")
	}

	// 删除旧的码值并重新保存
	if err := l.svcCtx.DictEnumModel.DeleteByDictId(l.ctx, id); err != nil {
		logx.Errorf("删除旧码值失败: %v", err)
	}
	for _, enum := range req.Enums {
		enumData := &dictmodel.DictEnum{
			DictId:     id,
			Code:       enum.Code,
			Value:      enum.Value,
			CreateUser: username,
		}
		if _, err := l.svcCtx.DictEnumModel.Insert(l.ctx, enumData); err != nil {
			logx.Errorf("保存码值失败: %v", err)
		}
	}

	// 删除旧的关联文件并重新保存
	if err := l.svcCtx.RelationDictFileModel.DeleteByDictId(l.ctx, id); err != nil {
		logx.Errorf("删除旧关联文件失败: %v", err)
	}
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
			logx.Errorf("保存关联文件失败: %v", err)
		}
	}

	// Step 5: 查询完整数据返回
	enums, _ := l.svcCtx.DictEnumModel.FindByDictId(l.ctx, id)
	catalogName := getCatalogName(req.CatalogId)
	deptName, deptPathNames, _ := getDeptInfo(req.DepartmentIds)

	result := buildDictVo(l.ctx, dictData, enums, catalogName, deptName, deptPathNames, false)

	logx.Infof("码表修改成功: id=%d, chName=%s", id, req.ChName)
	return &result, nil
}
