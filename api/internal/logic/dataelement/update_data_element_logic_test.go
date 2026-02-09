// Code scaffolded by speckit. Safe to edit.

package dataelement

import (
	"context"
	"testing"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/stretchr/testify/assert"
)

// TestUpdateDataElement_Success 测试更新数据元成功
func TestUpdateDataElement_Success(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewUpdateDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.UpdateDataElementReq{
		Id:            1,
		NameEn:        "updatedField",
		NameCn:        "更新后的字段",
		StdType:       1,
		DataType:      0,
		RelationType:  "none",
		CatalogId:     1,
		Description:   "更新后的描述",
		DepartmentIds: "1",
	}

	resp, err := logic.UpdateDataElement(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, "updatedField", resp.NameEn)
	assert.Equal(t, "更新后的字段", resp.NameCn)
}

// TestUpdateDataElement_NotFound 测试更新不存在的数据元
func TestUpdateDataElement_NotFound(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewUpdateDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.UpdateDataElementReq{
		Id:            99999,
		NameEn:        "updatedField",
		NameCn:        "更新后的字段",
		StdType:       1,
		DataType:      0,
		RelationType:  "none",
		CatalogId:     1,
		DepartmentIds: "1",
	}

	resp, err := logic.UpdateDataElement(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

// TestUpdateDataElement_VersionIncrement 测试关键字段变更时版本号递增
func TestUpdateDataElement_VersionIncrement(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewUpdateDataElementLogic(ctx, &svc.ServiceContext{})

	// 先获取当前版本
	getLogic := NewGetDataElementDetailLogic(ctx, &svc.ServiceContext{})
	detailReq := &types.DataElementDetailReq{Type: 1, Value: 1}
	oldDetail, _ := getLogic.GetDataElementDetail(detailReq)
	oldVersion := oldDetail.Version

	// 更新名称（触发版本递增）
	req := &types.UpdateDataElementReq{
		Id:            1,
		NameEn:        "newName",
		NameCn:        "新名称",
		StdType:       1,
		DataType:      0,
		RelationType:  "none",
		CatalogId:     1,
		DepartmentIds: "1",
	}

	resp, err := logic.UpdateDataElement(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Greater(t, resp.Version, oldVersion)
}

// TestUpdateDataElement_ChangeRelationType 测试更改关联类型
func TestUpdateDataElement_ChangeRelationType(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewUpdateDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.UpdateDataElementReq{
		Id:            1,
		NameEn:        "testField",
		NameCn:        "测试字段",
		StdType:       1,
		DataType:      0,
		RelationType:  "codeTable", // 从none改为codeTable
		DictCode:      "TEST_DICT",
		CatalogId:     1,
		DepartmentIds: "1",
	}

	resp, err := logic.UpdateDataElement(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "codeTable", resp.RelationType)
}

// TestUpdateDataElement_UpdateStdFiles 测试更新关联文件
func TestUpdateDataElement_UpdateStdFiles(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewUpdateDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.UpdateDataElementReq{
		Id:            1,
		NameEn:        "testField",
		NameCn:        "测试字段",
		StdType:       1,
		DataType:      0,
		RelationType:  "none",
		CatalogId:     1,
		DepartmentIds: "1",
		StdFiles:      []int64{1, 2, 3},
	}

	resp, err := logic.UpdateDataElement(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.StdFiles, 3)
}

// TestUpdateDataElement_CatalogNotExist 测试目录不存在
func TestUpdateDataElement_CatalogNotExist(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewUpdateDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.UpdateDataElementReq{
		Id:            1,
		NameEn:        "testField",
		NameCn:        "测试字段",
		StdType:       1,
		DataType:      0,
		RelationType:  "none",
		CatalogId:     99999, // 不存在的目录
		DepartmentIds: "1",
	}

	resp, err := logic.UpdateDataElement(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

// TestUpdateDataElement_WithDataLength 测试更新数据长度
func TestUpdateDataElement_WithDataLength(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewUpdateDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.UpdateDataElementReq{
		Id:            1,
		NameEn:        "testField",
		NameCn:        "测试字段",
		StdType:       1,
		DataType:      0,
		DataLength:    10,
		RelationType:  "none",
		CatalogId:     1,
		DepartmentIds: "1",
	}

	resp, err := logic.UpdateDataElement(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(10), resp.DataLength)
}

// TestUpdateDataElement_WithDataPrecision 测试更新数据精度
func TestUpdateDataElement_WithDataPrecision(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewUpdateDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.UpdateDataElementReq{
		Id:            1,
		NameEn:        "testField",
		NameCn:        "测试字段",
		StdType:       1,
		DataType:      1, // Decimal
		DataLength:    10,
		DataPrecision: 2,
		RelationType:  "none",
		CatalogId:     1,
		DepartmentIds: "1",
	}

	resp, err := logic.UpdateDataElement(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(10), resp.DataLength)
	assert.Equal(t, int32(2), resp.DataPrecision)
}
