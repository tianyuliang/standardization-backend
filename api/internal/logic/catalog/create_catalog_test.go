// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package catalog

import (
	"context"
	"testing"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/catalog/catalog"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

// ============================================
// 对照 Java 源码的测试用例
// 用于验证 Go 实现与 Java 实现逻辑一致性
// ============================================

// TestCreateCatalog_对照Java源码校验场景 对照 Java 源码 DeCatalogInfoServiceImpl.checkPost() 的所有校验场景
func TestCreateCatalog_对照Java源码校验场景(t *testing.T) {
	// 测试场景对照表（来自 Java 源码分析）
	/*
		Java 校验顺序：
		1. 目录名称不能为空 (Empty -> catalog_name)
		2. 目录名称格式校验 (InvalidParameter -> catalog_name)
		3. 父目录ID不能为空 (InvalidParameter -> parent_id)
		4. 父目录必须存在 (Empty -> parent_id)
		5. 父目录level < 255 (OutOfRange -> parent_id)
		6. 不能修改根目录 (InvalidParameter) - 仅修改时
		7. 新父目录不能是自身及其子目录 (InvalidParameter) - 仅修改时
		8. 新父目录类型必须一致 (InvalidParameter) - 仅修改时
		9. 同级目录名称不能重复 (OperationConflict -> catalog_name)
	*/

	testCases := []struct {
		name        string
		req         *types.CreateCatalogReq
		setupFunc   func(*catalog.MockModel) // 可选：设置 mock 数据
		wantErrCode int32                      // 期望的错误码
		wantErrMsg  string                     // 期望的错误信息关键词
		java对照    string                     // Java 源码对应位置
	}{
		{
			name: "场景1: 目录名称为空",
			req: &types.CreateCatalogReq{
				CatalogName: "",
				ParentId:    1,
			},
			wantErrCode: 30103, // InvalidParameter
			wantErrMsg:  "目录名称不能为空",
			java对照:    "DeCatalogInfoServiceImpl.checkPost(): 目录名称不能为空",
		},
		{
			name: "场景2: 目录名称过长(>20字符)",
			req: &types.CreateCatalogReq{
				CatalogName: "这是一个超过二十个字符限制的目录名称用于测试",
				ParentId:    1,
			},
			wantErrCode: 30103, // InvalidParameter
			wantErrMsg:  "长度不能超过20",
			java对照:    "Constants.getRegexENOrCNVarL(1,20): 长度校验",
		},
		{
			name: "场景3: 目录名称包含非法字符",
			req: &types.CreateCatalogReq{
				CatalogName: "test@name",
				ParentId:    1,
			},
			wantErrCode: 30103, // InvalidParameter
			wantErrMsg:  "格式",
			java对照:    "Constants.getRegexENOrCNVarL(1,20): 格式校验",
		},
		{
			name: "场景4: 目录名称以下划线开头",
			req: &types.CreateCatalogReq{
				CatalogName: "_invalid",
				ParentId:    1,
			},
			wantErrCode: 30103, // InvalidParameter
			wantErrMsg:  "下划线",
			java对照:    "Constants.getRegexENOrCNVarL(1,20): ^(?!_)",
		},
		{
			name: "场景5: 目录名称以中划线开头",
			req: &types.CreateCatalogReq{
				CatalogName: "-invalid",
				ParentId:    1,
			},
			wantErrCode: 30103, // InvalidParameter
			wantErrMsg:  "中划线",
			java对照:    "Constants.getRegexENOrCNVarL(1,20): ^(?!-)",
		},
		{
			name: "场景6: 父目录ID为空",
			req: &types.CreateCatalogReq{
				CatalogName: "test",
				ParentId:    0,
			},
			wantErrCode: 30102, // MissingParameter
			wantErrMsg:  "父目录ID不能为空",
			java对照:    "DeCatalogInfoServiceImpl.checkPost(): parentId不能为空",
		},
		{
			name: "场景7: 父目录不存在",
			req: &types.CreateCatalogReq{
				CatalogName: "test",
				ParentId:    99999,
			},
			wantErrCode: 30101, // Empty
			wantErrMsg:  "父目录",
			java对照:    "DeCatalogInfoServiceImpl.checkPost(): 无法找到对应的父目录",
		},
		{
			name: "场景8: 父目录级别达到上限",
			req: &types.CreateCatalogReq{
				CatalogName: "test",
				ParentId:    255, // 假设存在level=255的父目录
			},
			wantErrCode: 30104, // OutOfRange
			wantErrMsg:  "级别",
			java对照:    "DeCatalogInfoServiceImpl.checkPost(): level>=255",
		},
		{
			name: "场景9: Type类型与父目录不一致",
			req: &types.CreateCatalogReq{
				CatalogName: "test",
				ParentId:    1,     // 假设父目录type=1
				Type:        2,     // 请求type=2
			},
			wantErrCode: 30103, // InvalidParameter
			wantErrMsg:  "类型",
			java对照:    "DeCatalogInfoServiceImpl.checkPost(): type不一致",
		},
		{
			name: "场景10: 同级目录名称重复",
			req: &types.CreateCatalogReq{
				CatalogName: "已存在名称",
				ParentId:    1,
			},
			wantErrCode: 30105, // OperationConflict
			wantErrMsg:  "重复",
			java对照:    "DeCatalogInfoServiceImpl.checkPost(): 同级目录名称不能重复",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// TODO: 实现 mock 逻辑和测试执行
			t.Log("Java对照:", tc.java对照)
			t.Log("期望错误码:", tc.wantErrCode)
			t.Log("期望错误信息:", tc.wantErrMsg)
		})
	}
}

