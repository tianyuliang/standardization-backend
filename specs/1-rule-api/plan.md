# 编码规则管理 (rule-api) Technical Plan

> **Branch**: `1-rule-api`
> **Spec Path**: `specs/1-rule-api/`
> **Created**: 2026-02-05
> **Status**: Draft

---

## Summary

从 Java 迁移编码规则管理模块到 Go-Zero，实现18个API接口，支持REGEX和CUSTOM两种编码规则类型。技术方案核心决策：

1. **主键兼容性**：使用 Long 类型 BIGINT 主键（与Java保持一致，这是项目宪章例外情况）
2. **表结构保持**：完全复用 Java 表结构（f_ 前缀字段）
3. **接口兼容**：18个API路径、参数、响应格式、错误码100%一致
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
- **影响**：仅影响 rule 模块，其他模块仍遵循 UUID v7 规范
- **风险**：低，Java已稳定运行多年

---

## Technical Context

| Item | Value |
|------|-------|
| **Language** | Go 1.24+ |
| **Framework** | Go-Zero v1.9+ |
| **Storage** | MySQL 8.0 (复用Java表结构) |
| **Cache** | Redis 7.0 |
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
| 编码规则 | 30300-30399 | `internal/errorx/codes.go` |

**错误码映射**（与Java保持一致）：

| Java 错误码 | Go 错误码 | 说明 |
|-------------|-----------|------|
| Standardization.DATA_NOT_EXIST | 30301 | 数据不存在 |
| Standardization.PARAMETER_EMPTY | 30302 | 参数为空 |
| Standardization.InvalidParameter | 30303 | 参数无效 |

### 第三方库确认

| 库 | 原因 | 确认状态 |
|----|------|----------|
| github.com/google/uuid | UUID v7 生成（部门ID等） | ✅ 已批准 |
| regexp | 正则表达式校验（Go标准库） | ✅ 已批准 |
| encoding/json | JSON序列化（Go标准库） | ✅ 已批准 |

---

## Go-Zero 开发流程

按以下顺序完成技术设计和代码生成：

| Step | 任务 | 方式 | 产出 |
|------|------|------|------|
| 1 | 定义 API 文件 | AI 实现 | `api/doc/rule/rule.api` |
| 2 | 生成 Handler/Types | goctl 生成 | `api/internal/handler/rule/`, `types/` |
| 3 | 实现 Model 接口 | AI 手写 | `model/rule/rule/` |
| 4 | 实现 Logic 层 | AI 实现 | `api/internal/logic/rule/` |
| 5 | 更新路由注册 | 手动 | `api/doc/api.api` |

