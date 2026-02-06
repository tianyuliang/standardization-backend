# 标准文件管理 (std-file-api) Technical Plan

> **Branch**: `3-std-file-api`
> **Spec Path**: `specs/3-std-file-api/`
> **Created**: 2026-02-06
> **Status**: Draft

---

## Summary

从 Java 迁移标准文件管理模块到 Go-Zero，实现16个API接口，支持FILE和URL两种附件类型。技术方案核心决策：

1. **主键兼容性**：使用 Long 类型 BIGINT 主键（与Java保持一致）
2. **表结构保持**：完全复用 Java 表结构（f_ 前缀字段）
3. **接口兼容**：16个API路径、参数、响应格式、错误码100%一致
4. **纯 SQLx 策略**：所有数据库操作使用 SQLx，手工编写 SQL 查询

---

## Constitution Check

| 约束项 | 宪章要求 | 本方案 | 状态 | 理由 |
|--------|----------|--------|------|------|
| UUID v7 主键 | 必须使用 | **使用 BIGINT** | ⚠️ 例外 | 必须与Java实现100%兼容，前端已依赖Long类型ID |
| 分层架构 | Handler→Logic→Model | 严格遵循 | ✅ 遵守 | Handler仅解析参数，Logic处理业务，Model访问数据 |
| 函数行数 | ≤50行 | 遵守 | ✅ 遵守 | 复杂逻辑拆分到私有函数 |
| 错误处理 | 必须处理所有error | 使用 idrm-go-base errorx | ✅ 遵守 | 统一错误处理和错误码 |
| 测试覆盖 | ≥80% | 单元+集成测试 | ✅ 遵守 | 核心逻辑全覆盖 |
| 通用库 | 必须使用 | errorx, response, validator等 | ✅ 遵守 | 禁止自定义实现 |

**⚠️ 例外申请**：使用 BIGINT 主键而非 UUID v7
- **原因**：Java实现使用自增Long类型ID，前端已依赖此格式
- **影响**：仅影响 std-file 模块，其他模块仍遵循 UUID v7 规范
- **风险**：低，Java已稳定运行多年

---

## Technical Context

| Item | Value |
|------|-------|
| **Language** | Go 1.24+ |
| **Framework** | Go-Zero v1.9+ |
| **Storage** | MySQL 8.0 (复用Java表结构) |
| **File Storage** | OSS (阿里云对象存储) |
| **DB Access** | SQLx (纯 SQL) |
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
| 标准文件 | 30200-30299 | `internal/errorx/codes.go` |

**错误码映射**（与Java保持一致）：

| Java 错误码 | Go 错误码 | 说明 |
|-------------|-----------|------|
| Standardization.DATA_NOT_EXIST | 30201 | 数据不存在 |
| Standardization.PARAMETER_EMPTY | 30202 | 参数为空 |
| Standardization.InvalidParameter | 30203 | 参数无效 |
| Standardization.DATA_EXIST | 30204 | 数据已存在 |
| Standardization.FileDownloadFailed | 30205 | 文件下载失败 |

### 第三方库确认

| 库 | 原因 | 确认状态 |
|----|------|----------|
| github.com/google/uuid | UUID v7 生成（部门ID等） | ✅ 已批准 |

---

## Go-Zero 开发流程

按以下顺序完成技术设计和代码生成：

| Step | 任务 | 方式 | 产出 |
|------|------|------|------|
| 1 | 定义 API 文件 | AI 实现 | `api/doc/stdfile/stdfile.api` |
| 2 | 生成 Handler/Types | goctl 生成 | `api/internal/handler/stdfile/`, `types/` |
| 3 | 实现 Model 接口 | AI 手写 | `model/stdfile/stdfile/` |
| 4 | 实现 Logic 层 | AI 实现 | `api/internal/logic/stdfile/` |
| 5 | 更新路由注册 | 手动 | `api/doc/api.api` |

