# 数据元管理 (dataelement-api) Technical Plan

> **Branch**: `6-std-dataelement-api`
> **Spec Path**: `specs/6-std-dataelement-api/`
> **Created**: 2026-02-09
> **Status**: Draft

---

## Summary

从 Java 迁移数据元管理模块到 Go-Zero，实现19个API接口，支持数据元的完整生命周期管理。技术方案核心决策：

1. **主键兼容性**：使用 BIGINT 类型主键（与Java保持一致）
2. **表结构保持**：完全复用 Java 表结构（f_ 前缀字段）
3. **接口兼容**：19个API路径、参数、响应格式、错误码100%一致
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
- **影响**：仅影响 dataelement 模块，与其他模块保持一致
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
| 数据元管理 | 30800-30899 | `internal/errorx/dataelement.go` |

**错误标识映射**（与Java保持一致）：

| Java 异常类 | Go 错误标识 | 说明 |
|-------------|-------------|------|
| DATA_NOT_EXIST | DataElementNotExist | 数据不存在 |
| PARAMETER_EMPTY | ParameterEmpty | 参数为空 |
| InvalidParameter | InvalidParameter | 参数无效 |

---

## Go-Zero 开发流程

按以下顺序完成技术设计和代码生成：

| Step | 任务 | 方式 | 产出 |
|------|------|------|------|
| 1 | 定义 API 文件 | AI 实现 | `api/doc/dataelement/dataelement.api` |
| 2 | 生成 Handler/Types | goctl 生成 | `api/internal/handler/dataelement/`, `types/` |
| 3 | 实现 Model 接口 | AI 手写 | `model/dataelement/dataelement/` |
| 4 | 实现 Logic 层 | AI 实现 | `api/internal/logic/dataelement/` |
| 5 | 更新路由注册 | 手动 | `api/doc/api.api` |

**⚠️ 重要：goctl 命令**
```bash
# 永远使用以下命令生成代码
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

---

## File Structure

### 文件产出清单

| 序号 | 文件 | 生成方式 | 位置 |
|------|------|----------|------|
| 1 | API 文件 | AI 实现 | `api/doc/dataelement/dataelement.api` |
| 2 | DDL 文件 | 复用Java | `migrations/dataelement/raw/t_data_element_info.sql` |
| 3 | DDL 文件 | 复用Java | `migrations/dataelement/raw/t_relation_de_file.sql` |
| 4 | Model 接口 | AI 手写 | `model/dataelement/dataelement/interface.go` |
| 5 | Model 类型 | AI 手写 | `model/dataelement/dataelement/types.go` |
| 6 | Model 常量 | AI 手写 | `model/dataelement/dataelement/vars.go` |
| 7 | SQLx 实现 | AI 手写 | `model/dataelement/dataelement/sql_model.go` |
| 8 | Handler | goctl 生成 | `api/internal/handler/dataelement/` |
| 9 | Types | goctl 生成 | `api/internal/types/types.go` |
| 10 | Logic | AI 实现 | `api/internal/logic/dataelement/` |

### 代码结构

```
api/internal/
├── handler/dataelement/
│   ├── create_data_element_handler.go        # goctl 生成
│   ├── import_data_element_handler.go        # 批量导入
│   ├── export_data_element_handler.go        # 导出
│   ├── export_by_ids_handler.go              # ID导出
│   ├── list_data_element_handler.go          # 分页列表
│   ├── internal_list_handler.go              # 内部分页列表
│   ├── get_data_element_detail_handler.go    # 详情
│   ├── internal_get_detail_handler.go        # 内部详情
│   ├── update_data_element_handler.go        # 编辑
│   ├── delete_data_element_handler.go        # 删除
│   ├── update_state_handler.go               # 启用停用
│   ├── move_catalog_handler.go               # 移动目录
│   ├── query_by_std_file_catalog_handler.go  # 按文件目录查询
│   ├── query_by_std_file_handler.go          # 按文件查询
│   ├── query_std_file_handler.go             # 查询关联文件
│   ├── query_is_repeat_handler.go            # 检查重名
│   ├── delete_label_ids_handler.go           # 删除标签
│   ├── query_list_handler.go                 # 查询列表
│   ├── internal_query_list_handler.go        # 内部查询列表
│   └── routes.go
├── logic/dataelement/
│   ├── create_data_element_logic.go          # AI 实现
│   ├── import_data_element_logic.go          # 批量导入
│   ├── export_data_element_logic.go          # 导出
│   ├── export_by_ids_logic.go                # ID导出
│   ├── list_data_element_logic.go            # 分页列表
│   ├── internal_list_logic.go                # 内部分页列表
│   ├── get_data_element_detail_logic.go      # 详情
│   ├── internal_get_detail_logic.go          # 内部详情
│   ├── update_data_element_logic.go          # 编辑
│   ├── delete_data_element_logic.go          # 删除
│   ├── update_state_logic.go                 # 启用停用
│   ├── move_catalog_logic.go                 # 移动目录
│   ├── query_by_std_file_catalog_logic.go    # 按文件目录查询
│   ├── query_by_std_file_logic.go            # 按文件查询
│   ├── query_std_file_logic.go               # 查询关联文件
│   ├── query_is_repeat_logic.go              # 检查重名
│   ├── delete_label_ids_logic.go             # 删除标签
│   ├── query_list_logic.go                   # 查询列表
│   ├── internal_query_list_logic.go          # 内部查询列表
│   ├── mock/                                 # Mock 服务
│   │   └── service.go
│   └── common.go                             # 公共函数
├── types/
│   ├── dataelement.go                        # 数据元类型
│   └── types.go                              # goctl 生成（全局类型）
└── svc/
    └── servicecontext.go                      # 手动维护（注入Model、MQ等）

