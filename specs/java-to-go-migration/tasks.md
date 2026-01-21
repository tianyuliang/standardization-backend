# Java to Go Migration Tasks

> **Branch**: `feature/java-to-go-migration`
> **Spec Path**: `specs/java-to-go-migration/`
> **Created**: 2026-01-19
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
| `[US1]` | 关联 User Story 1（数据元素管理） |
| `[US2]` | 关联 User Story 2（目录管理） |
| `[US3]` | 关联 User Story 3（码表字典） |
| `[US4]` | 关联 User Story 4（编码规则） |
| `[US5]` | 关联 User Story 5（标准化任务） |
| `[US6]` | 关联 User Story 6（标准文件） |
| `[US7]` | 关联 User Story 7（消息队列） |
| `[US8]` | 关联 User Story 8（Redis） |
| `[TEST]` | 测试任务（必须完成） |

---

## Task Overview

| ID | Task | Story | Est. Lines |
|----|------|-------|------------|
| **Phase 0: 项目初始化** |
| T001 | 创建 Go 项目结构 | Setup | - |
| T002 | 初始化 go.mod 和依赖 | Setup | - |
| T003 | [P] 确认工具链安装 | Setup | - |
| **Phase 1: 基础设施** |
| T004 | 创建 base.api 基础类型 | Foundation | 30 |
| T005 | [P] 创建 api.doc/api.api 入口 | Foundation | 20 |
| T006 | 配置数据库连接 | Foundation | 40 |
| T007 | [P] 配置 Redis 连接 | Foundation | 30 |
| T008 | [P] 配置 idrm-go-base | Foundation | 20 |
| **Phase 2: 公共组件** |
| T009 | 实现 UUID v7 生成器 | Common | 40 |
| T010 | [P] 实现 Kafka 生产者 | US7 | 60 |
| T011 | [P] 实现 NSQ 生产者 | US7 | 50 |
| T012 | [P] 实现分布式锁 | US8 | 50 |
| **Phase 3: 数据元素模块 (US1)** |
| T013 | 创建 dataelement.api | US1 | 80 |
| T014 | 生成 Handler/Types | US1 | - |
| T015 | 创建 DDL | US1 | 40 |
| T016 | 实现 Model 层 | US1 | 80 |
| T017 | [TEST] Model 测试 | US1 | 100 |
| T018 | 实现 Logic 层 | US1 | 120 |
| T019 | [TEST] Logic 测试 | US1 | 150 |
| **Phase 4: 目录模块 (US2)** |
| T020 | 创建 catalog.api | US2 | 60 |
| T021 | 生成 Handler/Types | US2 | - |
| T022 | 创建 DDL | US2 | 30 |
| T023 | 实现 Model 层 | US2 | 70 |
| T024 | [TEST] Model 测试 | US2 | 80 |
| T025 | 实现 Logic 层 | US2 | 100 |
| T026 | [TEST] Logic 测试 | US2 | 120 |
| **Phase 5: 码表模块 (US3)** |
| T027 | 创建 dict.api | US3 | 50 |
| T028 | 生成 Handler/Types | US3 | - |
| T029 | 创建 DDL | US3 | 25 |
| T030 | 实现 Model 层 | US3 | 60 |
| T031 | [TEST] Model 测试 | US3 | 70 |
| T032 | 实现 Logic 层 | US3 | 90 |
| T033 | [P] [TEST] Excel 导入测试 | US3 | 60 |
| T034 | [P] [TEST] Excel 导出测试 | US3 | 50 |
| **Phase 6: 规则模块 (US4)** |
| T035 | 创建 rule.api | US4 | 50 |
| T036 | 生成 Handler/Types | US4 | - |
| T037 | 创建 DDL | US4 | 25 |
| T038 | 实现 Model 层 | US4 | 60 |
| T039 | [TEST] Model 测试 | US4 | 70 |
| T040 | 实现 Logic 层 | US4 | 90 |
| T041 | [TEST] Logic 测试 | US4 | 100 |
| **Phase 7: 任务模块 (US5)** |
| T042 | 创建 task.api | US5 | 40 |
| T043 | 生成 Handler/Types | US5 | - |
| T044 | 创建 DDL | US5 | 20 |
| T045 | 实现 Model 层 | US5 | 50 |
| T046 | [TEST] Model 测试 | US5 | 60 |
| T047 | 实现 Logic 层 | US5 | 80 |
| T048 | [TEST] Logic 测试 | US5 | 90 |
| **Phase 8: 文件模块 (US6)** |
| T049 | 创建 file.api | US6 | 50 |
| T050 | 生成 Handler/Types | US6 | - |
| T051 | 创建 DDL | US6 | 25 |
| T052 | 实现 Model 层 | US6 | 60 |
| T053 | [TEST] Model 测试 | US6 | 70 |
| T054 | 实现 Logic 层 | US6 | 90 |
| T055 | [P] [TEST] 文件上传测试 | US6 | 60 |
| T056 | [P] [TEST] 文件下载测试 | US6 | 50 |
| **Phase 9: 集成测试** |
| T057 | 集成测试框架搭建 | All | 40 |
| T058 | [TEST] API 端到端测试 | All | 200 |
| T059 | [P] [TEST] Kafka 消息测试 | US7 | 80 |
| T060 | [P] [TEST] NSQ 消息测试 | US7 | 80 |
| T061 | [P] [TEST] Redis 缓存测试 | US8 | 80 |
| T062 | [P] [TEST] 分布式锁测试 | US8 | 60 |
| **Phase 10: 对比验证** |
| T063 | 对比测试框架搭建 | All | 60 |
| T064 | [TEST] API 响应对比 | All | 200 |
| **Phase 11: 收尾** |
| T065 | 代码格式化 (gofmt) | All | - |
| T066 | 运行 golangci-lint | All | - |
| T067 | 生成 Swagger 文档 | All | - |
| T068 | 确认测试覆盖率 > 80% | All | - |