**goctl 命令**:
```bash
# 步骤1：在 api/doc/api.api 中 import stdfile 模块
# 步骤2：执行 goctl 生成代码（针对整个项目）
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

---

## File Structure

### 文件产出清单

| 序号 | 文件 | 生成方式 | 位置 |
|------|------|----------|------|
| 1 | API 文件 | AI 实现 | `api/doc/stdfile/stdfile.api` |
| 2 | DDL 文件 | 复用Java | `migrations/stdfile/raw/t_std_file.sql` |
| 3 | Model 接口 | AI 手写 | `model/stdfile/stdfile/interface.go` |
| 4 | Model 类型 | AI 手写 | `model/stdfile/stdfile/types.go` |
| 5 | Model 常量 | AI 手写 | `model/stdfile/stdfile/vars.go` |
| 6 | SQLx 实现 | AI 手写 | `model/stdfile/stdfile/sql_model.go` |
| 7 | Handler | goctl 生成 | `api/internal/handler/stdfile/` |
| 8 | Types | goctl 生成 | `api/internal/types/types.go` |
| 9 | Logic | AI 实现 | `api/internal/logic/stdfile/` |

### 代码结构

```
api/internal/
├── handler/stdfile/
│   ├── create_std_file_handler.go           # goctl 生成
│   ├── update_std_file_handler.go
│   ├── list_std_file_handler.go
│   ├── get_std_file_handler.go
│   ├── delete_std_file_handler.go
│   ├── update_std_file_state_handler.go
│   ├── remove_std_file_catalog_handler.go
│   ├── download_std_file_handler.go
│   ├── download_batch_std_file_handler.go
│   ├── query_relation_de_handler.go
│   ├── query_relation_dict_handler.go
│   ├── query_relation_rule_handler.go
│   ├── add_relation_handler.go
│   ├── query_relations_handler.go
│   ├── query_data_exists_handler.go
│   ├── batch_state_std_file_handler.go
│   └── routes.go
├── logic/stdfile/
│   ├── create_std_file_logic.go             # AI 实现
│   ├── update_std_file_logic.go
│   ├── list_std_file_logic.go
│   ├── get_std_file_logic.go
│   ├── delete_std_file_logic.go
│   ├── update_std_file_state_logic.go
│   ├── remove_std_file_catalog_logic.go
│   ├── download_std_file_logic.go
│   ├── download_batch_std_file_logic.go
│   ├── query_relation_de_logic.go
│   ├── query_relation_dict_logic.go
│   ├── query_relation_rule_logic.go
│   ├── add_relation_logic.go
│   ├── query_relations_logic.go
│   ├── query_data_exists_logic.go
│   ├── batch_state_std_file_logic.go
│   └── common.go                            # 公共函数
├── types/
│   └── types.go                             # goctl 生成（含所有请求/响应类型）
└── svc/
    └── servicecontext.go                    # 手动维护（注入Model、OSS等）

model/stdfile/stdfile/
├── interface.go                             # Model 接口定义
├── types.go                                 # 数据结构
├── vars.go                                  # 常量/错误定义
└── sql_model.go                             # SQLx 实现（所有查询）
```

---

## Architecture Overview

遵循 IDRM 分层架构：

```
HTTP Request → Handler → Logic → Model → Database
                    ↓         ↓
                  校验参数   业务逻辑
                  格式响应   调用Model
                            调用依赖服务
                            上传/下载文件
