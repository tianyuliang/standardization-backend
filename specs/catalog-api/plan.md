# 目录管理 API (Catalog API) 技术方案

> **Branch**: `feature/catalog-api`
> **Spec Path**: `specs/catalog-api/`
> **Created**: 2025-01-21
> **Status**: Draft

---

## Summary

将 Java Spring Boot 目录管理 API 转写为 Go-Zero 微服务架构，实现 6 个核心 CRUD 接口。技术方案保持与 Java 版本 100% 兼容，包括 API 路由、响应格式、业务逻辑和校验规则。使用 GORM 实现数据访问层，通过预留接口桩模块处理外部服务依赖（数据元、码表、编码规则、文件关联检查）。

---

## Technical Context

| Item | Value |
|------|-------|
| **Language** | Go 1.24+ |
| **Framework** | Go-Zero v1.9+ |
| **Storage** | MySQL 8.0 |
| **Cache** | Redis 7.0 (暂不使用) |
| **ORM** | GORM (主要) + SQLx (可选) |
| **Testing** | go test + 表驱动测试 |
| **Common Lib** | idrm-go-base v0.1.0+ |

---

## 通用库 (idrm-go-base)

**安装**:
```bash
go get github.com/jinguoxing/idrm-go-base@latest
```

### 模块初始化

| 模块 | 初始化方式 | 位置 |
|------|-----------|------|
| validator | `validator.Init()` 在 main.go | `api/main.go` |
| telemetry | `telemetry.Init(cfg)` 在 main.go | `api/main.go` |
| response | `httpx.SetErrorHandler(response.ErrorHandler)` | `api/main.go` |
| middleware | `rest.WithMiddlewares(...)` | `api/main.go` |

### 自定义错误码

| 功能 | 范围 | 位置 |
|------|------|------|
| 目录管理 | 30100-30199 | `internal/errorx/codes.go` |

**错误码定义**:
```go
package errorx

import "github.com/jinguoxing/idrm-go-base/errorx"

const (
    // 目录模块错误码 30100-30199
    ErrCatalogNotExist     = 30100
    ErrCatalogNameInvalid  = 30101
    ErrCatalogLevelInvalid = 30102
    ErrCatalogHasChildren  = 30103
)
```

### 第三方库确认

| 库 | 原因 | 确认状态 |
|----|------|----------|
| github.com/google/uuid | UUID v7 主键生成（但为兼容 Java 使用 string 存储 ID） | ✅ 已确认 |
| github.com/go-playground/validator | 参数校验（通过 validator 模块封装） | ✅ 已确认 |

---

## Go-Zero 开发流程

| Step | 任务 | 方式 | 产出 |
|------|------|------|------|
| 1 | 定义 API 文件 | AI 实现 | `api/doc/catalog/catalog.api` |
| 2 | 更新入口文件 | AI 手写 | `api/doc/api.api` |
| 3 | 生成 Handler/Types | goctl | `api/internal/handler/`, `types/` |
| 4 | 实现 Model 层 | AI 手写 | `model/catalog/catalog/` |
| 5 | 实现 Logic 层 | AI 实现 | `api/internal/logic/catalog/` |
| 6 | 实现桩模块 | AI 手写 | `model/catalog/stub/` |

**goctl 命令**:
```bash
# 在 api/doc/api.api 中 import "catalog/catalog.api"
# 然后执行 goctl 生成代码
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

---

## File Structure

### 文件产出清单

| 序号 | 文件 | 生成方式 | 位置 |
|------|------|----------|------|
| 1 | API 文件 | AI 手写 | `api/doc/catalog/catalog.api` |
| 2 | Handler | goctl 生成 | `api/internal/handler/catalog/` |
| 3 | Types | goctl 生成 | `api/internal/types/types.go` |
| 4 | Logic | AI 实现 | `api/internal/logic/catalog/` |
| 5 | Model | AI 实现 | `model/catalog/catalog/` |
| 6 | 桩模块 | AI 手写 | `model/catalog/stub/` |
| 7 | 错误码 | AI 手写 | `api/internal/errorx/codes.go` |

### 代码结构

```
api/internal/
├── handler/catalog/
│   ├── querytreehandler.go           # goctl 生成
│   ├── queryhandler.go               # goctl 生成
│   ├── createhandler.go              # goctl 生成
│   ├── updatehandler.go              # goctl 生成
│   ├── deletehandler.go              # goctl 生成
│   └── querywithfilehandler.go       # goctl 生成
├── logic/catalog/
│   ├── querytreelogic.go             # AI 实现
│   ├── querylogic.go                 # AI 实现
│   ├── createlogic.go                # AI 实现
│   ├── updatelogic.go                # AI 实现
│   ├── deletelogic.go                # AI 实现
│   └── querywithfilelogic.go         # AI 实现
├── types/
│   └── types.go                      # goctl 生成
├── errorx/
│   └── codes.go                      # 自定义错误码
└── svc/
    └── servicecontext.go             # 服务上下文

