# Quick Start: 标准任务管理 (task-api)

> **Feature**: task-api
> **Branch**: `5-std-task-api`
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
# 创建业务表标准创建池表
mysql -u root -p standardization < migrations/task/raw/t_business_table_std_create_pool.sql

# 创建标准创建任务表
mysql -u root -p standardization < migrations/task/raw/t_task_std_create.sql

# 创建任务结果表
mysql -u root -p standardization < migrations/task/raw/t_task_std_create_result.sql
```

---

## 开发流程

### Step 1: API 定义文件

```bash
# API 文件位置（已生成）
api/doc/task/task.api
```

### Step 2: 更新主 API 文件

在 `api/doc/api.api` 中添加导入：

```api
syntax = "v1"

import (
    "base.api"
    "task/task.api"  // 新增
)
```

### Step 3: 生成 Handler 和 Types 代码

```bash
# 从项目根目录执行
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

**生成文件**：
- `api/internal/handler/task/` - Handler 文件
- `api/internal/types/types.go` - 类型定义（追加）

### Step 4: 实现 Model 层

创建目录结构：

```bash
mkdir -p model/task/pool
mkdir -p model/task/task
```

**pool 模型文件清单**:
- `interface.go` - Model 接口定义
- `types.go` - 数据结构定义
- `vars.go` - 常量和错误定义
- `factory.go` - ORM 工厂函数
- `sql_model.go` - SQLx 实现

**task 模型文件清单**:
- `interface.go` - Model 接口定义
- `types.go` - 数据结构定义
- `vars.go` - 常量和错误定义
- `factory.go` - ORM 工厂函数
- `sql_model.go` - SQLx 实现

### Step 5: 实现 Logic 层

```bash
# Logic 文件位置（goctl 生成框架后需补充实现）
api/internal/logic/task/
├── accept_logic.go                      # 采纳
├── add_to_pending_logic.go              # 添加至待新建
├── cancel_field_logic.go                 # 撤销
├── create_task_logic.go                  # 创建任务
├── delete_field_logic.go                 # 移除字段
├── finish_task_logic.go                   # 完成任务
├── get_business_table_field_logic.go     # 业务表字段列表
├── get_business_table_logic.go           # 业务表列表
├── get_completed_tasks_logic.go          # 已完成任务列表
├── get_field_from_task_logic.go          # 任务关联字段
├── get_table_from_task_logic.go          # 任务关联业务表
├── get_task_by_id_logic.go               # 任务详情
├── get_uncompleted_tasks_logic.go        # 未处理任务列表
├── mock/                                 # Mock 服务层
│   └── service.go
├── query_task_process_logic.go           # 进度查询
├── query_task_state_logic.go             # 任务状态查询
├── rule_rec_logic.go                      # 编码规则推荐
├── std_create_logic.go                    # 标准关联暂存
├── std_rec_logic.go                       # 标准推荐
├── stand_rec_logic.go                     # 标准采纳推荐
├── submit_data_element_logic.go          # 提交选定数据元
├── submit_relation_logic.go               # 标准关联提交
├── update_description_logic.go            # 修改字段说明
├�── update_table_name_logic.go            # 修改表名称
└── common.go                              # 公共函数
```

### Step 6: 更新 ServiceContext

在 `api/internal/svc/servicecontext.go` 中注册 Model：

```go
type ServiceContext struct {
    Config config.Config

    // 注册 Task Models
    TaskStdCreateModel         model.TaskStdCreateModel
    TaskStdCreateResultModel   model.TaskStdCreateResultModel
    BusinessTablePoolModel     model.BusinessTablePoolModel
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

### 1. 添加至待新建标准池

```bash
curl -X POST http://localhost:8080/api/standardization/v1/dataelement/task/addToPending \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "businessTable": "T_USER",
    "businessTableDescription": "用户表",
    "tableField": "user_name",
    "createUserPhone": "13800138000"
  }'
```

### 2. 创建标准任务

```bash
curl -X POST http://localhost:8080/api/standardization/v1/dataelement/task/createTask \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "taskId": "550e8400-e29b-41d4-a716-446655440000",
    "ids": [
        "550e8400-e29b-41d4-a716-446655440001",
      "550e8400-e29b-41d4-a716-446655440002"
    ],
    "webhook": "http://example.com/webhook"
  }'
```

### 3. 查询未处理任务列表

```bash
curl -X GET "http://localhost:8080/api/standardization/v1/dataelement/task/std-create/uncompleted?offset=1&limit=20" \
  -H "Authorization: Bearer <token>"
```

### 4. 查询任务进度

```bash
curl -X POST http://localhost:8080/api/standardization/v1/dataelement/task/queryTaskProcess \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "taskId": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

### 5. 提交选定数据元

```bash
curl -X POST http://localhost:8080/api/standardization/v1/dataelement/task/submitDataElement \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "dataElementId": "1234567890"
  }'
```

### 6. 完成任务

```bash
curl -X POST http://localhost:8080/api/standardization/v1/dataelement/task/finishTask/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <token>"
```

### 7. 采纳标准

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/dataelement/task/acceptBusinessTableField \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "ids": [
      "12345678901234567890",
      "12345678901234567891"
    ]
  }'
```

### 8. 撤销

```bash
curl -X PUT http://localhost:8080/api/standardization/v1/dataelement/task/cancelBusinessTableField \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "ids": [
      "12345678901234567890"
    ]
  }'
```

---

## 错误响应示例

### 参数校验失败

```json
{
    "code": "30703",
    "description": "参数值校验不通过",
    "detail": [
        {"Key": "taskId", "Message": "taskId必须为36位UUID"}
    ],
    "solution": "请使用请求参数构造规范化的请求字符串。详细信息参见产品 API 文档。"
}
```

### 数据不存在

```json
{
    "code": "30701",
    "description": "数据不存在"
}
```

### 部分失败

```json
{
    "code": "30704",
    "description": "不能存在没有关联数据元的标准字段"
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

### Q4: TaskId 格式错误

**问题**: 创建任务时提示"taskId必须为36位UUID"

**解决**: 确保 taskId 是有效的 UUID v7 格式（36字符）

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