```

| 层级 | 职责 | 最大行数 | 关键操作 |
|------|------|----------|----------|
| Handler | 解析参数、格式化响应 | 30 | 参数绑定、调用Logic、返回结果 |
| Logic | 业务逻辑实现 | 50 | 业务规则、事务管理、调用依赖服务、文件操作 |
| Model | 数据访问 | 50 | CRUD操作、复杂查询、事务管理 |

### 依赖服务调用

Logic 层通过 ServiceContext 调用依赖服务：

```go
type ServiceContext struct {
    StdFileModel    model.StdFileModel
    CatalogModel    model.CatalogModel    // 目录服务
    DataElementModel model.DataElementModel // 数据元服务
    DictModel       model.DictModel       // 字典服务
    RuleModel       model.RuleModel       // 编码规则服务
    OSSClient       *oss.Client           // OSS客户端
}
```

---

## Interface Definitions

### StdFileModel 接口

```go
type StdFileModel interface {
    // 基础CRUD
    Insert(ctx context.Context, data *StdFile) (int64, error)
    FindOne(ctx context.Context, id int64) (*StdFile, error)
    Update(ctx context.Context, data *StdFile) error
    Delete(ctx context.Context, id int64) error

    // 查询方法
    FindByIds(ctx context.Context, ids []int64) ([]*StdFile, error)
    FindByNumber(ctx context.Context, number string) ([]*StdFile, error)
    FindByNameAndOrgType(ctx context.Context, name string, orgType int32) ([]*StdFile, error)
    FindByCatalogIds(ctx context.Context, opts *FindOptions) ([]*StdFile, int64, error)
    FindDataExists(ctx context.Context, filterId int64, number string, orgType int32, name string, deptIds string) (*StdFile, error)

    // 更新方法
    UpdateState(ctx context.Context, id int64, state int32, disableReason string) error
    RemoveCatalog(ctx context.Context, ids []int64, catalogId int64, updateUser string) error
    BatchUpdateState(ctx context.Context, ids []int64, state int32, disableReason string) error

    // 批量操作
    DeleteByIds(ctx context.Context, ids []int64) error
}

type FindOptions struct {
    CatalogId    *int64
    Keyword      string
    OrgType      *int32
    State        *int32
    DepartmentId string
    Page         int
    PageSize     int
    Sort         string
    Direction    string
}
```

---

## Data Model

### DDL（复用Java结构）

**位置**: `migrations/stdfile/raw/t_std_file.sql`

```sql
CREATE TABLE `t_std_file` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_number` varchar(64) DEFAULT NULL COMMENT '标准编号',
  `f_name` varchar(256) NOT NULL COMMENT '标准文件名称',
  `f_catalog_id` bigint(20) NOT NULL COMMENT '目录ID',
  `f_act_date` datetime DEFAULT NULL COMMENT '实施日期',
  `f_publish_date` datetime DEFAULT NULL COMMENT '发布日期',
  `f_disable_date` datetime DEFAULT NULL COMMENT '停用时间',
  `f_attachment_type` INT(2) NOT NULL DEFAULT 0 COMMENT '附件类型：0-文件附件，1-外置链接',
  `f_attachment_url` varchar(500) DEFAULT NULL COMMENT '链接地址',
  `f_file_name` varchar(256) DEFAULT NULL COMMENT '文件名',
  `f_org_type` INT(2) NOT NULL COMMENT '标准组织类型',
  `f_description` varchar(300) DEFAULT NULL COMMENT '说明',
  `f_state` INT(2) NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-停用',
  `f_disable_reason` varchar(800) DEFAULT NULL COMMENT '停用原因',
  `f_authority_id` varchar(100) DEFAULT NULL COMMENT '权限域',
  `f_department_ids` varchar(350) DEFAULT NULL COMMENT '部门ID',
  `f_third_dept_id` varchar(36) DEFAULT NULL COMMENT '第三方部门ID',
  `f_version` int(4) NOT NULL DEFAULT 1 COMMENT '版本号',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_create_user` varchar(128) DEFAULT NULL COMMENT '创建用户',
  `f_update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `f_update_user` varchar(128) DEFAULT NULL COMMENT '修改用户',
  `f_deleted` bigint(20) NOT NULL DEFAULT 0 COMMENT '逻辑删除标记',
  PRIMARY KEY (`f_id`),
  KEY `idx_catalog_id` (`f_catalog_id`),
  KEY `idx_org_type` (`f_org_type`),
  KEY `uk_number_deleted` (`f_number`,`f_deleted`),
  KEY `uk_name_orgtype_deleted` (`f_name`,`f_org_type`,`f_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标准文件管理表';
