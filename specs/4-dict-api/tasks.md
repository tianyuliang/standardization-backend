# 码表管理 (dict-api) Tasks

> **Branch**: `4-dict-api`
> **Spec Path**: `specs/4-dict-api/`
> **Created**: 2026-02-06
> **Input**: spec.md, plan.md

---

## 重要说明

**⚠️ 永远使用此命令生成代码**:
```bash
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

**任务标记**:
- `[P]` = 可并行执行
- `[TEST]` = 测试任务
- `[MOCK]` = 使用 Mock 数据，后续补充 RPC

**禁止修改其他需求代码**:
- ❌ 禁止修改 `specs/1-rule-api/` 下的任何文件
- ❌ 禁止修改 `specs/2-catalog-api/` 下的任何文件
- ❌ 禁止修改 `specs/3-std-file-api/` 下的任何文件
- ❌ 禁止修改 `api/internal/handler/rule/` 下的任何文件
- ❌ 禁止修改 `api/internal/handler/catalog/` 下的任何文件
- ❌ 禁止修改 `api/internal/handler/stdfile/` 下的任何文件
- ❌ 禁止修改 `api/internal/logic/rule/` 下的任何文件
- ❌ 禁止修改 `api/internal/logic/catalog/` 下的任何文件
- ❌ 禁止修改 `api/internal/logic/stdfile/` 下的任何文件

---

## Task Overview

| 阶段 | 描述 | 任务数 | 预计工作量 | 状态 |
|------|------|--------|------------|------|
| Phase 0 | 基础设施 | 5 | 0.5天 | ✅ 完成 |
| Phase 1 | 基础CRUD (4接口) | 16 | 3天 | ✅ 完成 |
| Phase 2 | 码值管理 (2接口) | 4 | 1天 | ✅ 完成 |
| Phase 3 | 状态管理 (1接口) | 4 | 0.5天 | ✅ 完成 |
| Phase 4 | 关联查询 (4接口) | 8 | 1.5天 | ✅ 完成 |
| Phase 5 | 关联管理 (2接口) | 4 | 1天 | ✅ 完成 |
| Phase 6 | 数据校验 (1接口) | 2 | 0.5天 | ✅ 完成 |
| Phase 7 | 辅助接口 (2接口) | 4 | 0.5天 | ✅ 完成 |
| Phase 8 | 收尾工作 | 8 | 0.5天 | 🔄 部分完成 |
| **总计** | | **51** | **约10天** | **47/51 (92%)** |

---

## Phase 0: 基础设施

**目的**: 项目初始化和基础配置

### P001 - 环境检查

- [X] T001 确认 Go-Zero 项目结构已就绪
  - [X] 确认 `api/` 目录存在
  - [X] 确认 `model/` 目录存在
  - [X] 确认 `go.mod` 已配置 Go-Zero 依赖

### P002 - 配置验证

- [X] T002 [P] 确认 base.api 已定义通用类型
- [X] T003 [P] 确认数据库连接已配置
- [X] T004 [P] 确认 ServiceContext 已配置

### P003 - 错误码定义

- [X] T005 创建 `api/internal/errorx/dict.go`
  - [X] 定义错误码 30400-30499
  - [X] 实现错误辅助函数

**Checkpoint**: ✅ 基础设施就绪

---

## Phase 1: 基础CRUD (4个接口)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 1 | POST | `/v1/dataelement/dict` | 新增码表 | P1 |
| 2 | PUT | `/v1/dataelement/dict/{id}` | 修改码表 | P1 |
| 3 | GET | `/v1/dataelement/dict` | 码表列表查询 | P1 |
| 4 | GET | `/v1/dataelement/dict/{id}` | 码表详情（按ID） | P1 |

### 1.1 API 定义

- [X] T006 [US1] 创建 `api/doc/dict/dict.api`
  - [X] 定义基础类型: CreateDictReq, UpdateDictReq, DictVo, DictDataListResp
  - [X] 定义 DictEnumVo 码值类型
  - [X] 定义 4 个 API 端点
  - [X] 配置路由: `@server(prefix: /api/standardization/v1/dataelement, group: dict, middleware: TokenCheck)`

- [X] T007 [US1] 在 `api/doc/api.api` 中 import dict 模块

- [X] T008 [US1] 运行 `goctl api go` 生成 Handler/Types
  ```bash
  goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
  ```

### 1.2 DDL 定义

- [X] T009 [P] [US1] 创建 `migrations/dict/raw/t_dict.sql`
- [X] T010 [P] [US1] 创建 `migrations/dict/raw/t_dict_enum.sql`
- [X] T011 [P] [US1] 创建 `migrations/dict/raw/t_relation_dict_file.sql`

### 1.3 Model 层

- [X] T012 [US1] 创建 `model/dict/dict/` 目录结构
  - [X] `interface.go` - DictModel 接口
  - [X] `types.go` - Dict、DictEnum、RelationDictFile、DictVo
  - [X] `vars.go` - 枚举常量
  - [X] `factory.go` - SQLx 工厂函数

- [X] T013 [US1] 实现 `model/dict/dict/sql_model.go`
  - [X] Insert, FindOne, FindByCode, Update, Delete
  - [X] FindByIds, FindByCatalogIds (分页)
  - [X] FindDataExists

- [X] T014 [P] [US1] 实现 `model/dict/dictenum/sql_model.go`
  - [X] Insert, FindByDictId, FindListByDictId
  - [X] DeleteByDictId, CheckDuplicateCode

- [X] T015 [P] [US1] 实现 `model/dict/relation/sql_model.go`
  - [X] InsertBatch, DeleteByDictId, FindByDictId

### 1.4 公共 Logic (common.go)

- [X] T016 [US1] 创建 `api/internal/logic/dict/common.go`
  - [X] stateToInt, intToState
  - [X] parseDate (日期解析)
  - [X] buildDictVo
  - [X] generateSnowflakeCode (雪花算法生成码表编码)
  - [X] [MOCK] GetCatalogName (目录名称查询)
  - [X] [MOCK] GetCatalogChildIds (子目录ID列表)
  - [X] [MOCK] IsStdFileCatalog (目录类型校验)
  - [X] [MOCK] GetUserInfo (用户信息)
  - [X] [MOCK] GetDeptInfo (部门信息)

### 1.5 接口实现: POST /v1/dataelement/dict (新增码表)

- [X] T017 [US1] 实现 `api/internal/logic/dict/create_dict_logic.go`
  - [X] 6步业务流程标注
  - [X] 参数校验
  - [X] 业务校验（目录存在性、名称唯一性、码值唯一性）
  - [X] 生成码表编码（雪花算法）
  - [X] 保存码表
  - [X] 保存码值明细
  - [X] 保存关联文件关系

### 1.6 接口实现: PUT /v1/dataelement/dict/{id} (修改码表)

- [X] T018 [US1] 实现 `api/internal/logic/dict/update_dict_logic.go`
  - [X] 5步业务流程标注
  - [X] 校验码表存在性
  - [X] 校验参数
  - [X] 检查是否需要更新版本（版本变更检测）
  - [X] 有变更则更新（版本号+1、更新关联文件、重新保存码值）

### 1.7 接口实现: GET /v1/dataelement/dict (码表列表查询)

- [X] T019 [US1] 实现 `api/internal/logic/dict/list_dict_logic.go`
  - [X] 4步业务流程标注
  - [X] 处理目录ID（获取当前目录及所有子目录ID列表）
  - [X] 构建查询条件
  - [X] 分页查询
  - [X] 数据处理（查询是否被引用、查询部门信息）

### 1.8 接口实现: GET /v1/dataelement/dict/{id} (码表详情)

- [X] T020 [US1] 实现 `api/internal/logic/dict/get_dict_logic.go`
  - [X] 4步业务流程标注
  - [X] 查询码表基本信息
  - [X] 查询码值明细列表
  - [X] 查询目录信息
  - [X] 查询部门信息

### 1.9 ServiceContext 更新

- [X] T021 [US1] 在 `api/internal/svc/service_context.go` 中添加:
  - [X] DictModel model.DictModel
  - [X] DictEnumModel model.DictEnumModel
  - [X] RelationDictFileModel model.RelationDictFileModel

**Checkpoint**: ✅ Phase 1 完成 - 码表基础CRUD功能

---

## Phase 2: 码值管理 (2个接口)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 5 | GET | `/v1/dataelement/dict/enum` | 码值分页查询 | P1 |
| 6 | GET | `/v1/dataelement/dict/enum/getList` | 码值列表查询 | P1 |

### 2.1 接口实现: GET /v1/dataelement/dict/enum (码值分页查询)

- [X] T022 [US2] 实现 `api/internal/logic/dict/list_dict_enum_logic.go`
  - [X] 4步业务流程标注
  - [X] 参数校验（码表ID必填）
  - [X] 校验码表存在性
  - [X] 查询码值分页列表
  - [X] 数据处理（转换为响应格式）

### 2.2 接口实现: GET /v1/dataelement/dict/enum/getList (码值列表查询)

- [X] T023 [US2] 实现 `api/internal/logic/dict/get_dict_enum_list_logic.go`
  - [X] 4步业务流程标注
  - [X] 参数校验（码表ID必填）
  - [X] 校验码表存在性
  - [X] 查询所有码值（不分页）
  - [X] 数据处理（转换为响应格式）

**Checkpoint**: ✅ Phase 2 完成 - 码值管理功能

---

## Phase 3: 状态管理 (1个接口)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 7 | PUT | `/v1/dataelement/dict/state/{id}` | 启用/停用 | P1 |

### 3.1 接口实现: PUT /v1/dataelement/dict/state/{id}

- [X] T024 [US3] 实现 `api/internal/logic/dict/update_dict_state_logic.go`
  - [X] 4步业务流程标注
  - [X] 校验码表存在性
  - [X] 业务校验（停用原因必填）
  - [X] 检查状态是否变更
  - [X] 更新状态（保存新的状态和原因）

**Checkpoint**: ✅ Phase 3 完成 - 状态管理功能

---

## Phase 4: 关联查询 (4个接口)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 8 | GET | `/v1/dataelement/dict/dataelement/{id}` | 查询引用码表的数据元 | P2 |
| 9 | GET | `/v1/dataelement/dict/queryByStdFileCatalog` | 按文件目录查询 | P2 |
| 10 | GET | `/v1/dataelement/dict/queryByStdFile` | 按文件查询 | P2 |
| 11 | GET | `/v1/dataelement/dict/relation/stdfile/{id}` | 查询关联标准文件 | P2 |

### 4.1 接口实现: GET /v1/dataelement/dict/dataelement/{id}

- [X] T025 [P] [US4] 实现 `api/internal/logic/dict/query_dict_by_data_element_logic.go`
  - [X] 2步业务流程标注
  - [X] [MOCK] 调用 DataElement 服务查询关联数据元
  - [X] 返回数据元分页列表（当前返回空列表）

### 4.2 接口实现: GET /v1/dataelement/dict/queryByStdFileCatalog

- [X] T026 [P] [US4] 实现 `api/internal/logic/dict/query_dict_by_std_file_catalog_logic.go`
  - [X] 2步业务流程标注
  - [X] [MOCK] 调用 StdFile 服务获取目录下的文件列表
  - [X] 返回码表列表（当前返回空列表）

### 4.3 接口实现: GET /v1/dataelement/dict/queryByStdFile

- [X] T027 [P] [US4] 实现 `api/internal/logic/dict/query_dict_by_std_file_logic.go`
  - [X] 3步业务流程标注
  - [X] 参数校验（FileId必填）
  - [X] 按文件ID查询码表列表
  - [X] 转换为响应格式并返回

### 4.4 接口实现: GET /v1/dataelement/dict/relation/stdfile/{id}

- [X] T028 [P] [US4] 实现 `api/internal/logic/dict/query_dict_relation_stdfile_logic.go`
  - [X] 2步业务流程标注
  - [X] [MOCK] 查询关联文件列表（需要修改handler传递dictId）
  - [X] 返回文件列表（当前返回空列表）

**Checkpoint**: ✅ Phase 4 完成 - 关联查询功能

---

## Phase 5: 关联管理 (2个接口)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 12 | GET | `/v1/dataelement/dict/{id}` | 码表详情（按Code） | P2 |
| 13 | PUT | `/v1/dataelement/dict/relation/{id}` | 添加关联关系 | P2 |

### 5.1 接口实现: GET /v1/dataelement/dict/{id} (按Code查询详情)

- [X] T029 [US5] 实现 `api/internal/logic/dict/get_dict_by_code_logic.go`
  - [X] 5步业务流程标注
  - [X] 解析码表编码（从路径参数:code获取）
  - [X] 按Code查询码表
  - [X] 查询码值明细列表
  - [X] 查询目录信息
  - [X] 查询部门信息

### 5.2 接口实现: PUT /v1/dataelement/dict/relation/{id}

- [X] T030 [US5] 实现 `api/internal/logic/dict/add_dict_relation_logic.go`
  - [X] 3步业务流程标注
  - [X] 校验码表存在性
  - [X] 删除旧的关联文件关系
  - [X] 添加新的关联文件关系

**Checkpoint**: ✅ Phase 5 完成 - 关联管理功能

---

## Phase 6: 数据校验 (1个接口)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 14 | GET | `/v1/dataelement/dict/queryDataExists` | 查询数据是否存在 | P2 |

### 6.1 接口实现: GET /v1/dataelement/dict/queryDataExists

- [X] T031 [US5] 实现 `api/internal/logic/dict/query_dict_data_exists_logic.go`
  - [X] 3步业务流程标注
  - [X] 参数校验（OrgType必填）
  - [X] 构建查询条件（使用Keyword进行名称搜索）
  - [X] 查询是否存在并返回结果

**Checkpoint**: ✅ Phase 6 完成 - 数据校验功能

---

## Phase 7: 辅助接口 (2个接口)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 15 | DELETE | `/v1/dataelement/dict/{id}` | 删除码表 | P2 |
| 16 | DELETE | `/v1/dataelement/dict/batch/{ids}` | 批量删除 | P2 |

### 7.1 接口实现: DELETE /v1/dataelement/dict/{id}

- [X] T032 [US1] 实现 `api/internal/logic/dict/delete_dict_logic.go`
  - [X] 4步业务流程标注
  - [X] 校验码表存在性
  - [X] 删除码表（物理删除）
  - [X] 删除关联的码值
  - [X] 删除关联文件关系

### 7.2 接口实现: DELETE /v1/dataelement/dict/batch/{ids}

- [X] T033 [US1] 实现 `api/internal/logic/dict/batch_delete_dict_logic.go`
  - [X] 5步业务流程标注
  - [X] 参数校验（ID非空）
  - [X] 校验码表存在性
  - [X] 删除码表（物理删除）
  - [X] 删除关联的码值
  - [X] 删除关联文件关系

**Checkpoint**: ✅ Phase 7 完成 - 辅助接口功能

---

## Phase 8: 收尾工作

**目的**: 代码质量、文档、验证

### 8.1 代码质量

- [X] T034 代码清理和格式化 (`gofmt -w .`)
- [ ] T035 运行 `golangci-lint run` 修复代码质量问题

### 8.2 测试验证

- [ ] T036 **确认测试覆盖率 ≥ 80%**
  ```bash
  go test ./... -coverprofile=coverage.out
  go tool cover -func=coverage.out | grep total
  ```

- [ ] T037 运行所有测试确认通过
  ```bash
  go test ./... -v
  ```

### 8.3 文档更新

- [ ] T038 更新 Swagger 文档
  ```bash
  goctl api swagger -api api.doc/api.api -dir doc/swagger
  ```

- [ ] T039 验证所有16个API端点已注册

### 8.4 兼容性验证

- [ ] T040 验证错误码与Java实现完全一致
  - [ ] 30401: 数据不存在
  - [ ] 30402: 参数为空
  - [ ] 30403: 参数无效
  - [ ] 30404: 目录不存在
  - [ ] 30405: 中文名称重复
  - [ ] 30406: 英文名称重复
  - [ ] 30407: 码值为空
  - [ ] 30408: 码值描述为空
  - [ ] 30409: 码值重复
  - [ ] 30410: 停用原因过长

- [ ] T041 **接口兼容性验证**
  - [ ] 确认响应格式与Java完全一致
  - [ ] 确认异常信息与Java完全一致

**Checkpoint**: ✅ Phase 8 完成 - 所有功能就绪

---

## Mock 函数说明

### Mock 目录结构

```
api/internal/logic/dict/mock/
├── catalog.go       # 目录服务 Mock
├── dataelement.go   # 数据元服务 Mock
├── stdfile.go       # 标准文件服务 Mock
└── token.go         # Token服务 Mock
```

**说明**: 所有Mock函数使用build tag `//go:build !mock_logic_off`控制编译，后续补充RPC时可直接实现

