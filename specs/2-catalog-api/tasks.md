# 目录管理 (catalog-api) Tasks

> **Branch**: `2-catalog-api`
> **Spec Path**: `specs/2-catalog-api/`
> **Created**: 2026-02-06
> **Input**: spec.md, plan.md
> **Status**: ✅ Completed

---

## 重要命令

### API 代码生成 (必读)

⚠️ **每次修改 `.api` 文件后，必须执行以下命令重新生成代码**：

```bash
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

**参数说明**:
- `--style=go_zero`: 使用 Go-Zero 官方代码风格（snake_case 命名）
- `--type-group`: 按类型分组生成文件（handler/logic/types 分离）

**生成文件命名规则**:
- Handler: `xxx_handler.go` (snake_case，如 `create_catalog_handler.go`)
- Logic: `xxx_logic.go` (snake_case，如 `create_catalog_logic.go`)
- Types: `types.go` (统一生成到 `api/internal/types/`)

⚠️ **不要使用其他 style 参数**，否则会生成 camelCase 文件导致重复定义错误。

---

## 任务组织说明

**组织方式**: 按接口增量定义维度，每个接口独立完成 API → Model → Logic → Test

**Mock 策略**: 对于依赖其他服务的逻辑（dataelement、dict、rule、stdfile），使用注释标记 TODO，后续补充 RPC 调用

**任务标记**:
- `[P]` = 可并行执行
- `[TEST]` = 测试任务
- `[MOCK]` = 使用 Mock 数据，后续补充 RPC

---

## Task Overview

| 阶段 | 描述 | 任务数 | 预计工作量 | 状态 |
|------|------|--------|------------|------|
| Phase 0 | 基础设施 | 5 | 1天 | ✅ 已完成 |
| Phase 1 | 查询接口 (2个) | 8 | 2天 | ✅ 已完成 |
| Phase 2 | CRUD接口 (3个) | 12 | 3天 | ✅ 已完成 |
| Phase 3 | 文件树接口 (1个) | 4 | 1天 | ✅ 已完成 |
| Phase 4 | 收尾工作 | 4 | 0.5天 | ✅ 已完成 |
| **总计** | | **33** | **约7.5天** | **100%** |

---

## Phase 0: 基础设施

**目的**: 项目初始化和基础配置

### P001 - 环境检查

- [x] T001 确认 Go-Zero 项目结构已就绪
  - [x] 确认 `api/` 目录存在
  - [x] 确认 `model/` 目录存在
  - [x] 确认 `go.mod` 已配置 Go-Zero 依赖

### P002 - 配置验证

- [x] T002 [P] 确认 base.api 已定义通用类型
- [x] T003 [P] 确认数据库连接已配置

### P003 - 错误码定义

- [x] T004 创建 `api/internal/errorx/codes.go`
  - [x] 定义错误码 30100-30199
  - [x] 实现错误辅助函数
  - [x] 30101: 数据不存在 (Empty)
  - [x] 30102: 参数缺失 (MissingParameter)
  - [x] 30103: 参数无效 (InvalidParameter)
  - [x] 30104: 超出范围 (OutOfRange)
  - [x] 30105: 操作冲突 (OperationConflict)
  - [x] 30106: 数据已存在 (DATA_EXIST)

**Checkpoint**: ✅ 基础设施就绪

---

## Phase 1: 查询接口 (2个)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 1 | GET | `/catalog/query_tree` | 查询目录树 | P1 |
| 2 | GET | `/catalog/query` | 检索目录 | P1 |

### 1.1 API 定义

- [x] T005 创建 `api/doc/catalog/catalog.api`
  - [x] 定义基础类型: CreateCatalogReq, UpdateCatalogReq, CatalogResp, CatalogTreeNodeVo
  - [x] 定义 QueryTreeReq, QueryReq
  - [x] 定义 6 个 API 端点
  - [x] 配置路由: `@server(prefix: /api/standardization/v1, group: catalog)`

- [x] T006 在 `api/doc/api.api` 中 import catalog 模块

- [x] T007 运行 `goctl api go` 生成 Handler/Types
  ```bash
  goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
  ```

### 1.2 DDL 定义

- [x] T008 [P] 创建 `migrations/catalog/raw/t_de_catalog_info.sql`

### 1.3 Model 层

- [x] T009 创建 `model/catalog/catalog/` 目录结构
  - [x] `interface.go` - CatalogModel 接口
  - [x] `types.go` - Catalog、CatalogTreeNodeVo、CatalogWithFileVo
  - [x] `vars.go` - 枚举常量、错误码
  - [x] `factory.go` - 工厂函数

- [x] T010 实现 `model/catalog/catalog/sql_model.go`
  - [x] Insert, FindOne, Update, Delete
  - [x] FindByType, FindByTypeAndLevel
  - [x] FindByName, FindByParentId
  - [x] FindByIds, FindTree
  - [x] FindChildren, DeleteByIds

- [x] T011 **[TEST]** `model/catalog/catalog/sql_model_test.go`
  - [x] Test Insert
  - [x] Test FindOne
  - [x] Test FindByType
  - [x] Test FindTree

### 1.4 公共 Logic (common.go)

- [x] T012 创建 `api/internal/logic/catalog/common.go`
  - [x] ValidateCatalogName (目录名称校验)
  - [x] BuildTree (构建树形结构)
  - [x] GetAllChildIds (递归获取子级ID)
  - [x] CheckCatalogDelete (删除前校验)

### 1.5 接口实现: GET /catalog/query_tree

- [x] T013 实现 `api/internal/logic/catalog/query_tree_logic.go`
  - [x] 业务流程标注 (对应 Java: DeCatalogInfoController.querySonTree)
  - [x] 按type查询: 校验type有效性，获取最小level，构建树
  - [x] 按id查询: 校验id存在性，获取该目录及其子集

- [x] T014 **[TEST]** `api/internal/logic/catalog/query_tree_logic_test.go`

### 1.6 接口实现: GET /catalog/query

- [x] T015 [P] 实现 `api/internal/logic/catalog/query_logic.go`
  - [x] 业务流程标注 (对应 Java: DeCatalogInfoController.queryParentTree)
  - [x] 校验type有效性
  - [x] 按关键字检索目录名称
  - [x] 返回平铺列表

- [x] T016 [P] **[TEST]** `api/internal/logic/catalog/query_logic_test.go`

### 1.7 ServiceContext 更新

- [x] T017 更新 `api/internal/svc/service_context.go`
  - [x] 添加 CatalogModel
  - [x] 初始化 DB 连接 (*sqlx.DB)
  - [x] 初始化 Model 实例
  - [x] TODO: 后续补充 DataElementModel、DictModel、RuleModel、StdFileModel

**Checkpoint**: ✅ Phase 1 完成 - 查询接口实现

---

## Phase 2: CRUD接口 (3个)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 3 | POST | `/catalog` | 创建目录 | P1 |
| 4 | PUT | `/catalog/{id}` | 修改目录 | P1 |
| 5 | DELETE | `/catalog/{id}` | 删除目录 | P1 |

### 2.1 接口实现: POST /catalog

- [x] T018 [P] 实现 `api/internal/logic/catalog/create_catalog_logic.go`
  - [x] 业务流程标注 (对应 Java: DeCatalogInfoServiceImpl.checkPost + save)
  - [x] 目录名称格式校验
  - [x] 父目录存在性校验
  - [x] 目录级别校验 (<= 255)
  - [x] 同级名称唯一性校验
  - [x] 继承父目录type和level+1

- [x] T019 [P] **[TEST]** `api/internal/logic/catalog/create_catalog_logic_test.go`

### 2.2 接口实现: PUT /catalog/{id}

- [x] T020 [P] 实现 `api/internal/logic/catalog/update_catalog_logic.go`
  - [x] 业务流程标注 (对应 Java: DeCatalogInfoServiceImpl.update + checkPost type=1)
  - [x] 目录存在性校验
  - [x] 不允许修改根目录
  - [x] 目录名称格式校验
  - [x] 同级名称唯一性校验（排除自身）

- [x] T021 [P] **[TEST]** `api/internal/logic/catalog/update_catalog_logic_test.go`

### 2.3 接口实现: DELETE /catalog/{id}

- [x] T022 [P] 实现 `api/internal/logic/catalog/delete_catalog_logic.go`
  - [x] 业务流程标注 (对应 Java: DeCatalogInfoServiceImpl.removeWithChildren + checkCatalogDelete)
  - [x] 目录存在性校验
  - [x] 不允许删除根目录
  - [x] 检查目录及子目录下是否存在数据
  - [x] 递归删除所有子级

- [x] T023 [P] **[TEST]** `api/internal/logic/catalog/delete_catalog_logic_test.go`

**Checkpoint**: ✅ Phase 2 完成 - CRUD接口实现

---

## Phase 3: 文件树接口 (1个)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 6 | GET | `/catalog/query/with_file` | 查询目录及文件树 | P2 |

### 3.1 接口实现: GET /catalog/query/with_file

- [x] T024 实现 `api/internal/logic/catalog/query_with_file_logic.go`
  - [x] 业务流程标注 (对应 Java: DeCatalogInfoController.queryParentTree with files)
  - [x] 按关键字检索目录
  - [x] 按关键字检索文件
  - [x] 返回目录列表和文件列表的平铺结构

- [x] T025 **[TEST]** `api/internal/logic/catalog/query_with_file_logic_test.go`

**Checkpoint**: ✅ Phase 3 完成 - 文件树接口实现

---

## Phase 4: 收尾工作

### 4.1 代码质量

- [x] T026 代码清理和格式化 (`gofmt -w .`)
- [x] T027 运行 `golangci-lint run` 修复代码质量问题

### 4.2 测试验证

- [x] T028 确认测试覆盖率 ≥ 80%
  ```bash
  go test ./... -coverprofile=coverage.out
  go tool cover -func=coverage.out | grep total
  ```

- [x] T029 运行所有测试确认通过
  ```bash
  go test ./... -v
  ```

### 4.3 文档更新

- [x] T030 运行 `make swagger` 生成 Swagger 文档

- [x] T031 验证所有6个API端点已注册

### 4.4 兼容性验证

- [x] T032 验证错误码与Java实现完全一致

- [x] T033 接口兼容性验证
  - [x] 确认响应格式与Java完全一致
  - [x] 确认异常信息与Java完全一致

**Checkpoint**: ✅ Phase 4 完成 - 所有测试和验证通过

---

## Mock 函数说明

### 需要后续补充 RPC 的场景

| 场景 | 当前Mock实现 | 后续补充 |
|------|-------------|----------|
| 数据元校验 | `getMockDataElementByCatalog(catalogId)` | DataElement RPC |
| 码表校验 | `getMockDictByCatalog(catalogId)` | Dict RPC |
| 规则校验 | `getMockRuleByCatalog(catalogId)` | Rule RPC |
| 文件查询 | `getMockFiles(keyword)` | StdFile RPC |
| 文件详情 | `getMockFileById(fileId)` | StdFile RPC |

### Mock 示例

```go
// ====== MOCK 函数 (后续补充 RPC) ======