```

### Go Struct

```go
// StdFile 标准文件
type StdFile struct {
    Id             int64      `db:"f_id" json:"id"`
    Number         string     `db:"f_number" json:"number"`
    Name           string     `db:"f_name" json:"name"`
    CatalogId      int64      `db:"f_catalog_id" json:"catalogId"`
    ActDate        time.Time  `db:"f_act_date" json:"actDate"`
    PublishDate    time.Time  `db:"f_publish_date" json:"publishDate"`
    DisableDate    time.Time  `db:"f_disable_date" json:"disableDate"`
    AttachmentType int        `db:"f_attachment_type" json:"attachmentType"`
    AttachmentUrl  string     `db:"f_attachment_url" json:"attachmentUrl"`
    FileName       string     `db:"f_file_name" json:"fileName"`
    OrgType        int        `db:"f_org_type" json:"orgType"`
    Description    string     `db:"f_description" json:"description"`
    State          int        `db:"f_state" json:"state"`
    DisableReason  string     `db:"f_disable_reason" json:"disableReason"`
    AuthorityId    string     `db:"f_authority_id" json:"authorityId"`
    DepartmentIds  string     `db:"f_department_ids" json:"departmentIds"`
    ThirdDeptId    string     `db:"f_third_dept_id" json:"thirdDeptId"`
    Version        int        `db:"f_version" json:"version"`
    CreateTime     time.Time  `db:"f_create_time" json:"createTime"`
    CreateUser     string     `db:"f_create_user" json:"createUser"`
    UpdateTime     time.Time  `db:"f_update_time" json:"updateTime"`
    UpdateUser     string     `db:"f_update_user" json:"updateUser"`
    Deleted        int64      `db:"f_deleted" json:"deleted"`
}

// StdFileVo 标准文件视图（包含关联信息）
type StdFileVo struct {
    StdFile
    CatalogName        string `json:"catalogName,omitempty"`
    DepartmentId       string `json:"departmentId,omitempty"`
    DepartmentName     string `json:"departmentName,omitempty"`
    DepartmentPathNames string `json:"departmentPathNames,omitempty"`
}
```

### 枚举定义

```go
// OrgTypeEnum 标准组织类型
const (
    OrgTypeGroup         int32 = 0  // 团体标准
    OrgTypeEnterprise    int32 = 1  // 企业标准
    OrgTypeIndustry      int32 = 2  // 行业标准
    OrgTypeLocal         int32 = 3  // 地方标准
    OrgTypeNational      int32 = 4  // 国家标准
    OrgTypeInternational int32 = 5  // 国际标准
    OrgTypeForeign       int32 = 6  // 国外标准
    OrgTypeOther         int32 = 99 // 其他标准
)

// AttachmentTypeEnum 附件类型
const (
    AttachmentTypeFile int32 = 0 // 文件附件
    AttachmentTypeURL  int32 = 1 // 外置链接
)

// EnableDisableStatusEnum 启用停用状态
const (
    StateDisable int32 = 0 // 停用
    StateEnable  int32 = 1 // 启用
)
```

---

## API Contract

**位置**: `api/doc/stdfile/stdfile.api`

### 完整 API 定义

```api
syntax = "v1"

import "../base.api"