model/catalog/catalog/
├── interface.go                      # Model 接口定义
├── types.go                          # 数据结构
├── vars.go                           # 常量和错误
├── factory.go                        # ORM 工厂
├── gorm_dao.go                       # GORM 实现
└── sqlx_model.go                     # SQLx 实现（可选）

model/catalog/stub/
└── external_checker.go               # 外部依赖检查桩模块
```

---

## Architecture Overview

遵循 IDRM 分层架构：

```
HTTP Request → Handler → Logic → Model → Database
```

| 层级 | 职责 | 最大行数 | 示例 |
|------|------|----------|------|
| Handler | 解析参数、格式化响应、调用 Logic | 30 | 参数绑定、错误码转换 |
| Logic | 业务逻辑、事务管理、调用 Model | 50 | 树结构构建、循环检测 |
| Model | 数据访问、SQL 执行 | 50 | CRUD 操作、复杂查询 |

### 层级交互

```
┌─────────────────────────────────────────────────────┐
│                      Handler                         │
│  - 参数解析 (httpx.Parse)                            │
│  - 调用 Logic                                       │
│  - 响应格式化 (response.Success)                     │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────┐
│                      Logic                           │
│  - 业务规则校验 (正则、层级、重复检查)                │
│  - 树结构构建 (递归查询)                             │
│  - 调用 Model 执行 CRUD                              │
│  - 调用桩模块检查外部依赖                            │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────┐
│                      Model                            │
│  - GORM/SQLx 数据库操作                              │
│  - 事务管理                                          │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────┐
│                   MySQL Database                     │
│              t_de_catalog_info                        │
└─────────────────────────────────────────────────────┘
```

---

## Interface Definitions

### Model 接口

```go
// CatalogModel 目录数据访问接口
type CatalogModel interface {
    // Insert 创建目录
    Insert(ctx context.Context, data *Catalog) (*Catalog, error)

    // FindOne 根据 ID 查询单个目录
    FindOne(ctx context.Context, id string) (*Catalog, error)

    // FindByParent 查询父目录下的所有子目录
    FindByParent(ctx context.Context, parentId string) ([]*Catalog, error)

    // FindByType 根据类型查询所有目录
    FindByType(ctx context.Context, catalogType int32) ([]*Catalog, error)

    // FindByName 模糊查询目录名称
    FindByName(ctx context.Context, name string, catalogType int32) ([]*Catalog, error)

    // Update 更新目录
    Update(ctx context.Context, data *Catalog) error

    // Delete 删除目录（单条）
    Delete(ctx context.Context, id string) error

    // DeleteBatch 批量删除
    DeleteBatch(ctx context.Context, ids []string) error

    // WithTx 创建事务
    WithTx(tx interface{}) CatalogModel

    // Trans 事务执行
    Trans(ctx context.Context, fn func(ctx context.Context, model CatalogModel) error) error
}
```

### 外部依赖检查接口（桩模块）

```go
// ExternalChecker 外部依赖检查接口（桩模块）
type ExternalChecker interface {
    // CheckDataElement 检查目录下是否存在数据元
    CheckDataElement(ctx context.Context, catalogId string) (bool, error)

    // CheckDict 检查目录下是否存在码表
    CheckDict(ctx context.Context, catalogId string) (bool, error)

    // CheckRule 检查目录下是否存在编码规则
    CheckRule(ctx context.Context, catalogId string) (bool, error)

    // CheckFile 检查目录下是否存在文件
    CheckFile(ctx context.Context, catalogId string) (bool, error)
}

// StubExternalChecker 桩模块实现（当前阶段所有返回 false）
type StubExternalChecker struct{}

