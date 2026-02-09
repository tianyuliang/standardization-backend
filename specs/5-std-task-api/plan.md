# 标准任务管理 (std-task-api) Technical Plan

> **Branch**: `5-std-task-api`
> **Spec Path**: `specs/5-std-task-api/`
> **Created**: 2026-02-09
> **Status**: Draft

---

## Summary

从 Java 迁移标准任务管理模块到 Go-Zero，实现24个API接口，支持标准创建任务管理、标准推荐、规则推荐、业务表管理。技术方案核心决策：

1. **主键兼容性**：使用 BIGINT 主键（与Java保持一致，这是项目宪章例外情况）
2. **表结构保持**：完全复用 Java 表结构（f_ 前缀字段）
3. **接口兼容**：24个API路径、参数、响应格式、错误码100%一致
4. **纯 SQLx 策略**：所有数据库操作使用 SQLx，手工编写 SQL 查询
5. **推荐服务集成**：调用外部推荐算法服务，使用 HTTP 客户端

---

## Constitution Check

| 约束项 | 宪章要求 | 本方案 | 状态 | 理由 |
|--------|----------|--------|------|------|
| UUID v7 主键 | 必须使用 | **使用 BIGINT** | ⚠️ 例外 | 必须与Java实现100%兼容 |
| 分层架构 | Handler→Logic→Model | 严格遵循 | ✅ 遵守 | Handler仅解析参数，Logic处理业务，Model访问数据 |
| 函数行数 | ≤50行 | 遵守 | ✅ 遵守 | 复杂逻辑拆分到私有函数 |
| 错误处理 | 必须处理所有error | 使用 idrm-go-base errorx | ✅ 遵守 | 统一错误处理和错误码 |
| 测试覆盖 | ≥80% | 单元+集成测试 | ✅ 遵守 | 核心逻辑全覆盖 |
| 通用库 | 必须使用 | errorx, response, validator等 | ✅ 遵守 | 禁止自定义实现 |

**⚠️ 例外申请**：使用 BIGINT 主键而非 UUID v7
- **原因**：Java实现使用自增BIGINT类型ID，与现有系统集成
- **影响**：仅影响 std-task 模块
- **风险**：低，Java已稳定运行

---

## Technical Context

| Item | Value |
|------|-------|
| **Language** | Go 1.24+ |
| **Framework** | Go-Zero v1.9+ |
| **Storage** | MySQL 8.0 (复用Java表结构) |
| **DB Access** | SQLx (纯 SQL) |
| **HTTP Client** | 标准库 net/http |
| **Testing** | go test |
| **Common Lib** | idrm-go-base v0.1.0+ |
| **Migration From** | Java Spring Boot + MyBatis |

---

## 通用库 (idrm-go-base)

### 模块使用

| 模块 | 用途 | 初始化 |
|------|------|--------|
| errorx | 统一错误处理和错误码 | - |
| response | 统一HTTP响应格式 | `httpx.SetErrorHandler(response.ErrorHandler)` |
| validator | 参数校验 | `validator.Init()` 在 main.go |
| telemetry | 日志和链路追踪 | `telemetry.Init(cfg)` 在 main.go |
| middleware | 认证、日志中间件 | `rest.WithMiddlewares(...)` |

### 自定义错误码

| 功能 | 范围 | 位置 |
|------|------|------|
| 标准任务管理 | 30700-30799 | `internal/errorx/task.go` |

**错误码映射**（与Java保持一致）：

| Go 错误码 | 说明 |
|-----------|------|
| 30701 | 数据不存在 |
| 30702 | 参数为空 |
| 30703 | 参数无效 |
| 30704 | 业务表不存在 |
| 30705 | 任务不存在 |

### 第三方库确认

| 库 | 原因 | 确认状态 |
|----|------|----------|
| net/http | HTTP客户端（调用推荐服务） | ✅ 已批准 |
| encoding/json | JSON序列化（Go标准库） | ✅ 已批准 |

---

## Go-Zero 开发流程

按以下顺序完成技术设计和代码生成：