type (
    // ==================== 请求类型 ====================

    // CreateStdFileReq 新增标准文件请求
    CreateStdFileReq {
        Number         string  `form:"number,optional"`
        Name           string  `form:"name,optional" validate:"required,max=300"`
        CatalogId      int64   `form:"catalog_id,optional,default=44"`
        OrgType        int32   `form:"org_type,optional" validate:"required"`
        ActDate        string  `form:"act_date,optional"`
        PublishDate    string  `form:"publish_date,optional"`
        Description    string  `form:"description,optional" validate:"max=300"`
        AttachmentType string  `form:"attachment_type,optional" validate:"required,oneof=FILE URL"`
        AttachmentUrl  string  `form:"attachment_url,optional"`
        State          string  `form:"state,optional" validate:"omitempty,oneof=enable disable"`
        DepartmentIds  string  `form:"department_ids,optional"`
    }

    // UpdateStdFileReq 修改标准文件请求
    UpdateStdFileReq {
        Number         string  `form:"number,optional"`
        Name           string  `form:"name,optional" validate:"required,max=300"`
        CatalogId      int64   `form:"catalog_id,optional" validate:"required"`
        OrgType        int32   `form:"org_type,optional" validate:"required"`
        ActDate        string  `form:"act_date,optional"`
        PublishDate    string  `form:"publish_date,optional"`
        Description    string  `form:"description,optional" validate:"max=300"`
        AttachmentType string  `form:"attachment_type,optional" validate:"required,oneof=FILE URL"`
        AttachmentUrl  string  `form:"attachment_url,optional"`
        State          string  `form:"state,optional" validate:"omitempty,oneof=enable disable"`
        DepartmentIds  string  `form:"department_ids,optional"`
    }

    // UpdateStdFileStateReq 修改文件状态请求
    UpdateStdFileStateReq {
        State  string `json:"state,optional" validate:"required,oneof=enable disable"`
        Reason string `json:"reason,optional"`
    }

    // RemoveCatalogReq 移动目录请求
    RemoveCatalogReq {
        Ids       []int64 `json:"ids,optional" validate:"required"`
        CatalogId int64   `json:"catalogId,optional" validate:"required"`
    }

    // BatchStateReq 批量状态更新请求
    BatchStateReq {
        Ids   []int64 `json:"ids,optional" validate:"required"`
        State string  `json:"state,optional" validate:"required,oneof=enable disable"`
        Reason string  `json:"reason,optional"`
    }

    // QueryByIdsReq 批量查询请求
    QueryByIdsReq {
        Ids []int64 `json:"ids,optional" validate:"required"`
    }

    // StdFileListQuery 列表查询参数
    StdFileListQuery {
        CatalogId    int64   `form:"catalog_id,optional"`
        Keyword      string  `form:"keyword,optional"`
        OrgType      int32   `form:"org_type,optional"`
        State        string  `form:"state,optional" validate:"omitempty,oneof=enable disable"`
        Offset       int     `form:"offset,optional,default=1"`
        Limit        int     `form:"limit,optional,default=20"`
        Sort         string  `form:"sort,optional,default=update_time"`
        Direction    string  `form:"direction,optional,default=desc"`
        DepartmentId string  `form:"department_id,optional"`
    }

    // StdFileRelationQuery 关联查询参数
    StdFileRelationQuery {
        Offset int `form:"offset,optional,default=1"`
        Limit  int `form:"limit,optional,default=20"`
    }

    // QueryDataExistsReq 检查数据存在请求
    QueryDataExistsReq {
        Number        string  `form:"number,optional"`
        OrgType       int32   `form:"org_type,optional"`
        Name          string  `form:"name,optional"`
        FilterId      int64   `form:"filter_id,optional"`
        DepartmentIds string  `form:"department_ids,optional"`
    }

    // StdFileRelationDto 关联关系请求
    StdFileRelationDto {
        DeIds    []int64 `json:"deIds,optional"`
        DictIds  []int64 `json:"dictIds,optional"`
        RuleIds  []int64 `json:"ruleIds,optional"`
    }

    // ==================== 响应类型 ====================

    // StdFileResp 标准文件响应
    StdFileResp {
        Id               int64  `json:"id"`
        Number           string `json:"number"`
        Name             string `json:"name"`
        CatalogId        int64  `json:"catalogId"`
        CatalogName      string `json:"catalogName,omitempty"`
        ActDate          string `json:"actDate,omitempty"`
        PublishDate      string `json:"publishDate,omitempty"`
        DisableDate      string `json:"disableDate,omitempty"`
        AttachmentType   string `json:"attachmentType"`
        AttachmentUrl    string `json:"attachmentUrl,omitempty"`
        FileName         string `json:"fileName,omitempty"`
        OrgType          int32  `json:"orgType"`
        Description      string `json:"description,omitempty"`
        State            string `json:"state"`
        DisableReason    string `json:"disableReason,omitempty"`
        Version          int    `json:"version"`
        DepartmentId     string `json:"departmentId,omitempty"`
        DepartmentName   string `json:"departmentName,omitempty"`
        DepartmentPathNames string `json:"departmentPathNames,omitempty"`
        CreateTime       string `json:"createTime,omitempty"`
        CreateUser       string `json:"createUser,omitempty"`
        UpdateTime       string `json:"updateTime,omitempty"`
        UpdateUser       string `json:"updateUser,omitempty"`
    }

    // StdFileListResp 文件列表响应
    StdFileListResp {
        TotalCount int64        `json:"totalCount"`
        Data       []StdFileResp `json:"data"`
    }

    // StdFileRelationResp 关联关系响应
    StdFileRelationResp {
        DeIds    []int64 `json:"deIds"`
        DictIds  []int64 `json:"dictIds"`
        RuleIds  []int64 `json:"ruleIds"`
    }

    // BaseResp 基础响应
    BaseResp {
        Code        string `json:"code"`
        Description string `json:"description"`
    }
)

