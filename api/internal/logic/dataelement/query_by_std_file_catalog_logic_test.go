// Code scaffolded by speckit. Safe to edit.

package dataelement

import (
	"context"
	"testing"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/stretchr/testify/assert"
)

// TestQueryByStdFileCatalog_Success 测试按文件目录查询成功
func TestQueryByStdFileCatalog_Success(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewQueryByStdFileCatalogLogic(ctx, &svc.ServiceContext{})

	req := &types.DataElementQueryByFileCatalogReq{
		CatalogId: 1,
		Offset:    1,
		Limit:     10,
		Direction: "desc",
		Sort:      "created_at",
	}

	resp, err := logic.QueryByStdFileCatalog(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.GreaterOrEqual(t, resp.TotalCount, int64(0))
}

// TestQueryByStdFileCatalog_Pagination 测试分页
func TestQueryByStdFileCatalog_Pagination(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewQueryByStdFileCatalogLogic(ctx, &svc.ServiceContext{})

	// 第一页
	req1 := &types.DataElementQueryByFileCatalogReq{
		CatalogId: 1,
		Offset:    1,
		Limit:     5,
		Direction: "desc",
		Sort:      "created_at",
	}

	resp1, err := logic.QueryByStdFileCatalog(req1)
	assert.NoError(t, err)
	assert.NotNil(t, resp1)
	assert.LessOrEqual(t, len(resp1.Entries), 5)

	// 第二页
	req2 := &types.DataElementQueryByFileCatalogReq{
		CatalogId: 1,
		Offset:    6,
		Limit:     5,
		Direction: "desc",
		Sort:      "created_at",
	}

	resp2, err := logic.QueryByStdFileCatalog(req2)
	assert.NoError(t, err)
	assert.NotNil(t, resp2)
}

// TestQueryByStdFileCatalog_SortAsc 测试正序排序
func TestQueryByStdFileCatalog_SortAsc(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewQueryByStdFileCatalogLogic(ctx, &svc.ServiceContext{})

	req := &types.DataElementQueryByFileCatalogReq{
		CatalogId: 1,
		Offset:    1,
		Limit:     10,
		Direction: "asc",
		Sort:      "created_at",
	}

	resp, err := logic.QueryByStdFileCatalog(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestQueryByStdFileCatalog_EmptyResult 测试空结果
func TestQueryByStdFileCatalog_EmptyResult(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewQueryByStdFileCatalogLogic(ctx, &svc.ServiceContext{})

	req := &types.DataElementQueryByFileCatalogReq{
		CatalogId: 99999, // 不存在的目录
		Offset:    1,
		Limit:     10,
	}

	resp, err := logic.QueryByStdFileCatalog(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(0), resp.TotalCount)
	assert.Empty(t, resp.Entries)
}
