# 目录管理 API (Catalog API) 转写 Tasks

> **Branch**: `feature/catalog-api`
> **Spec Path**: `specs/catalog-api/`
> **Created**: 2025-01-21
> **Input**: spec.md, plan.md

---

## Task Format

```
[ID] [P?] [Story] Description
```

| 标记 | 含义 |
|------|------|
| `T001` | 任务 ID |
| `[P]` | 可并行执行（不同文件，无依赖） |
| `[US1]` | 关联 User Story 1 |
| `[TEST]` | 测试任务（必须完成） |

---

## Task Overview

| ID | Task | Story | Status | Parallel | Est. Lines |
|----|------|-------|--------|----------|------------|
| T001 | 项目基础设置 | Setup | ✅ | - | - |
| T002 | 错误码定义 | Setup | ✅ | - | 20 |
| T003 | 桩模块定义 | Foundation | ✅ | - | 40 |
| T004 | API 文件定义 | US1 | ⏸️ | - | 300 |
| T005 | goctl 生成代码 | US1 | ⏸️ | - | - |
| T006 | Model 层实现 | US1 | ⏸️ | - | 120 |
| T007 | Logic 层实现 - 目录树查询 | US1 | ⏸️ | - | 80 |
| T008 | Logic 层实现 - 目录检索 | US2 | ⏸️ | - | 60 |
| T009 | Logic 层实现 - 创建目录 | US3 | ⏸️ | - | 70 |
| T010 | Logic 层实现 - 修改目录 | US4 | ⏸️ | - | 90 |
| T011 | Logic 层实现 - 删除目录 | US5 | ⏸️ | - | 80 |
| T012 | Logic 层实现 - 目录文件树查询 | US6 | ⏸️ | - | 50 |
| T013 | 单元测试 | All US | ⏸️ | [P] | 200 |
| T014 | 集成验证 | Polish | ⏸️ | - | - |

---

## Phase 1: Setup

**目的**: 项目初始化和基础配置

- [x] T001 确认 Go-Zero 项目结构已就绪
  - 验证 `api/` 目录存在
  - 验证 `model/` 目录存在
  - 验证 `go.mod` 文件存在

- [x] T002 [P] 创建目录管理错误码定义 `api/internal/errorx/codes.go`
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

- [x] T003 [P] 创建外部依赖检查桩模块 `model/catalog/stub/external_checker.go`
  ```go
  package stub

  import "context"

  // ExternalChecker 外部依赖检查接口（桩模块）
  type ExternalChecker interface {
      CheckDataElement(ctx context.Context, catalogId string) (bool, error)
      CheckDict(ctx context.Context, catalogId string) (bool, error)
      CheckRule(ctx context.Context, catalogId string) (bool, error)
      CheckFile(ctx context.Context, catalogId string) (bool, error)
  }

  // StubExternalChecker 桩模块实现（当前阶段所有返回 false）
  type StubExternalChecker struct{}

  func (s *StubExternalChecker) CheckDataElement(ctx context.Context, catalogId string) (bool, error) {
      return false, nil // 桩实现：暂不检查
  }
  func (s *StubExternalChecker) CheckDict(ctx context.Context, catalogId string) (bool, error) {
      return false, nil
  }
  func (s *StubExternalChecker) CheckRule(ctx context.Context, catalogId string) (bool, error) {
      return false, nil
  }
  func (s *StubExternalChecker) CheckFile(ctx context.Context, catalogId string) (bool, error) {
      return false, nil
  }
  ```

**Checkpoint**: ✅ 开发环境就绪

---

## Phase 2: Foundation

**目的**: 必须完成后才能开始 User Story 实现

- [x] T004 确认 base.api 已定义通用类型
  - 验证 `api/doc/base.api` 文件存在
  - 验证 PageBaseInfo, PageInfo 等类型已定义

- [x] T005 确认数据库连接配置
  - 验证 `api/etc/api.yaml` 数据库配置已设置
  - 确认数据库表 `t_de_catalog_info` 已存在

**Checkpoint**: ✅ 基础设施就绪，可开始 User Story 实现

---

## Phase 3: User Story 1 - 目录树查询 (P1) 🎯 MVP

**目标**: 实现通过类型或ID查询目录树结构功能

**独立测试**: 调用 `GET /v1/catalog/query_tree?type=1` 返回完整的目录树 JSON

### Step 1: 定义 API 文件

- [x] T006 [US1] 创建 `api/doc/catalog/catalog.api`
  - 定义所有 6 个接口的 Request/Response 类型
  - 定义 CatalogTreeNodeVo, CatalogInfoVo, CatalogListByFileVo 等类型
  - 定义服务配置 (@server prefix: /api/v1/catalog)