func (s *StubExternalChecker) CheckDataElement(ctx context.Context, catalogId string) (bool, error) {
    return false, nil // 桩实现：暂不检查
}
// ... 其他方法类似
```

---

## Data Model

### DDL

> **注意**: 数据库表 `t_de_catalog_info` 已存在，DDL 位于 `migrations/` 目录，无需重新创建。

**表结构参考**:
```sql
-- 目录信息表（已存在）
CREATE TABLE `t_de_catalog_info` (
    `f_id` BIGINT(20) NOT NULL COMMENT '目录唯一标识（雪花算法）',
    `f_catalog_name` VARCHAR(20) NOT NULL COMMENT '目录名称',
    `f_description` VARCHAR(255) DEFAULT NULL COMMENT '目录说明',
    `f_level` TINYINT(3) UNSIGNED NOT NULL COMMENT '目录级别（1-255，1为根目录）',
    `f_parent_id` BIGINT(20) NOT NULL COMMENT '父级标识',
    `f_type` TINYINT(3) UNSIGNED NOT NULL COMMENT '目录类型（1=数据元,2=码表,3=编码规则,4=文件）',
    `f_authority_id` BIGINT(20) DEFAULT NULL COMMENT '权限域（预留字段）',
    PRIMARY KEY (`f_id`),
    KEY `idx_type_level` (`f_type`, `f_level`),
    KEY `idx_parent_id` (`f_parent_id`),
    KEY `idx_type_parent_name` (`f_type`, `f_parent_id`, `f_catalog_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据标准目录信息表';
```

### Go Struct

```go
// Catalog 目录实体
type Catalog struct {
    Id          string    `gorm:"column:f_id;primaryKey;size:20" json:"id"`                           // 目录唯一标识（string 存储以兼容 JSON）
    CatalogName string    `gorm:"column:f_catalog_name;size:20;not null" json:"catalogName"`            // 目录名称
    Description string    `gorm:"column:f_description;size:255" json:"description"`                     // 目录说明
    Level       int32     `gorm:"column:f_level;unsigned;not null" json:"level"`                        // 目录级别
    ParentId    string    `gorm:"column:f_parent_id;size:20;not null;index:idx_parent_id" json:"parentId"` // 父级标识
    Type        int32     `gorm:"column:f_type;unsigned;not null;index:idx_type_level" json:"type"`         // 目录类型
    AuthorityId *string   `gorm:"column:f_authority_id;size:20" json:"authorityId,omitempty"`            // 权限域（可选）
    CreatedAt   time.Time `gorm:"column:f_created_at;not null" json:"createdAt,omitempty"`              // 创建时间
    UpdatedAt   time.Time `gorm:"column:f_updated_at;not null" json:"updatedAt,omitempty"`              // 更新时间
}

// TableName 指定表名
func (Catalog) TableName() string {
    return "t_de_catalog_info"
}
```

### 请求/响应类型

```go
// CatalogTreeNodeVo 目录树节点
type CatalogTreeNodeVo struct {
    Id          string                `json:"id"`
    CatalogName string                `json:"catalogName"`
    Description string                `json:"description"`
    Level       int32                 `json:"level"`
    ParentId    string                `json:"parentId"`
    Type        int32                 `json:"type"`
    AuthorityId *string               `json:"authorityId,omitempty"`
    Count       *int32                `json:"count,omitempty"`        // 目录下的数据量
    Children    []*CatalogTreeNodeVo  `json:"children,omitempty"`     // 递归子目录
    HaveChildren bool                 `json:"haveChildren"`           // 是否有子节点
}

// CatalogInfoVo 目录信息
type CatalogInfoVo struct {
    Id          string  `json:"id"`
    CatalogName string  `json:"catalogName"`
    Description string  `json:"description"`
    Level       int32   `json:"level"`
    ParentId    string  `json:"parentId"`
    Type        int32   `json:"type"`
    AuthorityId *string `json:"authorityId,omitempty"`
    Count       *int32  `json:"count,omitempty"`
}

// FileCountVo 文件统计
type FileCountVo struct {
    FileId     string `json:"fileId"`
    FileName   string `json:"fileName"`
    CatalogId  string `json:"catalogId"`
}

// CatalogListByFileVo 目录文件列表
type CatalogListByFileVo struct {
    Catalogs []*Catalog     `json:"catalogs"`
    Files    []*FileCountVo `json:"files"`
}
```

---

## API Contract

### API 定义文件

**位置**: `api/doc/catalog/catalog.api`

```api
syntax = "v1"

// 导入通用类型
import "../base.api"

// ============================================
// 目录管理 API 请求/响应类型
// ============================================

type (
    // QueryTreeReq 查询目录树请求
    QueryTreeReq {
        Type int32  `form:"type,required"`                  // 目录类型 1=数据元, 2=码表, 3=编码规则, 4=文件
        Id   *int64 `form:"id,optional"`                   // 目录ID（可选，指定时返回该节点子树）
    }

    // QueryTreeResp 查询目录树响应
    QueryTreeResp {
        Code        string             `json:"code"`
        Description string             `json:"description"`
        Data        *CatalogTreeNodeVo `json:"data,omitempty"`
    }

    // QueryReq 按名称检索目录请求
    QueryReq {
        Type        int32  `form:"type,required"`           // 目录类型
        CatalogName *string `form:"catalog_name,optional"`  // 目录名称（模糊查询）
    }

    // QueryResp 查询目录响应
    QueryResp {
        Code        string           `json:"code"`
        Description string           `json:"description"`
        Data        []*CatalogInfoVo `json:"data,omitempty"`
    }

    // CreateReq 创建目录请求
    CreateReq {
        CatalogName string `json:"catalogName,required,max=20"` // 目录名称（1-20字符）
        ParentId    string `json:"parentId,required"`            // 父目录ID
        Description string `json:"description,optional,max=255"`  // 目录说明（可选）
    }

    // CreateResp 创建目录响应
    CreateResp {
        Code        string `json:"code"`
        Description string `json:"description"`
    }

    // UpdateReq 修改目录请求
    UpdateReq {
        Id          string `path:"id,required"`                // 目录ID（路径参数）
        CatalogName string `json:"catalogName,required,max=20"` // 目录名称
        ParentId    string `json:"parentId,required"`            // 父目录ID
        Description string `json:"description,optional,max=255"`  // 目录说明（可选）
    }

    // UpdateResp 修改目录响应
    UpdateResp {
        Code        string `json:"code"`
        Description string `json:"description"`
    }

    // DeleteReq 删除目录请求
    DeleteReq {
        Id string `path:"id,required"` // 目录ID（路径参数）
    }

    // DeleteResp 删除目录响应
    DeleteResp {
        Code        string `json:"code"`
        Description string `json:"description"`
    }

    // QueryWithFileReq 查询目录与文件树列表请求
    QueryWithFileReq {
        CatalogName *string `form:"catalog_name,optional"` // 目录名称（模糊查询）
    }

    // QueryWithFileResp 查询目录与文件树列表响应
    QueryWithFileResp {
        Code        string                `json:"code"`
        Description string                `json:"description"`
        Data        *CatalogListByFileVo   `json:"data,omitempty"`
    }

    // ============================================
    // 目录节点类型定义
    // ============================================

    // CatalogTreeNodeVo 目录树节点
    CatalogTreeNodeVo {
        Id          string              `json:"id"`
        CatalogName string              `json:"catalogName"`
        Description string              `json:"description"`
        Level       int32               `json:"level"`
        ParentId    string              `json:"parentId"`
        Type        int32               `json:"type"`
        AuthorityId string              `json:"authorityId,omitempty"`
        Count       int32               `json:"count,omitempty"`
        Children    []*CatalogTreeNodeVo `json:"children,omitempty"`
        HaveChildren bool               `json:"haveChildren"`
    }

    // CatalogInfoVo 目录信息
    CatalogInfoVo {
        Id          string  `json:"id"`
        CatalogName string  `json:"catalogName"`
        Description string  `json:"description"`
        Level       int32   `json:"level"`
        ParentId    string  `json:"parentId"`
        Type        int32   `json:"type"`
        AuthorityId string  `json:"authorityId,omitempty"`
        Count       int32   `json:"count,omitempty"`
    }

    // Catalog 目录实体
    Catalog {
        Id          string `json:"id"`
        CatalogName string `json:"catalogName"`
        Description string `json:"description"`
        Level       int32  `json:"level"`
        ParentId    string `json:"parentId"`
        Type        int32  `json:"type"`
        AuthorityId string `json:"authorityId,omitempty"`
    }

    // FileCountVo 文件统计
    FileCountVo {
        FileId    string `json:"fileId"`
        FileName  string `json:"fileName"`
        CatalogId string `json:"catalogId"`
    }

    // CatalogListByFileVo 目录文件列表
    CatalogListByFileVo {
        Catalogs []*Catalog      `json:"catalogs"`
        Files    []*FileCountVo   `json:"files"`
    }
)

// ============================================
// 目录管理 API 服务定义
// ============================================

@server(
    prefix: /api/v1/catalog
    group: catalog
    middleware: Validator
)
service standardization-backend-api {
    @doc "查询目录树"
    @handler querytree
    get /query_tree (QueryTreeReq) returns (QueryTreeResp)

    @doc "按名称检索目录"
    @handler query
    get /query (QueryReq) returns (QueryResp)

    @doc "创建目录"
    @handler create
    post / (CreateReq) returns (CreateResp)

    @doc "修改目录"
    @handler update
    put /:id (UpdateReq) returns (UpdateResp)

    @doc "删除目录"
    @handler delete
    delete /:id (DeleteReq) returns (DeleteResp)

    @doc "查询目录与文件树列表"
    @handler querywithfile
    get /query/with_file (QueryWithFileReq) returns (QueryWithFileResp)
}
```

### 更新入口文件

**位置**: `api/doc/api.api` - 添加 import 语句

```api
syntax = "v1"

info (
    title:   "standardization-backend API"
    desc:    "standardization-backend API 服务"
    version: "v1"
)

// 引入通用类型
import "base.api"

// 引入目录管理模块
import "catalog/catalog.api"

// 定义 service 名称为 "api"
@server (
    prefix: /api/v1
)
service api {
    @handler HealthCheck
    get /health returns (HealthResp)
}

type HealthResp {
    Status string `json:"status"`
}
```

---

## Business Logic Implementation

### 1. 查询目录树 (QueryTree)

**业务逻辑**:
1. 校验 type 参数在 [1,2,3,4] 范围内
2. 如果提供了 id，查询该目录并返回其子树
3. 如果未提供 id，查询该类型所有目录，构建完整树
4. 树结构构建：递归查询子节点，按 level 分层

**关键算法**:
```go
// 构建树结构（递归）
func buildTree(catalogs []*Catalog, rootLevel int32) []*CatalogTreeNodeVo {
    // 1. 创建映射：id -> node
    nodeMap := make(map[string]*CatalogTreeNodeVo)
    for _, c := range catalogs {
        nodeMap[c.Id] = &CatalogTreeNodeVo{...}
    }

    // 2. 构建树：children 添加到父节点
    roots := []*CatalogTreeNodeVo{}
    for _, node := range nodeMap {
        if node.Level == rootLevel {
            roots = append(roots, node)
        } else {
            parent := nodeMap[node.ParentId]
            parent.Children = append(parent.Children, node)
        }
    }
    return roots
}
```

### 2. 按名称检索 (Query)

**业务逻辑**:
1. 校验 type 参数
2. 目录名称 SQL 特殊字符转义
3. 模糊查询：`WHERE f_catalog_name LIKE %name%` (小写不敏感)
4. 仅返回 level > 1 的目录（排除根目录）

### 3. 创建目录 (Create)

**业务逻辑**:
1. 校验目录名称格式：正则 `^[\u4e00-\u9fa5a-zA-Z0-9][\u4e00-\u9fa5a-zA-Z0-9_-]{0,19}$`
2. 校验父目录存在且 level < 255
3. 继承父目录的 type，设置 level = 父目录 level + 1
4. 检查同级目录名称唯一性
5. 插入数据库

### 4. 修改目录 (Update)

**业务逻辑**:
1. 校验目录存在
2. 校验不是根目录（level <= 1 不允许修改）
3. 校验新父目录存在且 level < 255
4. 循环检测：新父目录不能是自身子目录
5. 类型一致性：新父目录 type 必须与当前目录一致
6. 检查同级名称唯一性（排除自身）
7. 更新数据库

### 5. 删除目录 (Delete)

**业务逻辑**:
1. 校验目录存在
2. 校验不是根目录（level <= 1 不允许删除）
3. **递归获取所有子目录 ID**
4. **调用桩模块检查关联数据**（当前阶段跳过）
5. 批量删除目录及子目录

**递归获取子目录算法**:
```go
func getAllChildIds(ctx context.Context, model CatalogModel, catalogId string, ids *[]string) error {
    // 1. 查询直接子目录
    children, err := model.FindByParent(ctx, catalogId)
    if err != nil {
        return err
    }

    // 2. 递归处理子目录
    for _, child := range children {
        *ids = append(*ids, child.Id)
        if err := getAllChildIds(ctx, model, child.Id, ids); err != nil {
            return err
        }
    }
    return nil
}
```

### 6. 查询目录与文件树 (QueryWithFile)

**业务逻辑**:
1. 模糊查询目录名称
2. 查询文件列表（当前阶段返回空列表，待文件模块实现）
3. 组装 CatalogListByFileVo 响应

---

## Validation Rules

### 目录名称正则

```go
const catalogNameRegex = `^[\u4e00-\u9fa5a-zA-Z0-9][\u4e00-\u9fa5a-zA-Z0-9_-]{0,19}$`

// 验证函数
func ValidateCatalogName(name string) bool {
    matched, _ := regexp.MatchString(catalogNameRegex, name)
    return matched
}
```

### 参数校验（通过 validator 标签）

| 参数 | 校验规则 | 错误码 |
|------|----------|--------|
| type | required, range=[1:4] | ErrInvalidParam |
| catalog_name | max=20, regex | ErrInvalidParam |
| parent_id | required | ErrMissingParam |
| level | range=[1:255] | ErrOutOfRange |

---

## Testing Strategy

### 单元测试

| 测试类型 | 覆盖内容 | 方法 |
|----------|----------|------|
| Model 测试 | CRUD 操作、事务 | Mock DB, 表驱动测试 |
| Logic 测试 | 业务规则、树结构 | Mock Model, 表驱动测试 |
| Handler 测试 | HTTP 接口 | httptest, 表驱动测试 |

**测试覆盖率目标**: > 80%

### 测试用例示例

```go
func TestCreateCatalogLogic(t *testing.T) {
    tests := []struct {
        name    string
        req     *types.CreateReq
        wantErr bool
        errCode string
    }{
        {
            name: "正常创建",
            req: &types.CreateReq{
                CatalogName: "测试目录",
                ParentId:    "root-id",
            },
            wantErr: false,
        },
        {
            name: "目录名称为空",
            req: &types.CreateReq{
                CatalogName: "",
                ParentId:    "root-id",
            },
            wantErr: true,
            errCode: errorx.ErrInvalidParam,
        },
        {
            name: "目录名称格式错误",
            req: &types.CreateReq{
                CatalogName: "-invalid",
                ParentId:    "root-id",
            },
            wantErr: true,
            errCode: errorx.ErrInvalidParam,
        },
    }
    // ... 表驱动测试实现
}
```

### API 契约测试

**目标**: 与 Java 版本响应 100% 一致

```bash
# 运行契约测试
go test ./api/tests/contract -v -catalog-contract
```

---

## Error Handling

### 错误码映射

| Java ErrorCodeEnum | Go errorx.Error Code | HTTP Status |
|-------------------|---------------------|-------------|
| MissingParameter | errorx.ErrMissingParam | 400 |
| InvalidParameter | errorx.ErrInvalidParam | 400 |
| Empty | errorx.ErrNotFound | 404 |
| OutOfRange | errorx.ErrOutOfRange | 400 |
| OperationConflict | errorx.ErrConflict | 409 |
| DATA_EXIST | errorx.ErrDataExist | 409 |
| CatalogServiceError | errorx.ErrInternal (30100) | 500 |

### 错误响应格式

```json
{
    "code": "InvalidParameter",
    "description": "目录名称格式不正确",
    "detail": {
        "field": "catalogName",
        "error": "名称不能以_或-开头"
    },
    "solution": "请使用中文、英文、数字开头，可包含_和-"
}
```

---

## Logging & Observability

### 日志规范

```go
import "github.com/zeromicro/go-zero/core/logx"

// 结构化日志
logx.WithContext(ctx).Infow("创建目录",
    logx.Field("catalogId", catalog.Id),
    logx.Field("catalogName", catalog.CatalogName),
    logx.Field("parentId", catalog.ParentId),
)

// 错误日志
logx.WithContext(ctx).Errorw("创建目录失败",
    logx.Field("error", err.Error()),
    logx.Field("parentId", req.ParentId),
)
```

### Trace ID

通过 `telemetry` 模块自动注入，每个请求包含唯一 trace_id。

---

## Deployment Considerations

### 环境变量

```yaml
# api/etc/catalog.yaml
Name: standardization-backend-api
Host: 0.0.0.0
Port: 8888
Mode: dev

# 数据库配置
DataSource: root:password@tcp(127.0.0.1:3306)/af_std?charset=utf8mb4&parseTime=true

# Redis 配置（可选）
RedisHost: 127.0.0.1:6379
RedisType: node
RedisPass: ""

# 日志配置
Log:
    ServiceName: standardization-backend-api
    Mode: console
    Level: info
```

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-01-21 | - | 初始版本 - Java 目录管理 API 转写技术方案 |