---

## Phase 0: 项目初始化

**目的**: 创建 Go 项目基础结构

### Task T001: 创建 Go 项目结构 ✅ COMPLETED

**描述**: 按 plan.md 中定义的项目结构创建目录

**步骤**:
1. 创建 `api/` 目录结构
2. 创建 `model/` 目录结构
3. 创建 `migrations/` 目录结构
4. 创建 `pkg/` 目录结构
5. 创建 `api/internal/` 子目录

**验收**: 目录结构与 plan.md 一致

---

### Task T002: 初始化 go.mod 和依赖 ✅ COMPLETED

**描述**: 初始化 Go 模块并安装核心依赖

**步骤**:
```bash
go mod init github.com/dsg/standardization-backend
go get github.com/zeromicro/go-zero@latest
go get github.com/jinguoxing/idrm-go-base@latest
go get gorm.io/gorm@latest
go get gorm.io/driver/mysql@latest
go get github.com/go-redsync/redsync/v4@latest
go get github.com/IBM/sarama@latest
go get github.com/nsqio/go-nsq@latest
go get github.com/redis/go-redis/v9@latest
go get github.com/xuri/excelize/v2@latest
```

**验收**: go.mod 文件存在，依赖可正常编译

---

### Task T003: [P] 确认工具链安装 ✅ COMPLETED

**描述**: 确认开发工具已安装

**步骤**:
```bash
# 确认 Go 版本
go version  # 需要 1.24+

# 确认 goctl
goctl version

# 确认其他工具
golangci-lint version
```

**验收**: 所有工具命令可正常执行

---

## Phase 1: 基础设施

**目的**: 配置项目基础组件

---

### Task T004: 创建 base.api 基础类型 ✅ COMPLETED

**引用**: plan.md#File Structure

**文件**: `api/doc/base.api`

**步骤**:
1. 定义通用 Request 类型（BaseIdReq, BasePageReq）
2. 定义通用 Response 类型（BaseResp, BasePageResp）
3. 定义错误响应格式

**预计代码**: ~30 行

**验收**: base.api 可被其他 .api 文件 import

---

### Task T005: [P] 创建 api.doc/api.api 入口 ✅ COMPLETED

**引用**: plan.md#Go-Zero 开发流程

**文件**: `api/doc/api.api`

**步骤**:
1. 定义入口服务
2. import base.api
3. 预留各模块 import 位置