model/dataelement/dataelement/
├── interface.go                              # Model 接口定义
├── types.go                                  # 数据结构（DataElement, RelationDeFile）
├── vars.go                                   # 常量/错误定义
├── factory.go                                # SQLx 工厂函数
└── sql_model.go                              # SQLx 实现（所有查询）

model/dataelement/relation/
├── interface.go                              # 关联关系Model接口
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
    Config config.Config

    // DataElement Models
    DataElementModel      model.DataElementModel
    RelationDeFileModel   model.RelationDeFileModel

    // RPC Clients (后续替换Mock)
    CatalogRpc    catalog.Catalog          // 目录服务
    DictRpc       dict.Dict                 // 码表服务
    RuleRpc       rule.Rule                 // 编码规则服务
    StdFileRpc    stdfile.StdFile           // 标准文件服务
}
```

---

## Interface Definitions

### DataElementModel 接口

```go
type DataElementModel interface {
    // 基础CRUD
    Insert(ctx context.Context, data *DataElement) (int64, error)
    FindOne(ctx context.Context, id int64) (*DataElement, error)
    FindOneByCode(ctx context.Context, code int64) (*DataElement, error)
    Update(ctx context.Context, data *DataElement) error
    Delete(ctx context.Context, id int64) error

    // 查询方法
    FindByCatalogIds(ctx context.Context, opts *FindOptions) ([]*DataElement, int64, error)
    FindByIds(ctx context.Context, ids []int64) ([]*DataElement, error)
    FindByCodes(ctx context.Context, codes []int64) ([]*DataElement, error)
    FindByFileCatalog(ctx context.Context, opts *FindOptions) ([]*DataElement, int64, error)
    FindByFileId(ctx context.Context, fileId int64, opts *FindOptions) ([]*DataElement, int64, error)
    FindByStdFileCatalog(ctx context.Context, opts *FindOptions) ([]*DataElement, int64, error)
    FindByStdFileId(ctx context.Context, fileId int64, opts *FindOptions) ([]*DataElement, int64, error)

    // 检查方法
    CheckNameCnExists(ctx context.Context, nameCn string, stdType int32, excludeId int64, deptIds string) (bool, error)
    CheckNameEnExists(ctx context.Context, nameEn string, excludeId int64, deptIds string) (bool, error)
    FindDataExists(ctx context.Context, name string, stdType int32, excludeId int64, deptIds string) (*DataElement, error)

    // 更新方法
    UpdateState(ctx context.Context, ids []int64, state int32, reason string) error
    MoveCatalog(ctx context.Context, ids []int64, catalogId int64, updateUser string) error
    DeleteLabelIds(ctx context.Context, ids []int64) error

    // 批量操作
    DeleteByIds(ctx context.Context, ids []int64) error
}

