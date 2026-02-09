// Code scaffolded by speckit. Safe to edit.

package dataelement

import (
	"context"
	"testing"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/stretchr/testify/assert"
)

// TestIsRepeat_NameCnDuplicate 测试中文名称重复
func TestIsRepeat_NameCnDuplicate(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewIsRepeatLogic(ctx, &svc.ServiceContext{})

	req := &types.IsRepeatReq{
		NameCn:  "测试字段", // 假设已存在
		StdType: 1,
	}

	resp, err := logic.IsRepeat(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Result)
}

// TestIsRepeat_NameEnDuplicate 测试英文名称重复
func TestIsRepeat_NameEnDuplicate(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewIsRepeatLogic(ctx, &svc.ServiceContext{})

	req := &types.IsRepeatReq{
		NameEn: "testField", // 假设已存在
	}

	resp, err := logic.IsRepeat(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Result)
}

// TestIsRepeat_NoRepeat 测试名称不重复
func TestIsRepeat_NoRepeat(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewIsRepeatLogic(ctx, &svc.ServiceContext{})

	req := &types.IsRepeatReq{
		NameCn:  "新字段名称123",
		NameEn:  "newField123",
		StdType: 1,
	}

	resp, err := logic.IsRepeat(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Result)
}

// TestIsRepeat_WithExcludeId 测试排除指定ID
func TestIsRepeat_WithExcludeId(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewIsRepeatLogic(ctx, &svc.ServiceContext{})

	req := &types.IsRepeatReq{
		NameCn:  "测试字段",
		NameEn:  "testField",
		StdType: 1,
		Id:      1, // 排除ID=1的数据元
	}

	resp, err := logic.IsRepeat(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestIsRepeat_EmptyNames 测试空名称
func TestIsRepeat_EmptyNames(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewIsRepeatLogic(ctx, &svc.ServiceContext{})

	req := &types.IsRepeatReq{
		NameCn:  "",
		NameEn:  "",
		StdType: 1,
	}

	resp, err := logic.IsRepeat(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Result)
}

// TestIsRepeat_OnlyCheckNameCn 测试只检查中文名称
func TestIsRepeat_OnlyCheckNameCn(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewIsRepeatLogic(ctx, &svc.ServiceContext{})

	req := &types.IsRepeatReq{
		NameCn:  "测试字段",
		NameEn:  "",
		StdType: 1,
	}

	resp, err := logic.IsRepeat(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestIsRepeat_OnlyCheckNameEn 测试只检查英文名称
func TestIsRepeat_OnlyCheckNameEn(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewIsRepeatLogic(ctx, &svc.ServiceContext{})

	req := &types.IsRepeatReq{
		NameCn: "",
		NameEn: "testField",
	}

	resp, err := logic.IsRepeat(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestIsRepeat_DifferentStdType 测试不同标准类型
func TestIsRepeat_DifferentStdType(t *testing.T) {
	t.Skip("需要配置测试数据库")

	ctx := context.Background()
	logic := NewIsRepeatLogic(ctx, &svc.ServiceContext{})

	req := &types.IsRepeatReq{
		NameCn:  "测试字段",
		StdType: 2, // 不同的标准类型
	}

	resp, err := logic.IsRepeat(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