**goctl 命令**:
```bash
# 步骤1：在 api/doc/api.api 中 import rule 模块
# 步骤2：执行 goctl 生成代码（针对整个项目）
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

---

## File Structure

### 文件产出清单

| 序号 | 文件 | 生成方式 | 位置 |
|------|------|----------|------|
| 1 | API 文件 | AI 实现 | `api/doc/rule/rule.api` |
| 2 | DDL 文件 | 复用Java | `migrations/rule/raw/t_rule.sql` |
| 3 | Model 接口 | AI 手写 | `model/rule/rule/interface.go` |
| 4 | Model 类型 | AI 手写 | `model/rule/rule/types.go` |
| 5 | Model 常量 | AI 手写 | `model/rule/rule/vars.go` |
| 6 | SQLx 实现 | AI 手写 | `model/rule/rule/sql_model.go` |
| 8 | Handler | goctl 生成 | `api/internal/handler/rule/` |
| 9 | Types | goctl 生成 | `api/internal/types/types.go` |
| 10 | Logic | AI 实现 | `api/internal/logic/rule/` |

### 代码结构

```
api/internal/
├── handler/rule/
│   ├── create_rule_handler.go           # goctl 生成
│   ├── update_rule_handler.go
│   ├── get_rule_handler.go
│   ├── get_rule_id_handler.go           # internal接口
│   ├── get_rule_detail_by_data_id_handler.go  # internal接口
│   ├── get_rule_detail_by_data_code_handler.go # internal接口
│   ├── list_rule_handler.go
│   ├── delete_rule_handler.go
│   ├── update_rule_state_handler.go
│   ├── remove_rule_catalog_handler.go
│   ├── query_rule_used_data_element_handler.go
│   ├── query_rule_by_ids_handler.go
│   ├── query_internal_rule_by_ids_handler.go  # internal接口
│   ├── query_rule_by_std_file_catalog_handler.go
│   ├── query_rule_by_std_file_handler.go
│   ├── query_data_exists_handler.go
│   ├── query_std_files_by_rule_handler.go
│   ├── get_custom_date_format_handler.go
│   └── routes.go
├── logic/rule/
│   ├── create_rule_logic.go              # AI 实现
│   ├── update_rule_logic.go
│   ├── get_rule_logic.go
│   ├── get_rule_id_logic.go
│   ├── get_rule_detail_by_data_id_logic.go
│   ├── get_rule_detail_by_data_code_logic.go
│   ├── list_rule_logic.go
│   ├── delete_rule_logic.go
│   ├── update_rule_state_logic.go
│   ├── remove_rule_catalog_logic.go
│   ├── query_rule_used_data_element_logic.go
│   ├── query_rule_by_ids_logic.go
│   ├── query_internal_rule_by_ids_logic.go
│   ├── query_rule_by_std_file_catalog_logic.go
│   ├── query_rule_by_std_file_logic.go
│   ├── query_data_exists_logic.go
│   ├── query_std_files_by_rule_logic.go
│   ├── get_custom_date_format_logic.go
│   └── common.go                          # 公共函数（表达式校验、MQ消息等）
├── types/
│   └── types.go                           # goctl 生成（含所有请求/响应类型）
└── svc/
    └── servicecontext.go                   # 手动维护（注入Model、MQ等）

model/rule/rule/
├── interface.go                           # Model 接口定义
├── types.go                               # 数据结构（Rule、RelationRuleFile）
├── vars.go                                # 常量/错误定义
├── factory.go                             # SQLx 工厂函数
└── sql_model.go                           # SQLx 实现（所有查询）
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
```

| 层级 | 职责 | 最大行数 | 关键操作 |
|------|------|----------|----------|
| Handler | 解析参数、格式化响应 | 30 | 参数绑定、调用Logic、返回结果 |
| Logic | 业务逻辑实现 | 50 | 业务规则、事务管理、调用依赖服务 |
| Model | 数据访问 | 50 | CRUD操作、复杂查询、事务管理 |

### 依赖服务调用

Logic 层通过 ServiceContext 调用依赖服务：

```go
type ServiceContext struct {
    RuleModel       model.RuleModel
    CatalogModel    model.CatalogModel    // 目录服务
    DataElementModel model.DataElementModel // 数据元服务
    DictModel       model.DictModel       // 字典服务
    StdFileModel    model.StdFileModel    // 标准文件服务
    KafkaProducer   *kafka.Producer       // MQ消息
}
```

---

## Interface Definitions

### RuleModel 接口

```go
type RuleModel interface {
    // 基础CRUD
    Insert(ctx context.Context, data *Rule) (int64, error)
    FindOne(ctx context.Context, id int64) (*Rule, error)
    Update(ctx context.Context, data *Rule) error
    Delete(ctx context.Context, id int64) error

    // 查询方法
    FindByIds(ctx context.Context, ids []int64) ([]*Rule, error)
    FindByNameAndOrgType(ctx context.Context, name string, orgType int32) ([]*Rule, error)
    FindByCatalogIds(ctx context.Context, opts *FindOptions) ([]*Rule, int64, error)
    FindByStdFileCatalog(ctx context.Context, opts *FindOptions) ([]*Rule, int64, error)
    FindByFileId(ctx context.Context, fileId int64) ([]*Rule, error)
    FindNotUsedStdFile(ctx context.Context, opts *FindOptions) ([]*Rule, int64, error)
    FindByStdFile(ctx context.Context, fileId int64, opts *FindOptions) ([]*Rule, int64, error)
    FindDataExists(ctx context.Context, name string, filterId int64, deptIds string) (*Rule, error)

    // 更新方法
    UpdateState(ctx context.Context, id int64, state int32, reason string) error
    RemoveCatalog(ctx context.Context, ids []int64, catalogId int64, updateUser string) error
    UpdateVersionByIds(ctx context.Context, ids []int64, updateUser string) error

    // 批量操作
    DeleteByIds(ctx context.Context, ids []int64) error
}

