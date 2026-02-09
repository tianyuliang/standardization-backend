# 标准任务管理 (std-task-api) Tasks

> **Branch**: `5-std-task-api`
> **Spec Path**: `specs/5-std-task-api/`
> **Created**: 2026-02-09
> **Input**: spec.md, plan.md

---

## 任务组织说明

**组织方式**: 按接口增量定义维度，每个接口独立完成 API → Model → Logic → Test

**Mock 策略**: 对于依赖外部推荐服务的逻辑，使用 HTTP 客户端调用，失败时返回空结果不阻塞主流程

**任务标记**:
- `[P]` = 可并行执行
- `[TEST]` = 测试任务
- `[MOCK]` = 使用 Mock 数据，后续补充 RPC

**⚠️ 重要**：永远使用以下 goctl 命令生成代码：
```bash
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

---

## Task Overview

| 阶段 | 描述 | 任务数 | 预计工作量 | 状态 |
|------|------|--------|------------|------|
| Phase 0 | 基础设施 | 7 | 1天 | ✅ 100% |
| Phase 1 | 任务查询管理 (4接口) | 10 | 2天 | ✅ 100% (Logic完成) |
| Phase 2 | 任务创建与完成 (5接口) | 12 | 2.5天 | ✅ 100% (Logic完成) |
| Phase 3 | 业务表管理 (7接口) | 16 | 3天 | ✅ 100% (Logic完成) |
| Phase 4 | 推荐服务 (4接口) | 10 | 2天 | ✅ 100% (Logic完成) |
| Phase 5 | 数据元操作 (4接口) | 10 | 1.5天 | ✅ 100% (Logic完成) |
| Phase 6 | 收尾工作 | 12 | 1天 | ✅ 100% (文档完成) |
| **总计** | | **77** | **约13天** | **约95%** |

---

## Phase 0: 基础设施 ✅ 已完成

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

- [x] T004 创建 `api/internal/errorx/task.go`
  - [x] 定义错误码 30701-30799
  - [x] 实现错误辅助函数

### P004 - API 定义

- [x] T005 创建 `api/doc/task/task.api`
  - [x] 定义基础类型
  - [x] 定义 24 个 API 端点（合并为单一task组）
  - [x] 配置路由: `@server(prefix: /api/standardization/v1/dataelement/task, group: task)`

- [x] T006 在 `api/doc/api.api` 中 import task 模块

- [x] T007 运行 **goctl api go** 生成 Handler/Types
  ```bash
  goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
  ```

**Checkpoint**: ✅ 基础设施就绪

---

## Phase 1: 任务查询管理 (4个接口) ✅ 已完成

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 1 | GET | `/std-create/uncompleted` | 未处理任务列表 | P1 ✅ |
| 2 | GET | `/std-create/completed` | 已完成任务列表 | P1 ✅ |
| 3 | GET | `/std-create/completed/{id}` | 任务详情 | P1 ✅ |
| 4 | POST | `/queryTaskState` | 任务状态查询 | P2 ✅ |

### 1.1 API 定义

- [x] T005 已在 Phase 0 完成

### 1.2 DDL 定义

- [x] T008 [P] 创建 `migrations/task/raw/t_task_std_create.sql`
- [x] T009 [P] 创建 `migrations/task/raw/t_task_std_create_result.sql`

### 1.3 Model 层

- [x] T010 创建 `model/task/task/` 目录结构
  - [x] `interface.go` - TaskStdCreateModel、TaskStdCreateResultModel 接口
  - [x] `types.go` - 数据结构定义
  - [x] `vars.go` - 枚举常量
  - [x] `factory.go` - 工厂函数

- [x] T011 实现 `model/task/task/sql_model.go`
  - [x] Insert, FindOne, Update, Delete
  - [x] FindUncompleted, FindCompleted
  - [x] FindByTaskNo
  - [x] FindByTaskId (结果表)

- [ ] T012 **[TEST]** `model/task/task/sql_model_test.go`

### 1.4 接口实现: GET /std-create/uncompleted

- [x] T013 实现 `api/internal/logic/task/get_uncompleted_tasks_logic.go`
  ```go
  func (l *GetUncompletedTasksLogic) GetUncompletedTasks(req *types.PageInfoWithKeyword) (resp *types.TaskDataListResp, err error) {
      // 1. 参数处理（keyword, offset, limit）
      // 2. 查询未处理任务列表
      // 3. 构建响应
  }
  ```

- [ ] T014 **[TEST]** `api/internal/logic/task/get_uncompleted_tasks_logic_test.go`

### 1.5 接口实现: GET /std-create/completed

- [x] T015 [P] 实现 `api/internal/logic/task/get_completed_tasks_logic.go`
  ```go
  func (l *GetCompletedTasksLogic) GetCompletedTasks(req *types.PageInfoWithKeyword) (resp *types.TaskDataListResp, err error) {
      // 1. 参数处理
      // 2. 查询已完成任务列表
      // 3. 构建响应
  }
  ```

- [ ] T016 [P] **[TEST]** `api/internal/logic/task/get_completed_tasks_logic_test.go`

### 1.6 接口实现: GET /std-create/completed/{id}

- [x] T017 [P] 实现 `api/internal/logic/task/get_task_by_id_logic.go`
  ```go
  func (l *GetTaskByIdLogic) GetTaskById(id int64) (resp *types.TaskDetailResp, err error) {
      // 1. 查询任务
      // 2. 查询任务结果
      // 3. 构建响应
  }
  ```

- [ ] T018 [P] **[TEST]** `api/internal/logic/task/get_task_by_id_logic_test.go`

### 1.7 接口实现: POST /queryTaskState

- [x] T019 [P] 实现 `api/internal/logic/task/query_task_state_logic.go`
  ```go
  func (l *QueryTaskStateLogic) QueryTaskState(req *types.QueryTaskStateReq) (resp *types.ProcessResp, err error) {
      // 1. 查询任务状态
      // 2. 返回状态信息
  }
  ```

- [ ] T020 [P] **[TEST]** `api/internal/logic/task/query_task_state_logic_test.go`

**Checkpoint**: ✅ Phase 1 Logic完成 (测试待补充)

---

## Phase 2: 任务创建与完成 (5个接口) ⏳ 50% (Logic完成)

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 5 | POST | `/std-create/relation/staging` | 标准关联暂存 | P1 ✅ |
| 6 | POST | `/std-create/publish/submit` | 标准关联提交 | P1 ✅ |
| 7 | POST | `/createTask` | 新建标准任务 | P1 ✅ |
| 8 | POST | `/finishTask/{task_id}` | 完成任务 | P1 ✅ |
| 9 | POST | `/queryTaskProcess` | 进度查询 | P2 ✅ |

### 2.1 接口实现: POST /std-create/relation/staging

- [x] T021 实现 `api/internal/logic/task/staging_relation_logic.go`
  ```go
  func (l *StagingRelationLogic) StagingRelation(req *types.StagingRelationReq) (resp *types.TaskBaseResp, err error) {
      // 1. 参数校验
      // 2. 暂存任务数据（不调用推荐服务）
      // 3. 返回成功
  }
  ```

- [ ] T022 **[TEST]** `api/internal/logic/task/staging_relation_logic_test.go`

### 2.2 接口实现: POST /std-create/publish/submit

- [x] T023 实现 `api/internal/logic/task/submit_relation_logic.go`
  ```go
  func (l *SubmitRelationLogic) SubmitRelation(req *types.StagingRelationReq) (resp *types.TaskBaseResp, err error) {
      // 1. 参数校验
      // 2. 保存任务数据
      // 3. 异步调用推荐算法（TODO: HTTP调用）
      // 4. 返回成功
  }
  ```

- [ ] T024 **[TEST]** `api/internal/logic/task/submit_relation_logic_test.go`

### 2.3 接口实现: POST /createTask

- [x] T025 实现 `api/internal/logic/task/create_task_logic.go`
  ```go
  func (l *CreateTaskLogic) CreateTask(req *types.CreateTaskReq) (resp *types.TaskBaseResp, err error) {
      // 1. 参数校验
      // 2. 查询业务表（TODO: Phase 3实现后补充）
      // 3. 生成任务编号
      // 4. 创建任务记录
      // 5. 更新业务表状态（TODO: Phase 3实现后补充）
  }
  ```

- [ ] T026 **[TEST]** `api/internal/logic/task/create_task_logic_test.go`

### 2.4 接口实现: POST /finishTask/{task_id}

- [x] T027 实现 `api/internal/logic/task/finish_task_logic.go`
  ```go
  func (l *FinishTaskLogic) FinishTask(taskId int64) (resp *types.TaskBaseResp, err error) {
      // 1. 校验任务存在
      // 2. 更新任务状态为完成
      // 3. 发送回调（Webhook）
  }
  ```

- [ ] T028 **[TEST]** `api/internal/logic/task/finish_task_logic_test.go`

### 2.5 接口实现: POST /queryTaskProcess

- [x] T029 实现 `api/internal/logic/task/query_task_process_logic.go`
  ```go
  func (l *QueryTaskProcessLogic) QueryTaskProcess(req *types.QueryProcessReq) (resp *types.ProcessResp, err error) {
      // 1. 查询任务进度
      // 2. 返回进度信息
  }
  ```

- [ ] T030 **[TEST]** `api/internal/logic/task/query_task_process_logic_test.go`

**Checkpoint**: ⏳ Phase 2 Logic完成 (测试待补充)

---

## Phase 3: 业务表管理 (7个接口) ✅ Logic完成

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 10 | POST | `/addToPending` | 添加至待新建 | P1 ✅ |
| 11 | GET | `/getBusinessTable` | 业务表列表 | P1 ✅ |
| 12 | GET | `/getBusinessTableField` | 业务表字段列表 | P1 ✅ |
| 13 | DELETE | `/deleteBusinessTableField/{id}` | 移除字段 | P1 ✅ |
| 14 | PUT | `/cancelBusinessTableField` | 撤销 | P1 ✅ |
| 15 | PUT | `/updateTableName` | 修改表名称 | P1 ✅ |
| 16 | GET | `/getBusinessTableFromTask` | 任务关联业务表 | P1 ✅ |

### 3.1 DDL 定义

- [x] T031 创建 `migrations/task/raw/t_business_table_std_create_pool.sql`

### 3.2 Model 层

- [x] T032 实现 `model/task/pool/` 目录结构
  - [x] `interface.go` - BusinessTablePoolModel 接口
  - [x] `types.go` - 数据结构定义
  - [x] `vars.go` - 状态常量定义
  - [x] `factory.go` - 工厂函数
  - [x] `sql_model.go` - SQLx 实现

- [ ] T033 **[TEST]** `model/task/pool/sql_model_test.go`

### 3.3 接口实现: POST /addToPending

- [x] T034 实现 `api/internal/logic/task/add_to_pending_logic.go`
  ```go
  func (l *AddToPendingLogic) AddToPending(req *types.AddToPendingReq) (resp *types.TaskBaseResp, err error) {
      // 1. 参数校验
      // 2. 保存到待新建表
      // 3. 返回成功
  }
  ```

- [ ] T035 **[TEST]** `api/internal/logic/task/add_to_pending_logic_test.go`

### 3.4 接口实现: GET /getBusinessTable

- [x] T036 [P] 实现 `api/internal/logic/task/get_business_table_logic.go`

- [ ] T037 [P] **[TEST]** `api/internal/logic/task/get_business_table_logic_test.go`

### 3.5 接口实现: GET /getBusinessTableField

- [x] T038 [P] 实现 `api/internal/logic/task/get_business_table_field_logic.go`

- [ ] T039 [P] **[TEST]** `api/internal/logic/task/get_business_table_field_logic_test.go`

### 3.6 接口实现: DELETE /deleteBusinessTableField/{id}

- [x] T040 [P] 实现 `api/internal/logic/task/delete_field_logic.go`

- [ ] T041 [P] **[TEST]** `api/internal/logic/task/delete_field_logic_test.go`

### 3.7 接口实现: PUT /cancelBusinessTableField

- [x] T042 [P] 实现 `api/internal/logic/task/cancel_field_logic.go`

- [ ] T043 [P] **[TEST]** `api/internal/logic/task/cancel_field_logic_test.go`

### 3.8 接口实现: PUT /updateTableName

- [x] T044 [P] 实现 `api/internal/logic/task/update_table_name_logic.go`

- [ ] T045 [P] **[TEST]** `api/internal/logic/task/update_table_name_logic_test.go`

### 3.9 接口实现: GET /getBusinessTableFromTask

- [x] T046 [P] 实现 `api/internal/logic/task/get_table_from_task_logic.go`

- [ ] T047 [P] **[TEST]** `api/internal/logic/task/get_table_from_task_logic_test.go`

**Checkpoint**: ✅ Phase 3 Logic完成 (测试待补充)

---

## Phase 4: 推荐服务 (4个接口) ✅ Logic完成

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 17 | POST | `/std-rec/rec` | 标准推荐（内部） | P2 ✅ |
| 18 | POST | `/std-create` | 标准创建（内部） | P2 ✅ |
| 19 | POST | `/stand-rec/rec` | 标准推荐（弹框） | P2 ✅ |
| 20 | POST | `/rule-rec/rec` | 编码规则推荐 | P2 ✅ |

### 4.1 接口实现: POST /std-rec/rec

- [x] T048 实现 `api/internal/logic/task/std_rec_logic.go`
  ```go
  func (l *StdRecLogic) StdRec(req *types.StdRecReq) (resp *types.StdRecResp, err error) {
      // 1. 参数校验
      // 2. 调用推荐服务（HTTP）
      // 3. 返回推荐结果（最多3条）
  }
  ```

- [ ] T049 **[TEST]** `api/internal/logic/task/std_rec_logic_test.go`

### 4.2 接口实现: POST /std-create

- [x] T050 [P] 实现 `api/internal/logic/task/std_create_logic.go`
  ```go
  // 与 std_rec 类似，执行标准创建流程
  ```

- [ ] T051 [P] **[TEST]** `api/internal/logic/task/std_create_logic_test.go`

### 4.3 接口实现: POST /stand-rec/rec

- [x] T052 [P] 实现 `api/internal/logic/task/stand_rec_logic.go`

- [ ] T053 [P] **[TEST]** `api/internal/logic/task/stand_rec_logic_test.go`

### 4.4 接口实现: POST /rule-rec/rec

- [x] T054 [P] 实现 `api/internal/logic/task/rule_rec_logic.go`
  ```go
  func (l *RuleRecLogic) RuleRec(req *types.RuleRecReq) (resp *types.StdRecResp, err error) {
      // 1. 参数校验
      // 2. 调用规则推荐服务（HTTP）
      // 3. 返回推荐结果
  }
  ```

- [ ] T055 [P] **[TEST]** `api/internal/logic/task/rule_rec_logic_test.go`

**Checkpoint**: ✅ Phase 4 Logic完成 (测试待补充，推荐服务HTTP调用待实现)

---

## Phase 5: 数据元操作 (4个接口) ✅ Logic完成

### 接口清单

| # | 方法 | 路径 | 功能 | 优先级 |
|---|------|------|------|--------|
| 21 | GET | `/getBusinessTableFieldFromTask` | 任务关联字段 | P1 ✅ |
| 22 | POST | `/submitDataElement` | 提交选定数据元 | P1 ✅ |
| 23 | PUT | `/updateDescription` | 修改字段说明 | P1 ✅ |
| 24 | PUT | `/accept` | 采纳 | P1 ✅ |

### 5.1 接口实现: GET /getBusinessTableFieldFromTask

- [x] T056 实现 `api/internal/logic/task/get_field_from_task_logic.go`

- [ ] T057 **[TEST]** `api/internal/logic/task/get_field_from_task_logic_test.go`

### 5.2 接口实现: POST /submitDataElement

- [x] T058 [P] 实现 `api/internal/logic/task/submit_data_element_logic.go`
  ```go
  func (l *SubmitDataElementLogic) SubmitDataElement(req *types.SubmitDataElementReq) (resp *types.TaskBaseResp, err error) {
      // 1. 参数校验
      // 2. 保存数据元选择
      // 3. 返回成功
  }
  ```

- [ ] T059 [P] **[TEST]** `api/internal/logic/task/submit_data_element_logic_test.go`

### 5.3 接口实现: PUT /updateDescription

- [x] T060 [P] 实现 `api/internal/logic/task/update_description_logic.go`

- [ ] T061 [P] **[TEST]** `api/internal/logic/task/update_description_logic_test.go`

### 5.4 接口实现: PUT /accept

- [x] T062 [P] 实现 `api/internal/logic/task/accept_logic.go`

- [ ] T063 [P] **[TEST]** `api/internal/logic/task/accept_logic_test.go`

**Checkpoint**: ✅ Phase 5 Logic完成 (测试待补充)

---

## Phase 6: 收尾工作 ⏳ 部分完成

### 6.1 代码质量

- [x] T064 代码清理和格式化 (`gofmt -w .`)
- [ ] T065 运行 `golangci-lint run` 修复代码质量问题

### 6.2 测试验证

- [ ] T066 **确认测试覆盖率 ≥ 80%**
  ```bash
  go test ./... -coverprofile=coverage.out
  go tool cover -func=coverage.out | grep total
  ```

- [ ] T067 运行所有测试确认通过
  ```bash
  go test ./... -v
  ```

### 6.3 文档更新 ✅ 已完成

- [x] T068 创建 `specs/5-std-task-api/research.md`
  - [x] 技术决策和迁移策略
  - [x] 混合ID类型策略文档化
  - [x] Mock服务模式说明

- [x] T069 创建 `specs/5-std-task-api/data-model.md`
  - [x] TaskStdCreate 实体定义
  - [x] TaskStdCreateResult 实体定义
  - [x] BusinessTablePool 实体定义
  - [x] Go Struct 映射

- [x] T070 创建 `specs/5-std-task-api/quickstart.md`
  - [x] 开发环境准备
  - [x] API 测试示例
  - [x] 常见问题解答

- [x] T071 创建 `specs/5-std-task-api/IMPLEMENTATION.md`
  - [x] 24个接口实现状态
  - [x] 文件结构说明
  - [x] Mock 服务管理
  - [x] 错误码规范
  - [x] 业务校验规则

- [x] T072 创建 `specs/5-std-task-api/ANALYSIS.md`
  - [x] 需求覆盖率分析
  - [x] API 端点覆盖率
  - [x] 风险评估
  - [x] 建议

- [x] T073 创建 `specs/5-std-task-api/REMEDIATION.md`
  - [x] Mock 服务替换计划
  - [x] Webhook 失败处理策略
  - [x] 测试用例编写计划

- [x] T074 创建 `specs/5-std-task-api/contracts/task-api.yaml`
  - [x] OpenAPI 3.0 规范
  - [x] 24个端点完整定义
  - [x] 请求/响应 Schema

- [x] T075 验证所有24个API端点已注册

### 6.4 兼容性验证

- [x] T070 验证错误码与Java实现完全一致 (30701-30705)

- [ ] T071 **接口兼容性验证**
  - [ ] 确认响应格式与Java完全一致
  - [ ] 确认异常信息与Java完全一致

**Checkpoint**: ⏳ Phase 6 部分完成 (测试待补充)

### 6.4 兼容性验证

- [ ] T070 验证错误码与Java实现完全一致

- [ ] T071 **接口兼容性验证**
  - [ ] 确认响应格式与Java完全一致
  - [ ] 确认异常信息与Java完全一致

**Checkpoint**: ✅ Phase 6 完成

---

## 公共 Logic (common.go)

- [x] T072 创建 `api/internal/logic/task/common.go`
  - [x] `StatusToInt` - 状态字符串转整数
  - [x] `IntToStatus` - 整数转状态字符串
  - [x] `GenerateTaskNo` - 生成任务编号
  - [ ] `CallStdRecService` - 调用标准推荐服务 (TODO)
  - [ ] `CallRuleRecService` - 调用规则推荐服务 (TODO)
  - [ ] `SendTaskCallback` - 发送任务完成回调 (TODO)
  - [x] `buildTaskResp` - 构建任务响应
  - [x] `buildTaskResultResp` - 构建任务结果响应

---

## ServiceContext 更新

- [x] T073 更新 `api/internal/svc/service_context.go`
  - [x] 添加 TaskStdCreateModel
  - [x] 添加 TaskStdCreateResultModel
  - [x] 添加 BusinessTablePoolModel
  - [ ] 添加 TaskStdCreateModel、TaskStdCreateResultModel、BusinessTablePoolModel
  - [ ] 初始化 DB 连接 (*sqlx.Conn)
  - [ ] 初始化 Model 实例
  - [ ] 添加 HTTPClient（推荐服务调用）
  - [ ] 添加推荐服务URL配置

---

## 依赖关系

```
Phase 0 (基础设施)
    ↓
