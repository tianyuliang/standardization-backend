// Code scaffolded by speckit. Safe to edit.

package dataelement

import (
	"context"
	"testing"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/stretchr/testify/assert"
)

// TestCreateDataElement_Success 测试创建数据元成功场景
func TestCreateDataElement_Success(t *testing.T) {
	// TODO: 初始化测试数据库和依赖
	// 1. 创建测试数据库连接
	// 2. 初始化 ServiceContext
	// 3. 准备测试数据

	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewCreateDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.CreateDataElementReq{
		NameEn:        "testField",
		NameCn:        "测试字段",
		StdType:       1,
		DataType:      0,
		RelationType:  "none",
		CatalogId:     1,
		Description:   "测试描述",
		DepartmentIds: "1",
		ThirdDeptId:   "1",
	}

	resp, err := logic.CreateDataElement(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Greater(t, resp.Id, int64(0))
	assert.Equal(t, "testField", resp.NameEn)
	assert.Equal(t, "测试字段", resp.NameCn)
}

// TestCreateDataElement_DuplicateNameCn 测试中文名称重复
func TestCreateDataElement_DuplicateNameCn(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewCreateDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.CreateDataElementReq{
		NameEn:        "testField2",
		NameCn:        "测试字段", // 重复的中文名
		StdType:       1,
		DataType:      0,
		RelationType:  "none",
		CatalogId:     1,
		DepartmentIds: "1",
	}

	_, err := logic.CreateDataElement(req)
	assert.Error(t, err)
}

// TestCreateDataElement_CatalogNotExist 测试目录不存在
func TestCreateDataElement_CatalogNotExist(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewCreateDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.CreateDataElementReq{
		NameEn:        "testField3",
		NameCn:        "测试字段3",
		StdType:       1,
		DataType:      0,
		RelationType:  "none",
		CatalogId:     99999, // 不存在的目录
		DepartmentIds: "1",
	}

	_, err := logic.CreateDataElement(req)
	assert.Error(t, err)
}

// TestCreateDataElement_WithCodeTable 测试关联码表
func TestCreateDataElement_WithCodeTable(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewCreateDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.CreateDataElementReq{
		NameEn:        "testField4",
		NameCn:        "测试字段4",
		StdType:       1,
		DataType:      0,
		RelationType:  "codeTable",
		DictCode:      "TEST_DICT",
		CatalogId:     1,
		DepartmentIds: "1",
	}

	resp, err := logic.CreateDataElement(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "codeTable", resp.RelationType)
}

// TestCreateDataElement_WithRule 测试关联编码规则
func TestCreateDataElement_WithRule(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewCreateDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.CreateDataElementReq{
		NameEn:        "testField5",
		NameCn:        "测试字段5",
		StdType:       1,
		DataType:      0,
		RelationType:  "codeRule",
		RuleId:        1,
		CatalogId:     1,
		DepartmentIds: "1",
	}

	resp, err := logic.CreateDataElement(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "codeRule", resp.RelationType)
}

// TestCreateDataElement_WithStdFiles 测试关联标准文件
func TestCreateDataElement_WithStdFiles(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewCreateDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.CreateDataElementReq{
		NameEn:        "testField6",
		NameCn:        "测试字段6",
		StdType:       1,
		DataType:      0,
		RelationType:  "none",
		CatalogId:     1,
		DepartmentIds: "1",
		StdFiles:      []int64{1, 2},
	}

	resp, err := logic.CreateDataElement(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.StdFiles)
}
