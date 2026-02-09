# 数据元管理 (dataelement-api) Tasks

> **Branch**: `6-std-dataelement-api`
> **Spec Path**: `specs/6-std-dataelement-api/`
> **Created**: 2026-02-09
> **Input**: spec.md, plan.md

---

## ⚠️ 重要声明

**永远使用以下 goctl 命令生成代码**:
```bash
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

此命令必须在更新 API 定义后立即执行，以确保 Handler 和 Types 代码同步更新。

---

## 任务组织说明

**组织方式**: 按接口增量定义维度，每个接口独立完成 API → Model → Logic → Test

**Mock 策略**: 对于依赖其他服务的逻辑（catalog、dict、rule、stdfile），使用 `mock/service.go` 标记 TODO，后续补充 RPC 调用

**任务标记**:
- `[P]` = 可并行执行
- `[TEST]` = 测试任务
- `[MOCK]` = 使用 Mock 数据，后续补充 RPC

---

## Task Overview

| 阶段 | 描述 | 任务数 | 预计工作量 | 状态 |
|------|------|--------|------------|------|
| Phase 0 | 基础设施 | 7 | 1天 | ✅ 完成 |
| Phase 1 | CRUD基础 (6接口) | 20 | 3天 | ✅ 完成 (20/20) |
| Phase 2 | 状态管理 (1接口) | 2 | 0.5天 | ✅ 完成 (2/2) |
| Phase 3 | 目录管理 (1接口) | 2 | 0.5天 | ✅ 完成 (2/2) |
| Phase 4 | 关联查询 (3接口) | 6 | 1.5天 | ✅ 完成 (6/6) |
| Phase 5 | 导入导出 (2接口) | 4 | 2天 | ✅ 完成 (4/4) |
| Phase 6 | 其他接口 (4接口) | 8 | 1.5天 | ✅ 完成 (8/8) |
| Phase 4 | 关联查询 (3接口) | 8 | 1.5天 | ⏳ 待开始 |
| Phase 5 | 导入导出 (3接口) | 8 | 2天 | ⏳ 待开始 |
| Phase 6 | 其他接口 (5接口) | 10 | 1.5天 | ⏳ 待开始 |
| Phase 7 | 收尾工作 | 5 | 1天 | ⏳ 待开始 |
| **总计** | | **62** | **约11天** | **8%** |

---

## Phase 0: 基础设施 ✅ 完成

**目的**: 项目初始化和基础配置

### P001 - 环境检查

- [x] T001 确认 Go-Zero 项目结构已就绪
  - [x] 确认 `api/` 目录存在
  - [x] 确认 `model/` 目录存在
  - [x] 确认 `go.mod` 已配置 Go-Zero 依赖

### P002 - 配置验证

- [x] T002 [P] 确认 base.api 已定义通用类型
- [x] T003 [P] 确认数据库连接已配置
- [x] T004 [P] 确认 Kafka 配置已就绪

### P003 - 错误码定义

- [x] T005 创建 `api/internal/errorx/dataelement.go`
  - [x] 定义错误标识（不使用数字码值）
  - [x] 实现错误辅助函数
  - [x] 格式: `DataElementNotExist()`, `ParameterEmpty()`, `InvalidParameter(field, message)`

### P004 - API 定义

- [x] T006 创建 `api/doc/dataelement/dataelement.api`
  - [x] 定义基础请求/响应类型
  - [x] 定义 19 个 API 端点
  - [x] 配置路由: `@server(prefix: /v1, group: dataelement)`

- [x] T007 在 `api/doc/api.api` 中 import dataelement 模块

- [x] T008 **永远使用以下命令生成代码**:
  ```bash
  goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
  ```

**Checkpoint**: Phase 0 完成后进入 Phase 1

---

## Phase 1: CRUD基础 (6个接口) ⏳ 待开始

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 1 | POST | `/v1/dataelement` | 创建数据元 | P1 |
| 2 | GET | `/v1/dataelement` | 分页检索数据元 | P1 |
| 3 | GET | `/v1/dataelement/detail` | 查看数据元详情 | P1 |
| 4 | PUT | `/v1/dataelement/{id}` | 编辑数据元 | P1 |
| 5 | DELETE | `/v1/dataelement/{ids}` | 删除数据元 | P1 |
| 6 | GET | `/v1/dataelement/internal/list` | 分页检索（内部） | P2 |

### 1.1 DDL 定义

- [x] T009 [P] 创建 `migrations/dataelement/raw/t_data_element_info.sql`
- [x] T010 [P] 创建 `migrations/dataelement/raw/t_relation_de_file.sql`

### 1.2 Model 层

- [x] T011 创建 `model/dataelement/dataelement/` 目录结构
  - [x] `interface.go` - DataElementModel 接口
  - [x] `types.go` - DataElement、RelationDeFile 数据结构
  - [x] `vars.go` - 枚举常量定义
  - [x] `factory.go` - SQLx 工厂函数

- [x] T012 实现 `model/dataelement/dataelement/sql_model.go`
  - [x] Insert, FindOne, FindOneByCode, Update, Delete
  - [x] FindByCatalogIds (分页)
  - [x] FindByIds, FindByCodes
  - [x] CheckNameCnExists, CheckNameEnExists
  - [x] FindDataExists

- [x] T013 创建 `model/dataelement/relation/` 目录结构
  - [x] `interface.go` - RelationDeFileModel 接口
  - [x] `types.go` - 数据结构定义
  - [x] `vars.go` - 常量定义
  - [x] `factory.go` - SQLx 工厂函数
  - [x] `sql_model.go` - SQLx 实现

### 1.3 Mock 服务

- [x] T014 创建 `api/internal/logic/dataelement/mock/service.go`
  - [x] CheckCatalogExist() - 目录存在性校验
  - [x] GetDictInfo() - 码表信息查询
  - [x] GetRuleInfo() - 编码规则信息查询
  - [x] GetStdFileIds() - 标准文件ID验证

### 1.4 公共 Logic

- [x] T015 创建 `api/internal/logic/dataelement/common.go`
  - [x] `CalculateValueRange()` - 值域计算
  - [x] `CheckVersionChange()` - 版本变更检测
  - [x] `SendMqMessage()` - MQ消息发送
  - [x] `BuildDataElementVo()` - 构建数据元视图对象

### 1.5 接口实现: POST /v1/dataelement (创建数据元)

- [x] T016 实现 `api/internal/logic/dataelement/create_data_element_logic.go`
  ```go
  func (l *CreateDataElementLogic) CreateDataElement(req *types.CreateDataElementReq) (*types.DataElementDetailVo, error) {
      // 1. 参数校验 (Handler 已完成)
      // 2. 目录存在性校验 (mock)
      // 3. 中文名称唯一性校验 (同stdType)
      // 4. 英文名称唯一性校验 (同部门)
      // 5. 关联校验 (码表/规则是否存在) (mock)
      // 6. 文件ID有效性校验 (mock)
      // 7. 生成ID和Code
      // 8. 保存数据元基本信息
      // 9. 保存关联文件关系
      // 10. 发送MQ消息
      // 11. 返回详情
  }
  ```

- [x] T017 **[TEST]** `api/internal/logic/dataelement/create_data_element_logic_test.go`

### 1.6 接口实现: GET /v1/dataelement (分页检索)

- [x] T018 [P] 实现 `api/internal/logic/dataelement/list_data_element_logic.go`
  ```go
  func (l *ListDataElementLogic) ListDataElement(req *types.PageInfoWithKeyword) (*types.DataElementListResp, error) {
      // 1. 处理目录ID（获取当前目录及所有子目录ID列表）
      // 2. 构建查询条件
      // 3. 分页查询
      // 4. 查询关联的码表、规则、文件、部门信息
      // 5. 构建响应
  }
  ```

- [x] T019 [P] **[TEST]** `api/internal/logic/dataelement/list_data_element_logic_test.go`

### 1.7 接口实现: GET /v1/dataelement/detail (查看详情)

- [x] T020 [P] 实现 `api/internal/logic/dataelement/get_data_element_detail_logic.go`
  ```go
  func (l *GetDataElementDetailLogic) GetDataElementDetail(req *types.DataElementDetailReq) (*types.DataElementDetailVo, error) {
      // 1. 校验查询参数
      // 2. 查询数据元基本信息（按ID或Code）
      // 3. 查询关联码表信息 (mock)
      // 4. 查询关联编码规则信息 (mock)
      // 5. 查询关联目录信息 (mock)
      // 6. 计算值域
      // 7. 查询部门信息
      // 8. 查询分级标签信息
      // 9. 构建响应
  }
  ```

- [x] T021 [P] **[TEST]** `api/internal/logic/dataelement/get_data_element_detail_logic_test.go`

### 1.8 接口实现: PUT /v1/dataelement/{id} (编辑数据元)

- [x] T022 [P] 实现 `api/internal/logic/dataelement/update_data_element_logic.go`
  ```go
  func (l *UpdateDataElementLogic) UpdateDataElement(req *types.UpdateDataElementReq) (*types.DataElementDetailVo, error) {
      // 1. 校验数据元存在性
      // 2. 目录存在性校验 (mock)
      // 3. 关联校验 (mock)
      // 4. 判断版本是否变更
      //    - 关键属性变更（名称/关联类型/关联对象）：版本号+1
      // 5. 更新数据元基本信息
      // 6. 更新关联文件关系
      // 7. 发送MQ消息
      // 8. 返回详情
  }
  ```

- [x] T023 [P] **[TEST]** `api/internal/logic/dataelement/update_data_element_logic_test.go`

### 1.9 接口实现: DELETE /v1/dataelement/{ids} (删除数据元)

- [x] T024 [P] 实现 `api/internal/logic/dataelement/delete_data_element_logic.go`
  ```go
  func (l *DeleteDataElementLogic) DeleteDataElement(ids string) error {
      // 1. 校验ID列表
      // 2. 物理删除数据元
      // 3. 同步删除关联文件关系
      // 4. 发送MQ消息（delete类型）
      // 5. 返回成功
  }
  ```

- [x] T025 [P] **[TEST]** `api/internal/logic/dataelement/delete_data_element_logic_test.go`

### 1.10 接口实现: GET /v1/dataelement/internal/list (内部分页检索)

- [x] T026 [P] 实现 `api/internal/logic/dataelement/list_data_element_internal_logic.go`

- [x] T027 [P] **[TEST]** `api/internal/logic/dataelement/list_data_element_internal_logic_test.go`

**Checkpoint**: Phase 1 CRUD基础完成

---

## Phase 2: 状态管理 (1个接口) ⏳ 待开始

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 7 | PUT | `/v1/dataelement/state/{ids}` | 启用/停用 | P1 |

### 2.1 接口实现: PUT /v1/dataelement/state/{ids}

- [x] T028 实现 `api/internal/logic/dataelement/update_state_logic.go`
  ```go
  func (l *UpdateStateLogic) UpdateState(ids string, req *types.UpdateStateReq) (*types.EmptyResp, error) {
      // 1. 校验ID列表，转换为int64数组
      // 2. state=disable时校验停用原因（必填且<=800字符）
      // 3. 批量更新状态
      // 4. 发送MQ消息
      // 5. 返回成功
  }
  ```

- [x] T029 **[TEST]** `api/internal/logic/dataelement/update_state_logic_test.go`

**Checkpoint**: Phase 2 状态管理完成

---

## Phase 3: 目录管理 (1个接口) ⏳ 待开始

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 8 | POST | `/v1/dataelement/catalog/remove` | 移动目录 | P2 |

### 3.1 接口实现: POST /v1/dataelement/catalog/remove

- [x] T030 实现 `api/internal/logic/dataelement/remove_catalog_logic.go`
  ```go
  func (l *RemoveCatalogLogic) RemoveCatalog(req *types.RemoveCatalogReq) (*types.EmptyResp, error) {
      // 1. 校验ID列表
      // 2. 校验目标目录存在性 (mock)
      // 3. 批量更新目录ID（版本号+1）
      // 4. 发送MQ消息
      // 5. 返回成功
  }
  ```

- [x] T031 **[TEST]** `api/internal/logic/dataelement/remove_catalog_logic_test.go`

**Checkpoint**: Phase 3 目录管理完成

---

## Phase 4: 关联查询 (3个接口) ⏳ 待开始

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 9 | GET | `/v1/dataelement/query/byStdFileCatalog` | 按文件目录检索 | P2 |
| 10 | GET | `/v1/dataelement/query/byStdFile` | 按文件检索 | P2 |
| 11 | GET | `/v1/dataelement/query/stdFile` | 查询关联文件 | P2 |

### 4.1 接口实现: GET /v1/dataelement/query/byStdFileCatalog

- [x] T032 实现 `api/internal/logic/dataelement/query_by_std_file_catalog_logic.go`

- [x] T033 **[TEST]** `api/internal/logic/dataelement/query_by_std_file_catalog_logic_test.go`

### 4.2 接口实现: GET /v1/dataelement/query/byStdFile

- [x] T034 [P] 实现 `api/internal/logic/dataelement/query_by_std_file_logic.go`

- [x] T035 [P] **[TEST]** `api/internal/logic/dataelement/query_by_std_file_logic_test.go`

### 4.3 接口实现: GET /v1/dataelement/query/stdFile

- [x] T036 [P] 实现 `api/internal/logic/dataelement/query_std_file_logic.go`

- [x] T037 [P] **[TEST]** `api/internal/logic/dataelement/query_std_file_logic_test.go`

**Checkpoint**: Phase 4 关联查询完成

---

## Phase 5: 导入导出 (3个接口) ⏳ 待开始

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 12 | POST | `/v1/dataelement/import` | 批量导入 | P2 |
| 13 | POST | `/v1/dataelement/export` | 导出(支持目录/ID) | P2 |

### 5.1 接口实现: POST /v1/dataelement/import

- [x] T038 实现 `api/internal/logic/dataelement/import_data_element_logic.go`
  ```go
  func (l *ImportDataElementLogic) ImportDataElement(req *types.ImportDataElementReq) (*types.ImportResultVo, error) {
      // 1. 校验目录存在性 (mock)
      // 2. 校验文件类型和大小
      // 3. 解析Excel文件
      // 4. 逐行校验数据
      // 5. 批量保存校验通过的数据
      // 6. 返回导入结果（成功/失败列表）
      // 7. 发送MQ消息
  }
  ```

- [x] T039 **[TEST]** `api/internal/logic/dataelement/import_data_element_logic_test.go`

### 5.2 接口实现: POST /v1/dataelement/export (支持按目录和按ID导出)

- [x] T040 [P] 实现 `api/internal/logic/dataelement/export_data_element_logic.go`

- [x] T041 [P] **[TEST]** `api/internal/logic/dataelement/export_data_element_logic_test.go`

**注意**: 单个导出接口同时支持按目录导出(CatalogId)和按ID导出(Ids参数)

**Checkpoint**: Phase 5 导入导出完成

---

## Phase 6: 其他接口 (5个接口) ⏳ 待开始

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 15 | GET | `/v1/dataelement/query/isRepeat` | 检查重名 | P2 |
| 16 | DELETE | `/v1/dataelement/label/:id` | 删除标签 | P3 |
| 17 | GET | `/v1/dataelement/internal/detail/:id` | 内部详情 | P3 |
| 18 | POST | `/v1/dataelement/internal/getDataElementPageByRuleId` | 按规则ID查询 | P3 |

### 6.1 接口实现: GET /v1/dataelement/query/isRepeat

- [x] T044 实现 `api/internal/logic/dataelement/is_repeat_logic.go`

- [x] T045 **[TEST]** `api/internal/logic/dataelement/is_repeat_logic_test.go`

### 6.2 接口实现: DELETE /v1/dataelement/label/:id

- [x] T046 [P] 实现 `api/internal/logic/dataelement/delete_label_logic.go`

- [x] T047 [P] **[TEST]** `api/internal/logic/dataelement/delete_label_logic_test.go`

### 6.3 接口实现: GET /v1/dataelement/internal/detail/:id

- [x] T048 [P] 实现 `api/internal/logic/dataelement/get_data_element_internal_logic.go`

- [x] T049 [P] **[TEST]** `api/internal/logic/dataelement/get_data_element_internal_logic_test.go`

### 6.4 接口实现: POST /v1/dataelement/internal/getDataElementPageByRuleId

- [x] T050 [P] 实现 `api/internal/logic/dataelement/get_page_by_rule_id_logic.go`

- [x] T051 [P] **[TEST]** `api/internal/logic/dataelement/get_page_by_rule_id_logic_test.go`

**注意**: query/list 和 internal/query/list 接口功能已由 Phase 1 的 list_data_element_logic.go 和 list_data_element_internal_logic.go 实现

**Checkpoint**: Phase 6 其他接口完成

---

## Phase 7: 收尾工作 ⏳ 待开始

### 7.1 代码质量

- [ ] T054 代码清理和格式化 (`gofmt -w .`)
- [ ] T055 运行 `golangci-lint run` 修复代码质量问题

### 7.2 测试验证

- [ ] T056 **确认测试覆盖率 ≥ 80%**
  ```bash
  go test ./... -coverprofile=coverage.out
  go tool cover -func=coverage.out | grep total
  ```

- [ ] T057 运行所有测试确认通过
  ```bash
  go test ./... -v
  ```

### 7.3 文档更新

- [ ] T058 创建 `specs/6-std-dataelement-api/research.md`
  - [ ] 技术决策和迁移策略
  - [ ] Java兼容性说明

- [ ] T059 创建 `specs/6-std-dataelement-api/data-model.md`
  - [ ] 数据库表结构
  - [ ] Go Struct 定义

- [ ] T060 创建 `specs/6-std-dataelement-api/quickstart.md`
  - [ ] 开发环境准备
  - [ ] API 测试示例

- [ ] T061 创建 `specs/6-std-dataelement-api/IMPLEMENTATION.md`
  - [ ] 19个接口实现状态
  - [ ] 文件结构说明

- [ ] T062 创建 `specs/6-std-dataelement-api/contracts/dataelement-api.yaml`
  - [ ] OpenAPI 3.0 规范

### 7.4 兼容性验证

- [ ] T063 验证错误码与Java实现完全一致
- [ ] T064 **接口兼容性验证**
  - [ ] 确认响应格式与Java完全一致
  - [ ] 确认异常信息与Java完全一致

**Checkpoint**: Phase 7 收尾工作完成

---

## ServiceContext 更新

- [ ] T065 更新 `api/internal/svc/service_context.go`
  - [ ] 添加 DataElementModel
  - [ ] 添加 RelationDeFileModel
  - [ ] 初始化 DB 连接 (*sqlx.Conn)
  - [ ] 初始化 Model 实例

---

## 依赖关系

```
Phase 0 (基础设施)
    ↓
Phase 1 (CRUD基础)
    ↓
Phase 2 (状态管理)
    ↓
Phase 3 (目录管理)
    ↓
Phase 4 (关联查询)
    ↓
Phase 5 (导入导出)
    ↓
Phase 6 (其他接口)
    ↓
Phase 7 (收尾工作)
```

---

## 测试要求 🧪

| 要求 | 标准 |
|------|------|
| **单元测试覆盖率** | ≥ 80% |
| **关键路径测试** | 100% 覆盖（创建、查询、编辑、删除） |
| **边界测试** | 必须包含 |
| **错误处理测试** | 必须包含 |

---

## 实施进度跟踪

| Week | Phase | 内容 | 完成度 |
|------|-------|------|--------|
| 1 | Phase 0-3 | 基础设施 + CRUD + 状态 + 目录 | 0% |

---

## 文档变更历史

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-02-09 | 初始版本 |