| Step | 任务 | 方式 | 产出 |
|------|------|------|------|
| 1 | 定义 API 文件 | AI 实现 | `api/doc/task/task.api` |
| 2 | 生成 Handler/Types | goctl 生成 | `api/internal/handler/task/`, `types/` |
| 3 | 实现 Model 接口 | AI 手写 | `model/task/task/` |
| 4 | 实现 Logic 层 | AI 实现 | `api/internal/logic/task/` |
| 5 | 更新路由注册 | 手动 | `api/doc/api.api` |

**⚠️ 重要 goctl 命令（永远使用此指令）**:
```bash
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

---

## File Structure

### 文件产出清单

| 序号 | 文件 | 生成方式 | 位置 |
|------|------|----------|------|
| 1 | API 文件 | AI 实现 | `api/doc/task/task.api` |
| 2 | DDL 文件 | 复用Java | `migrations/task/raw/t_task_std_create.sql` |
| 3 | Model 接口 | AI 手写 | `model/task/task/interface.go` |
| 4 | Model 类型 | AI 手写 | `model/task/task/types.go` |
| 5 | Model 常量 | AI 手写 | `model/task/task/vars.go` |
| 6 | SQLx 实现 | AI 手写 | `model/task/task/sql_model.go` |
| 7 | Handler | goctl 生成 | `api/internal/handler/task/` |
| 8 | Types | goctl 生成 | `api/internal/types/types.go` |
| 9 | Logic | AI 实现 | `api/internal/logic/task/` |

### 代码结构

```
api/internal/
├── handler/task/
│   ├── get_uncompleted_tasks_handler.go
│   ├── get_completed_tasks_handler.go
│   ├── get_task_by_id_handler.go
│   ├── staging_relation_handler.go
│   ├── submit_relation_handler.go
│   ├── add_to_pending_handler.go
│   ├── get_business_table_handler.go
│   ├── get_business_table_field_handler.go
│   ├── delete_field_handler.go
│   ├── create_task_handler.go
│   ├── cancel_field_handler.go
│   ├── get_table_from_task_handler.go
│   ├── get_field_from_task_handler.go
│   ├── submit_data_element_handler.go
│   ├── finish_task_handler.go
│   ├── query_task_process_handler.go
│   ├── query_task_state_handler.go
│   ├── update_description_handler.go
│   ├── accept_handler.go
│   ├── update_table_name_handler.go
│   ├── std_rec_handler.go
│   ├── std_create_handler.go
│   ├── stand_rec_handler.go
│   ├── rule_rec_handler.go
│   └── routes.go
├── logic/task/
│   ├── get_uncompleted_tasks_logic.go
│   ├── get_completed_tasks_logic.go
│   ├── get_task_by_id_logic.go
│   ├── staging_relation_logic.go
│   ├── submit_relation_logic.go
│   ├── add_to_pending_logic.go
│   ├── get_business_table_logic.go
│   ├── get_business_table_field_logic.go
│   ├── delete_field_logic.go
│   ├── create_task_logic.go
│   ├── cancel_field_logic.go
│   ├── get_table_from_task_logic.go
│   ├── get_field_from_task_logic.go
│   ├── submit_data_element_logic.go
│   ├── finish_task_logic.go
│   ├── query_task_process_logic.go
│   ├── query_task_state_logic.go
│   ├── update_description_logic.go
│   ├── accept_logic.go
│   ├── update_table_name_logic.go
│   ├── std_rec_logic.go
│   ├── std_create_logic.go
│   ├── stand_rec_logic.go
│   ├── rule_rec_logic.go
│   └── common.go
├── types/
│   └── types.go
└── svc/
    └── servicecontext.go

model/task/task/
├── interface.go
├── types.go
├── vars.go
├── factory.go
└── sql_model.go
```

---

## Architecture Overview

遵循 IDRM 分层架构：

```
HTTP Request → Handler → Logic → Model → Database
                    ↓         ↓
                  校验参数   业务逻辑
                  格式响应   调用Model
                            调用推荐服务
