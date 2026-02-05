# Quick Start: 编码规则管理 (rule-api)

> **Feature**: rule-api
> **Branch**: `1-rule-api`
> **Date**: 2026-02-05

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
api/doc/rule/rule.api
```

**内容**：见 `plan.md` 的 API Contract 章节

### Step 2: 更新主 API 文件

在 `api/doc/api.go` 中添加导入：

```go
import (
    "handler/base"
    "handler/rule"
    // ... 其他 imports
)
```

在 `api/doc/api.api` 中添加导入：

```api
syntax = "v1"

import (
    "base.api"
    "rule/rule.api"  // 新增
)

// ... 其他内容
```

### Step 3: 生成 Handler 和 Types 代码

```bash
# 从项目根目录执行
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

**生成文件**：
- `api/internal/handler/rule/` - Handler 文件
- `api/internal/types/types.go` - 类型定义（追加）

### Step 4: 实现 Model 层

创建目录结构：

```bash
mkdir -p model/rule/rule
```

**文件清单**：
- `interface.go` - Model 接口定义
- `types.go` - 数据结构定义
- `vars.go` - 常量和错误定义
- `factory.go` - ORM 工厂函数
- `gorm_dao.go` - GORM 实现
- `sqlx_model.go` - SQLx 实现

### Step 5: 实现 Logic 层

```bash
# Logic 文件位置（goctl 生成框架后需补充实现）
api/internal/logic/rule/
├── create_rule_logic.go
├── update_rule_logic.go
├── get_rule_logic.go
├── get_rule_id_logic.go
├── get_rule_detail_by_data_id_logic.go
├── get_rule_detail_by_data_code_logic.go
├── list_rule_logic.go
├── delete_rule_logic.go
├── update_rule_state_logic.go
├── remove_rule_catalog_logic.go
├── query_rule_used_data_element_logic.go
├── query_rule_by_ids_logic.go
├── query_internal_rule_by_ids_logic.go
├── query_rule_by_std_file_catalog_logic.go
├── query_rule_by_std_file_logic.go
├── query_data_exists_logic.go
├── query_std_files_by_rule_logic.go
├── get_custom_date_format_logic.go
└── common.go  # 公共函数
```

### Step 6: 更新 ServiceContext

在 `api/internal/svc/servicecontext.go` 中注册 Model：

```go
type ServiceContext struct {
    Config config.Config

    // 注册 Rule Model
    RuleModel model.RuleModel
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

### 1. 创建规则 (REGEX)

```bash
curl -X POST http://localhost:8080/api/standardization/v1/rule \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "测试规则",
    "orgType": 0,
    "ruleType": "REGEX",
    "regex": "[A-Z]{4}\\d{6}",
    "description": "测试规则说明",
    "catalogId": 33,
    "state": "enable"
  }'
```

### 2. 创建规则 (CUSTOM)

```bash
curl -X POST http://localhost:8080/api/standardization/v1/rule \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "自定义规则",
    "orgType": 0,
    "ruleType": "CUSTOM",
    "custom": [
      {
        "segment_length": 4,
        "type": "number",
        "value": ""
      },
      {
        "segment_length": 6,
        "type": "date",
        "value": "yyyyMMdd"
      }
    ],
    "catalogId": 33
  }'
```

### 3. 查询规则列表

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/rule?offset=1&limit=20" \
  -H "Authorization: Bearer <token>"
```

### 4. 查询规则详情

```bash
curl -X GET http://localhost:8080/api/standardization/v1/rule/1234567890 \
  -H "Authorization: Bearer <token>"
```

### 5. 修改规则

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/rule/1234567890 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "修改后的规则名称",
    "orgType": 0,
    "ruleType": "REGEX",
    "regex": "[A-Z]{4}\\d{8}",
    "description": "修改后的规则说明"
  }'
```

### 6. 停用规则

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/rule/state/1234567890 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "state": "disable",
    "reason": "测试停用原因"
  }'
```

### 7. 启用规则

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/rule/state/1234567890 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "state": "enable"
  }'
```

### 8. 删除规则

```bash
curl -X DELETE http://localhost:8080/api/standardization/v1/rule/1234567890,1234567891 \
  -H "Authorization: Bearer <token>"
```

### 9. 批量查询规则

```bash
curl -X POST http://localhost:8080/api/standardization/v1/rule/queryByIds \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "ids": [1234567890, 1234567891, 1234567892]
  }'
```

### 10. 获取自定义日期格式

```bash
curl -X GET http://localhost:8080/api/standardization/v1/rule/getCustomDateFormat
```

---

## 错误响应示例

### 参数校验失败

```json
{
  "code": "30303",
  "description": "参数值校验不通过",
  "detail": [
    {"Key": "name", "Message": "规则名称已存在"}
  ],
  "solution": "请使用请求参数构造规范化的请求字符串。详细信息参见产品 API 文档。"
}
```

### 数据不存在

```json
{
  "code": "30301",
  "description": "数据不存在"
}
```

### 参数为空

```json
{
  "code": "30302",
  "description": "参数值校验不通过",
  "detail": [
    {"Key": "disable_reason", "Message": "停用必须填写停用原因"}
  ],
  "solution": "请使用请求参数构造规范化的请求字符串。详细信息参见产品 API 文档。"
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

**解决**: 在 `ServiceContext` 中注册 Model，并在 Logic 中使用 `l.svcCtx.XxxModel` 访问

### Q3: 数据库连接失败

**问题**: 运行时报数据库连接错误

**解决**: 检查 `api/etc/api.yaml` 中的数据库配置是否正确

### Q4: 正则表达式校验失败

**问题**: 创建 REGEX 规则时提示"正则表达式非法"

**解决**: 确保正则表达式语法正确，可以使用在线工具验证：https://regex101.com/

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
