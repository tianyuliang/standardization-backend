# 目录管理 (catalog-api) Tasks

> **Branch**: `2-catalog-api`
> **Spec Path**: `specs/2-catalog-api/`
> **Created**: 2026-02-06
> **Input**: spec.md, plan.md
> **Status**: Draft

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
| Phase 2 | CRUD接口 (3个) | 12 | 3天 | ⏳ 进行中 |
| Phase 3 | 文件树接口 (1个) | 4 | 1天 | ⏳ 待开始 |
| Phase 4 | 收尾工作 | 4 | 0.5天 | ⏳ 待开始 |
| **总计** | | **33** | **约7.5天** | **约30%** |

---

## Phase 0: 基础设施

**目的**: 项目初始化和基础配置

### P001 - 环境检查

- [ ] T001 确认 Go-Zero 项目结构已就绪
  - [ ] 确认 `api/` 目录存在
  - [ ] 确认 `model/` 目录存在
  - [ ] 确认 `go.mod` 已配置 Go-Zero 依赖

### P002 - 配置验证

- [ ] T002 [P] 确认 base.api 已定义通用类型
- [ ] T003 [P] 确认数据库连接已配置

### P003 - 错误码定义

- [ ] T004 创建 `api/internal/errorx/codes.go`
  - [ ] 定义错误码 30100-30199
  - [ ] 实现错误辅助函数
  - [ ] 30101: 数据不存在 (Empty)
  - [ ] 30102: 参数缺失 (MissingParameter)
  - [ ] 30103: 参数无效 (InvalidParameter)
  - [ ] 30104: 超出范围 (OutOfRange)
  - [ ] 30105: 操作冲突 (OperationConflict)
  - [ ] 30106: 数据已存在 (DATA_EXIST)

**Checkpoint**: ✅ 基础设施就绪

---

## Phase 1: 查询接口 (2个)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 1 | GET | `/catalog/query_tree` | 查询目录树 | P1 |
| 2 | GET | `/catalog/query` | 检索目录 | P1 |

### 1.1 API 定义

- [ ] T005 创建 `api/doc/catalog/catalog.api`
  - [ ] 定义基础类型: CreateCatalogReq, UpdateCatalogReq, CatalogResp, CatalogTreeNodeVo
  - [ ] 定义 QueryTreeReq, QueryReq
  - [ ] 定义 6 个 API 端点
  - [ ] 配置路由: `@server(prefix: /api/standardization/v1, group: catalog)`

- [ ] T006 在 `api/doc/api.api` 中 import catalog 模块

- [ ] T007 运行 `goctl api go` 生成 Handler/Types
  ```bash
  goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
  ```

### 1.2 DDL 定义

- [ ] T008 [P] 创建 `migrations/catalog/raw/t_de_catalog_info.sql`

### 1.3 Model 层

- [ ] T009 创建 `model/catalog/catalog/` 目录结构
  - [ ] `interface.go` - CatalogModel 接口
  - [ ] `types.go` - Catalog、CatalogTreeNodeVo、CatalogWithFileVo
  - [ ] `vars.go` - 枚举常量、错误码
  - [ ] `factory.go` - 工厂函数

- [ ] T010 实现 `model/catalog/catalog/sql_model.go`
  - [ ] Insert, FindOne, Update, Delete
  - [ ] FindByType, FindByTypeAndLevel
  - [ ] FindByName, FindByParentId
  - [ ] FindByIds, FindTree
  - [ ] FindChildren, DeleteByIds

- [ ] T011 **[TEST]** `model/catalog/catalog/sql_model_test.go`
  - [ ] Test Insert
  - [ ] Test FindOne
  - [ ] Test FindByType
  - [ ] Test FindTree

### 1.4 公共 Logic (common.go)

- [ ] T012 创建 `api/internal/logic/catalog/common.go`
  - [ ] ValidateCatalogName (目录名称校验)
  - [ ] BuildTree (构建树形结构)
  - [ ] getAllChildIds (递归获取子级ID)
  - [ ] CheckCatalogDelete (删除前校验)
  ```go
  // 目录名称校验
  func ValidateCatalogName(name string) error

  // 构建树形结构
  func BuildTree(catalogs []*Catalog, rootLevel int32) []*Catalog

  // 递归获取子级ID列表
  func getAllChildIds(ctx context.Context, model CatalogModel, parentId int64) ([]int64, error)

  // 删除前校验
  func CheckCatalogDelete(ctx context.Context, catalog *Catalog, svcCtx *svc.ServiceContext) error
  ```

### 1.5 接口实现: GET /catalog/query_tree