type FindOptions struct {
    CatalogId    *int64
    StdType      *int32
    State        *int32
    DataType     *int32
    Keyword      string
    DepartmentId string
    FileCatalogId *int64
    FileId       *int64
    Page         int
    PageSize     int
    Sort         string
    Direction    string
}
```

### RelationDeFileModel 接口

```go
type RelationDeFileModel interface {
    InsertBatch(ctx context.Context, data []*RelationDeFile) error
    DeleteByDeId(ctx context.Context, deId int64) error
    DeleteByDeIds(ctx context.Context, deIds []int64) error
    DeleteByFileId(ctx context.Context, fileId int64) error
    FindByDeId(ctx context.Context, deId int64, opts *PageOptions) ([]*RelationDeFile, int64, error)
    FindByDeIds(ctx context.Context, deIds []int64) (map[int64][]int64, error)
}

type PageOptions struct {
    Page    int
    PageSize int
}
```

---

## Data Model

### DDL（复用Java结构）

**位置**: `migrations/dataelement/raw/t_data_element_info.sql`

```sql
CREATE TABLE `t_data_element_info` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_de_id` bigint(20) NOT NULL COMMENT '关联标识',
  `f_name_en` varchar(128) NOT NULL COMMENT '英文名称',
  `f_name_cn` varchar(128) NOT NULL COMMENT '中文名称',
  `f_synonym` varchar(300) DEFAULT NULL COMMENT '同义词',
  `f_std_type` INT(2) NOT NULL COMMENT '标准分类',
  `f_data_type` INT(2) NOT NULL COMMENT '数据类型',
  `f_data_length` INT(4) DEFAULT NULL COMMENT '数据长度',
  `f_data_precision` INT(4) DEFAULT NULL COMMENT '数据精度',
  `f_dict_code` bigint(20) DEFAULT NULL COMMENT '关联码表编码',
  `f_rule_id` bigint(20) DEFAULT NULL COMMENT '关联编码规则ID',
  `f_relation_type` VARCHAR(50) DEFAULT NULL COMMENT '关联类型：no/codeTable/codeRule',
  `f_catalog_id` bigint(20) NOT NULL COMMENT '目录ID',
  `f_label_id` bigint(20) DEFAULT NULL COMMENT '数据分级标签ID',
  `f_description` varchar(300) DEFAULT NULL COMMENT '说明',
  `f_version` INT(4) NOT NULL DEFAULT 1 COMMENT '版本号',
  `f_state` INT(2) NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-停用',
  `f_authority_id` varchar(100) DEFAULT NULL COMMENT '权限域',
  `f_department_ids` varchar(350) DEFAULT NULL COMMENT '部门ID',
  `f_third_dept_id` varchar(36) DEFAULT NULL COMMENT '第三方部门ID',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_create_user` varchar(128) DEFAULT NULL COMMENT '创建用户',
  `f_update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `f_update_user` varchar(128) DEFAULT NULL COMMENT '修改用户',
  `f_deleted` bigint(20) NOT NULL DEFAULT 0 COMMENT '逻辑删除标记',
  PRIMARY KEY (`f_id`),
  KEY `idx_de_id` (`f_de_id`),
  KEY `idx_dict_code` (`f_dict_code`),
  KEY `idx_rule_id` (`f_rule_id`),
  KEY `idx_catalog_id` (`f_catalog_id`),
  KEY `idx_state` (`f_state`),
  KEY `idx_std_type` (`f_std_type`),
  KEY `idx_deleted` (`f_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据元基本信息表';
