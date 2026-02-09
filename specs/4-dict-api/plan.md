# 码表管理 (dict-api) Technical Plan

> **Branch**: `4-dict-api`
> **Spec Path**: `specs/4-dict-api/`
> **Created**: 2026-02-06
> **Status**: Draft

---

## Summary

从 Java 迁移码表管理模块到 Go-Zero，实现16个API接口，支持码表及码值的完整管理功能。技术方案核心决策：

1. **主键兼容性**：使用 Long 类型 BIGINT 主键（与Java保持一致，这是项目宪章例外情况）
2. **表结构保持**：完全复用 Java 表结构（f_ 前缀字段）
3. **接口兼容**：16个API路径、参数、响应格式、错误码100%一致
4. **纯 SQLx 策略**：所有数据库操作使用 SQLx，手工编写 SQL 查询
5. **雪花算法**：码表code使用雪花算法生成，全局唯一

**废弃接口**（不实现）：
- POST /v1/dataelement/dict/import - 导入
- POST /v1/dataelement/dict/export - 导出

**实际实现**：16个接口

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
| **Snowflake** | 雪花算法生成码表code |

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
| 码表管理 | 30400-30499 | `api/internal/errorx/dict.go` |

**错误码映射**（与Java保持一致）：

| 错误码常量 | 值 | 说明 |
|-----------|-----|------|
| ErrCodeDictDataNotExist | 30401 | 数据不存在 |
| ErrCodeDictParamEmpty | 30402 | 参数为空 |
| ErrCodeDictInvalidParam | 30403 | 参数无效 |
| ErrCodeDictCatalogNotExist | 30404 | 目录不存在 |
| ErrCodeDictChNameDuplicate | 30405 | 中文名称重复 |
| ErrCodeDictEnNameDuplicate | 30406 | 英文名称重复 |
| ErrCodeDictEnumCodeEmpty | 30407 | 码值为空 |
| ErrCodeDictEnumValueEmpty | 30408 | 码值描述为空 |
| ErrCodeDictEnumCodeDuplicate | 30409 | 码值重复 |
| ErrCodeDictReasonTooLong | 30410 | 停用原因过长 |

**⚠️ 重要**：错误码单独文件 `api/internal/errorx/dict.go`，不包含在 codes.go 中

### 第三方库确认

| 库 | 原因 | 确认状态 |
|----|------|----------|
| github.com/google/uuid | UUID v7 生成（部门ID等） | ✅ 已批准 |
| encoding/json | JSON序列化（Go标准库） | ✅ 已批准 |

---

## Go-Zero 开发流程

按以下顺序完成技术设计和代码生成：

| Step | 任务 | 方式 | 产出 |
|------|------|------|------|
| 1 | 定义 API 文件 | AI 实现 | `api/doc/dict/dict.api` |
| 2 | 生成 Handler/Types | goctl 生成 | `api/internal/handler/dict/`, `types/` |
| 3 | 实现 Model 接口 | AI 手写 | `model/dict/dict/` |
| 4 | 实现 Logic 层 | AI 实现 | `api/internal/logic/dict/` |
| 5 | 更新路由注册 | 手动 | `api/doc/api.api` |