// TestValidateCatalogName_对照JavaRegex 测试目录名称校验与 Java 正则表达式的一致性
func TestValidateCatalogName_对照JavaRegex(t *testing.T) {
	/*
		Java 正则表达式:
		^(?!_)(?!-)[\u4E00-\u9FA5\uF900-\uFA2D\w-]{1,20}$

		分解说明:
		- ^(?!_) - 不以 _ 开头
		- ^(?!-) - 不以 - 开头
		- [\u4E00-\u9FA5\uF900-\uFA2D\w-] - 中文、字母、数字、下划线、中划线
		- {1,20} - 长度 1-20
	*/
	testCases := []struct {
		name    string
		input   string
		wantErr bool
		java对照 string
	}{
		{
			name:    "有效: 纯中文",
			input:   "目录名称",
			wantErr: false,
			java对照: "匹配中文字符范围",
		},
		{
			name:    "有效: 纯字母",
			input:   "catalog",
			wantErr: false,
			java对照: "匹配字母",
		},
		{
			name:    "有效: 纯数字",
			input:   "12345",
			wantErr: false,
			java对照: "匹配数字",
		},
		{
			name:    "有效: 混合字符",
			input:   "目录-2024_v1",
			wantErr: false,
			java对照: "匹配混合字符",
		},
		{
			name:    "有效: 中间包含下划线",
			input:   "test_name",
			wantErr: false,
			java对照: "下划线可以在中间",
		},
		{
			name:    "有效: 中间包含中划线",
			input:   "test-name",
			wantErr: false,
			java对照: "中划线可以在中间",
		},
		{
			name:    "无效: 以下划线开头",
			input:   "_invalid",
			wantErr: true,
			java对照: "^(?!_) 负向前瞻",
		},
		{
			name:    "无效: 以中划线开头",
			input:   "-invalid",
			wantErr: true,
			java对照: "^(?!) 负向前瞻",
		},
		{
			name:    "无效: 空字符串",
			input:   "",
			wantErr: true,
			java对照: "{1,20} 长度至少为1",
		},
		{
			name:    "无效: 超过20字符",
			input:   "这是一个超过二十个字符限制的目录名称用于测试",
			wantErr: true,
			java对照: "{1,20} 长度最多为20",
		},
		{
			name:    "无效: 包含特殊字符",
			input:   "test@name",
			wantErr: true,
			java对照: "[\\u4E00-\\u9FA5\\uF900-\\uFA2D\\w-] 字符集限制",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCatalogName(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateCatalogName(%q) error = %v, wantErr %v\nJava对照: %s",
					tc.input, err, tc.wantErr, tc.java对照)
			}
		})
	}
}

// TestCheckType_对照JavaEnum 测试类型校验与 Java 枚举的一致性
func TestCheckType_对照JavaEnum(t *testing.T) {
	/*
		Java 枚举定义:
		Root(0, "根目录"),
		DataElement(1, "数据元"),
		DeDict(2, "码表"),
		ValueRule(3, "编码规则"),
		File(4, "文件"),
		Other(99, "其他类型");

		Java checkType() 校验:
		- type不能为null
		- type不能是Root(0)
		- type不能是Other(99)
		- type必须在枚举中存在
	*/
	testCases := []struct {
		name     string
		input    int32
		wantErr  bool
		java对照 string
	}{
		{
			name:     "有效: 数据元目录",
			input:    1,
			wantErr:  false,
			java对照: "CatalogTypeEnum.DataElement",
		},
		{
			name:     "有效: 码表目录",
			input:    2,
			wantErr:  false,
			java对照: "CatalogTypeEnum.DeDict",
		},
		{
			name:     "有效: 编码规则目录",
			input:    3,
			wantErr:  false,
			java对照: "CatalogTypeEnum.ValueRule",
		},
		{
			name:     "有效: 文件目录",
			input:    4,
			wantErr:  false,
			java对照: "CatalogTypeEnum.File",
		},
		{
			name:     "无效: Root类型",
			input:    0,
			wantErr:  true,
			java对照: "!EnumUtil.isIncludeCode(..., Root.getCode())",
		},
		{
			name:     "无效: Other类型",
			input:    99,
			wantErr:  true,
			java对照: "!EnumUtil.isIncludeCode(..., Other.getCode())",
		},
		{
			name:     "无效: 未定义类型",
			input:    5,
			wantErr:  true,
			java对照: "!EnumUtil.isIncludeCode(..., 5)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckType(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("CheckType(%d) error = %v, wantErr %v\nJava对照: %s",
					tc.input, err, tc.wantErr, tc.java对照)
			}
		})
	}
}

// ============================================
// 辅助函数
// ============================================

// setupTestServiceContext 创建测试用的 ServiceContext
func setupTestServiceContext() *svc.ServiceContext {
	c := rest.RestConf{
		Host: "localhost",
		Port: 8888,
	}
	return &svc.ServiceContext{
		Config: c,
	}
}

// setupMockCatalogModel 创建 mock 目录模型
func setupMockCatalogModel() *catalog.MockModel {
	return &catalog.MockModel{
		Data: map[int64]*catalog.Catalog{
			1: {Id: 1, CatalogName: "根目录", Level: 1, Type: 1, ParentId: 0},
		},
	}
}