**预计代码**: ~20 行

**验收**: api.api 文件格式正确

---

### Task T006: 配置数据库连接 ✅ COMPLETED

**引用**: plan.md#Technical Context

**文件**: `api/internal/svc/servicecontext.go`

**步骤**:
1. 创建数据库配置结构体
2. 使用 GORM 连接 MariaDB
3. 验证连接成功

**预计代码**: ~40 行

**验收**: 可成功连接数据库

---

### Task T007: [P] 配置 Redis 连接 ✅ COMPLETED

**引用**: plan.md#Redis 配置

**文件**: `pkg/cache/redis/redis.go`

**步骤**:
1. 创建 Redis 配置结构体
2. 实现连接逻辑
3. 实现 Get/Set/Del 方法

**预计代码**: ~30 行

**验收**: 可成功连接 Redis 并读写数据

---

### Task T008: [P] 配置 idrm-go-base ✅ COMPLETED

**引用**: plan.md#通用库

**文件**: `api/internal/svc/servicecontext.go`

**步骤**:
1. 初始化 validator
2. 初始化 telemetry
3. 设置 response 错误处理器

**预计代码**: ~20 行

**验收**: idrm-go-base 各模块正常工作

---

## Phase 2: 公共组件

**目的**: 实现跨模块共享的公共组件

---

### Task T009: 实现 UUID v7 生成器 ✅ COMPLETED

**引用**: spec.md#Business Rules > BR-01

**文件**: `pkg/idgen/uuidv7.go`

**步骤**:
1. 实现 UUID v7 生成算法
2. 添加单元测试

**预计代码**: ~40 行

**验收**: 生成的 UUID 符合 v7 规范且唯一

---

### Task T010: [P] 实现 Kafka 生产者 ✅ COMPLETED

**引用**: plan.md#Message Queue Integration

**文件**: `pkg/mq/kafka/producer.go`

**步骤**:
1. 创建 KafkaProducer 结构体
2. 实现 SendMessage 方法
3. 添加连接池配置
4. 定义 Topic 常量

**预计代码**: ~60 行

**验收**: 可成功发送 Kafka 消息

---

### Task T011: [P] 实现 NSQ 生产者 ✅ COMPLETED

**引用**: plan.md#Message Queue Integration

**文件**: `pkg/mq/nsq/producer.go`

**步骤**:
1. 创建 NSQProducer 结构体
2. 实现 Publish 方法
3. 添加连接池配置

**预计代码**: ~50 行

**验收**: 可成功发送 NSQ 消息

---

### Task T012: [P] 实现分布式锁 ✅ COMPLETED

**引用**: plan.md#Distributed Lock

**文件**: `pkg/cache/lock/distributed_lock.go`

**步骤**:
1. 定义 DistributedLock 接口
2. 实现 RedsyncLock
3. 实现 Lock/Unlock 方法
4. 添加单元测试

**预计代码**: ~50 行

**验收**: 分布式锁可正常工作

---

## Phase 3: 数据元素模块 (US1)

**目标**: 实现数据元管理 CRUD、历史记录、业务表字段管理

**独立测试**: 对比 Java 版本 API 响应

---

### Task T013: 创建 dataelement.api

**引用**: plan.md#Module 1

**文件**: `api/doc/dataelement/dataelement.api`

**步骤**:
1. 定义 CreateDataElementReq/Resp
2. 定义 GetDataElementReq/Resp
3. 定义 UpdateDataElementReq/Resp
4. 定义 DeleteDataElementReq/Resp
5. 定义 ListDataElementReq/Resp
6. 在 api.api 中 import 此模块

**预计代码**: ~80 行

**验收**: .api 文件语法正确

---

### Task T014: 生成 Handler/Types

**引用**: plan.md#Go-Zero 开发流程

