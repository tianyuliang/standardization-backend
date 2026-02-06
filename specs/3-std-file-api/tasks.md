# 标准文件管理 (std-file-api) Tasks

> **Branch**: `3-std-file-api`
> **Spec Path**: `specs/3-std-file-api/`
> **Created**: 2026-02-06
> **Input**: spec.md, plan.md

---

## 任务组织说明

**组织方式**: 按接口增量定义维度，每个接口独立完成 API → Model → Logic → Test

**Mock 策略**: 对于依赖其他服务的逻辑（catalog、dataelement、dict、rule），使用 Mock 函数统一收集在 `logic/stdfile/mock/` 目录下，方便后续更改为正常逻辑。

**永远使用的 goctl 命令** (重要！！！):
```bash
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

**任务标记**:
- `[P]` = 可并行执行
- `[TEST]` = 测试任务
- `[MOCK]` = 使用 Mock 数据，后续补充 RPC

---

## Task Overview

| 阶段 | 描述 | 任务数 | 预计工作量 | 状态 |
|------|------|--------|------------|------|
| Phase 0 | 基础设施 | 5 | 0.5天 | ⏳ 待开始 |
| Phase 1 | 基础CRUD (4接口) | 16 | 3天 | ⏳ 待开始 |
| Phase 2 | 状态管理 (2接口) | 4 | 1天 | ⏳ 待开始 |
| Phase 3 | 目录移动 (1接口) | 2 | 0.5天 | ⏳ 待开始 |
| Phase 4 | 文件下载 (2接口) | 4 | 1天 | ⏳ 待开始 |
| Phase 5 | 关联查询 (3接口) | 6 | 1.5天 | ⏳ 待开始 |
| Phase 6 | 关联管理 (2接口) | 4 | 1天 | ⏳ 待开始 |
| Phase 7 | 辅助接口 (1接口) | 2 | 0.5天 | ⏳ 待开始 |
| Phase 8 | 收尾工作 | 4 | 0.5天 | ⏳ 待开始 |
| **总计** | | **47** | **约10天** | **0%** |

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
- [ ] T004 [P] 确认 OSS 配置已就绪

### P003 - 错误码定义

- [ ] T005 创建 `api/internal/errorx/codes.go`
  - [ ] 定义错误码 30200-30299
  - [ ] 实现错误辅助函数
  - [ ] 参考规则模块错误码定义

**Checkpoint**: ⏳ 基础设施待就绪

---

## Phase 1: 基础 CRUD (4个接口)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 1 | POST | `/v1/std-file` | 新增标准文件 | P1 |
| 2 | PUT | `/v1/std-file/{id}` | 修改标准文件 | P1 |
| 3 | GET | `/v1/std-file` | 分页列表查询 | P1 |
| 4 | GET | `/v1/std-file/{id}` | 详情查询 | P1 |

### 1.1 API 定义

- [ ] T006 创建 `api/doc/stdfile/stdfile.api`
  - [ ] 定义基础类型: CreateStdFileReq, UpdateStdFileReq, StdFileResp, StdFileListResp
  - [ ] 定义 StdFileRelationDto 关联类型
  - [ ] 定义 16 个 API 端点
  - [ ] 配置路由: `@server(prefix: /api/standardization/v1, group: stdfile)`

- [ ] T007 在 `api/doc/api.api` 中 import stdfile 模块

- [ ] T008 **永远使用这个命令** 运行 goctl 生成 Handler/Types
  ```bash
  goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
  ```

### 1.2 DDL 定义

- [ ] T009 [P] 创建 `migrations/stdfile/raw/t_std_file.sql`
  - [ ] 复用 Java 表结构

### 1.3 Model 层

- [ ] T010 创建 `model/stdfile/stdfile/` 目录结构
  - [ ] `interface.go` - StdFileModel 接口
  - [ ] `types.go` - StdFile、StdFileVo
  - [ ] `vars.go` - 枚举常量、错误码
  - [ ] `factory.go` - 工厂函数

- [ ] T011 实现 `model/stdfile/stdfile/sql_model.go`
  - [ ] Insert, FindOne, Update, Delete
  - [ ] FindByIds, FindByNumber, FindByNameAndOrgType
  - [ ] FindByCatalogIds (分页)
  - [ ] FindDataExists

- [ ] T012 **[TEST]** `model/stdfile/stdfile/sql_model_test.go`
  - [ ] Test Insert
  - [ ] Test FindOne
  - [ ] Test Update
  - [ ] Test FindByCatalogIds

### 1.4 公共 Logic (common.go)

- [ ] T013 创建 `api/internal/logic/stdfile/common.go`
  - [ ] dateToStr、timeToStr
  - [ ] buildStdFileResp
  - [ ] CheckNumberUnique (标准编号唯一性校验)
  - [ ] CheckNameUnique (文件名称唯一性校验)
  - [ ] CheckCatalogIdExist (目录存在性校验)
  - [ ] CheckVersionChange (版本变更检测)
  - [ ] isAllowedFileType (文件类型校验)

### 1.5 接口实现: POST /v1/std-file (新增标准文件)

- [ ] T014 实现 `api/internal/logic/stdfile/create_std_file_logic.go`
  - [ ] 8步业务流程标注
  - [ ] 标准编号唯一性校验
  - [ ] 文件名称唯一性校验
  - [ ] 目录存在性校验
  - [ ] 文件类型校验
  - [ ] 文件上传到OSS
  - [ ] TODO: 部门ID处理

- [ ] T015 **[TEST]** `api/internal/logic/stdfile/create_std_file_logic_test.go`

### 1.6 接口实现: PUT /v1/std-file/{id} (修改标准文件)

- [ ] T016 [P] 实现 `api/internal/logic/stdfile/update_std_file_logic.go`
  - [ ] 9步业务流程标注
  - [ ] 校验存在性
  - [ ] 标准编号唯一性校验（排除自身）
  - [ ] 文件名称唯一性校验（排除自身）
  - [ ] 目录存在性校验
  - [ ] 版本变更检测
  - [ ] 文件更新处理

- [ ] T017 [P] **[TEST]** `api/internal/logic/stdfile/update_std_file_logic_test.go`

### 1.7 接口实现: GET /v1/std-file (分页列表查询)

- [ ] T018 [P] 实现 `api/internal/logic/stdfile/list_std_file_logic.go`
  - [ ] 6步业务流程标注
  - [ ] 调用 Catalog RPC 获取子目录列表
  - [ ] 查询列表
  - [ ] 批量查询目录名称、部门信息

- [ ] T019 [P] **[TEST]** `api/internal/logic/stdfile/list_std_file_logic_test.go`

### 1.8 接口实现: GET /v1/std-file/{id} (详情查询)

- [ ] T020 [P] 实现 `api/internal/logic/stdfile/get_std_file_logic.go`
  - [ ] 5步业务流程标注
  - [ ] 查询文件
  - [ ] 查询目录名称
  - [ ] 查询部门信息

- [ ] T021 [P] **[TEST]** `api/internal/logic/stdfile/get_std_file_logic_test.go`

### 1.9 ServiceContext 更新

- [ ] T022 更新 `api/internal/svc/service_context.go`
  - [ ] 添加 StdFileModel
  - [ ] 初始化 DB 连接 (*sqlx.DB)
  - [ ] 初始化 Model 实例
  - [ ] TODO: 后续补充 OSS Client
  - [ ] TODO: 后续补充 RPC 客户端

**Checkpoint**: ⏳ Phase 1 待开始

---

## Phase 2: 状态管理 (2个接口)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 5 | PUT | `/v1/std-file/state/{id}` | 启用/停用 | P1 |
| 6 | PUT | `/v1/std-file/batchState` | 批量启用/停用 | P2 |

### 2.1 接口实现: PUT /v1/std-file/state/{id}

- [ ] T023 实现 `api/internal/logic/stdfile/update_std_file_state_logic.go`
  - [ ] 5步业务流程标注
  - [ ] 校验存在性
  - [ ] 停用时必须填写原因
  - [ ] 停用原因长度校验
  - [ ] 更新状态

- [ ] T024 **[TEST]** `api/internal/logic/stdfile/update_std_file_state_logic_test.go`

### 2.2 接口实现: PUT /v1/std-file/batchState

- [ ] T025 实现 `api/internal/logic/stdfile/batch_state_std_file_logic.go`
  - [ ] 5步业务流程标注
  - [ ] 批量校验存在性
  - [ ] 停用时必须填写原因
  - [ ] 批量更新状态

- [ ] T026 **[TEST]** `api/internal/logic/stdfile/batch_state_std_file_logic_test.go`

**Checkpoint**: ⏳ Phase 2 待开始

---

## Phase 3: 目录移动 (1个接口)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 7 | POST | `/v1/std-file/catalog/remove` | 目录移动 | P2 |

### 3.1 接口实现: POST /v1/std-file/catalog/remove

- [ ] T027 实现 `api/internal/logic/stdfile/remove_std_file_catalog_logic.go`
  - [ ] 4步业务流程标注
  - [ ] 校验目录存在性
  - [ ] 校验文件存在性
  - [ ] 批量更新 catalog_id

- [ ] T028 **[TEST]** `api/internal/logic/stdfile/remove_std_file_catalog_logic_test.go`

**Checkpoint**: ⏳ Phase 3 待开始

---

## Phase 4: 文件下载 (2个接口)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 8 | GET | `/v1/std-file/download/{id}` | 下载附件 | P1 |
| 9 | POST | `/v1/std-file/downloadBatch` | 批量下载 | P2 |

### 4.1 接口实现: GET /v1/std-file/download/{id}

- [ ] T029 实现 `api/internal/logic/stdfile/download_std_file_logic.go`
  - [ ] 4步业务流程标注
  - [ ] 校验文件存在
  - [ ] URL类型返回错误
  - [ ] FILE类型从OSS下载

- [ ] T030 **[TEST]** `api/internal/logic/stdfile/download_std_file_logic_test.go`

### 4.2 接口实现: POST /v1/std-file/downloadBatch

- [ ] T031 实现 `api/internal/logic/stdfile/download_batch_std_file_logic.go`
  - [ ] 6步业务流程标注
  - [ ] 校验文件存在
  - [ ] 过滤URL类型
  - [ ] 从OSS下载文件
  - [ ] 处理文件名重复
  - [ ] 打包成ZIP

- [ ] T032 **[TEST]** `api/internal/logic/stdfile/download_batch_std_file_logic_test.go`

**Checkpoint**: ⏳ Phase 4 待开始

---

## Phase 5: 关联查询 (3个接口)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 10 | GET | `/v1/std-file/relation/de/{id}` | 查询关联数据元 | P2 |
| 11 | GET | `/v1/std-file/relation/dict/{id}` | 查询关联码表 | P2 |
| 12 | GET | `/v1/std-file/relation/rule/{id}` | 查询关联编码规则 | P2 |

### 5.1 接口实现: GET /v1/std-file/relation/de/{id}

- [ ] T033 [P] 实现 `api/internal/logic/stdfile/query_relation_de_logic.go`
  - [ ] 3步业务流程标注
  - [ ] 校验文件存在
  - [ ] 调用 DataElement RPC 查询关联数据元

- [ ] T034 [P] **[TEST]** `api/internal/logic/stdfile/query_relation_de_logic_test.go`

### 5.2 接口实现: GET /v1/std-file/relation/dict/{id}

- [ ] T035 [P] 实现 `api/internal/logic/stdfile/query_relation_dict_logic.go`
  - [ ] 3步业务流程标注
  - [ ] 校验文件存在
  - [ ] 调用 Dict RPC 查询关联码表

- [ ] T036 [P] **[TEST]** `api/internal/logic/stdfile/query_relation_dict_logic_test.go`

### 5.3 接口实现: GET /v1/std-file/relation/rule/{id}

- [ ] T037 [P] 实现 `api/internal/logic/stdfile/query_relation_rule_logic.go`
  - [ ] 3步业务流程标注
  - [ ] 校验文件存在
  - [ ] 调用 Rule RPC 查询关联编码规则

- [ ] T038 [P] **[TEST]** `api/internal/logic/stdfile/query_relation_rule_logic_test.go`

**Checkpoint**: ⏳ Phase 5 待开始

---

## Phase 6: 关联管理 (2个接口)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 13 | PUT | `/v1/std-file/relation/{id}` | 添加关联关系 | P2 |
| 14 | GET | `/v1/std-file/relation/{id}` | 查询关联关系 | P2 |

### 6.1 接口实现: PUT /v1/std-file/relation/{id}

- [ ] T039 [P] 实现 `api/internal/logic/stdfile/add_relation_logic.go`
  - [ ] 5步业务流程标注
  - [ ] 校验文件存在
  - [ ] 调用 DataElement RPC 添加关联
  - [ ] 调用 Dict RPC 添加关联
  - [ ] 调用 Rule RPC 添加关联

- [ ] T040 [P] **[TEST]** `api/internal/logic/stdfile/add_relation_logic_test.go`

### 6.2 接口实现: GET /v1/std-file/relation/{id}

- [ ] T041 [P] 实现 `api/internal/logic/stdfile/query_relations_logic.go`
  - [ ] 4步业务流程标注
  - [ ] 校验文件存在
  - [ ] 调用 DataElement RPC 查询关联ID
  - [ ] 调用 Dict RPC 查询关联ID
  - [ ] 调用 Rule RPC 查询关联ID

- [ ] T042 [P] **[TEST]** `api/internal/logic/stdfile/query_relations_logic_test.go`

**Checkpoint**: ⏳ Phase 6 待开始

---

## Phase 7: 辅助接口 (1个接口)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 15 | GET | `/v1/std-file/queryDataExists` | 检查数据是否存在 | P2 |

### 7.1 接口实现: GET /v1/std-file/queryDataExists

- [ ] T043 [P] 实现 `api/internal/logic/stdfile/query_data_exists_logic.go`
  - [ ] 3步业务流程标注
  - [ ] 部门ID路径处理
  - [ ] 检查是否存在 (支持 filter_id 排除自身)

- [ ] T044 [P] **[TEST]** `api/internal/logic/stdfile/query_data_exists_logic_test.go`

**Checkpoint**: ⏳ Phase 7 待开始

---

## Phase 8: 收尾工作

### 8.1 代码质量

- [ ] T045 代码清理和格式化 (`gofmt -w .`)
- [ ] T046 运行 `golangci-lint run` 修复代码质量问题

### 8.2 测试验证

- [ ] T047 **确认测试覆盖率 ≥ 80%**
  ```bash
  go test ./... -coverprofile=coverage.out
  go tool cover -func=coverage.out | grep total
  ```

- [ ] T048 运行所有测试确认通过
  ```bash
  go test ./... -v
  ```

### 8.3 文档更新

- [ ] T049 更新 Swagger 文档
  ```bash
  make swagger
  ```

- [ ] T050 验证所有16个API端点已注册

### 8.4 兼容性验证

- [ ] T051 验证错误码与Java实现完全一致

- [ ] T052 **接口兼容性验证**
  - [ ] 确认响应格式与Java完全一致
  - [ ] 确认异常信息与Java完全一致

**Checkpoint**: ⏳ Phase 8 待开始

---

## Mock 函数说明

### Mock 目录结构

```
api/internal/logic/stdfile/mock/
├── catalog.go    # 目录服务 Mock
├── dataelement.go # 数据元服务 Mock
├── dict.go       # 码表服务 Mock
└── rule.go       # 编码规则服务 Mock
```

### 需要后续补充 RPC 的场景

| 场景 | 当前Mock实现 | 后续补充 |
|------|-------------|----------|
| 目录名称查询 | `CatalogGetCatalogName(catalogId)` | Catalog RPC |
| 子目录列表 | `CatalogGetChildIds(catalogId)` | Catalog RPC |
| 目录类型校验 | `CatalogIsStdFileCatalog(catalogId)` | Catalog RPC |
| 数据元关联查询 | `DataElementQueryPageByFileId(fileId)` | DataElement RPC |
| 数据元关联添加 | `DataElementAddRelation(fileId, deIds)` | DataElement RPC |
| 码表关联查询 | `DictQueryPageByFileId(fileId)` | Dict RPC |
| 码表关联添加 | `DictAddRelation(fileId, dictIds)` | Dict RPC |
| 编码规则关联查询 | `RuleQueryPageByFileId(fileId)` | Rule RPC |
| 编码规则关联添加 | `RuleAddRelation(fileId, ruleIds)` | Rule RPC |
| 部门信息 | `GetDeptInfo(deptId)` | 部门服务/Token |

### Mock 示例

```go
//go:build !mock_logic_off
// +build !mock_logic_off

