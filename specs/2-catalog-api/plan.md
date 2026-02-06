# 目录管理 (catalog-api) Technical Plan

> **Branch**: `2-catalog-api`
> **Spec Path**: `specs/2-catalog-api/`
> **Created**: 2026-02-06
> **Status**: Draft

---

## Summary

从 Java 迁移目录管理模块到 Go-Zero，实现8个API接口（排除2个已废弃接口），支持多级目录树的CRUD操作。技术方案核心决策：

1. **主键兼容性**：使用 Long 类型 BIGINT 主键（与Java保持一致）
2. **表结构保持**：完全复用 Java 表结构（f_ 前缀字段）
3. **接口兼容**：8个API路径、参数、响应格式、错误码100%一致
4. **纯 SQLx 策略**：所有数据库操作使用 SQLx，手工编写 SQL 查询

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
- **原因**：Java实现使用自增Long类型ID
- **影响**：仅影响 catalog 模块
- **风险**：低，Java已稳定运行多年

---

## Technical Context

| Item | Value |
|------|-------|
| **Language** | Go 1.24+ |
| **Framework** | Go-Zero v1.9+ |
| **Storage** | MySQL 8.0 (复用Java表结构) |
| **DB Access** | SQLx (纯 SQL) |
| **Testing** | go test |
| **Common Lib** | idrm-go-base v0.1.0+ |
| **Migration From** | Java Spring Boot + MyBatis |

---

## 通用库 (idrm-go-base)

### 自定义错误码

| 功能 | 范围 | 位置 |
|------|------|------|
| 目录管理 | 30100-30199 | `internal/errorx/codes.go` |

**错误码映射**（与Java保持一致）：

| Java 错误码 | Go 错误码 | 说明 |
|-------------|-----------|------|
| Standardization.Empty | 30101 | 数据不存在 |
| Standardization.MissingParameter | 30102 | 参数缺失 |
| Standardization.InvalidParameter | 30103 | 参数无效 |
| Standardization.OutOfRange | 30104 | 超出范围 |
| Standardization.OperationConflict | 30105 | 操作冲突 |
| Standardization.DATA_EXIST | 30106 | 数据已存在 |
| Standardization.CatalogServiceError | 30100 | 目录服务错误 |

---

## Go-Zero 开发流程

按以下顺序完成技术设计和代码生成：

| Step | 任务 | 方式 | 产出 |
|------|------|------|------|
| 1 | 定义 API 文件 | AI 实现 | `api/doc/catalog/catalog.api` |
| 2 | 生成 Handler/Types | goctl 生成 | `api/internal/handler/catalog/`, `types/` |
| 3 | 实现 Model 接口 | AI 手写 | `model/catalog/catalog/` |
| 4 | 实现 Logic 层 | AI 实现 | `api/internal/logic/catalog/` |
| 5 | 更新路由注册 | 手动 | `api/doc/api.api` |

---

## File Structure

### 代码结构