- [ ] T013 实现 `api/internal/logic/catalog/query_tree_logic.go`
  - [ ] 业务流程标注 (对应 Java: DeCatalogInfoController.querySonTree)
  - [ ] 按type查询: 校验type有效性，获取最小level，构建树
  - [ ] 按id查询: 校验id存在性，获取该目录及其子集
  ```go
  func (l *QueryTreeLogic) QueryTree(req *types.QueryTreeReq) (resp *types.CatalogTreeNodeVo, err error) {
      // 1. 如果指定id，查询指定目录的子集树
      if req.Id != 0 {
          // 1.1 校验目录存在
          // 1.2 获取目录及其所有子级
          // 1.3 构建树形结构
          // 1.4 可选：添加数据统计 (getCatalogCountMap)
          return buildTreeWithRoot(catalog)
      }
      // 2. 如果指定type，查询该类型的完整目录树
      if req.Type != 0 {
          // 2.1 校验type有效性
          // 2.2 查询该类型的所有目录
          // 2.3 获取最小level作为根节点
          // 2.4 构建树形结构
          // 2.5 可选：添加数据统计
          return buildFullTree(req.Type)
      }
      // 3. 都未指定，返回空
      return emptyTree()
  }
  ```

- [ ] T014 **[TEST]** `api/internal/logic/catalog/query_tree_logic_test.go`

### 1.6 接口实现: GET /catalog/query

- [ ] T015 [P] 实现 `api/internal/logic/catalog/query_logic.go`
  - [ ] 业务流程标注 (对应 Java: DeCatalogInfoController.queryParentTree)
  - [ ] 校验type有效性
  - [ ] 按关键字检索目录名称
  - [ ] 返回平铺列表
  ```go
  func (l *QueryLogic) Query(req *types.QueryReq) (resp *[]types.CatalogInfoVo, err error) {
      // 1. 校验type有效性
      if err := checkType(req.Type); err != nil { return }
      // 2. XSS转义关键字
      keyword := escapeSqlSpecialChars(req.Keyword)
      // 3. 按名称模糊查询 (level > 1)
      // 4. 返回平铺列表
  }
  ```

- [ ] T016 [P] **[TEST]** `api/internal/logic/catalog/query_logic_test.go`

### 1.7 ServiceContext 更新

- [ ] T017 更新 `api/internal/svc/service_context.go`
  - [ ] 添加 CatalogModel
  - [ ] 初始化 DB 连接 (*sqlx.DB)
  - [ ] 初始化 Model 实例
  - [ ] TODO: 后续补充 DataElementModel、DictModel、RuleModel、StdFileModel

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

- [ ] T018 [P] 实现 `api/internal/logic/catalog/create_catalog_logic.go`
  - [ ] 业务流程标注 (对应 Java: DeCatalogInfoServiceImpl.checkPost + save)
  - [ ] 目录名称格式校验
  - [ ] 父目录存在性校验
  - [ ] 目录级别校验 (<= 255)
  - [ ] 同级名称唯一性校验
  - [ ] 继承父目录type和level+1
  ```go
  func (l *CreateCatalogLogic) Create(req *types.CreateCatalogReq) (resp *types.CatalogResp, err error) {
      // 1. 目录名称格式校验
      if err := ValidateCatalogName(req.CatalogName); err != nil { return }
      // 2. 父目录存在性校验
      parent, err := l.svcCtx.CatalogModel.FindOne(l.ctx, req.ParentId)
      if err != nil { return errorx.CatalogParentNotExist() }
      // 3. 目录级别校验
      if parent.Level >= 255 { return errorx.CatalogLevelOutOfRange() }
      // 4. 根目录不允许修改/创建（level=1为根目录）
      // 5. 继承父目录type
      catalogType := parent.Type
      // 6. 同级名称唯一性校验
      siblings, _ := l.svcCtx.CatalogModel.FindByParentId(l.ctx, req.ParentId)
      for _, sib := range siblings {
          if sib.CatalogName == req.CatalogName && sib.Type == catalogType {
              return errorx.CatalogNameDuplicate()
          }
      }
      // 7. 保存数据 (level = parent.level + 1)
      // 8. 返回结果
  }
  ```

- [ ] T019 [P] **[TEST]** `api/internal/logic/catalog/create_catalog_logic_test.go`

### 2.2 接口实现: PUT /catalog/{id}

- [ ] T020 [P] 实现 `api/internal/logic/catalog/update_catalog_logic.go`
  - [ ] 业务流程标注 (对应 Java: DeCatalogInfoServiceImpl.update + checkPost type=1)
  - [ ] 目录存在性校验
  - [ ] 不允许修改根目录
  - [ ] 目录名称格式校验
  - [ ] 同级名称唯一性校验（排除自身）
  ```go
  func (l *UpdateCatalogLogic) Update(id int64, req *types.UpdateCatalogReq) (err error) {
      // 1. 校验目录存在
      catalog, err := l.svcCtx.CatalogModel.FindOne(l.ctx, id)
      if err != nil { return errorx.CatalogNotExist() }
      // 2. 不允许修改根目录
      if catalog.Level <= 1 { return errorx.CannotModifyRootCatalog() }
      // 3. 目录名称格式校验
      if err := ValidateCatalogName(req.CatalogName); err != nil { return }
      // 4. 同级名称唯一性校验（排除自身）
      siblings, _ := l.svcCtx.CatalogModel.FindByParentId(l.ctx, catalog.ParentId)
      for _, sib := range siblings {
          if sib.CatalogName == req.CatalogName && sib.Id != id {
              return errorx.CatalogNameDuplicate()
          }
      }
      // 5. 更新数据
      // 6. 返回结果
  }
  ```