```

| 层级 | 职责 | 最大行数 | 关键操作 |
|------|------|----------|----------|
| Handler | 解析参数、格式化响应 | 30 | 参数绑定、调用Logic、返回结果 |
| Logic | 业务逻辑实现 | 50 | 业务规则、事务管理、调用推荐服务 |
| Model | 数据访问 | 50 | CRUD操作、复杂查询 |

### 依赖服务调用

Logic 层通过 ServiceContext 调用依赖服务：

```go
type ServiceContext struct {
    TaskStdCreateModel        model.TaskStdCreateModel
    TaskStdCreateResultModel  model.TaskStdCreateResultModel
    BusinessTablePoolModel    model.BusinessTablePoolModel
    RecServiceURL            string  // 推荐服务URL
    RecRuleServiceURL        string  // 规则推荐服务URL
    HTTPClient               *http.Client  // HTTP客户端
}
```

---

## Interface Definitions

### TaskStdCreateModel 接口

```go
type TaskStdCreateModel interface {
    // 基础CRUD
    Insert(ctx context.Context, data *TaskStdCreate) (int64, error)
    FindOne(ctx context.Context, id int64) (*TaskStdCreate, error)
    Update(ctx context.Context, data *TaskStdCreate) error
    Delete(ctx context.Context, id int64) error

    // 查询方法
    FindUncompleted(ctx context.Context, keyword string, page, pageSize int) ([]*TaskStdCreate, int64, error)
    FindCompleted(ctx context.Context, keyword string, page, pageSize int) ([]*TaskStdCreate, int64, error)
    FindByStatus(ctx context.Context, status int32) ([]*TaskStdCreate, error)
    FindByTaskNo(ctx context.Context, taskNo string) (*TaskStdCreate, error)
}

type FindTaskOptions struct {
    Keyword  string
    Page     int
    PageSize int
}
```

### TaskStdCreateResultModel 接口

```go
type TaskStdCreateResultModel interface {
    Insert(ctx context.Context, data *TaskStdCreateResult) (int64, error)
    FindByTaskId(ctx context.Context, taskId int64) ([]*TaskStdCreateResult, error)
    DeleteByTaskId(ctx context.Context, taskId int64) error
}
```

### BusinessTablePoolModel 接口

```go
type BusinessTablePoolModel interface {
    Insert(ctx context.Context, data *BusinessTablePool) (int64, error)
    FindOne(ctx context.Context, id int64) (*BusinessTablePool, error)
    Update(ctx context.Context, data *BusinessTablePool) error
    Delete(ctx context.Context, id int64) error
    FindAll(ctx context.Context, keyword string, page, pageSize int) ([]*BusinessTablePool, int64, error)
    FindByTableName(ctx context.Context, tableName string) (*BusinessTablePool, error)
    FindFieldsByTableName(ctx context.Context, tableName string) ([]string, error)
    UpdateStatus(ctx context.Context, id int64, status int32) error
}
```

---

## Data Model

### DDL（复用Java结构）

**位置**: `migrations/task/raw/`

```sql
-- t_task_std_create 标准创建任务表
CREATE TABLE `t_task_std_create` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_task_no` varchar(64) NOT NULL COMMENT '任务编号',
  `f_table` varchar(128) DEFAULT NULL COMMENT '业务表名称',
  `f_table_description` varchar(256) DEFAULT NULL COMMENT '业务表描述',
  `f_table_field` varchar(1024) DEFAULT NULL COMMENT '表字段名称',
  `f_status` INT(2) NOT NULL DEFAULT 0 COMMENT '任务状态：0-未处理，1-处理中，2-处理完成',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_create_user` varchar(128) DEFAULT NULL COMMENT '创建用户（ID）',
  `f_create_user_phone` varchar(32) DEFAULT NULL COMMENT '创建用户联系方式',
  `f_update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `f_update_user` varchar(128) DEFAULT NULL COMMENT '修改用户（ID）',
  `f_webhook` varchar(256) DEFAULT NULL COMMENT 'AF回调地址',
  `f_deleted` bigint(20) NOT NULL DEFAULT 0 COMMENT '逻辑删除标记',
  PRIMARY KEY (`f_id`),
  KEY `idx_task_no` (`f_task_no`),
  KEY `idx_status_deleted` (`f_status`,`f_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标准创建任务表';