type FindOptions struct {
    OrgType     *int32
    State       *int32
    RuleType    *int32
    Keyword     string
    DepartmentId string
    Page        int
    PageSize    int
    Sort        string
    Direction   string
}
```

### RelationRuleFileModel 接口

```go
type RelationRuleFileModel interface {
    InsertBatch(ctx context.Context, data []*RelationRuleFile) error
    DeleteByRuleId(ctx context.Context, ruleId int64) error
    DeleteByFileId(ctx context.Context, fileId int64) error
    FindByRuleId(ctx context.Context, ruleId int64) ([]*RelationRuleFile, error)
    DeleteByRuleIds(ctx context.Context, ruleIds []int64) error
}
```

---

## Data Model

### DDL（复用Java结构）

**位置**: `migrations/rule/raw/t_rule.sql`

```sql
CREATE TABLE `t_rule` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_name` varchar(128) NOT NULL COMMENT '规则名称',
  `f_catalog_id` bigint(20) NOT NULL COMMENT '所属目录id',
  `f_org_type` INT(2) NOT NULL COMMENT '标准组织类型',
  `f_description` varchar(300) DEFAULT NULL COMMENT '说明',
  `f_rule_type` INT(2) NOT NULL DEFAULT 0 COMMENT '规则类型：0-正则表达式，1-自定义配置',
  `f_version` int(4) NOT NULL DEFAULT 1 COMMENT '版本号，从1开始',
  `f_expression` varchar(1024) NOT NULL COMMENT '表达式：正则表达式或JSON配置',
  `f_state` INT(2) NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-停用',
  `f_disable_reason` VARCHAR(1024) NULL COMMENT '停用原因',
  `f_authority_id` varchar(100) DEFAULT NULL COMMENT '权限域',
  `f_department_ids` varchar(350) NULL COMMENT '部门ID',
  `f_third_dept_id` varchar(36) NULL COMMENT '第三方部门ID',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_create_user` varchar(128) DEFAULT NULL COMMENT '创建用户',
  `f_update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `f_update_user` varchar(128) DEFAULT NULL COMMENT '修改用户',
  `f_deleted` bigint(20) NOT NULL DEFAULT 0 COMMENT '逻辑删除标记',
  PRIMARY KEY (`f_id`),
  KEY `uk_name_orgtype_deleted` (`f_name`,`f_org_type`,`f_deleted`,`f_department_ids`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='编码规则';

CREATE TABLE `t_relation_rule_file` (
  `f_id` bigint(20) NOT NULL COMMENT '主键',
  `f_rule_id` bigint(20) NOT NULL COMMENT '规则ID',
  `f_file_id` bigint(20) NOT NULL COMMENT '文件ID',
  UNIQUE KEY `uk_ruleid_fileid` (`f_rule_id`,`f_file_id`),
  PRIMARY KEY (`f_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='编码规则-文件关系表';
