# Quick Start: 码表管理 (dict-api)

> **Feature**: dict-api
> **Branch**: `4-dict-api`
> **Date**: 2026-02-06

---

## 开发环境准备

### 1. 安装依赖

```bash
# 安装 Go 依赖
go mod tidy

# 安装 goctl (如果尚未安装)
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 安装 golang-migrate (如果尚未安装)
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 2. 配置数据库

确保 MySQL 数据库已创建并运行：

```bash
# 创建数据库（如果尚未创建）
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS standardization CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

---

## 开发流程

### Step 1: 创建 API 定义文件

```bash
# API 文件位置（已生成）
api/doc/dict/dict.api
```

**内容**：见 `plan.md` 的 API Contract 章节

### Step 2: 更新主 API 文件

在 `api/doc/api.api` 中添加导入：

```api
syntax = "v1"

import (
    "base.api"
    "rule/rule.api"
    "catalog/catalog.api"
    "stdfile/stdfile.api"
    "dict/dict.api"  // 新增
)
```

### Step 3: 生成 Handler 和 Types 代码

**⚠️ 永远使用此命令**:

```bash
# 从项目根目录执行
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

**生成文件**：
- `api/internal/handler/dict/` - Handler 文件
- `api/internal/types/dict.go` - 类型定义

### Step 4: 实现 Model 层

创建目录结构：

```bash
mkdir -p model/dict/dict
```

**文件清单**：
- `interface.go` - Model 接口定义
- `types.go` - 数据结构定义
- `vars.go` - 常量和错误定义
- `factory.go` - SQLx 工厂函数
- `sql_model.go` - SQLx 实现

### Step 5: 实现 Logic 层

```bash
# Logic 文件位置（goctl 生成框架后需补充实现）
api/internal/logic/dict/
├── list_dict_logic.go
├── get_dict_logic.go
├── get_dict_by_code_logic.go
├── create_dict_logic.go
├── update_dict_logic.go
├── delete_dict_logic.go
├── batch_delete_dict_logic.go
├── update_dict_state_logic.go
├── list_dict_enum_logic.go
├── get_dict_enum_list_logic.go
├── query_dict_by_data_element_logic.go
├── query_dict_by_std_file_catalog_logic.go
├── query_dict_by_std_file_logic.go
├── query_dict_relation_stdfile_logic.go
├── add_dict_relation_logic.go
├── query_dict_data_exists_logic.go
└── common.go  # 公共函数
```

### Step 6: 更新 ServiceContext

在 `api/internal/svc/servicecontext.go` 中注册 Model：

```go
type ServiceContext struct {
    Config config.Config

    // 注册 Dict Model
    DictModel            model.DictModel
    DictEnumModel        model.DictEnumModel
    RelationDictFileModel model.RelationDictFileModel
}
```

### Step 7: 创建错误码文件

创建 `api/internal/errorx/dict.go`：

```go
package errorx

import (
    "github.com/jinguoxing/idrm-go-base/errorx"
)

// ========== dict-api 错误码定义 (30400-30499) ==========

const (
    ErrCodeDictDataNotExist     = 30401 // 数据不存在
    ErrCodeDictParamEmpty       = 30402 // 参数为空
    ErrCodeDictInvalidParam     = 30403 // 参数无效
    ErrCodeDictCatalogNotExist   = 30404 // 目录不存在
    ErrCodeDictChNameDuplicate   = 30405 // 中文名称重复
    ErrCodeDictEnNameDuplicate   = 30406 // 英文名称重复
    ErrCodeDictEnumCodeEmpty     = 30407 // 码值为空
    ErrCodeDictEnumValueEmpty    = 30408 // 码值描述为空
    ErrCodeDictEnumCodeDuplicate = 30409 // 码值重复
    ErrCodeDictReasonTooLong     = 30410 // 停用原因过长
)

// 辅助函数...
```

### Step 8: 创建 Mock 目录

创建 Mock 目录和文件：

```bash
mkdir -p api/internal/logic/dict/mock
```

**Mock 文件**:
- `catalog.go` - 目录服务 Mock
- `dataelement.go` - 数据元服务 Mock
- `stdfile.go` - 标准文件服务 Mock
- `token.go` - Token 服务 Mock

### Step 9: 运行和测试

```bash
# 运行服务
go run api.go -f api/etc/api.yaml

# 运行测试
go test -v ./...
```

---

## API 测试示例

### 1. 创建码表

```bash
curl -X POST http://localhost:8080/api/standardization/v1/dataelement/dict \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "chName": "性别码表",
    "enName": "Gender",
    "description": "性别字典",
    "orgType": 0,
    "catalogId": 34,
    "enums": [
      {"code": "M", "value": "男"},
      {"code": "F", "value": "女"}
    ],
    "stdFiles": [1, 2, 3],
    "departmentIds": "a/ab"
  }'
