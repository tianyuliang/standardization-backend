# 标准任务管理 (std-task-api) Requirements Checklist

> **Branch**: `5-std-task-api`
> **Spec Path**: `specs/5-std-task-api/`
> **Created**: 2026-02-09
> **Status**: Draft

---

## 检查清单说明

本文档用于跟踪标准任务管理模块的所有需求，确保实现完整且符合规范。

---

## Phase 0: 基础设施

### 环境检查
- [ ] Go-Zero 项目结构已就绪
- [ ] 数据库连接已配置

### 错误码定义
- [ ] 错误码文件 `api/internal/errorx/task.go` 已创建
- [ ] 错误码范围 30701-30799 已定义
- [ ] 错误辅助函数已实现

### API 定义
- [ ] `api/doc/task/task.api` 已创建
- [ ] 24个 API 端点已定义
- [ ] 路由配置正确：`/api/standardization/v1/dataelement/task`
- [ ] `api.doc/api.api` 已 import task 模块
- [ ] goctl 命令已执行：`goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group`

---

## Phase 1: 任务查询管理 (4接口)

### DDL 定义
- [ ] `t_task_std_create.sql` 已创建
- [ ] `t_task_std_create_result.sql` 已创建

### Model 层
- [ ] `model/task/task/interface.go` 已创建
- [ ] `model/task/task/types.go` 已创建
- [ ] `model/task/task/vars.go` 已创建
- [ ] `model/task/task/factory.go` 已创建
- [ ] `model/task/task/sql_model.go` 已实现
- [ ] 单元测试已编写

### 接口实现
- [ ] GET `/std-create/uncompleted` - 未处理任务列表
- [ ] GET `/std-create/completed` - 已完成任务列表
- [ ] GET `/std-create/completed/{id}` - 任务详情
- [ ] POST `/queryTaskState` - 任务状态查询

---

## Phase 2: 任务创建与完成 (5接口)

### 接口实现
- [ ] POST `/std-create/relation/staging` - 标准关联暂存
- [ ] POST `/std-create/publish/submit` - 标准关联提交
- [ ] POST `/createTask` - 新建标准任务
- [ ] POST `/finishTask/{task_id}` - 完成任务
- [ ] POST `/queryTaskProcess` - 进度查询

---

## Phase 3: 业务表管理 (7接口)

### DDL 定义
- [ ] `t_business_table_std_create_pool.sql` 已创建

### Model 层
- [ ] `model/task/pool/` 目录结构已创建
- [ ] BusinessTablePoolModel 已实现
- [ ] 单元测试已编写

### 接口实现
- [ ] POST `/addToPending` - 添加至待新建
- [ ] GET `/getBusinessTable` - 业务表列表
- [ ] GET `/getBusinessTableField` - 业务表字段列表
- [ ] DELETE `/deleteBusinessTableField/{id}` - 移除字段
- [ ] PUT `/cancelBusinessTableField` - 撤销
- [ ] PUT `/updateTableName` - 修改表名称
- [ ] GET `/getBusinessTableFromTask` - 任务关联业务表

---

## Phase 4: 推荐服务 (4接口)

### 接口实现
- [ ] POST `/std-rec/rec` - 标准推荐（内部）
- [ ] POST `/std-create` - 标准创建（内部）
- [ ] POST `/stand-rec/rec` - 标准推荐（弹框）
- [ ] POST `/rule-rec/rec` - 编码规则推荐

---

## Phase 5: 数据元操作 (4接口)

### 接口实现
- [ ] GET `/getBusinessTableFieldFromTask` - 任务关联字段
- [ ] POST `/submitDataElement` - 提交选定数据元
- [ ] PUT `/updateDescription` - 修改字段说明
- [ ] PUT `/accept` - 采纳

---

## Phase 6: 收尾工作

### 代码质量
- [ ] 代码格式化（gofmt）已完成
- [ ] golangci-lint 检查已通过

### 测试验证
- [ ] 测试覆盖率 ≥ 80%
- [ ] 所有测试已通过

### 文档更新
- [ ] Swagger 文档已生成
- [ ] 24个 API 端点已验证

### 兼容性验证
- [ ] 错误码与Java实现一致
- [ ] 响应格式与Java实现一致
- [ ] 异常信息与Java实现一致

---

## 技术约束检查

### 主键类型
- [ ] ⚠️ 使用 BIGINT 主键（与Java兼容的例外情况）

### 表结构
- [ ] 完全复用 Java 表结构（f_ 前缀字段）

### 接口兼容性
- [ ] 24个API路径100%一致
- [ ] 请求参数100%一致
- [ ] 响应格式100%一致
- [ ] 错误信息100%一致

### 数据访问
- [ ] 使用纯 SQLx 实现
- [ ] 手工编写 SQL 查询

### Logic 层
- [ ] 每个函数 ≤ 50 行
- [ ] 处理步骤已标注

---

## 依赖服务集成

- [ ] 推荐算法服务（HTTP 客户端）
- [ ] 规则推荐服务（HTTP 客户端）
- [ ] AF系统回调（HTTP POST）

---

## 统计信息

| 项目 | 数量 |
|------|------|
| API 接口 | 24 |
| 数据表 | 3 |
| 错误码 | 30701-30799 |
| 任务总数 | 73 |

---

## 签署确认

| 角色 | 姓名 | 日期 | 签名 |
|------|------|------|------|
| 开发者 | | | |
| 审核者 | | | |
| 测试者 | | | |