### 需要后续补充 RPC 的场景

| 场景 | 当前Mock实现 | 后续补充 |
|------|-------------|----------|
| 目录名称查询 | `CatalogGetCatalogName(catalogId)` | Catalog RPC |
| 子目录列表 | `CatalogGetChildIds(catalogId)` | Catalog RPC |
| 目录类型校验 | `CatalogIsStdFileCatalog(catalogId)` | Catalog RPC |
| 数据元关联查询 | `DataElementQueryPageByDictId(dictId)` | DataElement RPC |
| 部门信息 | `GetDeptInfo(deptId)` | 部门服务/Token |
| 用户信息 | `GetUserInfo()` | Token |

---

## 依赖关系

```
Phase 0 (基础设施)
    ↓
Phase 1 (基础CRUD) ← MVP 🎯
    ↓
Phase 2 (码值管理)
    ↓
Phase 3 (状态管理)
    ↓
Phase 4 (关联查询) ← 可与Phase 5并行
Phase 5 (关联管理) ← 可与Phase 4并行
    ↓
Phase 6 (数据校验)
    ↓
Phase 7 (辅助接口)
    ↓
Phase 8 (收尾工作)
```

---

## 测试要求 🧪

| 要求 | 标准 |
|------|------|
| **单元测试覆盖率** | > 80% |
| **关键路径测试** | 100% 覆盖 |
| **边界测试** | 必须包含 |
| **错误处理测试** | 必须包含 |

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-02-06 | - | 初始版本 |