```

**位置**: `migrations/dataelement/raw/t_relation_de_file.sql`

```sql
CREATE TABLE `t_relation_de_file` (
  `f_id` bigint(20) NOT NULL COMMENT '主键',
  `f_de_id` bigint(20) NOT NULL COMMENT '数据元ID',
  `f_file_id` bigint(20) NOT NULL COMMENT '文件ID',
  UNIQUE KEY `uk_deid_fileid` (`f_de_id`,`f_file_id`),
  PRIMARY KEY (`f_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据元-文件关系表';
```

### Go Struct

```go
// DataElement 数据元实体
type DataElement struct {
    Id              int64     `db:"f_id" json:"id"`
    Code            int64     `db:"f_de_id" json:"code"`
    NameEn          string    `db:"f_name_en" json:"nameEn"`
    NameCn          string    `db:"f_name_cn" json:"nameCn"`
    Synonym         string    `db:"f_synonym" json:"synonyms"`
    StdType         int32     `db:"f_std_type" json:"stdType"`
    DataType        int32     `db:"f_data_type" json:"dataType"`
    DataLength      *int      `db:"f_data_length" json:"dataLength"`
    DataPrecision   *int      `db:"f_data_precision" json:"dataPrecision"`
    DictCode        *int64    `db:"f_dict_code" json:"dictCode"`
    RuleId          *int64    `db:"f_rule_id" json:"ruleId"`
    RelationType    string    `db:"f_relation_type" json:"relationType"`
    CatalogId       int64     `db:"f_catalog_id" json:"catalogId"`
    LabelId         *int64    `db:"f_label_id" json:"labelId"`
    Description     string    `db:"f_description" json:"description"`
    Version         int       `db:"f_version" json:"version"`
    State           int32     `db:"f_state" json:"state"`
    AuthorityId     string    `db:"f_authority_id" json:"authorityId"`
    DepartmentIds   string    `db:"f_department_ids" json:"departmentIds"`
    ThirdDeptId     string    `db:"f_third_dept_id" json:"thirdDeptId"`
    CreateTime      string    `db:"f_create_time" json:"createTime"`
    CreateUser      string    `db:"f_create_user" json:"createUser"`
    UpdateTime      string    `db:"f_update_time" json:"updateTime"`
    UpdateUser      string    `db:"f_update_user" json:"updateUser"`
    Deleted         int64     `db:"f_deleted" json:"deleted"`
}

// RelationDeFile 数据元-文件关系
type RelationDeFile struct {
    Id     int64 `db:"f_id" json:"id"`
    DeId   int64 `db:"f_de_id" json:"deId"`
    FileId int64 `db:"f_file_id" json:"fileId"`
}
```

---

## Enums

### 标准分类枚举 (StdTypeEnum)

| 枚举值 | code | 说明 |
|--------|------|------|
| GROUP | 0 | 团体标准 |
| ENTERPRISE | 1 | 企业标准 |
| INDUSTRY | 2 | 行业标准 |
| LOCAL | 3 | 地方标准 |
| NATIONAL | 4 | 国家标准 |
| INTERNATIONAL | 5 | 国际标准 |
| FOREIGN | 6 | 国外标准 |
| OTHER | 99 | 其他标准 |

### 数据类型枚举 (DataTypeEnum)

| 枚举值 | code | 说明 |
|--------|------|------|
| Number | 0 | 数字型 |
| Char | 1 | 字符型 |
| Date | 2 | 日期型 |
| DateTime | 3 | 日期时间型 |
| Boolean | 5 | 布尔型 |
| Decimal | 7 | 高精度型 |
| Binary | 8 | 小数型 |
| Time | 9 | 时间型 |
| Integer | 10 | 整数型 |
| Unknown | 99 | 未知 |

### 启用停用状态枚举 (StateEnum)

| 枚举值 | code | 说明 |
|--------|------|------|
| DISABLE | 0 | 停用 |
| ENABLE | 1 | 启用 |

### 关联类型枚举 (RelationTypeEnum)

| 枚举值 | 说明 |
|--------|------|
| no | 无限制 |
| codeTable | 码表关联 |
| codeRule | 编码规则关联 |

---

## Key Implementation Details

### 版本控制规则

修改以下字段时版本号+1：
- relationType（关联类型）
- dictCode（码表编码）
- ruleId（编码规则ID）
- nameCn（中文名称）
- nameEn（英文名称）

修改以下字段时版本号不变：
- catalogId（目录ID）
- dataType（数据类型）
- dataLength（数据长度）
- dataPrecision（数据精度）
- description（说明）

### 值域计算规则

```go
// CalculateValueRange 计算值域
func CalculateValueRange(de *DataElement, dictData string) string {
    if de.RelationType == "codeTable" && dictData != "" {
        return dictData // 返回码表值JSON: [M,F]
    }

    // 根据数据类型计算
    if de.DataType == 0 || de.DataType == 7 { // Number or Decimal
        length := de.DataLength ?? 38
        precision := de.DataPrecision ?? 0
        min := -int(math.Pow10(float64(length-precision))) + 1
        max := int(math.Pow10(float64(length-precision))) - 1
        return fmt.Sprintf("(%d,%d)", min, max)
    }

    return ""
}
```

### MQ 消息格式

```go
// SendMqMessage 发送MQ消息
func SendMqMessage(ctx context.Context, producer *kafka.Producer, entities []*DataElement, operation string) error {
    mqDto := map[string]interface{}{
        "header": map[string]interface{}{},
        "payload": map[string]interface{}{
            "type": "smart-recommendation-graph",
            "content": map[string]interface{}{
                "type":       operation, // insert/update/delete
                "table_name": "t_data_element_info",
                "entities":   entities,
            },
        },
    }

    jsonData, _ := json.Marshal(mqDto)
    return producer.SendMessage("MQ_MESSAGE_SAILOR", string(jsonData))
}
```

---

## API Contract

### 基础路径

```
/api/standardization/v1/dataelement
```

### API 端点清单

| 序号 | 方法 | 路径 | 功能 | 标签 |
|------|------|------|------|------|
| 1 | POST | `/v1/dataelement` | 创建数据元 | open |
| 2 | POST | `/v1/dataelement/import` | 批量导入数据元 | - |
| 3 | POST | `/v1/dataelement/export` | 导出数据元 | - |
| 4 | POST | `/v1/dataelement/export/{ids}` | 通过ID集合导出 | - |
| 5 | GET | `/v1/dataelement` | 分页检索数据元 | open |
| 6 | GET | `/v1/dataelement/internal/list` | 分页检索（内部） | - |
| 7 | GET | `/v1/dataelement/detail` | 查看数据元详情 | open |
| 8 | GET | `/v1/dataelement/internal/detail` | 查看详情（内部） | - |
| 9 | PUT | `/v1/dataelement/{id}` | 编辑数据元 | open |
| 10 | DELETE | `/v1/dataelement/{ids}` | 删除数据元 | - |
| 11 | PUT | `/v1/dataelement/state/{ids}` | 启用/停用 | open |
| 12 | PUT | `/v1/dataelement/move_catalog/{ids}` | 移动目录 | - |
| 13 | GET | `/v1/dataelement/query/byStdFileCatalog` | 按文件目录检索 | open |
| 14 | GET | `/v1/dataelement/query/byStdFile` | 按文件检索 | open |
| 15 | GET | `/v1/dataelement/query/stdFile` | 查询关联文件 | - |
| 16 | GET | `/v1/dataelement/query/isRepeat` | 检查重名 | - |
| 17 | DELETE | `/v1/dataelement/labelIds/{ids}` | 删除标签 | - |
| 18 | GET | `/v1/dataelement/query/list` | 查询数据元列表 | - |
| 19 | GET | `/v1/dataelement/internal/query/list` | 查询列表（内部） | - |

---

## Testing Strategy

### 单元测试

| 组件 | 覆盖率 | 重点 |
|------|--------|------|
| Model 层 | 90%+ | SQL查询正确性、事务处理 |
| Logic 层 | 85%+ | 业务规则、异常处理 |
| Handler 层 | 70%+ | 参数校验、响应格式 |

### 集成测试

- 创建→查询→编辑→删除完整流程
- 批量操作（导入、删除、状态更新）
- 关联操作（文件关联、目录移动）
- 异常场景（重复检查、权限校验）

---

## Rollout Plan

| Phase | 内容 | 预计工作量 |
|-------|------|------------|
| Phase 0 | API定义、代码生成 | 1天 |
| Phase 1 | Model层实现 | 1.5天 |
| Phase 2 | CRUD接口实现（6个） | 3天 |
| Phase 3 | 状态管理接口（1个） | 0.5天 |
| Phase 4 | 目录管理接口（1个） | 0.5天 |
| Phase 5 | 关联查询接口（3个） | 1.5天 |
| Phase 6 | 导入导出接口（3个） | 2天 |
| Phase 7 | 其他接口（5个） | 1.5天 |
| Phase 8 | 测试和文档 | 2天 |
| **总计** | | **约13.5天** |

---

## Open Questions

无