- [x] T007 [US1] 在 `api/doc/api.api` 入口文件中 import 目录管理模块
  ```api
  import "catalog/catalog.api"
  ```

### Step 2: 生成代码

- [x] T008 [US1] 运行 `goctl api go` 生成 Handler/Types
  ```bash
  goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
  ```

- [ ] T009 [US1] [P] 运行 `goctl api swagger` 生成 Swagger 文档
  ```bash
  make swagger
  ```

### Step 3: 实现 Model 层

- [x] T010 [P] [US1] 创建 `model/catalog/catalog/interface.go`
  - 定义 CatalogModel 接口
  - 定义 Insert, FindOne, FindByParent, FindByType, FindByName, Update, Delete, DeleteBatch 方法签名

- [x] T011 [P] [US1] 创建 `model/catalog/catalog/types.go`
  - 定义 Catalog 实体结构
  - 定义 GORM 标签和 JSON 标签

- [x] T012 [P] [US1] 创建 `model/catalog/catalog/vars.go`
  - 定义目录名称正则: `^[\u4e00-\u9fa5a-zA-Z0-9][\u4e00-\u9fa5a-zA-Z0-9_-]{0,19}$`
  - 定义目录类型常量

- [x] T013 [P] [US1] 创建 `model/catalog/catalog/factory.go`
  - 定义 NewCatalogModel 工厂函数
  - 支持 GORM 和 SQLx 双 ORM

- [x] T014 [US1] 实现 `model/catalog/catalog/gorm_dao.go`
  - 实现 CatalogModel 接口的所有方法
  - 使用 GORM 进行数据库操作

- [ ] T015 [P] [US1] **[TEST]** 创建 `model/catalog/catalog/gorm_dao_test.go`
  - 测试 CRUD 操作
  - 测试边界情况和错误处理

### Step 4: 实现 Logic 层

- [x] T016 [US1] 实现 `api/internal/logic/catalog/querytreelogic.go`
  - 校验 type 参数在 [1,2,3,4] 范围内
  - 如果提供 id，查询该目录并返回其子树
  - 如果未提供 id，查询该类型所有目录，构建完整树
  - 实现树结构构建算法（递归查询 + 父子映射）
  - 预计代码: ~80 行

- [ ] T017 [US1] **[TEST]** 测试 `api/internal/logic/catalog/querytreelogic_test.go`
  - 测试按类型查询完整树
  - 测试按 ID 查询子树
  - 测试参数校验
  - 测试树结构构建逻辑

**Checkpoint**: ✅ User Story 1 已完成，代码 + 测试 全部通过

---

## Phase 4: User Story 2 - 目录检索 (P1)

**目标**: 实现通过目录名称模糊检索目录功能

**独立测试**: 调用 `GET /v1/catalog/query?catalog_name=测试&type=1` 返回匹配的目录列表

### Implementation + Test

- [x] T018 [US2] 实现 `api/internal/logic/catalog/querylogic.go`
  - 校验 type 参数
  - 目录名称 SQL 特殊字符转义
  - 模糊查询：`WHERE f_catalog_name LIKE %name%` (小写不敏感)
  - 仅返回 level > 1 的目录（排除根目录）
  - 预计代码: ~60 行

- [ ] T019 [US2] **[TEST]** 测试 `api/internal/logic/catalog/querylogic_test.go`
  - 测试模糊查询功能
  - 测试 SQL 特殊字符转义
  - 测试空结果场景

**Checkpoint**: ✅ User Story 2 已完成，代码 + 测试 全部通过

---

## Phase 5: User Story 3 - 创建目录 (P1)

**目标**: 实现创建新的目录节点功能

**独立测试**: 调用 `POST /v1/catalog` 创建目录，返回成功

### Implementation + Test

- [x] T020 [US3] 实现 `api/internal/logic/catalog/createlogic.go`
  - 校验目录名称格式：正则 `^[\u4e00-\u9fa5a-zA-Z0-9][\u4e00-\u9fa5a-zA-Z0-9_-]{0,19}$`
  - 校验父目录存在且 level < 255
  - 继承父目录的 type，设置 level = 父目录 level + 1
  - 检查同级目录名称唯一性
  - 插入数据库
  - 预计代码: ~70 行

- [ ] T021 [US3] **[TEST]** 测试 `api/internal/logic/catalog/createlogic_test.go`
  - 测试正常创建
  - 测试目录名称为空
  - 测试目录名称格式错误
  - 测试父目录不存在
  - 测试父目录级别超限
  - 测试同级名称重复

