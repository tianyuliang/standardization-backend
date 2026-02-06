# Quick Start: 标准文件管理 (std-file-api)

> **Feature**: std-file-api
> **Branch**: `3-std-file-api`
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
api/doc/stdfile/stdfile.api
```

**内容**：见 `plan.md` 的 API Contract 章节

### Step 2: 更新主 API 文件

在 `api/doc/api.api` 中添加导入：

```api
syntax = "v1"

import (
    "base.api"
    "catalog/catalog.api"
    "stdfile/stdfile.api"  // 新增
)
```

### Step 3: 生成 Handler 和 Types 代码

```bash
# 从项目根目录执行
# 永远使用这个命令！！！
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

**生成文件**：
- `api/internal/handler/stdfile/` - Handler 文件
- `api/internal/types/types.go` - 类型定义（追加）

### Step 4: 实现 Model 层

创建目录结构：

```bash
mkdir -p model/stdfile/stdfile
```

**文件清单**：
- `interface.go` - Model 接口定义
- `types.go` - 数据结构定义
- `vars.go` - 常量和错误定义
- `factory.go` - ORM 工厂函数
- `sql_model.go` - SQLx 实现

### Step 5: 实现 Logic 层

```bash
# Logic 文件位置（goctl 生成框架后需补充实现）
api/internal/logic/stdfile/
├── create_std_file_logic.go
├── update_std_file_logic.go
├── list_std_file_logic.go
├── get_std_file_logic.go
├── delete_std_file_logic.go
├── update_std_file_state_logic.go
├── batch_state_std_file_logic.go
├── remove_std_file_catalog_logic.go
├── download_std_file_logic.go
├── download_batch_std_file_logic.go
├── query_relation_de_logic.go
├── query_relation_dict_logic.go
├── query_relation_rule_logic.go
├── add_relation_logic.go
├── query_relations_logic.go
├── query_data_exists_logic.go
└── common.go  # 公共函数
```

### Step 6: 创建 Mock 函数

创建 Mock 目录结构：

```bash
mkdir -p api/internal/logic/stdfile/mock
```

**Mock 文件清单**：
- `catalog.go` - 目录服务 Mock
- `dataelement.go` - 数据元服务 Mock
- `dict.go` - 码表服务 Mock
- `rule.go` - 编码规则服务 Mock

### Step 7: 更新 ServiceContext

在 `api/internal/svc/servicecontext.go` 中注册 Model：

```go
type ServiceContext struct {
    Config config.Config

    // 注册 StdFile Model
    StdFileModel model.StdFileModel
    CatalogModel model.CatalogModel  // 目录服务
}
```

### Step 8: 运行和测试

```bash
# 运行服务
go run api.go -f api/etc/api.yaml

# 运行测试
go test -v ./...
```

---

## API 测试示例

### 1. 创建文件 (FILE类型)

```bash
curl -X POST http://localhost:8080/api/standardization/v1/std-file \
  -H "Content-Type: multipart/form-data" \
  -H "Authorization: Bearer <token>" \
  -F "file=@test.pdf" \
  -F "name=测试标准文件" \
  -F "catalogId=44" \
  -F "org_type=0" \
  -F "attachment_type=FILE" \
  -F "act_date=2026-01-01" \
  -F "state=enable"
```

### 2. 创建文件 (URL类型)

```bash
curl -X POST http://localhost:8080/api/standardization/v1/std-file \
  -H "Content-Type: multipart/form-data" \
  -H "Authorization: Bearer <token>" \
  -F "name=外部链接文件" \
  -F "catalogId=44" \
  -F "org_type=0" \
  -F "attachment_type=URL" \
  -F "attachment_url=https://example.com/file.pdf"
```

### 3. 查询文件列表

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/std-file?offset=1&limit=20" \
  -H "Authorization: Bearer <token>"
```

### 4. 查询文件详情

```bash
curl -X GET http://localhost:8080/api/standardization/v1/std-file/1234567890 \
  -H "Authorization: Bearer <token>"