// ==================== 外部API（open标签）====================

@server(
    prefix: /api/standardization/v1
    group: stdfile
    middleware: TokenCheck
)
service standardization-api {
    @doc "新增标准文件"
    @handler CreateStdFile
    post /std-file (CreateStdFileReq) returns (StdFileResp)

    @doc "根据ID修改标准文件"
    @handler UpdateStdFile
    put /std-file/:id (UpdateStdFileReq) returns (StdFileResp)

    @doc "标准文件-列表查询"
    @handler ListStdFile
    get /std-file (StdFileListQuery) returns (StdFileListResp)

    @doc "根据ID查询详情"
    @handler GetStdFile
    get /std-file/:id returns (StdFileResp)

    @doc "根据文件ID启用/停用"
    @handler UpdateStdFileState
    put /std-file/state/:id (UpdateStdFileStateReq) returns (BaseResp)

    @doc "根据id列表查询"
    @handler QueryStdFileByIds
    post /std-file/queryByIds (QueryByIdsReq) returns (StdFileListResp)

    @doc "标准文件-根据标准文件目录结构分页查询"
    @handler QueryStdFileByCatalog
    get /std-file/queryByCatalog (StdFileListQuery) returns (StdFileListResp)
}

// ==================== 其他API（无认证）====================

@server(
    prefix: /api/standardization/v1
    group: stdfile
)
service standardization-api {
    @doc "标准文件管理-批量删除"
    @handler DeleteStdFile
    delete /std-file/delete/:ids returns (BaseResp)

    @doc "移动到指定目录-目录移动"
    @handler RemoveStdFileCatalog
    post /std-file/catalog/remove (RemoveCatalogReq) returns (BaseResp)

    @doc "根据文件ID下载标准文件附件"
    @handler DownloadStdFile
    get /std-file/download/:id

    @doc "标准文件附件下载（批量）"
    @handler DownloadBatchStdFile
    post /std-file/downloadBatch (QueryByIdsReq)

    @doc "查询关联数据元"
    @handler QueryRelationDe
    get /std-file/relation/de/:id (StdFileRelationQuery) returns (StdFileListResp)

    @doc "查询关联的码表"
    @handler QueryRelationDict
    get /std-file/relation/dict/:id (StdFileRelationQuery) returns (StdFileListResp)

    @doc "查询关联编码规则"
    @handler QueryRelationRule
    get /std-file/relation/rule/:id (StdFileRelationQuery) returns (StdFileListResp)

    @doc "根据标准文件ID添加关联关系"
    @handler AddRelation
    put /std-file/relation/:id (StdFileRelationDto) returns (BaseResp)

    @doc "标准文件关联关系查询"
    @handler QueryRelations
    get /std-file/relation/:id returns (StdFileRelationResp)

    @doc "查询数据是否存在"
    @handler QueryDataExists
    get /std-file/queryDataExists (QueryDataExistsReq) returns (BaseResp)

    @doc "批量启用/停用"
    @handler BatchStateStdFile
    put /std-file/batchState (BatchStateReq) returns (BaseResp)
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
| P1 | 创建FILE/URL类型文件、修改文件、版本递增 |
| P1 | 停用文件（带原因）、启用文件 |
| P1 | 标准编号唯一性校验、目录存在性校验 |
| P1 | 文件类型校验、文件大小校验 |
| P2 | 列表查询（多条件筛选）、批量查询 |
| P2 | 目录移动、批量删除 |
| P2 | 关联查询（数据元、码表、编码规则） |
| P3 | 文件下载、批量下载 |

---

## Key Implementation Details

### 1. 文件上传处理

```go
// 处理文件上传
func (l *CreateStdFileLogic) handleFileUpload(file *multipart.FileHeader, fileId int64) error {
    if file == nil {
        return errorx.NewWithMsg(ErrCodeStdFileParamEmpty, "文件不能为空")
    }

    // 校验文件类型
    allowedTypes := []string{".doc", ".pdf", ".docx", ".txt", ".ppt", ".pptx", ".xls", ".xlsx"}
    if !isAllowedFileType(file.Filename, allowedTypes) {
        return errorx.NewWithMsg(ErrCodeStdFileInvalidParam, "不支持的文件类型")
    }

    // 校验文件大小
    if file.Size > 30*1024*1024 {
        return errorx.NewWithMsg(ErrCodeStdFileInvalidParam, "文件不能超过30M")
    }

    // 上传到OSS
    return l.svcCtx.OSSClient.UploadFile(file, strconv.FormatInt(fileId, 10))
}
```

### 2. 版本变更检测

```go
// CheckVersionChange 检测是否需要递增版本号
func CheckVersionChange(old *StdFile, new *UpdateStdFileReq) bool {
    // 比较字段：number, name, catalog_id, department_ids, org_type, description, act_date, attachment_type, attachment_url
    // 任何字段变化都返回 true
}
```

### 3. 批量下载文件名处理

```go
// getZipFileName 获取ZIP中的文件名（处理重复）
func getZipFileName(file *StdFile, duplicates map[string]int) string {
    baseName := strings.TrimSuffix(file.FileName, filepath.Ext(file.FileName))
    ext := filepath.Ext(file.FileName)

    if count := duplicates[file.FileName]; count > 1 {
        return fmt.Sprintf("%s(%s)(%s)%s", baseName, file.OrgType, file.Name, ext)
    }
    return file.FileName
}
```

### 4. 名称唯一性校验

```go
// CheckNumberUnique 校验标准编号唯一性
func (l *CreateStdFileLogic) CheckNumberUnique(number string) error {
    if number == "" {
        return nil
    }
    files, err := l.svcCtx.StdFileModel.FindByNumber(l.ctx, number)
    if err != nil {
        return err
    }
    if len(files) > 0 {
        return errorx.NewWithMsg(ErrCodeStdFileDataExist, "标准编号重复")
    }
    return nil
}

// CheckNameUnique 校验文件名称唯一性（同一orgType下）
func (l *CreateStdFileLogic) CheckNameUnique(name string, orgType int32) error {
    files, err := l.svcCtx.StdFileModel.FindByNameAndOrgType(l.ctx, name, orgType)
    if err != nil {
        return err
    }
    if len(files) > 0 {
        return errorx.NewWithMsg(ErrCodeStdFileDataExist, "标准文件名称重复")
    }
    return nil
}
```

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-02-06 | - | 初始版本 |
