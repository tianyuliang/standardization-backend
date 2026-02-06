# Quick Start: 目录管理 (catalog-api)

> **Feature**: catalog-api
> **Branch**: `2-catalog-api`
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
# API 文件位置
api/doc/catalog/catalog.api
```

**内容**：包含 6 个 API 端点定义

### Step 2: 更新主 API 文件

在 `api/doc/api.api` 中添加导入：

```api
syntax = "v1"

import (
    "base.api"
    "catalog/catalog.api"  // 新增
)

// ... 其他内容
```

### Step 3: 生成 Handler 和 Types 代码

⚠️ **每次修改 .api 文件后必须执行**：

```bash
# 从项目根目录执行
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

**生成文件**：
- `api/internal/handler/catalog/` - Handler 文件
- `api/internal/types/types.go` - 类型定义（追加）
- `api/internal/logic/catalog/` - Logic 框架文件

### Step 4: 实现 Model 层

创建目录结构：

```bash
mkdir -p model/catalog/catalog
```

**文件清单**：
- `interface.go` - Model 接口定义
- `types.go` - 数据结构定义
- `vars.go` - 枚举常量定义
- `factory.go` - SQLx 工厂函数

### Step 5: 实现 Logic 层

```bash
# Logic 文件位置
api/internal/logic/catalog/
├── create_catalog_logic.go
├── update_catalog_logic.go
├── delete_catalog_logic.go
├── query_tree_logic.go
├── query_logic.go
├── query_with_file_logic.go
└── common.go  # 公共函数
```

### Step 6: 更新 ServiceContext

在 `api/internal/svc/service_context.go` 中注册 Model：

```go
type ServiceContext struct {
    Config       config.Config
    db           *sqlx.DB
    RuleModel    rulemodel.RuleModel
    CatalogModel catalogmodel.CatalogModel  // 新增
}
```

### Step 7: 运行和测试

```bash
# 运行服务
go run api.go -f api/etc/api.yaml

# 运行测试
go test -v ./...
```

---

## API 测试示例

### 1. 查询目录树（按类型）

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/catalog/query_tree?type=1" \
  -H "Authorization: Bearer <token>"
```

### 2. 查询目录树（按ID）

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/catalog/query_tree?id=100" \
  -H "Authorization: Bearer <token>"
```

### 3. 检索目录（带关键字）

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/catalog/query?type=1&keyword=测试" \
  -H "Authorization: Bearer <token>"
```

### 4. 检索目录（无关键字）

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/catalog/query?type=1" \
  -H "Authorization: Bearer <token>"
```

### 5. 创建目录

```bash
curl -X POST http://localhost:8080/api/standardization/v1/catalog \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "catalogName": "测试目录",
    "parentId": 100,
    "type": 1,
    "description": "测试目录说明"
  }'
```

### 6. 修改目录

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/catalog/100 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "catalogName": "修改后的目录名称",
    "description": "修改后的目录说明"
  }'
```

### 7. 删除目录

```bash
curl -X DELETE http://localhost:8080/api/standardization/v1/catalog/100 \
  -H "Authorization: Bearer <token>"
```

### 8. 查询目录及文件树

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/catalog/query/with_file" \
  -H "Authorization: Bearer <token>"
```

### 9. 查询目录及文件树（带关键字）

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/catalog/query/with_file?keyword=测试" \
  -H "Authorization: Bearer <token>"
```

---

## 错误响应示例

### 参数校验失败 - 目录名称为空

```json
{
    "code": "30103",
    "description": "参数值校验不通过",
    "detail": [
        {"Key": "catalogName", "Message": "目录名称不能为空"}
    ],
    "solution": "请使用请求参数构造规范化的请求字符串。"
}
```

### 参数校验失败 - 目录名称过长

```json
{
    "code": "30103",
    "description": "参数值校验不通过",
    "detail": [
        {"Key": "catalogName", "Message": "目录名称长度不能超过20个字符"}
    ],
    "solution": "请使用请求参数构造规范化的请求字符串。"
}
```

### 参数校验失败 - 目录名称格式错误

```json
{
    "code": "30103",
    "description": "参数值校验不通过",
    "detail": [
        {"Key": "catalogName", "Message": "目录名称只能由中英文、数字、下划线、中划线组成"}
    ],
    "solution": "请使用请求参数构造规范化的请求字符串。"
}
```

### 参数校验失败 - 同级目录名称重复

```json
{
    "code": "30105",
    "description": "同级目录名称不能重复"
}
```

### 参数校验失败 - 父目录不存在

```json
{
    "code": "30101",
    "description": "无法找到对应的父目录"
}
```

### 参数校验失败 - 目录级别超限

```json
{
    "code": "30104",
    "description": "目录级别取值范围(1-255)"
}
```

### 参数校验失败 - 不允许删除根目录

```json
{
    "code": "30103",
    "description": "不允许删除根目录"
}
```

### 参数校验失败 - 目录或子目录下已存在数据

```json
{
    "code": "30106",
    "description": "目录或子目录下已存在数据，不允许删除"
}
```

### 参数校验失败 - 目录类型无效

```json
{
    "code": "30103",
    "description": "此类型不在有效值范围内"
}
```

### 参数校验失败 - 不能修改根目录

```json
{
    "code": "30103",
    "description": "不能修改根目录"
}
```

---

## 常见问题

### Q1: goctl 生成代码失败

**问题**: 执行 `goctl api go` 时报错

**解决**:
1. 检查 `api.doc/api.api` 语法是否正确
2. 确保所有 import 的文件存在
3. 检查 group 名称是否合法（不能包含连字符）

### Q2: 编译错误 - undefined: model.xxx

**问题**: Handler 或 Logic 中引用 Model 未定义

**解决**: 在 `ServiceContext` 中注册 Model，并在 Logic 中使用 `l.svcCtx.CatalogModel` 访问

### Q3: 数据库连接失败

**问题**: 运行时报数据库连接错误

**解决**: 检查 `api/etc/api.yaml` 中的数据库配置是否正确

### Q4: 树形结构构建错误

**问题**: 查询目录树时子节点未正确关联

**解决**:
- 检查 `level` 分组逻辑
- 确认 `parentId` 匹配逻辑
- 验证递归终止条件

### Q5: 逻辑删除后数据仍被查询

**问题**: 删除后的目录仍然出现在查询结果中

**解决**: 确保所有 SQL 查询都包含 `WHERE f_deleted = 0` 条件

---

## 下一步

完成开发和测试后，执行以下命令：

```bash
# 运行完整测试
go test -v -cover ./...

# 生成 Swagger 文档
goctl api swagger -api api/doc/api.api -o api/doc/swagger/swagger.json

# 构建镜像
make docker-build

# 部署到 K8s
make k8s-deploy
```