```

### Go Struct

```go
// Rule 编码规则
type Rule struct {
    Id             int64     `db:"f_id" json:"id"`
    Name           string    `db:"f_name" json:"name"`
    CatalogId      int64     `db:"f_catalog_id" json:"catalogId"`
    OrgType        int32     `db:"f_org_type" json:"orgType"`
    Description    string    `db:"f_description" json:"description"`
    RuleType       int32     `db:"f_rule_type" json:"ruleType"`
    Version        int32     `db:"f_version" json:"version"`
    Expression     string    `db:"f_expression" json:"-"`
    State          int32     `db:"f_state" json:"state"`
    DisableReason  string    `db:"f_disable_reason" json:"disableReason"`
    AuthorityId    string    `db:"f_authority_id" json:"authorityId"`
    DepartmentIds  string    `db:"f_department_ids" json:"departmentIds"`
    ThirdDeptId    string    `db:"f_third_dept_id" json:"thirdDeptId"`
    CreateTime     time.Time `db:"f_create_time" json:"createTime"`
    CreateUser     string    `db:"f_create_user" json:"createUser"`
    UpdateTime     time.Time `db:"f_update_time" json:"updateTime"`
    UpdateUser     string    `db:"f_update_user" json:"updateUser"`
    Deleted        int64     `db:"f_deleted" json:"deleted"`
}

// RuleCustom 自定义规则配置
type RuleCustom struct {
    SegmentLength int    `json:"segment_length"`
    Name          string `json:"name"`
    Type          string `json:"type"` // dict, number, english_letters, chinese_characters, any_characters, date, split_str
    Value         string `json:"value"`
}

// RuleVo 编码规则视图（包含关联信息）
type RuleVo struct {
    Rule
    CatalogName        string   `json:"catalogName"`
    FullCatalogName    string   `json:"fullCatalogName"`
    Regex              string   `json:"regex,omitempty"`     // REGEX类型时返回
    Custom             []RuleCustom `json:"custom,omitempty"` // CUSTOM类型时返回
    StdFiles           []int64  `json:"stdFiles,omitempty"`
    UsedFlag           bool     `json:"usedFlag"`
    DepartmentId       string   `json:"departmentId"`
    DepartmentName     string   `json:"departmentName"`
    DepartmentPathNames string  `json:"departmentPathNames"`
}