**Checkpoint**: ✅ User Story 3 已完成，代码 + 测试 全部通过

---

## Phase 6: User Story 4 - 修改目录 (P1)

**目标**: 实现修改目录名称和父目录功能

**独立测试**: 调用 `PUT /v1/catalog/{id}` 更新目录，返回成功

### Implementation + Test

- [ ] T022 [US4] 实现 `api/internal/logic/catalog/updatelogic.go`
  - 校验目录存在
  - 校验不是根目录（level <= 1 不允许修改）
  - 校验新父目录存在且 level < 255
  - 循环检测：新父目录不能是自身子目录
  - 类型一致性：新父目录 type 必须与当前目录一致
  - 检查同级名称唯一性（排除自身）
  - 更新数据库
  - 预计代码: ~90 行

- [ ] T023 [US4] **[TEST]** 测试 `api/internal/logic/catalog/updatelogic_test.go`
  - 测试正常修改
  - 测试修改根目录（应失败）
  - 测试循环父子关系
  - 测试类型不一致
  - 测试同级名称重复

**Checkpoint**: ✅ User Story 4 已完成，代码 + 测试 全部通过

---

## Phase 7: User Story 5 - 删除目录 (P1)

**目标**: 实现删除目录及其所有子目录功能

**独立测试**: 调用 `DELETE /v1/catalog/{id}` 删除目录及其子目录

### Implementation + Test

- [ ] T024 [US5] 实现 `api/internal/logic/catalog/deletelogic.go`
  - 校验目录存在
  - 校验不是根目录（level <= 1 不允许删除）
  - 递归获取所有子目录 ID
  - 调用桩模块检查关联数据（当前阶段跳过）
  - 批量删除目录及子目录
  - 实现递归获取子目录算法
  - 预计代码: ~80 行

- [ ] T025 [US5] **[TEST]** 测试 `api/internal/logic/catalog/deletelogic_test.go`
  - 测试正常删除
  - 测试删除根目录（应失败）
  - 测试级联删除
  - 测试桩模块调用

**Checkpoint**: ✅ User Story 5 已完成，代码 + 测试 全部通过

---

## Phase 8: User Story 6 - 目录文件树查询 (P2)

**目标**: 实现查询包含文件的目录树结构功能

**独立测试**: 调用 `GET /v1/catalog/query/with_file` 返回目录和文件列表

### Implementation + Test

- [ ] T026 [US6] 实现 `api/internal/logic/catalog/querywithfilelogic.go`
  - 模糊查询目录名称
  - 查询文件列表（当前阶段返回空列表，待文件模块实现）
  - 组装 CatalogListByFileVo 响应
  - 预计代码: ~50 行

- [ ] T027 [US6] **[TEST]** 测试 `api/internal/logic/catalog/querywithfilelogic_test.go`
  - 测试模糊查询
  - 测试空结果场景
  - 测试文件列表组装

**Checkpoint**: ✅ User Story 6 已完成，代码 + 测试 全部通过

---

## Phase 9: Polish & 验证

**目的**: 收尾工作和质量保证

- [ ] T028 [P] 代码清理和格式化
  ```bash
  gofmt -w .
  goimports -w .
  ```

- [ ] T029 [P] 运行 golangci-lint 代码检查
  ```bash
  make lint
  ```

- [ ] T030 **确认测试覆盖率 > 80%**
  ```bash
  go test ./... -coverprofile=coverage.out
  go tool cover -func=coverage.out | grep -E "(total|model/catalog|logic/catalog)"
  ```

- [ ] T031 验证 API 与 Java 版本契约一致性
  - 对比响应 JSON 结构
  - 验证错误码映射
  - 验证业务规则一致性

---

## Dependencies

```
Phase 1 (Setup)
    ↓
Phase 2 (Foundation)
    ↓
Phase 3 (US1: 目录树查询 + Tests) → Phase 4 (US2: 目录检索 + Tests)
    ↓                                   ↓
Phase 5 (US3: 创建目录 + Tests) → Phase 6 (US4: 修改目录 + Tests)
    ↓                                   ↓
Phase 7 (US5: 删除目录 + Tests) → Phase 8 (US6: 目录文件树查询 + Tests)
    ↓
Phase 9 (Polish & 验证)
```

### 并行执行说明

- `[P]` 标记的任务可与同 Phase 内其他 `[P]` 任务并行
- `[TEST]` 标记的任务必须与对应实现任务同步完成
- Model 层文件（T010-T014）可并行开发
- 不同 Logic 层文件（T016, T018, T020, T022, T024, T026）可并行开发
- 单元测试（T017, T019, T021, T023, T025, T027）必须与对应 Logic 实现同步