```
api/internal/
├── handler/catalog/
│   ├── query_tree_handler.go           # goctl 生成
│   ├── query_handler.go
│   ├── create_catalog_handler.go
│   ├── update_catalog_handler.go
│   ├── delete_catalog_handler.go
│   └── query_with_file_handler.go
├── logic/catalog/
│   ├── query_tree_logic.go              # AI 实现
│   ├── query_logic.go
│   ├── create_catalog_logic.go
│   ├── update_catalog_logic.go
│   ├── delete_catalog_logic.go
│   ├── query_with_file_logic.go
│   └── common.go                        # 公共函数（名称校验、树构建等）
├── types/
│   └── types.go                         # goctl 生成（含所有请求/响应类型）
└── svc/
    └── servicecontext.go                 # 手动维护（注入Model）

model/catalog/catalog/
├── interface.go                           # Model 接口定义
├── types.go                               # 数据结构（Catalog）
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

### 依赖服务调用

Logic 层通过 ServiceContext 调用依赖服务：

```go
type ServiceContext struct {
    CatalogModel     model.CatalogModel
    DataElementModel model.DataElementModel  // 数据元服务（删除校验）
    DictModel        model.DictModel         // 字典服务（删除校验）
    RuleModel        model.RuleModel         // 规则服务（删除校验）
    StdFileModel     model.StdFileModel      // 文件服务（删除校验、文件树）
}
```

---

## Interface Definitions

### CatalogModel 接口

```go
type CatalogModel interface {
    // 基础CRUD
    Insert(ctx context.Context, data *Catalog) (int64, error)
    FindOne(ctx context.Context, id int64) (*Catalog, error)
    Update(ctx context.Context, data *Catalog) error
    Delete(ctx context.Context, id int64) error

    // 查询方法
    FindByType(ctx context.Context, catalogType int32) ([]*Catalog, error)
    FindByTypeAndLevel(ctx context.Context, catalogType int32, minLevel int32) ([]*Catalog, error)
    FindByName(ctx context.Context, name string, catalogType int32) ([]*Catalog, error)
    FindByParentId(ctx context.Context, parentId int64) ([]*Catalog, error)
    FindByIds(ctx context.Context, ids []int64) ([]*Catalog, error)

    // 树形结构
    FindTree(ctx context.Context, catalogType int32, rootLevel int32) ([]*Catalog, error)
    FindChildren(ctx context.Context, parentId int64, catalogType int32) ([]*Catalog, error)

    // 批量操作
    DeleteByIds(ctx context.Context, ids []int64) error
}
```

---

## Data Model

### DDL（复用Java结构）

**位置**: `migrations/catalog/raw/t_de_catalog_info.sql`

```sql
CREATE TABLE `t_de_catalog_info` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '目录唯一标识',
  `f_catalog_name` varchar(20) NOT NULL COMMENT '目录名称',
  `f_description` varchar(300) DEFAULT NULL COMMENT '目录说明',
  `f_level` INT(4) NOT NULL DEFAULT 1 COMMENT '目录级别',
  `f_parent_id` bigint(20) DEFAULT NULL COMMENT '父级标识',
  `f_type` INT(2) NOT NULL COMMENT '目录类型',
  `f_authority_id` varchar(100) DEFAULT NULL COMMENT '权限域',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_create_user` varchar(128) DEFAULT NULL COMMENT '创建用户',
  `f_update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `f_update_user` varchar(128) DEFAULT NULL COMMENT '修改用户',
  `f_deleted` bigint(20) NOT NULL DEFAULT 0 COMMENT '逻辑删除标记',
  PRIMARY KEY (`f_id`),
  KEY `idx_parent_id` (`f_parent_id`),
  KEY `idx_type_deleted` (`f_type`,`f_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='目录信息表';
```

### Go Struct

```go
// Catalog 目录实体
type Catalog struct {
    Id           int64     `db:"f_id" json:"id"`
    CatalogName  string    `db:"f_catalog_name" json:"catalogName"`
    Description  string    `db:"f_description" json:"description"`
    Level        int32     `db:"f_level" json:"level"`
    ParentId     int64     `db:"f_parent_id" json:"parentId"`
    Type         int32     `db:"f_type" json:"type"`
    AuthorityId  string    `db:"f_authority_id" json:"authorityId"`
    CreateTime   time.Time `db:"f_create_time" json:"createTime"`
    CreateUser   string    `db:"f_create_user" json:"createUser"`
    UpdateTime   time.Time `db:"f_update_time" json:"updateTime"`
    UpdateUser   string    `db:"f_update_user" json:"updateUser"`
    Deleted      int64     `db:"f_deleted" json:"deleted"`
    Children     []*Catalog `db:"-" json:"children,omitempty"` // 递归子级
    Count        int32     `db:"-" json:"count,omitempty"`     // 数据统计
    HaveChildren bool      `db:"-" json:"haveChildren"`        // 是否有子级
}

// CatalogTreeNodeVo 目录树节点
type CatalogTreeNodeVo struct {
    Id           int64                `json:"id"`
    CatalogName  string               `json:"catalogName"`
    Level        int32                `json:"level"`
    ParentId     int64                `json:"parentId"`
    Type         int32                `json:"type"`
    Children     []*CatalogTreeNodeVo `json:"children,omitempty"`
    Count        int32                `json:"count,omitempty"`
    HaveChildren bool                 `json:"haveChildren"`
}

// CatalogWithFileVo 目录及文件树节点
type CatalogWithFileVo struct {
    Id           int64                `json:"id"`
    CatalogName  string               `json:"catalogName"`
    Level        int32                `json:"level"`
    ParentId     int64                `json:"parentId"`
    Type         int32                `json:"type"`
    Children     []*CatalogWithFileVo `json:"children,omitempty"`
    Files        []*FileCountVo       `json:"files,omitempty"`
    HaveChildren bool                 `json:"haveChildren"`
}

// FileCountVo 文件计数
type FileCountVo struct {
    FileId    int64  `json:"fileId"`
    FileName  string `json:"fileName"`
    CatalogId int64  `json:"catalogId"`
}

// CatalogListByFileVo 目录及文件列表响应
type CatalogListByFileVo struct {
    Catalogs []*Catalog `json:"catalogs"`
    Files    []*FileCountVo `json:"files"`
}
```

### 枚举定义

```go
// CatalogTypeEnum 目录类型枚举
const (
    CatalogTypeDataElement int32 = 1 // 数据元目录
    CatalogTypeDict        int32 = 2 // 码表目录
    CatalogTypeValueRule   int32 = 3 // 编码规则目录
    CatalogTypeFile        int32 = 4 // 文件目录
)

// 获取目录类型字符串
func GetCatalogTypeString(catalogType int32) string {
    switch catalogType {
    case CatalogTypeDataElement:
        return "DataElement"
    case CatalogTypeDict:
        return "Dict"
    case CatalogTypeValueRule:
        return "ValueRule"
    case CatalogTypeFile:
        return "File"
    default:
        return "Other"
    }
}

// 获取目录类型int
func GetCatalogTypeInt(catalogType string) int32 {
    switch catalogType {
    case "DataElement":
        return CatalogTypeDataElement
    case "Dict":
        return CatalogTypeDict
    case "ValueRule":
        return CatalogTypeValueRule
    case "File":
        return CatalogTypeFile
    default:
        return CatalogTypeDataElement
    }
}
```

---

## API Contract

**位置**: `api/doc/catalog/catalog.api`

### 接口清单（8个，排除2个废弃接口）

| 序号 | 方法 | 路径 | 功能 |
|------|------|------|------|
| 1 | GET | `/query_tree` | 查询目录树 |
| 2 | GET | `/query` | 检索目录 |
| 3 | POST | `/` | 创建目录 |
| 4 | PUT | `/{id}` | 修改目录 |
| 5 | DELETE | `/{id}` | 删除目录 |
| 6 | GET | `/query/with_file` | 查询目录及文件树 |

### 完整 API 定义

```api
syntax = "v1"

import "../base.api"

type (
    // ==================== 请求类型 ====================

    // CreateCatalogReq 创建目录请求
    CreateCatalogReq {
        CatalogName string `json:"catalogName" validate:"required,max=20"`
        ParentId    int64  `json:"parentId" validate:"required"`
        Type        int32  `json:"type" validate:"required,min=1,max=4"`
        Description string `json:"description,optional" validate:"max=300"`
    }

    // UpdateCatalogReq 修改目录请求
    UpdateCatalogReq {
        CatalogName string `json:"catalogName" validate:"required,max=20"`
        Description string `json:"description,optional" validate:"max=300"`
    }

    // QueryTreeReq 查询目录树请求
    QueryTreeReq {
        Type int32 `form:"type,optional" validate:"min=1,max=4"`
        Id   int64 `form:"id,optional"`
    }

    // QueryReq 检索目录请求
    QueryReq {
        Type    int32  `form:"type" validate:"required,min=1,max=4"`
        Keyword string `form:"keyword,optional"`
    }

    // QueryWithFileReq 查询目录及文件树请求
    QueryWithFileReq {
        Keyword string `form:"keyword,optional"`
    }

    // ==================== 响应类型 ====================

    // CatalogResp 目录响应
    CatalogResp {
        Id          int64          `json:"id"`
        CatalogName string         `json:"catalogName"`
        Description string         `json:"description,omitempty"`
        Level       int32          `json:"level"`
        ParentId    int64          `json:"parentId"`
        Type        int32          `json:"type"`
        Children    []*CatalogResp `json:"children,omitempty"`
        Count       int32          `json:"count,omitempty"`
        HaveChildren bool          `json:"haveChildren"`
    }

    // CatalogWithFileResp 目录及文件响应
    CatalogWithFileResp {
        Id           int64                  `json:"id"`
        CatalogName  string                 `json:"catalogName"`
        Level        int32                  `json:"level"`
        ParentId     int64                  `json:"parentId"`
        Type         int32                  `json:"type"`
        Children     []*CatalogWithFileResp `json:"children,omitempty"`
        Files        []*FileCountVo         `json:"files,omitempty"`
        HaveChildren bool                   `json:"haveChildren"`
    }

    // FileCountVo 文件计数响应
    FileCountVo {
        FileId    int64  `json:"fileId"`
        FileName  string `json:"fileName"`
        CatalogId int64  `json:"catalogId"`
    }

    // CatalogListByFileResp 目录及文件列表响应
    CatalogListByFileResp {
        Catalogs []*CatalogInfoVo `json:"catalogs"`
        Files    []*FileCountVo   `json:"files"`
    }

    // CatalogInfoVo 目录信息响应（平铺列表）
    CatalogInfoVo {
        Id          int64  `json:"id"`
        CatalogName string `json:"catalogName"`
        Level       int32  `json:"level"`
        ParentId    int64  `json:"parentId"`
        Type        int32  `json:"type"`
    }
)

// ==================== API 定义 ====================

@server(
    prefix: /api/standardization/v1
    group: catalog
)
service api {
    @doc "查询目录树"
    @handler QueryTree
    get /catalog/query_tree (QueryTreeReq) returns (CatalogResp)

    @doc "检索目录"
    @handler Query
    get /catalog/query (QueryReq) returns ([]CatalogInfoVo)

    @doc "创建目录"
    @handler CreateCatalog
    post /catalog (CreateCatalogReq) returns (CatalogResp)

    @doc "修改目录"
    @handler UpdateCatalog
    put /catalog/:id (UpdateCatalogReq) returns (EmptyResp)

    @doc "删除目录"
    @handler DeleteCatalog
    delete /catalog/:id returns (EmptyResp)

    @doc "查询目录及文件树"
    @handler QueryWithFile
    get /catalog/query/with_file (QueryWithFileReq) returns (CatalogListByFileResp)
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
| P1 | 创建目录、查询目录树、删除目录 |
| P1 | 目录名称格式校验、同级名称唯一性 |
| P1 | 父目录存在性校验、级别继承 |
| P1 | 根目录保护（不允许修改和删除） |
| P2 | 删除前数据校验、递归删除 |
| P2 | 目录检索、文件树查询 |
| P3 | 边界条件（level=255、空关键字等） |

---

## Key Implementation Details

### 1. 目录名称校验（Logic/common.go）

```go
// ValidateCatalogName 校验目录名称格式
func ValidateCatalogName(name string) error {
    name = strings.TrimSpace(name)

    // 长度校验
    if len(name) == 0 {
        return errorx.CatalogNameEmpty()
    }
    if len(name) > 20 {
        return errorx.CatalogNameTooLong()
    }

    // 格式校验：中英文、数字、_、-
    matched, _ := regexp.MatchString(`^[\p{Han}a-zA-Z0-9_-]+$`, name)
    if !matched {
        return errorx.CatalogNameInvalidFormat()
    }

    // 首字符校验：不能以_或-开头
    if strings.HasPrefix(name, "_") || strings.HasPrefix(name, "-") {
        return errorx.CatalogNameInvalidPrefix()
    }

    return nil
}
```

### 2. 树形结构构建（Logic/common.go）

```go
// BuildTree 构建目录树
func BuildTree(catalogs []*Catalog, rootLevel int32) []*Catalog {
    // 按level分组
    levelMap := make(map[int32][]*Catalog)
    for _, c := range catalogs {
        levelMap[c.Level] = append(levelMap[c.Level], c)
    }

    // 从rootLevel开始递归构建
    roots := levelMap[rootLevel]
    for _, root := range roots {
        buildChildren(root, levelMap)
    }

    return roots
}

func buildChildren(parent *Catalog, levelMap map[int32][]*Catalog) {
    children := levelMap[parent.Level+1]
    if len(children) == 0 {
        return
    }

    for _, child := range children {
        if child.ParentId == parent.Id {
            parent.Children = append(parent.Children, child)
            parent.HaveChildren = true
            buildChildren(child, levelMap)
        }
    }
}
```

### 3. 递归删除

```go
// DeleteWithChildren 递归删除目录及所有子级
func DeleteWithChildren(ctx context.Context, model CatalogModel, catalogId int64) error {
    // 1. 获取所有子级ID
    childIds, err := getAllChildIds(ctx, model, catalogId)
    if err != nil {
        return err
    }

    // 2. 批量删除
    allIds := append(childIds, catalogId)
    return model.DeleteByIds(ctx, allIds)
}

func getAllChildIds(ctx context.Context, model CatalogModel, parentId int64) ([]int64, error) {
    children, err := model.FindByParentId(ctx, parentId)
    if err != nil {
        return nil, err
    }

    ids := make([]int64, 0)
    for _, child := range children {
        ids = append(ids, child.Id)
        // 递归获取子级的子级
        childChildIds, err := getAllChildIds(ctx, model, child.Id)
        if err != nil {
            return nil, err
        }
        ids = append(ids, childChildIds...)
    }

    return ids, nil
}
```

### 4. 删除前数据校验

```go
// CheckCatalogDelete 校验目录是否可以删除
func CheckCatalogDelete(ctx context.Context, catalog *Catalog, svcCtx *svc.ServiceContext) error {
    // 1. 根目录不允许删除
    if catalog.Level <= 1 {
        return errorx.CannotDeleteRootCatalog()
    }

    // 2. 根据类型校验是否存在数据
    switch catalog.Type {
    case CatalogTypeDataElement:
        // TODO: 调用 DataElementModel 检查是否存在数据元
    case CatalogTypeDict:
        // TODO: 调用 DictModel 检查是否存在码表
    case CatalogTypeValueRule:
        // TODO: 调用 RuleModel 检查是否存在规则
    case CatalogTypeFile:
        // TODO: 调用 StdFileModel 检查是否存在文件
    }

    return nil
}
```

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-02-06 | - | 初始版本 |