**goctl 命令**:
```bash
# 步骤1：在 api/doc/api.api 中 import dict 模块
# 步骤2：执行 goctl 生成代码（针对整个项目）
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

**⚠️ 永远使用此命令**：`goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group`

---

## File Structure

### 文件产出清单

| 序号 | 文件 | 生成方式 | 位置 |
|------|------|----------|------|
| 1 | API 文件 | AI 实现 | `api/doc/dict/dict.api` |
| 2 | DDL 文件 | 复用Java | `migrations/dict/raw/*.sql` |
| 3 | Model 接口 | AI 手写 | `model/dict/dict/interface.go` |
| 4 | Model 类型 | AI 手写 | `model/dict/dict/types.go` |
| 5 | Model 常量 | AI 手写 | `model/dict/dict/vars.go` |
| 6 | SQLx 实现 | AI 手写 | `model/dict/dict/sql_model.go` |
| 7 | Handler | goctl 生成 | `api/internal/handler/dict/` |
| 8 | Types | goctl 生成 | `api/internal/types/dict.go` |
| 9 | Logic | AI 实现 | `api/internal/logic/dict/` |
| 10 | 错误码 | AI 手写 | `api/internal/errorx/dict.go` |
| 11 | Mock | AI 实现 | `api/internal/logic/dict/mock/` |

### 代码结构

```
api/internal/
├── handler/dict/
│   ├── list_dict_handler.go              # goctl 生成
│   ├── get_dict_handler.go               # goctl 生成
│   ├── get_dict_by_code_handler.go       # goctl 生成
│   ├── create_dict_handler.go            # goctl 生成
│   ├── update_dict_handler.go            # goctl 生成
│   ├── delete_dict_handler.go            # goctl 生成
│   ├── batch_delete_dict_handler.go      # goctl 生成
│   ├── update_dict_state_handler.go      # goctl 生成
│   ├── list_dict_enum_handler.go         # goctl 生成
│   ├── get_dict_enum_list_handler.go     # goctl 生成
│   ├── query_dict_by_data_element_handler.go  # goctl 生成
│   ├── query_dict_by_std_file_catalog_handler.go  # goctl 生成
│   ├── query_dict_by_std_file_handler.go # goctl 生成
│   ├── query_dict_relation_stdfile_handler.go     # goctl 生成
│   ├── add_dict_relation_handler.go      # goctl 生成
│   └── query_dict_data_exists_handler.go # goctl 生成
├── logic/dict/
│   ├── list_dict_logic.go                # AI 实现
│   ├── get_dict_logic.go                 # AI 实现
│   ├── get_dict_by_code_logic.go         # AI 实现
│   ├── create_dict_logic.go              # AI 实现
│   ├── update_dict_logic.go              # AI 实现
│   ├── delete_dict_logic.go              # AI 实现
│   ├── batch_delete_dict_logic.go        # AI 实现
│   ├── update_dict_state_logic.go        # AI 实现
│   ├── list_dict_enum_logic.go           # AI 实现
│   ├── get_dict_enum_list_logic.go       # AI 实现
│   ├── query_dict_by_data_element_logic.go  # AI 实现
│   ├── query_dict_by_std_file_catalog_logic.go  # AI 实现
│   ├── query_dict_by_std_file_logic.go   # AI 实现
│   ├── query_dict_relation_stdfile_logic.go     # AI 实现
│   ├── add_dict_relation_logic.go        # AI 实现
│   ├── query_dict_data_exists_logic.go   # AI 实现
│   └── common.go                         # 公共函数
├── types/
│   └── dict.go                           # goctl 生成
└── svc/
    └── servicecontext.go                  # 手动维护（注入Model）

model/dict/dict/
├── interface.go                           # Model 接口定义
├── types.go                               # 数据结构
├── vars.go                                # 常量/错误定义
├── factory.go                             # SQLx 工厂函数
└── sql_model.go                           # SQLx 实现

api/internal/logic/dict/mock/
├── catalog.go                             # 目录服务 Mock
├── dataelement.go                         # 数据元服务 Mock
├── stdfile.go                             # 标准文件服务 Mock
└── token.go                               # Token 服务 Mock

api/internal/errorx/
└── dict.go                                # 码表错误码定义（30400-30499）
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
| Logic | 业务逻辑实现 | 50 | 业务规则、事务管理、调用依赖服务、标记处理步骤 |
| Model | 数据访问 | 50 | CRUD操作、复杂查询、事务管理 |

---

## Interface Definitions

### DictModel 接口

```go
type DictModel interface {
    // 基础CRUD
    Insert(ctx context.Context, data *Dict) (int64, error)
    FindOne(ctx context.Context, id int64) (*Dict, error)
    FindByCode(ctx context.Context, code int64) (*Dict, error)
    Update(ctx context.Context, data *Dict) error
    Delete(ctx context.Context, id int64) error

    // 查询方法
    FindByIds(ctx context.Context, ids []int64) ([]*Dict, error)
    FindByCatalogIds(ctx context.Context, opts *FindOptions) ([]*Dict, int64, error)
    FindByStdFileCatalog(ctx context.Context, opts *FindOptions) ([]*Dict, int64, error)
    FindByFileId(ctx context.Context, fileId int64) ([]*Dict, error)
    FindNotUsedStdFile(ctx context.Context, opts *FindOptions) ([]*Dict, int64, error)
    FindByStdFile(ctx context.Context, fileId int64, opts *FindOptions) ([]*Dict, int64, error)
    FindDataExists(ctx context.Context, chName, enName string, orgType int32, filterId int64, deptIds string) (*Dict, error)

    // 更新方法
    UpdateState(ctx context.Context, id int64, state int32, reason string) error
    UpdateVersionByIds(ctx context.Context, ids []int64, updateUser string) error

    // 批量操作
    DeleteByIds(ctx context.Context, ids []int64) error
}
```

---

## Data Model

### DDL（复用Java结构）

**位置**: `migrations/dict/raw/`

```sql
CREATE TABLE `t_dict` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_code` bigint(20) NOT NULL COMMENT '码表编码',
  `f_ch_name` varchar(128) NOT NULL COMMENT '中文名称',
  `f_en_name` varchar(128) NOT NULL COMMENT '英文名称',
  `f_description` varchar(300) DEFAULT NULL COMMENT '业务含义',
  `f_catalog_id` bigint(20) NOT NULL COMMENT '目录ID',
  `f_org_type` INT(2) NOT NULL COMMENT '所属组织类型',
  `f_version` INT(4) NOT NULL DEFAULT 1 COMMENT '版本号',
  `f_state` INT(2) NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-停用',
  `f_disable_reason` VARCHAR(1024) DEFAULT NULL COMMENT '停用原因',
  `f_authority_id` varchar(100) DEFAULT NULL COMMENT '权限域',
  `f_department_ids` varchar(350) DEFAULT NULL COMMENT '部门ID',
  `f_third_dept_id` varchar(36) DEFAULT NULL COMMENT '第三方部门ID',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_create_user` varchar(128) DEFAULT NULL COMMENT '创建用户',
  `f_update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `f_update_user` varchar(128) DEFAULT NULL COMMENT '修改用户',
  `f_deleted` bigint(20) NOT NULL DEFAULT 0 COMMENT '逻辑删除标记',
  PRIMARY KEY (`f_id`),
  KEY `idx_code` (`f_code`),
  KEY `idx_orgtype_deleted` (`f_org_type`,`f_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='码表主表';

CREATE TABLE `t_dict_enum` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_dict_id` bigint(20) NOT NULL COMMENT '码表ID',
  `f_code` varchar(50) NOT NULL COMMENT '码值',
  `f_value` varchar(128) NOT NULL COMMENT '码值描述',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_create_user` varchar(128) DEFAULT NULL COMMENT '创建用户',
  PRIMARY KEY (`f_id`),
  KEY `idx_dict_id` (`f_dict_id`),
  KEY `idx_code` (`f_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='码表明细表';

CREATE TABLE `t_relation_dict_file` (
  `f_id` bigint(20) NOT NULL COMMENT '主键',
  `f_dict_id` bigint(20) NOT NULL COMMENT '码表ID',
  `f_file_id` bigint(20) NOT NULL COMMENT '文件ID',
  UNIQUE KEY `uk_dictid_fileid` (`f_dict_id`,`f_file_id`),
  PRIMARY KEY (`f_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='码表-文件关系表';
```

### Go Struct

```go
// Dict 码表
type Dict struct {
    Id             int64     `db:"f_id" json:"id"`
    Code           int64     `db:"f_code" json:"code"`
    ChName         string    `db:"f_ch_name" json:"chName"`
    EnName         string    `db:"f_en_name" json:"enName"`
    Description    string    `db:"f_description" json:"description"`
    CatalogId      int64     `db:"f_catalog_id" json:"catalogId"`
    OrgType        int32     `db:"f_org_type" json:"orgType"`
    Version        int32     `db:"f_version" json:"version"`
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
```

---

## 重要约束

### 永远使用的 goctl 命令

```bash
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

**说明**：
- 必须针对 `api/doc/api.api` 入口文件执行
- 不能针对单个功能文件执行
- 这会生成整个项目的代码

### Logic 层处理步骤标记

每个 Logic 函数必须标记处理步骤：

```go
func (l *CreateDictLogic) CreateDict(req *types.CreateDictReq) (*types.DictVo, error) {
    // Step 1: 参数校验
    // Step 2: 业务校验（目录存在性、名称唯一性）
    // Step 3: 生成码表编码（雪花算法）
    // Step 4: 保存码表
    // Step 5: 保存码值明细
    // Step 6: 保存关联文件关系
}
```

### 禁止修改其他需求代码

- 禁止修改 `specs/1-rule-api/` 下的任何文件
- 禁止修改 `specs/2-catalog-api/` 下的任何文件
- 禁止修改 `specs/3-std-file-api/` 下的任何文件
- 禁止修改 `api/internal/handler/rule/` 下的任何文件
- 禁止修改 `api/internal/handler/catalog/` 下的任何文件
- 禁止修改 `api/internal/handler/stdfile/` 下的任何文件
- 禁止修改 `api/internal/logic/rule/` 下的任何文件
- 禁止修改 `api/internal/logic/catalog/` 下的任何文件
- 禁止修改 `api/internal/logic/stdfile/` 下的任何文件

---

## Testing Strategy

| 类型 | 方法 | 覆盖率 |
|------|------|--------|
| 单元测试 | 表驱动测试，Mock Model | > 80% |
| 集成测试 | 测试数据库 | 核心流程 |

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-02-06 | - | 初始版本 |