```

### 2. 查询码表列表

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/dataelement/dict?offset=1&limit=20" \
  -H "Authorization: Bearer <token>"
```

### 3. 查询码表详情

```bash
curl -X GET http://localhost:8080/api/standardization/v1/dataelement/dict/1234567890 \
  -H "Authorization: Bearer <token>"
```

### 4. 按Code查询详情

```bash
curl -X GET http://localhost:8080/api/standardization/v1/dataelement/dict/code/9876543210 \
  -H "Authorization: Bearer <token>"
```

### 5. 修改码表

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/dataelement/dict/1234567890 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "chName": "性别码表（修改）",
    "enName": "Gender",
    "orgType": 0,
    "catalogId": 34,
    "enums": [
      {"code": "M", "value": "男"},
      {"code": "F", "value": "女"},
      {"code": "U", "value": "未知"}
    ]
  }'
```

### 6. 停用码表

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/dataelement/dict/state/1234567890 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "state": "disable",
    "reason": "测试停用原因"
  }'
```

### 7. 启用码表

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/dataelement/dict/state/1234567890 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "state": "enable"
  }'
```

### 8. 删除码表

```bash
curl -X DELETE http://localhost:8080/api/standardization/v1/dataelement/dict/1234567890 \
  -H "Authorization: Bearer <token>"
```

### 9. 批量删除

```bash
curl -X DELETE http://localhost:8080/api/standardization/v1/dataelement/dict/batch/1234567890,1234567891 \
  -H "Authorization: Bearer <token>"
```

### 10. 码值分页查询

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/dataelement/dict/enum?dict_id=1234567890&offset=1&limit=20" \
  -H "Authorization: Bearer <token>"
```

### 11. 码值列表查询

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/dataelement/dict/enum/getList?dict_id=1234567890" \
  -H "Authorization: Bearer <token>"
```

### 12. 查询引用的数据元

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/dataelement/dict/dataelement/1234567890?offset=1&limit=20" \
  -H "Authorization: Bearer <token>"
```

### 13. 按文件目录查询

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/dataelement/dict/queryByStdFileCatalog?catalog_id=34&offset=1&limit=20" \
  -H "Authorization: Bearer <token>"
```

### 14. 查询关联的标准文件

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/dataelement/dict/relation/stdfile/1234567890?offset=1&limit=20" \
  -H "Authorization: Bearer <token>"
```

### 15. 添加关联关系

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/dataelement/dict/relation/1234567890 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "stdFiles": [1, 2, 3]
  }'
```

### 16. 数据存在检查

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/dataelement/dict/queryDataExists?org_type=0&ch_name=性别码表" \
  -H "Authorization: Bearer <token>"
```

---

## 错误响应示例

### 参数校验失败 - 码值为空

```json
{
    "code": "30407",
    "description": "码值输入不能为空"
}
```

### 参数校验失败 - 中文名称重复

```json
{
    "code": "30405",
    "description": "码表中文名称、标准分类不能全部重复"
}
```

### 数据不存在

```json
{
    "code": "30401",
    "description": "数据不存在"
}
```

### 停用原因不能为空

```json
{
    "code": "30402",
    "description": "停用原因不能为空"
}
```

---

## 常见问题

### Q1: goctl 生成代码失败

**问题**: 执行 `goctl api go` 时报错

**解决**:
1. 检查 `api/doc/api.api` 语法是否正确
2. 确保所有 import 的文件存在
3. 检查 group 名称是否合法（不能包含连字符）

### Q2: 编译错误 - undefined: model.xxx

**问题**: Handler 或 Logic 中引用 Model 未定义

**解决**: 在 `ServiceContext` 中注册 Model

### Q3: 数据库连接失败

**问题**: 运行时报数据库连接错误

**解决**: 检查 `api/etc/api.yaml` 中的数据库配置是否正确

### Q4: 码值重复校验失败

**问题**: 创建码表时提示"码值出现重复记录"

**解决**: 确保同一码表内 code 值唯一

---

## 下一步

完成开发和测试后，执行以下命令：

```bash
# 运行完整测试
go test -v -cover ./...

# 生成 Swagger 文档
goctl api swagger -api api/doc/api.api -dir doc/swagger

# 代码格式化
gofmt -w .

# 代码质量检查
golangci-lint run

# 构建镜像
make docker-build

# 部署到 K8s
make k8s-deploy
```