// getMockDataElementByCatalog 检查目录下是否存在数据元
func getMockDataElementByCatalog(catalogId int64) bool {
    // TODO: 调用 DataElement RPC 检查
    // 当前返回false表示无数据
    return false
}

// getMockDictByCatalog 检查目录下是否存在码表
func getMockDictByCatalog(catalogId int64) bool {
    // TODO: 调用 Dict RPC 检查
    // 当前返回false表示无数据
    return false
}

// getMockRuleByCatalog 检查目录下是否存在规则
func getMockRuleByCatalog(catalogId int64) bool {
    // TODO: 调用 Rule RPC 检查
    // 当前返回false表示无数据
    return false
}

// getMockStdFileByCatalog 检查目录下是否存在文件
func getMockStdFileByCatalog(catalogId int64) bool {
    // TODO: 调用 StdFile RPC 检查
    // 当前返回false表示无数据
    return false
}

// getMockFiles 按关键字查询文件
func getMockFiles(keyword string) []*FileCountVo {
    // TODO: 调用 StdFile RPC 查询
    // 当前返回空列表
    return []*FileCountVo{}
}
```

---

## 依赖关系

```
Phase 0 (基础设施)
    ↓
Phase 1 (查询接口) ← MVP 🎯
    ↓
Phase 2 (CRUD接口)
    ↓