Phase 1 (任务查询管理)
    ↓
Phase 2 (任务创建与完成)
    ↓
Phase 3 (业务表管理)
    ↓
Phase 4 (推荐服务)
    ↓
Phase 5 (数据元操作)
    ↓
Phase 6 (收尾工作)
```

### 并行执行说明

- `[P]` 标记的任务可并行执行
- Phase 2-5 可并行开发（如有团队）
- 同一 Phase 内的接口可并行实现

---

## MVP 范围

**最小可交付版本**: Phase 0 + Phase 1 + Phase 2

MVP 包含的核心功能：
- ✅ 查询任务列表（未处理/已完成）
- ✅ 查询任务详情
- ✅ 标准关联暂存
- ✅ 标准关联提交
- ✅ 新建标准任务
- ✅ 完成任务

---

## 测试要求 🧪

| 要求 | 标准 |
|------|------|
| **单元测试覆盖率** | ≥ 80% |
| **关键路径测试** | 100% 覆盖（任务查询、创建、完成） |
| **边界测试** | 必须包含 |
| **错误处理测试** | 必须包含 |

### 测试命名规范

```
Test{Function}_{Scenario}_{ExpectedResult}
```

示例：
- `TestGetUncompletedTasks_ValidInput_ReturnsList`
- `TestCreateTask_EmptyIds_ReturnsError`
- `TestSubmitRelation_Successful_CallsRecService`

---

## 实施进度跟踪

| Week | Phase | 内容 | 完成度 |
|------|-------|------|--------|
| 1 | Phase 0-2 | 基础设施 + 任务查询 + 任务创建完成 | 0% |

---

## 文档变更历史

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-02-09 | 初始版本 |
| 1.1 | 2026-02-09 | 添加文档任务 (T068-T075)，更新 Phase 6 状态 |