---

## 测试要求 🧪

| 要求 | 标准 |
|------|------|
| **单元测试覆盖率** | > 80% |
| **关键路径测试** | 100% 覆盖 |
| **边界测试** | 必须包含 |
| **错误处理测试** | 必须包含 |

### 测试命名规范

```
Test{Function}_{Scenario}_{ExpectedResult}
```

示例：
- `TestQueryTree_ValidType_ReturnsTree`
- `TestQueryTree_InvalidType_ReturnsError`
- `TestCreateCatalog_DuplicateName_ReturnsConflict`

---

## Notes

- 每个 Task 完成后运行 `go test ./...` 验证
- **实现和测试必须同时提交**
- 每完成一个 User Story 运行完整测试验证
- 遇到问题及时记录到 spec.md 的 Open Questions
- API 响应格式必须与 Java 版本 100% 一致

---

## 实现问题记录

### 问题 1: CatalogModel 未初始化

**现象**: Logic 层调用 `svcCtx.CatalogModel.Insert()` 时，CatalogModel 为 nil

**原因**:
- `service_context.go` 中 `CatalogModel` 设置为 `nil`
- 没有在 `NewServiceContext` 中初始化
- 缺少 `main.go` 入口文件

**解决方案**:
1. 更新 `config/config.go` 添加 `DB *gorm.DB` 字段
2. 创建 `api/main.go` 初始化数据库连接
3. 在 `NewServiceContext` 中调用 `catalog.NewCatalogModel()` 初始化

**相关提交**: `121fdf0` - fix(catalog-api): initialize CatalogModel in ServiceContext

---

### 问题 2: goctl 生成代码风格问题

**现象**:
- goctl 使用 `--style=go_zero` 生成 `querytreelogic.go` (驼峰合并)
- Handler 调用 `l.Querytree()` (小写方法名)
- 与 Go 命名规范不一致

**原因**: `--style=go_zero` 风格不适合大型项目

**解决方案**:
1. 使用 `--style=default` 生成带下划线的文件名
2. 保留 goctl 生成的 `_logic.go` 文件结构（含 `logx.Logger`）
3. 在 goctl 生成的文件中直接实现业务逻辑
4. 不要创建额外的 `xxxlogic.go` 文件

**正确结构**:
```
goctl 生成: querytree_logic.go (模板，含 logx.Logger)
开发者实现: 在 querytree_logic.go 中添加业务逻辑
Handler 调用: NewQuerytreeLogic(ctx, svcCtx).Querytree(&req)
```

**相关提交**: `e91d459` - fix(catalog-api): move logic from incorrectly named files

---

### 问题 3: Handler 方法名与 Logic 层不匹配

**现象**: Handler 调用 `QueryTree()` 但 Logic 定义的是 `Querytree()`

**原因**: 手动修改了 Handler 的方法名为驼峰，但 goctl 生成的是小写

**解决方案**:
- 保持 goctl 生成的方法名不变 (`Querytree` 而不是 `QueryTree`)
- Go-Zero 的命名风格是方法名小写开头 (如 `Querytree`, `Create`, `Update`)

**相关提交**: `542e470` - fix(catalog-api): fix handler method naming

---

### 问题 4: ValidateCatalogName 引用不存在的变量

**现象**: 代码引用 `catalog.CatalogNamePatternRegexp.MatchString()` 但变量不存在

**原因**: `vars.go` 只定义了字符串常量，没有定义 regexp 对象

**解决方案**:
- 简化为直接校验首字符规则
- 或在 `vars.go` 中初始化 `regexp.MustCompile(catalog.CatalogNamePattern)`

---

### 经验总结

1. **goctl 使用建议**:
   - 始终使用 `--style=default` 生成带下划线的文件名
   - 在 goctl 生成的文件中直接实现业务逻辑
   - 不要创建额外的文件

2. **初始化顺序**:
   - Config (数据库配置)
   - ServiceContext (初始化 Model)
   - Logic (调用 Model)
   - Handler (调用 Logic)

3. **命名规范**:
   - goctl 生成的文件名: `xxx_logic.go`, `xxx_handler.go`
   - goctl 生成的方法名: `Querytree`, `Create`, `Update` (首字母大写，内部小写)
   - 不要修改 goctl 生成的方法签名

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-01-21 | - | 初始版本 - Java 目录管理 API 转写任务拆分 |
| 1.1 | 2025-01-21 | - | 添加实现问题记录和经验总结 |