Phase 3 (文件树接口)
    ↓
Phase 4 (收尾工作)
```

### 并行执行说明

- `[P]` 标记的任务可并行执行
- Phase 1-3 的接口可并行开发（如有团队）

---

## MVP 范围

**最小可交付版本**: Phase 0 + Phase 1

MVP 包含的核心功能：
- ✅ 查询目录树（按类型或按ID）
- ✅ 检索目录（按关键字）

---

## 测试要求 🧪

| 要求 | 标准 |
|------|------|
| **单元测试覆盖率** | ≥ 80% |
| **关键路径测试** | 100% 覆盖（查询、创建、删除） |
| **边界测试** | 必须包含 |
| **错误处理测试** | 必须包含 |

### 测试命名规范

```
Test{Function}_{Scenario}_{ExpectedResult}
```

示例：
- `TestQueryTree_ByType_ReturnsTree`
- `TestQueryTree_ById_ReturnsSubTree`
- `TestCreateCatalog_ValidInput_ReturnsCatalog`
- `TestCreateCatalog_DuplicateName_ReturnsError`
- `TestDeleteCatalog_RootCatalog_ReturnsError`

---

## 实施进度跟踪

| Week | Phase | 内容 | 完成度 |
|------|-------|------|--------|
| 1 | Phase 0-4 | 基础设施 + 查询 + CRUD + 文件树 + 收尾 | 100% |

---

## 文档变更历史

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-02-06 | 初始版本 - 参考rule-api任务结构 |
