// Code scaffolded by speckit. Safe to edit.

package dataelement

import (
	"context"
	"testing"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/stretchr/testify/assert"
)

// TestListDataElement_Success 测试分页查询成功
func TestListDataElement_Success(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewListDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.DataElementListReq{
		CatalogId:     0,
		StdType:       0,
		State:         "",
		DataType:      0,
		RelationType:  "",
		DepartmentIds: "",
		Offset:        1,
		Limit:         10,
		Direction:     "desc",
		Sort:          "created_at",
	}

	resp, err := logic.ListDataElement(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.GreaterOrEqual(t, resp.TotalCount, int64(0))
}

// TestListDataElement_WithCatalog 测试按目录查询
func TestListDataElement_WithCatalog(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewListDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.DataElementListReq{
		CatalogId:     1,
		StdType:       0,
		State:         "",
		DataType:      0,
		RelationType:  "",
		DepartmentIds: "",
		Offset:        1,
		Limit:         10,
		Direction:     "desc",
		Sort:          "created_at",
	}

	resp, err := logic.ListDataElement(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestListDataElement_WithKeyword 测试关键字搜索
func TestListDataElement_WithKeyword(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewListDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.DataElementListReq{
		CatalogId:     0,
		StdType:       0,
		State:         "",
		DataType:      0,
		RelationType:  "",
		DepartmentIds: "",
		Keyword:       "测试",
		Offset:        1,
		Limit:         10,
		Direction:     "desc",
		Sort:          "created_at",
	}

	resp, err := logic.ListDataElement(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestListDataElement_WithState 测试按状态过滤
func TestListDataElement_WithState(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewListDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.DataElementListReq{
		CatalogId:     0,
		StdType:       0,
		State:         "enable",
		DataType:      0,
		RelationType:  "",
		DepartmentIds: "",
		Offset:        1,
		Limit:         10,
		Direction:     "desc",
		Sort:          "created_at",
	}

	resp, err := logic.ListDataElement(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// 验证所有结果都是启用状态
	for _, entry := range resp.Entries {
		assert.Equal(t, "enable", entry.State)
	}
}

// TestListDataElement_WithDataType 测试按数据类型过滤
func TestListDataElement_WithDataType(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewListDataElementLogic(ctx, &svc.ServiceContext{})

	req := &types.DataElementListReq{
		CatalogId:     0,
		StdType:       0,
		State:         "",
		DataType:      0, // Number类型
		RelationType:  "",
		DepartmentIds: "",
		Offset:        1,
		Limit:         10,
		Direction:     "desc",
		Sort:          "created_at",
	}

	resp, err := logic.ListDataElement(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// 验证所有结果都是Number类型
	for _, entry := range resp.Entries {
		assert.Equal(t, int32(0), entry.DataType)
	}
}

// TestListDataElement_Pagination 测试分页
func TestListDataElement_Pagination(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewListDataElementLogic(ctx, &svc.ServiceContext{})

	// 第一页
	req1 := &types.DataElementListReq{
		CatalogId: 0,
		Offset:    1,
		Limit:     5,
		Direction: "desc",
		Sort:      "created_at",
	}

	resp1, err := logic.ListDataElement(req1)
	assert.NoError(t, err)
	assert.NotNil(t, resp1)
	assert.LessOrEqual(t, len(resp1.Entries), 5)

	// 第二页
	req2 := &types.DataElementListReq{
		CatalogId: 0,
		Offset:    6,
		Limit:     5,
		Direction: "desc",
		Sort:      "created_at",
	}

	resp2, err := logic.ListDataElement(req2)
	assert.NoError(t, err)
	assert.NotNil(t, resp2)
}