**命令**:
```bash
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

**验收**: Handler 和 Types 文件已生成

---

### Task T015: 创建 DDL

**引用**: plan.md#Module 1 > DDL 设计

**文件**: `migrations/dataelement/data_element_info.sql`

**步骤**:
1. 按计划创建 DDL
2. 添加索引
3. 添加注释

**预计代码**: ~40 行

**验收**: DDL 可在数据库中执行

---

### Task T016: 实现 Model 层

**引用**: plan.md#Module 1 > Model 接口

**文件**: `model/dataelement/dataelement/`

**步骤**:
1. 创建 interface.go
2. 创建 types.go (DataElement 结构体)
3. 创建 vars.go (常量和错误)
4. 实现 gorm_dao.go

**预计代码**: ~80 行

**验收**: Model 接口完整实现

---

### Task T017: [TEST] Model 测试

**引用**: tasks-template.md#测试要求

**文件**: `model/dataelement/dataelement/gorm_dao_test.go`

**步骤**:
1. 测试 Insert 方法
2. 测试 FindOne 方法
3. 测试 FindByCode 方法
4. 测试 Update 方法
5. 测试 Delete 方法
6. 测试 List 方法
7. 测试边界情况和错误处理

**预计代码**: ~100 行

**验收**: 测试覆盖率 > 80%

---

### Task T018: 实现 Logic 层

**引用**: spec.md#AC-01 至 AC-06

**文件**: `api/internal/logic/dataelement/`

**步骤**:
1. 实现 CreateDataElementLogic
2. 实现 GetDataElementLogic
3. 实现 UpdateDataElementLogic
4. 实现 DeleteDataElementLogic
5. 实现 ListDataElementLogic
6. 集成 Kafka 消息发送
7. 集成 Redis 缓存

**预计代码**: ~120 行

**验收**: Logic 层实现完整业务逻辑

---

### Task T019: [TEST] Logic 测试

**文件**: `api/internal/logic/dataelement/*_test.go`

**步骤**:
1. Mock Model 层
2. 测试正常流程
3. 测试异常处理
4. 测试边界情况

**预计代码**: ~150 行

**验收**: 测试覆盖率 > 80%

---

## Phase 4: 目录模块 (US2)

**目标**: 实现目录树管理、移动、统计

**独立测试**: 对比 Java 版本 API 响应

---

### Task T020-T026: 目录模块实现

**说明**: 按照 Phase 3 相同步骤实现目录模块

**核心差异**:
- 树形结构查询（递归）
- 目录移动逻辑（tree_path 更新）
- 统计查询

---

## Phase 5: 码表模块 (US3)

**目标**: 实现码表 CRUD、Excel 导入导出

**独立测试**: 对比 Java 版本 API 响应

---

### Task T027-T034: 码表模块实现

**说明**: 按照相同模式实现

**核心差异**:
- Excel 导入使用 excelize 库
- Excel 导出生成 .xlsx 文件

---

## Phase 6: 规则模块 (US4)

**目标**: 实现规则管理、自定义规则、规则推荐

**独立测试**: 对比 Java 版本 API 响应

---

### Task T035-T041: 规则模块实现

**说明**: 按照相同模式实现

---

## Phase 7: 任务模块 (US5)

**目标**: 实现任务管理、状态流转

**独立测试**: 对比 Java 版本 API 响应

---

### Task T042-T048: 任务模块实现

**说明**: 按照相同模式实现

**核心差异**:
- 状态机验证
- 任务状态流转规则

---

## Phase 8: 文件模块 (US6)

**目标**: 实现文件管理、附件管理

**独立测试**: 对比 Java 版本 API 响应

---

### Task T049-T056: 文件模块实现

**说明**: 按照相同模式实现

**核心差异**:
- 文件上传处理
- 文件下载响应

---

## Phase 9: 集成测试

**目的**: 验证模块间协作

---

### Task T057: 集成测试框架搭建

**文件**: `api/test/integration/setup.go`

**步骤**:
1. 使用 testcontainers 启动测试数据库
2. 初始化测试数据
3. 创建测试 HTTP 服务器

**预计代码**: ~40 行

---

### Task T058: [TEST] API 端到端测试

**文件**: `api/test/integration/api_test.go`

**步骤**:
1. 测试完整业务流程
2. 验证事务回滚
3. 验证并发处理

**预计代码**: ~200 行

---

### Task T059-T062: 中间件测试

**步骤**:
1. 测试 Kafka 消息发送和消费
2. 测试 NSQ 消息发送和消费
3. 测试 Redis 缓存读写
4. 测试分布式锁获取和释放

---

## Phase 10: 对比验证

**目的**: 确认 Go 版本与 Java 版本功能一致

---

### Task T063: 对比测试框架搭建

**文件**: `api/test/comparison/framework.go`

**步骤**:
1. 创建并行调用 Java/Go API 的框架
2. 实现响应对比逻辑
3. 生成对比报告

**预计代码**: ~60 行

---

### Task T064: [TEST] API 响应对比

**文件**: `api/test/comparison/compare_test.go`

**步骤**:
1. 对比所有 API 端点
2. 验证响应结构一致
3. 验证业务逻辑一致
4. 生成差异报告

**预计代码**: ~200 行

---

## Phase 11: 收尾

**目的**: 代码质量和文档

---

### Task T065: 代码格式化

```bash
gofmt -w .
```

---

### Task T066: 运行 golangci-lint

```bash
golangci-lint run
```

---

### Task T067: 生成 Swagger 文档

```bash
goctl api plugin -plugin goctl-swagger="swagger -filename api.json" -api api/doc/api.api -dir .
```

---

### Task T068: 确认测试覆盖率 > 80%

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

---

## Dependencies

```
Phase 0 (项目初始化)
    ↓
Phase 1 (基础设施)
    ↓
Phase 2 (公共组件)
    ↓
Phase 3 (US1: 数据元素) ←─┐
    ↓                      │
Phase 4 (US2: 目录)       │
    ↓                      │
Phase 5 (US3: 码表)       │ 并行执行
Phase 6 (US4: 规则)       │ (如有团队)
Phase 7 (US5: 任务)       │
Phase 8 (US6: 文件) ──────┘
    ↓
Phase 9 (集成测试)
    ↓
Phase 10 (对比验证)
    ↓
Phase 11 (收尾)
```

---

## 测试要求 🧪

| 要求 | 标准 |
|------|------|
| **单元测试覆盖率** | > 80% |
| **关键路径测试** | 100% 覆盖 |
| **边界测试** | 必须包含 |
| **错误处理测试** | 必须包含 |
| **并发测试** | 分布式锁相关 |

---

## Notes

- 每个 Task 完成后提交代码
- **实现和测试必须同时提交**
- 每个 Checkpoint 运行 `go test ./...` 验证
- 遇到问题及时记录到 spec.md#Open Questions

---

## Java 到 Go 迁移经验总结 📚

> **基于目录模块 (Catalog) 迁移实践总结**
> **目的**: 避免后续迁移重复相同错误

### 常见错误模式

| 错误 | 说明 | 解决方案 |
|------|------|----------|
| **Base Path 不一致** | Go 用 `/c`，Java 用 `/v1/catalog` | `.api` 文件中 `prefix` 必须与 Java `@RequestMapping` 完全一致 |
| **路径参数缺失** | `put /:id` 但 Request 无 `Id` 字段 | 路径参数必须在 Request 类型中添加对应字段 |
| **参数名不匹配** | Java 用 `catalog_name`，Go 用 `name` | 完全匹配 Java 参数名（包括大小写） |
| **条件应用位置** | `LIKE "%%"` 在 name 为空时仍执行 | 只在参数非空时才应用条件 |
| **逻辑反转** | 拒绝在根目录下创建子目录 | 创建允许在 level=1 下创建，更新不允许修改 level<=1 |
| **删除逻辑错误** | 有子目录就拒绝删除 | Java 使用级联删除 `removeBatchByIds(getIDList())` |
| **循环引用缺失** | 更新时未检查循环父目录 | 必须检查新父目录不是自身或子孙节点 |
| **类型不匹配** | 返回自定义类型而非 Java VO | 使用与 Java VO 完全一致的响应类型 |
| **语言不一致** | Java 中文错误消息，Go 变成英文 | 所有错误消息必须使用中文，与 Java 版本保持一致 |
| **goctl 文件覆盖** | 修改 .api 后旧文件仍保留 | 修改 .api 后必须重新运行 goctl，让其覆盖生成新文件 |
| **Handler 类型错误** | Handler 请求类型与 Logic 不一致 | 修改 .api 文件后需手动检查并更新 Handler |
| **子包导出缺失** | model 子包类型无法导入 | 创建父包文件（如 catalog.go）导出子包类型 |
| **ServiceContext 错误** | 忽略 NewServiceContext 的 error 返回值 | 正确处理: `ctx, err := svc.NewServiceContext(c)` |
| **路径参数重复** | `/:id` 路径参数且请求体含 `Id` 字段 | 路径参数用 `r.PathValue()` 提取，请求体不含 id |
| **Type 字段冗余** | Create 请求包含 `Type` 字段 | 创建时从父目录继承类型，不需要客户端传递 |
| **验证长度不匹配** | API 验证 `max=64`，数据库 `varchar(20)` | 验证长度必须与数据库约束一致 |
| **重复导入 base.api** | 模块 .api 导入 base.api 导致重复定义 | 只在 api.api 中导入，模块 .api 不导入 |
| **错误处理不当** | 用 `fmt.Errorf(variable)` 处理预定义错误 | 使用 `errors.New(variable)` 处理 string 类型错误常量 |

### Model 层最佳实践

1. **辅助方法**:
   - `GetIDList(catalogId)` - 获取目录及所有子孙的 ID（用于级联操作和循环校验）
   - `DeleteBatch(ids)` - 批量删除（用于级联删除）
   - `CheckNameExists(name, parentId, excludeId)` - 重复校验（支持排除自身）

2. **软删除**:
   - 始终使用 `deleted_at` 字段，而非硬删除
   - 查询条件必须包含 `deleted_at IS NULL`

3. **ID 类型**:
   - 使用 UUID v7 替代 Java 的 Snowflake ID
   - 字段类型: `string` (非 `int64`)

4. **子包结构**:
   - Model 实际位于 `model/{module}/{module}/` 子目录
   - 需要在 `model/{module}/` 创建父包文件导出类型
   ```go
   // model/catalog/catalog.go
   package catalog
   import "github.com/dsg/standardization-backend/model/catalog/catalog"
   type Model = catalog.Model
   type Catalog = catalog.Catalog
   // ... 导出常量和错误
   ```

### Logic 层最佳实践

1. **创建操作**:
   ```go
   // ✅ 正确: 允许在根目录下创建
   parent.Level + 1  // 即使 parent.Level = 1 也允许
   ```

2. **更新操作**:
   ```go
   // ✅ 正确: 多重校验
   if c.Level <= 1 { return error }  // 根目录不能修改
   if parent.Type != c.Type { return error }  // 类型必须一致
   if contains(newParentId, getIDList(c.Id)) { return error }  // 防止循环
   ```

3. **删除操作**:
   ```go
   // ✅ 正确: 级联删除
   idList := GetIDList(c.Id)  // 获取自身及所有子孙
   DeleteBatch(idList)         // 批量删除
   ```

4. **查询操作**:
   ```go
   // ✅ 正确: 条件只在非空时应用
   query := Where("type = ? AND level > 1", type)
   if name != "" {  // 只在参数非空时应用 LIKE
       query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+name+"%")
   }
   ```

### API 定义最佳实践

1. **路径参数** (正确做法):
   ```api
   put /:id (UpdateCatalogReq) returns (BaseResp)
   // UpdateCatalogReq 不包含 Id 字段（从路径参数获取）
   type UpdateCatalogReq {
       CatalogName string `json:"catalogName" validate:"required,max=20"`
       ParentId    string `json:"parentId" validate:"required"`
       Description string `json:"description"`
   }

   // Handler 中提取路径参数
   id := r.PathValue("id")
   l.UpdateCatalog(&req, id)

   // Logic 中接收独立参数
   func (l *UpdateCatalogLogic) UpdateCatalog(req *types.UpdateCatalogReq, id string)
   ```
   - Java 使用 `@PathVariable` 和 `@RequestBody` 分开传递
   - Go 使用 `r.PathValue()` 从路径提取，不放入请求体

2. **查询参数**:
   ```api
   get /query (QueryCatalogReq) returns (QueryCatalogResp)
   // 使用 form 标签，非 json
   type QueryCatalogReq {
       Type        int32  `form:"type" validate:"required"`
       CatalogName string `form:"catalog_name" optional:"`
   }
   ```

3. **响应类型**:
   - 使用与 Java VO 完全一致的字段名
   - 包含所有 Java VO 中定义的字段
   - 额外字段使用 `optional` 或标记为 TODO

### 项目结构最佳实践

1. **API Base Path 定义**:
   ```api
   @server(
       prefix: /v1/catalog  # 必须与 Java @RequestMapping 完全一致
       group: catalog
   )
   ```
   - Java: `@RequestMapping("/v1/catalog")`
   - Go: `prefix: /v1/catalog`（不能简写为 `/c`）

2. **goctl 文件管理**:
   - goctl 生成下划线命名文件（如 `service_context.go`、`create_catalog_logic.go`）
   - 修改 `.api` 文件后，需要重新运行 `goctl api go`，goctl 会覆盖已存在的文件
   - **正确做法**：
     1. 修改 `.api` 文件
     2. 运行 `goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group`
     3. goctl 会覆盖生成所有文件（handler、logic、types、routes、service_context）
     4. 手动编辑生成的文件添加业务逻辑（保留 goctl 生成的结构）
   - **错误做法**：手动保留旧文件，会导致重复声明错误

3. **Handler 与 Logic 同步**:
   - 修改 `.api` 文件后，goctl 不会自动更新 Handler
   - 需手动检查 Handler 中的请求类型是否与 Logic 一致
   ```go
   // Handler 中必须使用正确的请求类型
   var req types.QueryWithFileReq  // 非 QueryTreeReq
   ```

4. **ServiceContext 初始化**:
   ```go
   // main.go 中正确处理错误
   ctx, err := svc.NewServiceContext(c)
   if err != nil {
       panic(fmt.Sprintf("Failed to create service context: %v", err))
   }
   ```

4. **go.mod 依赖管理**:
   - 初始 go.mod 可能只有 indirect 依赖
   - 运行 `go mod tidy` 整理直接依赖
   ```bash
   go mod tidy
   ```

5. **避免导入不存在类型**:
   - Model 中导出的类型必须实际存在
   - VO 类型（如 `CatalogTreeNodeVo`）在 `api/internal/types` 中
   - 不要在 Model 包中导出 VO 类型

6. **API 模块导入规范**:
   - 只在 `api.doc.api.api` 中导入 `base.api`
   - 各模块 .api 文件（如 `catalog/catalog.api`）不要重复导入 `base.api`
   - 避免类型重复定义错误

7. **验证长度与数据库一致**:
   - API 请求验证必须与数据库字段约束一致
   - 示例：数据库 `varchar(20)` → API `validate:"max=20"`
   - 避免验证过长导致数据库插入失败

8. **错误处理规范**:
   - 预定义错误常量是 `string` 类型，使用 `errors.New()` 而非 `fmt.Errorf()`
   - 只有需要格式化的错误才使用 `fmt.Errorf()`
   ```go
   // ✅ 正确
   return nil, errors.New(catalog.ErrInvalidParent)

   // ❌ 错误 (non-constant format string)
   return nil, fmt.Errorf(catalog.ErrInvalidParent)

   // ✅ 需要格式化时
   return nil, fmt.Errorf("无效的目录ID: %w", err)
   ```

### 测试检查要点

- [ ] 测试在根目录下创建子目录（level=1 -> level=2）
- [ ] 测试修改根目录应被拒绝
- [ ] 测试循环父目录校验（A->B，尝试 B->A）
- [ ] 测试级联删除（删除父目录，子目录也被删除）
- [ ] 测试空参数搜索（返回所有，而非报错）
- [ ] 测试模糊搜索（大小写不敏感）

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-01-19 | - | 初始版本 - Java to Go 迁移任务拆分 |
| 1.1 | 2026-01-20 | Claude | 添加 Java 到 Go 迁移经验总结，基于目录模块实践 |
| 1.2 | 2026-01-20 | Claude | 添加项目结构最佳实践和 goctl 相关问题解决 |
| 1.3 | 2026-01-20 | Claude | 添加 API Base Path 一致性检查和修复说明 |
| 1.4 | 2026-01-20 | Claude | 添加路径参数处理、Type 字段冗余、验证长度等新问题和解决方案 |
