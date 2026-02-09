# Quick Start: 数据元管理 (dataelement-api)

> **Feature**: dataelement-api
> **Branch**: `6-std-dataelement-api`
> **Date**: 2026-02-09

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

### 3. 执行 DDL 迁移

```bash
# 创建数据元表
mysql -u root -p standardization < migrations/dataelement/raw/t_data_element_info.sql

# 创建数据元-文件关系表
mysql -u root -p standardization < migrations/dataelement/raw/t_relation_de_file.sql
```

---

## 开发流程

### Step 1: API 定义文件

```bash
# API 文件位置（已生成）
api/doc/dataelement/dataelement.api
```

### Step 2: 更新主 API 文件

在 `api/doc/api.api` 中添加导入：

```api
syntax = "v1"

import (
    "base.api"
    "dataelement/dataelement.api"  # 新增
)
```

### Step 3: 生成 Handler 和 Types 代码

```bash
# 从项目根目录执行
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

**⚠️ 重要：永远使用此命令生成代码**

**生成文件**：
- `api/internal/handler/dataelement/` - Handler 文件
- `api/internal/types/types.go` - 类型定义（追加）

### Step 4: 实现 Model 层

创建目录结构：

```bash
mkdir -p model/dataelement/dataelement
mkdir -p model/dataelement/relation
```

**dataelement 模型文件清单**:
- `interface.go` - Model 接口定义
- `types.go` - 数据结构定义
- `vars.go` - 常量和错误定义
- `factory.go` - ORM 工厂函数
- `sql_model.go` - SQLx 实现

**relation 模型文件清单**:
- `interface.go` - Model 接口定义
- `types.go` - 数据结构定义
- `vars.go` - 常量和错误定义
- `factory.go` - ORM 工厂函数
- `sql_model.go` - SQLx 实现

### Step 5: 实现 Logic 层

```bash
# Logic 文件位置（goctl 生成框架后需补充实现）
api/internal/logic/dataelement/
├── accept_logic.go                      # (无此接口)
├── create_data_element_logic.go        # 创建数据元
├── delete_data_element_logic.go        # 删除数据元
├── delete_label_ids_logic.go           # 删除标签
├── export_by_ids_logic.go              # ID导出
├── export_data_element_logic.go        # 导出
├── get_data_element_detail_logic.go    # 详情
├── import_data_element_logic.go        # 批量导入
├── internal_get_detail_logic.go        # 内部详情
├── internal_list_logic.go                # 内部分页列表
├── internal_query_list_logic.go        # 内部查询列表
├── list_data_element_logic.go          # 分页列表
├── mock/                                 # Mock 服务层
│   └── service.go
├── move_catalog_logic.go               # 移动目录
├── query_by_std_file_catalog_logic.go # 按文件目录查询
├── query_by_std_file_logic.go          # 按文件查询
├── query_is_repeat_logic.go             # 检查重名
├── query_list_logic.go                  # 查询列表
├── query_std_file_logic.go             # 查询关联文件
├── common.go                             # 公共函数
├── update_data_element_logic.go        # 编辑数据元
└── update_state_logic.go               # 启用/停用
```

### Step 6: 更新 ServiceContext

在 `api/internal/svc/servicecontext.go` 中注册 Model：

```go
type ServiceContext struct {
    Config config.Config

    // 注册 DataElement Models
    DataElementModel      model.DataElementModel
    RelationDeFileModel   model.RelationDeFileModel
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

### 1. 创建数据元

```bash
curl -X POST http://localhost:8080/api/standardization/v1/dataelement \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "nameEn": "userName",
    "nameCn": "用户名",
    "synonyms": "账户,账号",
    "stdType": 4,
    "dataType": 1,
    "dataLength": 50,
    "relationType": "codeTable",
    "catalogId": 33,
    "dictId": 123,
    "description": "用户名称"
}'
```

### 2. 分页查询数据元

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/dataelement?offset=1&limit=20" \
  -H "Authorization: Bearer <token>"
```

### 3. 查询数据元详情

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/dataelement/detail?type=1&value=1234567890" \
  -H "Authorization: Bearer <token>"
```

### 4. 编辑数据元

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/dataelement/1234567890 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "nameEn": "userName",
    "nameCn": "用户名",
    "stdType": 4,
    "dataType": 1,
    "relationType": "codeTable",
    "catalogId": 33,
    "description": "用户名称"
  }'
```

### 5. 删除数据元

```bash
curl -X DELETE http://localhost:8080/api/standardization/v1/dataelement/1234567890,1234567891 \
  -H "Authorization: Bearer <token>"
```

### 6. 启用/停用数据元

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/dataelement/state/1234567890,1234567891 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "state": "disable",
    "reason": "停用原因"
  }'
```

### 7. 移动目录

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/dataelement/move_catalog/1234567890,1234567891 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "catalogId": 34
  }'
```

### 8. 检查重名

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/dataelement/query/isRepeat?id=1234567890&name=用户名&stdType=4&repeatType=DE_NAME_CN" \
  -H "Authorization: Bearer <token>"
```

### 9. 批量导入

```bash
curl -X POST http://localhost:8080/api/standardization/v1/dataelement/import \
  -H "Authorization: Bearer <token>" \
  -F "file=@dataelement.xlsx" \
  -F "catalogId=33"
```

### 10. 导出数据元

```bash
curl -X POST http://localhost:8080/api/standardization/v1/dataelement/export \
  -H "Authorization: Bearer <token>" \
  -d '{
    "catalogId": 33,
    "keyword": "用户",
    "state": "enable",
    "stdType": 4
  }'
```

---

## 错误响应示例

### 参数校验失败

```json
{
    "code": "InvalidParameter",
    "description": "参数值校验不通过",
    "detail": [
        {"Key": "catalog_id", "Message": "数据元对应的目录不存在"}
    ]
}
```

### 数据不存在

```json
{
    "code": "DataElementNotExist",
    "description": "数据不存在"
}
```

### 名称重复

```json
{
    "code": "InvalidParameter",
    "description": "参数值校验不通过",
    "detail": [
        {"Key": "name_cn", "Message": "中文名称+标准分类不能全部重复"}
    ]
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

### Q4: 值域计算错误

**问题**: 返回的值域格式不正确

**解决**: 检查 `common.go` 中的 `CalculateValueRange` 函数实现

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

---

**文档版本**: v1.0
**更新时间**: 2026-02-09
**维护人**: AI Assistant