// RelationRuleFile 规则-文件关系
type RelationRuleFile struct {
    Id     int64 `db:"f_id" json:"id"`
    RuleId int64 `db:"f_rule_id" json:"ruleId"`
    FileId int64 `db:"f_file_id" json:"fileId"`
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

// RuleTypeEnum 规则类型
const (
    RuleTypeRegex  int32 = 0 // 正则表达式
    RuleTypeCustom int32 = 1 // 自定义配置
)

// EnableDisableStatusEnum 启用停用状态
const (
    StateDisable int32 = 0 // 停用
    StateEnable  int32 = 1 // 启用
)

// RuleCustomTypeEnum 自定义规则类型
const (
    CustomTypeDict             int32 = 1 // 码表
    CustomTypeNumber           int32 = 2 // 数字
    CustomTypeEnglishLetters   int32 = 3 // 英文字母
    CustomTypeChineseCharacters int32 = 4 // 汉字
    CustomTypeAnyCharacters    int32 = 5 // 任意字符
    CustomTypeDate             int32 = 6 // 日期
    CustomTypeSplitStr         int32 = 7 // 分割字符串
)

// CustomDateFormat 支持的日期格式
var CustomDateFormat = []string{
    "yyyy", "yyyyMM", "yyyy-MM-dd", "yyyyMMdd",
    "yyyy-MM-dd HH:mm:ss", "yyyyMMddHHmmss",
    "yyyy-MM-dd'T'HH:mm:ss", "yyyyMMdd'T'HHmmss",
    "HHmmss", "HH:mm:ss",
}
```

---

## API Contract

**位置**: `api/doc/rule/rule.api`

### 完整 API 定义

```api
syntax = "v1"

import "../base.api"

type (
    // ==================== 请求类型 ====================

    // CreateRuleReq 新增编码规则请求
    CreateRuleReq {
        Name         string        `json:"name,optional" validate:"required,max=128"`
        OrgType      int32         `json:"orgType,optional" validate:"required"`
        RuleType     string        `json:"ruleType,optional" validate:"required,oneof=REGEX CUSTOM"`
        Regex        string        `json:"regex,optional"`
        Custom       []RuleCustom  `json:"custom,optional"`
        Description  string        `json:"description,optional" validate:"max=300"`
        CatalogId    int64         `json:"catalogId,optional,default=33"`
        StdFiles     []int64       `json:"stdFiles,optional"`
        State        string        `json:"state,optional" validate:"omitempty,oneof=enable disable"`
        DepartmentIds string       `json:"departmentIds,optional"`
    }

    // UpdateRuleReq 修改编码规则请求
    UpdateRuleReq {
        Name         string        `json:"name,optional" validate:"required,max=128"`
        OrgType      int32         `json:"orgType,optional" validate:"required"`
        RuleType     string        `json:"ruleType,optional" validate:"required,oneof=REGEX CUSTOM"`
        Regex        string        `json:"regex,optional"`
        Custom       []RuleCustom  `json:"custom,optional"`
        Description  string        `json:"description,optional" validate:"max=300"`
        CatalogId    int64         `json:"catalogId,optional,default=33"`
        StdFiles     []int64       `json:"stdFiles,optional"`
        State        string        `json:"state,optional" validate:"omitempty,oneof=enable disable"`
        DepartmentIds string       `json:"departmentIds,optional"`
    }

    // UpdateRuleStateReq 修改规则状态请求
    UpdateRuleStateReq {
        State string `json:"state,optional" validate:"required,oneof=enable disable"`
        Reason string `json:"reason,optional"`
    }

    // RemoveCatalogReq 移动目录请求
    RemoveCatalogReq {
        Ids       []int64 `json:"ids,optional" validate:"required"`
        CatalogId int64   `json:"catalogId,optional" validate:"required"`
    }

    // QueryByIdsReq 批量查询请求
    QueryByIdsReq {
        Ids []int64 `json:"ids,optional" validate:"required"`
    }

    // RuleListQuery 列表查询参数
    RuleListQuery {
        CatalogId    int64   `form:"catalogId,optional"`
        Keyword      string  `form:"keyword,optional"`
        OrgType      int32   `form:"org_type,optional"`
        State        string  `form:"state,optional" validate:"omitempty,oneof=enable disable"`
        RuleType     string  `form:"rule_type,optional" validate:"omitempty,oneof=REGEX CUSTOM"`
        Offset       int     `form:"offset,optional,default=1"`
        Limit        int     `form:"limit,optional,default=20"`
        Sort         string  `form:"sort,optional,default=create_time"`
        Direction    string  `form:"direction,optional,default=desc"`
        DepartmentId string  `form:"department_id,optional"`
    }

    // QueryDataExistsReq 检查数据存在请求
    QueryDataExistsReq {
        Name          string `form:"name,optional" validate:"required"`
        FilterId      int64  `form:"filter_id,optional"`
        DepartmentIds string `form:"department_ids,optional"`
    }

    // RuleCustom 自定义规则配置
    RuleCustom {
        SegmentLength int    `json:"segment_length,optional" validate:"required"`
        Name          string `json:"name,optional"`
        Type          string `json:"type,optional" validate:"required"`
        Value         string `json:"value,optional" validate:"required"`
    }

    // ==================== 响应类型 ====================

    // RuleResp 编码规则响应
    RuleResp {
        Id                 int64         `json:"id"`
        Name               string        `json:"name"`
        CatalogId          int64         `json:"catalogId"`
        CatalogName        string        `json:"catalogName,omitempty"`
        FullCatalogName    string        `json:"fullCatalogName,omitempty"`
        OrgType            int32         `json:"orgType"`
        Description        string        `json:"description,omitempty"`
        RuleType           string        `json:"ruleType"`
        Version            int32         `json:"version"`
        Regex              string        `json:"regex,omitempty"`
        Custom             []RuleCustom  `json:"custom,omitempty"`
        State              string        `json:"state"`
        DisableReason      string        `json:"disableReason,omitempty"`
        StdFiles           []int64       `json:"stdFiles,omitempty"`
        UsedFlag           bool          `json:"usedFlag"`
        DepartmentId       string        `json:"departmentId,omitempty"`
        DepartmentName     string        `json:"departmentName,omitempty"`
        DepartmentPathNames string       `json:"departmentPathNames,omitempty"`
        CreateTime         string        `json:"createTime,omitempty"`
        CreateUser         string        `json:"createUser,omitempty"`
        UpdateTime         string        `json:"updateTime,omitempty"`
        UpdateUser         string        `json:"updateUser,omitempty"`
    }

    // RuleListResp 规则列表响应
    RuleListResp {
        TotalCount int64      `json:"totalCount"`
        Data       []RuleResp `json:"data"`
    }

    // DataElementResp 数据元响应（用于引用查询）
    DataElementResp {
        Id       int64  `json:"id"`
        DataCode string `json:"dataCode"`
        DataName string `json:"dataName"`
    }

    // StdFileResp 标准文件响应
    StdFileResp {
        Id       int64  `json:"id"`
        FileName string `json:"fileName"`
    }

    // CustomDateFormatResp 日期格式响应
    CustomDateFormatResp {
        Data []string `json:"data"`
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
    group: rule
    middleware: TokenCheck
)
service standardization-api {
    @doc "新增编码规则"
    @handler CreateRule
    post /rule (CreateRuleReq) returns (RuleResp)

    @doc "根据ID修改编码规则"
    @handler UpdateRule
    put /rule/:id (UpdateRuleReq) returns (RuleResp)

    @doc "编码规则-详情查看"
    @handler GetRule
    get /rule/:id (RuleListQuery) returns (RuleResp)

    @doc "编码规则-列表查询"
    @handler ListRule
    get /rule (RuleListQuery) returns (RuleListResp)

    @doc "根据编码规则id停用/启用"
    @handler UpdateRuleState
    put /rule/state/:id (UpdateRuleStateReq) returns (BaseResp)

    @doc "根据id列表查询"
    @handler QueryRuleByIds
    post /rule/queryByIds (QueryByIdsReq) returns (RuleListResp)

    @doc "编码规则-根据标准文件目录结构分页查询"
    @handler QueryRuleByStdFileCatalog
    get /rule/queryByStdFileCatalog (RuleListQuery) returns (RuleListResp)

    @doc "编码规则-根据标准文件分页查询"
    @handler QueryRuleByStdFile
    get /rule/queryByStdFile (RuleListQuery) returns (RuleListResp)

    @doc "编码规则-根据文件ID分页查询编码规则关联的标准文件"
    @handler QueryStdFilesByRule
    get /rule/relation/stdfile/:id (RuleListQuery) returns (RuleListResp)
}

// ==================== 内部API（无认证）====================

@server(
    prefix: /v1/rule
    group: rule_internal
)
service standardization-api {
    @doc "编码规则-根据ID查看详情（内部）"
    @handler GetRuleInternal
    get /internal/getId/:id returns (RuleResp)

    @doc "编码规则-根据数据元ID查看详情（内部）"
    @handler GetRuleByDataId
    get /internal/getDetailByDataId/:dataId returns (RuleResp)

    @doc "编码规则-根据数据元编码查看详情（内部）"
    @handler GetRuleByDataCode
    get /internal/getDetailByDataCode/:dataCode returns (RuleResp)

    @doc "根据id列表查询（内部接口）"
    @handler QueryRuleByIdsInternal
    post /internal/queryByIds (QueryByIdsReq) returns (RuleListResp)
}

// ==================== 其他API（无认证）====================

@server(
    prefix: /api/standardization/v1
    group: rule
)
service standardization-api {
    @doc "编码规则-删除&批量删除"
    @handler DeleteRule
    delete /rule/:ids returns (BaseResp)

    @doc "编码规则-移动到指定目录"
    @handler RemoveRuleCatalog
    post /rule/catalog/remove (RemoveCatalogReq) returns (BaseResp)

    @doc "编码规则-分页查询引用编码规则的标准数据元列表"
    @handler QueryRuleUsedDataElement
    get /rule/relation/de/:id (RuleListQuery) returns (RuleListResp)

    @doc "编码规则-查询数据是否存在"
    @handler QueryDataExists
    get /rule/queryDataExists (QueryDataExistsReq) returns (BaseResp)

    @doc "编码规则-获取自定义规则日期格式化字符串"
    @handler GetCustomDateFormat
    get /rule/getCustomDateFormat returns (CustomDateFormatResp)
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
| P1 | 创建REGEX/CUSTOM规则、修改规则、版本递增 |
| P1 | 停用规则（带原因）、启用规则 |
| P1 | 名称唯一性校验、目录存在性校验 |
| P1 | 正则表达式校验、自定义配置校验 |
| P2 | 列表查询（多条件筛选）、批量查询 |
| P2 | 目录移动、批量删除 |
| P2 | 关联查询（数据元、标准文件） |
| P3 | 内部接口（按数据元查询） |

---

## Key Implementation Details

### 1. 表达式校验（Logic/common.go）

```go
// ValidateExpression 校验规则表达式
func ValidateExpression(ruleType string, regex string, custom []RuleCustom) error {
    if ruleType == "REGEX" {
        if regex == "" {
            return errorx.NewWithCode(30303, "regex", "正则表达式为空")
        }
        if _, err := regexp.Compile(regex); err != nil {
            return errorx.NewWithCode(30303, "regex", "正则表达式非法")
        }
    } else {
        if len(custom) == 0 {
            return errorx.NewWithCode(30303, "custom", "不能为空")
        }
        for i, c := range custom {
            prefix := fmt.Sprintf("custom[%d].", i+1)
            if c.SegmentLength <= 0 {
                return errorx.NewWithCode(30303, prefix+"segment_length", "值必须为正整数")
            }
            // ... 其他校验
        }
    }
    return nil
}
```

### 2. 版本变更检测（Logic/common.go）

```go
// CheckVersionChange 检测是否需要递增版本号
func CheckVersionChange(old *Rule, new *UpdateRuleReq, oldFiles []*RelationRuleFile) bool {
    // 比较字段：name, catalog_id, department_ids, org_type, description, rule_type, expression, 关联文件
    // 任何字段变化都返回 true
}
```

### 3. MQ消息发送

```go
// SendRuleMQMessage 发送MQ消息
func SendRuleMQMessage(producer *kafka.Producer, rules []*Rule, operation string) error {
    dto := DataMqDto{
        Header:  make(map[string]interface{}),
        Payload: Payload{
            Type: "smart-recommendation-graph",
            Content: Content{
                Type:      operation, // insert/update/delete
                TableName: "t_rule",
                Entities:  rules,
            },
        },
    }
    // 序列化并发送到 Kafka
}
```

### 4. 名称唯一性校验

```go
// CheckNameUnique 校验名称唯一性（同一orgType下）
func (l *CreateRuleLogic) CheckNameUnique(name string, orgType int32, deptIds string) error {
    rules, err := l.svcCtx.RuleModel.FindByNameAndOrgType(l.ctx, name, orgType)
    if err != nil {
        return err
    }
    for _, r := range rules {
        if r.DepartmentIds == deptIds {
            return errorx.NewWithCode(30303, "name", "规则名称已存在")
        }
    }
    return nil
}
```

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-02-05 | - | 初始版本 |