```

### 5. 修改文件

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/std-file/1234567890 \
  -H "Content-Type: multipart/form-data" \
  -H "Authorization: Bearer <token>" \
  -F "name=修改后的文件名称" \
  -F "catalogId=44" \
  -F "org_type=0" \
  -F "attachment_type=FILE"
```

### 6. 停用文件

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/std-file/state/1234567890 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "state": "disable",
    "reason": "测试停用原因"
  }'
```

### 7. 启用文件

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/std-file/state/1234567890 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "state": "enable"
  }'
```

### 8. 删除文件

```bash
curl -X DELETE http://localhost:8080/api/standardization/v1/std-file/delete/1234567890,1234567891 \
  -H "Authorization: Bearer <token>"
```

### 9. 批量启用/停用

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/std-file/batchState \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "ids": [1234567890, 1234567891],
    "state": "disable",
    "reason": "批量停用"
  }'
```

### 10. 移动目录

```bash
curl -X POST http://localhost:8080/api/standardization/v1/std-file/catalog/remove \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "ids": [1234567890, 1234567891],
    "catalogId": 45
  }'
```

### 11. 下载文件

```bash
curl -X GET http://localhost:8080/api/standardization/v1/std-file/download/1234567890 \
  -H "Authorization: Bearer <token>" \
  --output downloaded_file.pdf
```

### 12. 批量下载

```bash
curl -X POST http://localhost:8080/api/standardization/v1/std-file/downloadBatch \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "ids": [1234567890, 1234567891]
  }' \
  --output files.zip
```

### 13. 查询关联数据元

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/std-file/relation/de/1234567890?offset=1&limit=20" \
  -H "Authorization: Bearer <token>"
```

### 14. 添加关联关系

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/std-file/relation/1234567890 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "deIds": [1, 2, 3],
    "dictIds": [1, 2, 3],
    "ruleIds": [1, 2, 3]
  }'
```

### 15. 查询关联关系

```bash
curl -X GET http://localhost:8080/api/standardization/v1/std-file/relation/1234567890 \
  -H "Authorization: Bearer <token>"
```

### 16. 检查数据是否存在

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/std-file/queryDataExists?number=DB3502/T%20035-2022" \
  -H "Authorization: Bearer <token>"
```

---

## 错误响应示例

### 参数校验失败

```json
{
  "code": "30203",
  "description": "参数值校验不通过",
  "detail": [
    {"Key": "number", "Message": "标准编号重复"}
  ],
  "solution": "请使用请求参数构造规范化的请求字符串。详细信息参见产品 API 文档。"
}
```

### 数据不存在

```json
{
  "code": "30201",
  "description": "数据不存在"
}
```

### 参数为空

```json
{
  "code": "30202",
  "description": "参数值校验不通过",
  "detail": [
    {"Key": "disable_reason", "Message": "停用必须填写停用原因"}
  ],
  "solution": "请使用请求参数构造规范化的请求字符串。详细信息参见产品 API 文档。"
}
```

### 文件下载失败

```json
{
  "code": "30205",
  "description": "[URL]类型没有文件附件"
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

**解决**: 在 `ServiceContext` 中注册 Model，并在 Logic 中使用 `l.svcCtx.XxxModel` 访问

### Q3: 数据库连接失败

**问题**: 运行时报数据库连接错误

**解决**: 检查 `api/etc/api.yaml` 中的数据库配置是否正确

### Q4: 文件上传失败

**问题**: 上传文件时提示"不支持的文件类型"

**解决**: 确保文件类型为 doc、pdf、docx、txt、ppt、pptx、xls、xlsx 中的一种

### Q5: 批量下载返回错误

**问题**: 批量下载时提示"[URL]类型没有文件附件"

**解决**: 这是因为所有选中的文件都是 URL 类型，只有 FILE 类型的文件才能下载

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