- [ ] T021 [P] **[TEST]** `api/internal/logic/catalog/update_catalog_logic_test.go`

### 2.3 接口实现: DELETE /catalog/{id}

- [ ] T022 [P] 实现 `api/internal/logic/catalog/delete_catalog_logic.go`
  - [ ] 业务流程标注 (对应 Java: DeCatalogInfoServiceImpl.removeWithChildren + checkCatalogDelete)
  - [ ] 目录存在性校验
  - [ ] 不允许删除根目录
  - [ ] 检查目录及子目录下是否存在数据
  - [ ] 递归删除所有子级
  ```go
  func (l *DeleteCatalogLogic) Delete(id int64) (err error) {
      // 1. 校验目录存在
      catalog, err := l.svcCtx.CatalogModel.FindOne(l.ctx, id)
      if err != nil { return errorx.CatalogNotExist() }
      // 2. 不允许删除根目录
      if catalog.Level <= 1 { return errorx.CannotDeleteRootCatalog() }
      // 3. 删除前校验（检查目录及子目录下是否存在数据）
      if err := CheckCatalogDelete(l.ctx, catalog, l.svcCtx); err != nil { return }
      // 4. 递归删除所有子级目录
      childIds, err := getAllChildIds(l.ctx, l.svcCtx.CatalogModel, id)
      allIds := append(childIds, id)
      return l.svcCtx.CatalogModel.DeleteByIds(l.ctx, allIds)
  }
  ```

- [ ] T023 [P] **[TEST]** `api/internal/logic/catalog/delete_catalog_logic_test.go`

**Checkpoint**: ✅ Phase 2 完成 - CRUD接口实现

---

## Phase 3: 文件树接口 (1个)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 6 | GET | `/catalog/query/with_file` | 查询目录及文件树 | P2 |

### 3.1 接口实现: GET /catalog/query/with_file

- [ ] T024 实现 `api/internal/logic/catalog/query_with_file_logic.go`
  - [ ] 业务流程标注 (对应 Java: DeCatalogInfoController.queryParentTree with files)
  - [ ] 按关键字检索目录
  - [ ] 按关键字检索文件
  - [ ] 返回目录列表和文件列表的平铺结构
  ```go
  func (l *QueryWithFileLogic) QueryWithFile(req *types.QueryWithFileReq) (resp *types.CatalogListByFileVo, err error) {
      // 1. XSS转义关键字
      keyword := escapeSqlSpecialChars(req.Keyword)
      // 2. 查询目录列表 (type=4, level>1)
      // TODO: 调用 StdFileModel 查询文件
      catalogs, _ := l.svcCtx.CatalogModel.FindByName(keyword, CatalogTypeFile)
      // 3. 查询文件列表
      // TODO: 调用 StdFileModel 查询文件
      files := getMockFiles(keyword)
      // 4. 构建响应
      return &types.CatalogListByFileVo{
          Catalogs: catalogs,
          Files:    files,
      }, nil
  }
  ```

- [ ] T025 **[TEST]** `api/internal/logic/catalog/query_with_file_logic_test.go`

**Checkpoint**: ✅ Phase 3 完成 - 文件树接口实现

---

## Phase 4: 收尾工作

### 4.1 代码质量

- [ ] T026 代码清理和格式化 (`gofmt -w .`)
- [ ] T027 运行 `golangci-lint run` 修复代码质量问题

### 4.2 测试验证

- [ ] T028 确认测试覆盖率 ≥ 80%
  ```bash
  go test ./... -coverprofile=coverage.out
  go tool cover -func=coverage.out | grep total
  ```

- [ ] T029 运行所有测试确认通过
  ```bash
  go test ./... -v
  ```

### 4.3 文档更新

- [ ] T030 运行 `make swagger` 生成 Swagger 文档

- [ ] T031 验证所有6个API端点已注册

### 4.4 兼容性验证

- [ ] T032 验证错误码与Java实现完全一致

- [ ] T033 接口兼容性验证
  - [ ] 确认响应格式与Java完全一致
  - [ ] 确认异常信息与Java完全一致

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
| 1 | Phase 0-4 | 基础设施 + 查询 + CRUD + 文件树 + 收尾 | 0% |

---

## 文档变更历史

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-02-06 | 初始版本 - 参考rule-api任务结构 |