package mock

import (
    "context"
    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
)

// CatalogGetCatalogName 获取目录名称
// 对应 Java: deCatalogInfoService.getById(catalogId).getCatalogName()
func CatalogGetCatalogName(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) string {
    // MOCK: 返回格式化的目录名称
    // TODO: 调用 Catalog RPC 获取目录名称
    if catalogId == 44 {
        return "全部目录"
    }
    return ""
}
```

---

## 依赖关系

```
Phase 0 (基础设施)
    ↓
Phase 1 (基础CRUD) ← MVP 🎯
    ↓
Phase 2 (状态管理)
    ↓
Phase 3 (目录移动)
    ↓
Phase 4 (文件下载)
    ↓
Phase 5 (关联查询)
    ↓
Phase 6 (关联管理)
    ↓
Phase 7 (辅助接口)
    ↓
Phase 8 (收尾工作)
```

### 并行执行说明

- `[P]` 标记的任务可并行执行
- Phase 2-7 可并行开发（如有团队）
- 同一 Phase 内的接口可并行实现

---

## MVP 范围

**最小可交付版本**: Phase 0 + Phase 1

MVP 包含的核心功能：
- ⏳ 创建标准文件（FILE/URL）
- ⏳ 修改标准文件（版本控制）
- ⏳ 查询文件详情
- ⏳ 列表查询（多条件筛选）

---

## 测试要求 🧪

| 要求 | 标准 |
|------|------|
| **单元测试覆盖率** | ≥ 80% |
| **关键路径测试** | 100% 覆盖（CRUD、状态管理） |
| **边界测试** | 必须包含 |
| **错误处理测试** | 必须包含 |

### 测试命名规范

```
Test{Function}_{Scenario}_{ExpectedResult}
```

示例：
- `TestCreateStdFile_ValidInput_ReturnsFile`
- `TestCreateStdFile_DuplicateNumber_ReturnsError`
- `TestUpdateStdFile_WithChanges_VersionIncremented`

---

## 实施进度跟踪

| Week | Phase | 内容 | 完成度 |
|------|-------|------|--------|
| 1 | Phase 0-8 | 基础设施 + CRUD + 状态 + 目录 + 下载 + 关联 + 管理 + 辅助 + 收尾 | 0% (待开始) |

---

## 文档变更历史

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-02-06 | 按接口增量维度创建任务 |