-- t_task_std_create_result 任务结果表
CREATE TABLE `t_task_std_create_result` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_task_id` bigint(20) NOT NULL COMMENT '标准推荐任务ID',
  `f_table_field` varchar(64) DEFAULT NULL COMMENT '表字段名称',
  `f_table_field_description` varchar(256) DEFAULT NULL COMMENT '表字段描述',
  `f_std_ref_file` varchar(256) DEFAULT NULL COMMENT '参考标准文件',
  `f_std_code` varchar(64) DEFAULT NULL COMMENT '标准编码',
  `f_rec_std_codes` varchar(512) DEFAULT NULL COMMENT '推荐算法结果标准编码',
  `f_std_ch_name` varchar(128) DEFAULT NULL COMMENT '标准中文名称',
  `f_std_en_name` varchar(256) DEFAULT NULL COMMENT '标准英文名称',
  PRIMARY KEY (`f_id`),
  KEY `idx_task_id` (`f_task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标准创建任务结果表';

-- t_business_table_std_create_pool 待新建标准表
CREATE TABLE `t_business_table_std_create_pool` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_table_name` varchar(128) NOT NULL COMMENT '业务表名称',
  `f_table_description` varchar(256) DEFAULT NULL COMMENT '业务表描述',
  `f_table_field` varchar(2048) DEFAULT NULL COMMENT '表字段',
  `f_status` INT(2) NOT NULL DEFAULT 0 COMMENT '状态：0-待新建，1-处理中，2-已完成',
  `f_create_user` varchar(128) DEFAULT NULL COMMENT '创建用户',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_deleted` bigint(20) NOT NULL DEFAULT 0 COMMENT '逻辑删除标记',
  PRIMARY KEY (`f_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='待新建标准表';
```

### Go Struct

```go
// TaskStdCreate 标准创建任务
type TaskStdCreate struct {
    Id                int64  `db:"f_id" json:"id"`
    TaskNo            string `db:"f_task_no" json:"taskNo"`
    Table             string `db:"f_table" json:"table"`
    TableDescription  string `db:"f_table_description" json:"tableDescription"`
    TableField        string `db:"f_table_field" json:"tableField"`
    Status            int32  `db:"f_status" json:"status"`
    CreateTime        string `db:"f_create_time" json:"createTime"`
    CreateUser        string `db:"f_create_user" json:"createUser"`
    CreateUserPhone   string `db:"f_create_user_phone" json:"createUserPhone"`
    UpdateTime        string `db:"f_update_time" json:"updateTime"`
    UpdateUser        string `db:"f_update_user" json:"updateUser"`
    Webhook           string `db:"f_webhook" json:"webhook"`
    Deleted           int64  `db:"f_deleted" json:"deleted"`
}

// TaskStdCreateResult 任务结果
type TaskStdCreateResult struct {
    Id                   int64  `db:"f_id" json:"id"`
    TaskId               int64  `db:"f_task_id" json:"taskId"`
    TableField           string `db:"f_table_field" json:"tableField"`
    TableFieldDescription string `db:"f_table_field_description" json:"tableFieldDescription"`
    StdRefFile           string `db:"f_std_ref_file" json:"stdRefFile"`
    StdCode              string `db:"f_std_code" json:"stdCode"`
    RecStdCodes          string `db:"f_rec_std_codes" json:"recStdCodes"`
    StdChName            string `db:"f_std_ch_name" json:"stdChName"`
    StdEnName            string `db:"f_std_en_name" json:"stdEnName"`
}

// BusinessTablePool 待新建标准表
type BusinessTablePool struct {
    Id               int64  `db:"f_id" json:"id"`
    TableName        string `db:"f_table_name" json:"tableName"`
    TableDescription string `db:"f_table_description" json:"tableDescription"`
    TableField       string `db:"f_table_field" json:"tableField"`
    Status           int32  `db:"f_status" json:"status"`
    CreateUser       string `db:"f_create_user" json:"createUser"`
    CreateTime       string `db:"f_create_time" json:"createTime"`
    Deleted          int64  `db:"f_deleted" json:"deleted"`
}
```

### 枚举定义

```go
// TaskStatusEnum 任务状态
const (
    TaskStatusUnhandled  int32 = 0 // 未处理
    TaskStatusProcessing int32 = 1 // 处理中
    TaskStatusCompleted  int32 = 2 // 处理完成
)

// PoolStatusEnum 业务表池状态
const (
    PoolStatusPending    int32 = 0 // 待新建
    PoolStatusProcessing int32 = 1 // 处理中
    PoolStatusCompleted  int32 = 2 // 已完成
)

// 字符串状态
const (
    TaskStateUnhandled  = "unhandled"
    TaskStateProcessing = "processing"
    TaskStateCompleted  = "completed"
)
```

---

## API Contract

**位置**: `api/doc/task/task.api`

**⚠️ 重要**：基础路径必须为 `/api/standardization/v1/dataelement/task`

### 完整 API 定义（部分核心接口）

```api
syntax = "v1"

import "../base.api"

type (
    // ==================== 核心请求类型 ====================

    // StagingRelationReq 标准关联暂存请求
    StagingRelationReq {
        TaskId        int64       `json:"taskId,optional" validate:"required"`
        BusinessTable string      `json:"businessTable,optional" validate:"required"`
        Fields        []FieldInfo `json:"fields,optional"`
    }

    // FieldInfo 字段信息
    FieldInfo {
        FieldName        string `json:"fieldName,optional" validate:"required"`
        FieldDescription string `json:"fieldDescription,optional"`
        DataType         string `json:"dataType,optional"`
    }

    // SubmitDataElementReq 提交数据元请求
    SubmitDataElementReq {
        TaskId   int64            `json:"taskId,optional" validate:"required"`
        Elements []DataElementInfo `json:"elements,optional"`
    }

    // DataElementInfo 数据元信息
    DataElementInfo {
        FieldId   int64  `json:"fieldId,optional"`
        StdCode   string `json:"stdCode,optional"`
        StdChName string `json:"stdChName,optional"`
        StdEnName string `json:"stdEnName,optional"`
    }

    // StdRecReq 标准推荐请求
    StdRecReq {
        TaskId               int64  `json:"taskId,optional" validate:"required"`
        BusinessTable        string `json:"businessTable,optional" validate:"required"`
        TableField           string `json:"tableField,optional" validate:"required"`
        TableFieldDescription string `json:"tableFieldDescription,optional"`
    }

    // RuleRecReq 规则推荐请求
    RuleRecReq {
        BusinessTable        string `json:"businessTable,optional" validate:"required"`
        TableField           string `json:"tableField,optional" validate:"required"`
        TableFieldDescription string `json:"tableFieldDescription,optional"`
        DataType             string `json:"dataType,optional"`
    }

    // CreateTaskReq 创建任务请求
    CreateTaskReq {
        BusinessTableIds []int64 `json:"businessTableIds,optional" validate:"required"`
        Webhook          string  `json:"webhook,optional"`
    }

    // UpdateDescriptionReq 修改字段说明请求
    UpdateDescriptionReq {
        FieldId              int64  `json:"fieldId,optional" validate:"required"`
        FieldDescription     string `json:"fieldDescription,optional"`
    }

    // UpdateTableNameReq 修改表名请求
    UpdateTableNameReq {
        TableId          int64  `json:"tableId,optional" validate:"required"`
        TableName        string `json:"tableName,optional" validate:"required"`
    }

    // AddToPendingReq 添加到待新建请求
    AddToPendingReq {
        BusinessTable          string `json:"businessTable,optional" validate:"required"`
        BusinessTableDescription string `json:"businessTableDescription,optional"`
        TableField             string `json:"tableField,optional"`
        CreateUserPhone        string `json:"createUserPhone,optional"`
    }

    // AcceptReq 采纳请求
    AcceptReq {
        TaskId int64 `json:"taskId,optional" validate:"required"`
    }

    // ==================== 响应类型 ====================

    // TaskResp 任务响应
    TaskResp {
        Id                int64  `json:"id"`
        TaskNo            string `json:"taskNo"`
        Table             string `json:"table"`
        TableDescription  string `json:"tableDescription"`
        TableField        string `json:"tableField"`
        Status            int32  `json:"status"`
        CreateTime        string `json:"createTime"`
        CreateUser        string `json:"createUser"`
        CreateUserPhone   string `json:"createUserPhone"`
        Webhook           string `json:"webhook"`
    }

    // TaskDataListResp 任务列表响应
    TaskDataListResp {
        TotalCount int64     `json:"totalCount"`
        Data       []TaskResp `json:"data"`
    }

    // TaskDetailResp 任务详情响应
    TaskDetailResp {
        Task    TaskResp            `json:"task"`
        Results []TaskStdCreateResult `json:"results"`
    }

    // BusinessTableResp 业务表响应
    BusinessTableResp {
        Id               int64  `json:"id"`
        TableName        string `json:"tableName"`
        TableDescription string `json:"tableDescription"`
        TableField       string `json:"tableField"`
        Status           int32  `json:"status"`
        CreateUser       string `json:"createUser"`
        CreateTime       string `json:"createTime"`
    }

    // BusinessTableListResp 业务表列表响应
    BusinessTableListResp {
        TotalCount int64              `json:"totalCount"`
        Data       []BusinessTableResp `json:"data"`
    }

    // FieldResp 字段响应
    FieldResp {
        Id               int64  `json:"id"`
        FieldName        string `json:"fieldName"`
        FieldDescription string `json:"fieldDescription"`
        DataType         string `json:"dataType"`
    }

    // FieldListResp 字段列表响应
    FieldListResp {
        TotalCount int64      `json:"totalCount"`
        Data       []FieldResp `json:"data"`
    }

    // StdRecResp 标准推荐响应
    StdRecResp {
        Code        string          `json:"code"`
        Description string          `json:"description"`
        Data        []StdRecItem    `json:"data"`
    }

    // StdRecItem 标准推荐项
    StdRecItem {
        StdCode   string  `json:"stdCode"`
        StdName   string  `json:"stdName"`
        MatchRate float64 `json:"matchRate"`
    }

    // TaskBaseResp 任务基础响应
    TaskBaseResp {
        Code        string `json:"code"`
        Description string `json:"description"`
    }

    // ProcessResp 进度响应
    ProcessResp {
        TaskId  int64  `json:"taskId"`
        Status  string `json:"status"`
        Message string `json:"message"`
    }
)

// ==================== 外部API ====================

@server(
    prefix: /api/standardization/v1/dataelement/task
    group: task
    middleware: TokenCheck
)
service standardization-api {
    @doc "未处理任务列表"
    @handler GetUncompletedTasks
    get /std-create/uncompleted returns (TaskDataListResp)

    @doc "已完成任务列表"
    @handler GetCompletedTasks
    get /std-create/completed returns (TaskDataListResp)

    @doc "任务详情"
    @handler GetTaskById
    get /std-create/completed/:id returns (TaskDetailResp)

    @doc "标准关联暂存"
    @handler StagingRelation
    post /std-create/relation/staging (StagingRelationReq) returns (TaskBaseResp)

    @doc "标准关联提交"
    @handler SubmitRelation
    post /std-create/publish/submit (StagingRelationReq) returns (TaskBaseResp)

    @doc "添加至待新建"
    @handler AddToPending
    post /addToPending (AddToPendingReq) returns (TaskBaseResp)

    @doc "业务表列表"
    @handler GetBusinessTable
    get /getBusinessTable returns (BusinessTableListResp)

    @doc "业务表字段列表"
    @handler GetBusinessTableField
    get /getBusinessTableField returns (FieldListResp)

    @doc "移除字段"
    @handler DeleteField
    delete /deleteBusinessTableField/:id returns (TaskBaseResp)

    @doc "新建标准任务"
    @handler CreateTask
    post /createTask (CreateTaskReq) returns (TaskBaseResp)

    @doc "撤销"
    @handler CancelField
    put /cancelBusinessTableField returns (TaskBaseResp)

    @doc "任务关联业务表"
    @handler GetTableFromTask
    get /getBusinessTableFromTask returns (BusinessTableListResp)

    @doc "任务关联字段"
    @handler GetFieldFromTask
    get /getBusinessTableFieldFromTask returns (FieldListResp)

    @doc "提交选定数据元"
    @handler SubmitDataElement
    post /submitDataElement (SubmitDataElementReq) returns (TaskBaseResp)

    @doc "完成任务"
    @handler FinishTask
    post /finishTask/:task_id returns (TaskBaseResp)

    @doc "进度查询"
    @handler QueryTaskProcess
    post /queryTaskProcess returns (ProcessResp)

    @doc "任务状态查询"
    @handler QueryTaskState
    post /queryTaskState returns (ProcessResp)

    @doc "修改字段说明"
    @handler UpdateDescription
    put /updateDescription (UpdateDescriptionReq) returns (TaskBaseResp)

    @doc "采纳"
    @handler Accept
    put /accept (AcceptReq) returns (TaskBaseResp)

    @doc "修改表名称"
    @handler UpdateTableName
    put /updateTableName (UpdateTableNameReq) returns (TaskBaseResp)
}

// ==================== 内部API ====================

@server(
    prefix: /api/standardization/v1/dataelement/task
    group: task_internal
)
service standardization-api {
    @doc "标准推荐（内部）"
    @handler StdRec
    post /std-rec/rec (StdRecReq) returns (StdRecResp)

    @doc "标准创建（内部）"
    @handler StdCreate
    post /std-create (StdRecReq) returns (TaskBaseResp)
}

@server(
    prefix: /api/standardization/v1/dataelement/task
    group: task_stand
)
service standardization-api {
    @doc "标准推荐（弹框）"
    @handler StandRec
    post /stand-rec/rec (StdRecReq) returns (StdRecResp)
}

@server(
    prefix: /api/standardization/v1/dataelement/task
    group: task_rule
)
service standardization-api {
    @doc "编码规则推荐"
    @handler RuleRec
    post /rule-rec/rec (RuleRecReq) returns (StdRecResp)
}
```

---

## Testing Strategy

| 类型 | 方法 | 覆盖率 | 重点 |
|------|------|--------|------|
| 单元测试 | 表驱动测试，Mock Model | > 80% | Logic层业务逻辑 |
| 集成测试 | 测试数据库 | 核心流程 | 端到端API测试 |

### 测试用例优先级

| 优先级 | 测试场景 |
|--------|----------|
| P1 | 创建任务、查询任务列表、完成任务 |
| P1 | 标准关联暂存、标准关联提交 |
| P1 | 添加业务表到待新建、查询业务表 |
| P2 | 标准推荐、规则推荐 |
| P2 | 提交数据元、修改字段说明 |

---

## Key Implementation Details

### 1. 推荐服务调用

```go
// CallStdRecService 调用标准推荐服务
func CallStdRecService(client *http.Client, serviceURL string, req *StdRecReq) (*StdRecResp, error) {
    jsonData, _ := json.Marshal(req)
    httpReq, _ := http.NewRequest("POST", serviceURL, bytes.NewBuffer(jsonData))
    httpReq.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(httpReq)
    if err != nil {
        logx.Errorf("调用推荐服务失败: %v", err)
        return &StdRecResp{Data: []StdRecItem{}}, nil // 返回空结果
    }
    defer resp.Body.Close()

    var result StdRecResp
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}
```

### 2. 任务编号生成

```go
// GenerateTaskNo 生成任务编号
func GenerateTaskNo() string {
    return fmt.Sprintf("TASK%s%04d", time.Now().Format("20060102"), rand.Intn(10000))
}
```

### 3. 完成任务回调

```go
// SendTaskCallback 发送任务完成回调
func SendTaskCallback(client *http.Client, webhook string, task *TaskStdCreate, results []*TaskStdCreateResult) error {
    if webhook == "" {
        return nil // Webhook为空时不发送
    }

    payload := map[string]interface{}{
        "taskId": task.Id,
        "status": "completed",
        "results": results,
    }

    jsonData, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", webhook, bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        logx.Errorf("发送任务回调失败: webhook=%s, error=%v", webhook, err)
        return err // 记录日志但不阻塞主流程
    }
    defer resp.Body.Close()
    return nil
}
```

### 4. 状态转换

```go
// StatusToInt 状态字符串转整数
func StatusToInt(state string) int32 {
    switch state {
    case "unhandled":
        return TaskStatusUnhandled
    case "processing":
        return TaskStatusProcessing
    case "completed":
        return TaskStatusCompleted
    default:
        return TaskStatusUnhandled
    }
}

// IntToStatus 整数转状态字符串
func IntToStatus(status int32) string {
    switch status {
    case TaskStatusUnhandled:
        return "unhandled"
    case TaskStatusProcessing:
        return "processing"
    case TaskStatusCompleted:
        return "completed"
    default:
        return "unhandled"
    }
}
```

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-02-09 | - | 初始版本 |
