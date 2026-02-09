# 标准任务管理接口实现文档

## 一、实现概述

### 1.1 项目背景

**源系统**: Java Spring Boot
**目标系统**: Go-Zero 微服务框架
**兼容性要求**: 100% 保持接口兼容（路径、参数、响应、异常信息）

### 1.2 实施范围

- **接口数量**: 24个
- **实现时间**: 2026-02-09
- **代码路径**: `api/internal/logic/task/`
- **Mock路径**: `api/internal/logic/task/mock/`

---

## 二、处理流程

### 2.1 工作流程图

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ Java源码分析 │ -> │ 接口框架搭建 │ -> │ Mock服务拆分 │ -> │ 业务逻辑实现 │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
                                                               |
                                                               v
                                                         ┌─────────────┐
                                                         │ 构建验证   │
                                                         └─────────────┘
```

### 2.2 关键里程碑

| 阶段 | 说明 | 状态 |
|------|------|------|
| Phase 0 | API 定义与代码生成 | ✅ 完成 |
| Phase 1-8 | 24个接口框架搭建 | ✅ 完成 |
| Phase 9 | 接口细节补充（业务逻辑） | ✅ 完成 |
| Phase 10 | Mock代码拆分与优化 | ✅ 完成 |

---

## 三、接口清单与实现状态

### 3.1 任务管理接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 1 | GET | `/api/standardization/v1/dataelement/task/std-create/uncompleted` | [get_uncompleted_tasks_logic.go](../../api/internal/logic/task/get_uncompleted_tasks_logic.go) | ✅ |
| 2 | GET | `/api/standardization/v1/dataelement/task/std-create/completed` | [get_completed_tasks_logic.go](../../api/internal/logic/task/get_completed_tasks_logic.go) | ✅ |
| 3 | GET | `/api/standardization/v1/dataelement/task/std-create/completed/{id}` | [get_task_by_id_logic.go](../../api/internal/logic/task/get_task_by_id_logic.go) | ✅ |
| 4 | POST | `/api/standardization/v1/dataelement/task/createTask` | [create_task_logic.go](../../api/internal/logic/task/create_task_logic.go) | ✅ |
| 5 | POST | `/api/standardization/v1/dataelement/task/finishTask/{task_id}` | [finish_task_logic.go](../../api/internal/logic/task/finish_task_logic.go) | ✅ |
| 6 | POST | `/api/standardization/v1/dataelement/task/queryTaskProcess` | [query_task_process_logic.go](../../api/internal/logic/task/query_task_process_logic.go) | ✅ |
| 7 | POST | `/api/standardization/v1/dataelement/task/queryTaskState` | [query_task_state_logic.go](../../api/internal/logic/task/query_task_state_logic.go) | ✅ |

### 3.2 标准关联接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 8 | POST | `/api/standardization/v1/dataelement/task/std-create/relation/staging` | [std_create_logic.go](../../api/internal/logic/task/std_create_logic.go) | ✅ |
| 9 | POST | `/api/standardization/v1/dataelement/task/std-create/publish/submit` | [submit_relation_logic.go](../../api/internal/logic/task/submit_relation_logic.go) | ✅ |
| 10 | POST | `/api/standardization/v1/dataelement/task/submitDataElement` | [submit_data_element_logic.go](../../api/internal/logic/task/submit_data_element_logic.go) | ✅ |

### 3.3 业务表管理接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 11 | POST | `/api/standardization/v1/dataelement/task/addToPending` | [add_to_pending_logic.go](../../api/internal/logic/task/add_to_pending_logic.go) | ✅ |
| 12 | GET | `/api/standardization/v1/dataelement/task/getBusinessTable` | [get_business_table_logic.go](../../api/internal/logic/task/get_business_table_logic.go) | ✅ |
| 13 | GET | `/api/standardization/v1/dataelement/task/getBusinessTableField` | [get_business_table_field_logic.go](../../api/internal/logic/task/get_business_table_field_logic.go) | ✅ |
| 14 | GET | `/api/standardization/v1/dataelement/task/getBusinessTableFromTask` | [get_table_from_task_logic.go](../../api/internal/logic/task/get_table_from_task_logic.go) | ✅ |
| 15 | GET | `/api/standardization/v1/dataelement/task/getBusinessTableFieldFromTask` | [get_field_from_task_logic.go](../../api/internal/logic/task/get_field_from_task_logic.go) | ✅ |
| 16 | DELETE | `/api/standardization/v1/dataelement/task/deleteBusinessTableField/{id}` | [delete_field_logic.go](../../api/internal/logic/task/delete_field_logic.go) | ✅ |
| 17 | PUT | `/api/standardization/v1/dataelement/task/updateDescription` | [update_description_logic.go](../../api/internal/logic/task/update_description_logic.go) | ✅ |
| 18 | PUT | `/api/standardization/v1/dataelement/task/updateTableName` | [update_table_name_logic.go](../../api/internal/logic/task/update_table_name_logic.go) | ✅ |
| 19 | PUT | `/api/standardization/v1/dataelement/task/cancelBusinessTableField` | [cancel_field_logic.go](../../api/internal/logic/task/cancel_field_logic.go) | ✅ |
| 20 | PUT | `/api/standardization/v1/dataelement/task/accept` | [accept_logic.go](../../api/internal/logic/task/accept_logic.go) | ✅ |

### 3.4 推荐接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 21 | POST | `/api/standardization/v1/dataelement/task/std-rec/rec` | [std_rec_logic.go](../../api/internal/logic/task/std_rec_logic.go) | ✅ |
| 22 | POST | `/api/standardization/v1/dataelement/task/std-create` | [std_create_logic.go](../../api/internal/logic/task/std_create_logic.go) | ✅ |
| 23 | POST | `/api/standardization/v1/dataelement/task/stand-rec/rec` | [stand_rec_logic.go](../../api/internal/logic/task/stand_rec_logic.go) | ✅ |
| 24 | POST | `/api/standardization/v1/dataelement/task/rule-rec/rec` | [rule_rec_logic.go](../../api/internal/logic/task/rule_rec_logic.go) | ✅ |

---

## 四、文件结构

### 4.1 核心实现文件

```
api/internal/logic/task/
├── common.go                          # 共享辅助函数、业务校验
├── mock/                              # Mock 服务层
│   └── service.go                     # 统一 Mock 服务
├── add_to_pending_logic.go            # 添加至待新建
├── accept_logic.go                    # 采纳
├── cancel_field_logic.go              # 撤销
├── create_task_logic.go               # 创建任务
├── delete_field_logic.go              # 移除字段
├── finish_task_logic.go               # 完成任务
├── get_business_table_field_logic.go  # 业务表字段列表
├── get_business_table_logic.go        # 业务表列表
├── get_completed_tasks_logic.go       # 已完成任务列表
├── get_field_from_task_logic.go       # 任务关联字段
├── get_table_from_task_logic.go       # 任务关联业务表
├── get_task_by_id_logic.go            # 任务详情
├── get_uncompleted_tasks_logic.go     # 未处理任务列表
├── query_task_process_logic.go        # 进度查询
├── query_task_state_logic.go          # 任务状态查询
├── rule_rec_logic.go                  # 编码规则推荐
├── stand_rec_logic.go                 # 标准采纳推荐
├── std_create_logic.go                # 标准关联暂存
├── std_rec_logic.go                   # 标准推荐
├── submit_data_element_logic.go       # 提交选定数据元
├── submit_relation_logic.go           # 标准关联提交
├── update_description_logic.go        # 修改字段说明
└── update_table_name_logic.go         # 修改表名称
```

### 4.2 类型定义文件

- `api/internal/types/task.go` - 任务模块类型定义
- `api/internal/types/types.go` - 全局类型定义（包含 DataElementDetailVo、TaskDetailDto）

### 4.3 错误码文件

- `api/internal/errorx/task.go` - task-api 错误码定义（30700-30799）

---

## 五、Mock 服务管理

### 5.1 Mock 文件清单

| 文件 | 服务 | 说明 |
|------|------|------|
| [mock/service.go](../../api/internal/logic/task/mock/service.go) | 统一Mock | 所有外部依赖的Mock函数 |

### 5.2 Mock 函数清单

| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `GetTaskDetailDto()` | 获取任务详情 | `AfService.getTaskDetailDto()` |
| `GetDataElementInfo()` | 获取数据元信息 | `DataElementInfoService.getById()` |
| `GetDataElementDetailVo()` | 获取数据元详情 | `DataElementInfoService.getDetailVo()` |
| `CallStdRecService()` | 调用标准推荐服务 | HTTP POST to 推荐服务 |
| `SendTaskCallback()` | 发送任务回调 | Webhook HTTP POST |

### 5.3 Mock 替换搜索方法

```bash
# 搜索所有 Mock 标记
grep -r "MOCK:" api/internal/logic/task/

# 搜索需要替换的函数
grep -r "TODO:" api/internal/logic/task/mock/
```

---

## 六、错误码规范

### 6.1 错误码清单（30700-30799）

| 错误码 | 说明 | 错误函数 | Java 对应 |
|--------|------|----------|----------|
| 30701 | 数据不存在 | `TaskDataNotExist()` | DATA_NOT_EXIST |
| 30702 | 参数为空 | `TaskParamEmpty()` | PARAMETER_EMPTY |
| 30703 | 参数无效 | `TaskInvalidParam()` | InvalidParameter |
| 30704 | 部分失败 | `TaskPartFail()` | PARTAIL_FAILURE |

---

## 七、业务校验规则

### 7.1 TaskId 格式校验

```go
// TaskId 必须为36位UUID v7
if len(taskId) != 36 {
    return errorx.TaskInvalidParam("taskId必须为36位UUID")
}
```

### 7.2 ID 格式校验（BIGINT）

```go
// PoolId 格式校验 (BIGINT, 19位)
if len(idStr) != 19 {
    return errorx.TaskInvalidParam(fmt.Sprintf("ID格式错误: %s, 必须为19位", idStr))
}
id, err := strconv.ParseInt(idStr, 10, 64)
if err != nil || id <= 0 {
    return errorx.TaskInvalidParam(fmt.Sprintf("ID格式错误: %s", idStr))
}
```

### 7.3 批量操作校验

```go
// 批量操作ID列表校验
for _, idStr := range req.Ids {
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil || id <= 0 {
        return errorx.TaskInvalidParam(fmt.Sprintf("ID格式错误: %s", idStr))
    }
    ids = append(ids, id)
}
```

### 7.4 任务完成前置条件

```go
// 对应 Java: BusinessTableStdCreatePoolServiceImpl.finishTask()
// 1. 查询所有关联池记录
// 2. 验证每个记录都有关联的数据元
// 3. 验证数据元标准分类与任务组织类型一致
// 4. 批量更新状态为 CREATED

for _, pool := range poolRecords {
    if pool.DataElementId == 0 {
        return errorx.TaskPartFail("不能存在没有关联数据元的标准字段")
    }
}
```

---

## 八、未完成项 (TODO)

### 8.1 Mock 服务替换

| 服务 | 函数数 | 优先级 |
|------|--------|--------|
| AfService | 1 | 高 |
| DataElementInfoService | 2 | 高 |
| StdRecService | 1 | 中 |
| Webhook | 1 | 中 |

### 8.2 事务处理

当前代码中标记了需要开启事务的位置：
- [create_task_logic.go](../../api/internal/logic/task/create_task_logic.go) - 创建任务
- [finish_task_logic.go](../../api/internal/logic/task/finish_task_logic.go) - 完成任务

---

## 九、注意事项

### 9.1 命名规范

```go
// 文件命名: snake_case.go
add_to_pending_logic.go
query_task_process_logic.go

// 函数命名: PascalCase
func CreateTask() {}
func QueryTaskProcess() {}

// 变量命名: camelCase
businessTableId
dataElementId
```

### 9.2 导入顺序

```go
import (
    // 标准库
    "context"
    "strconv"
    "strings"

    // 项目内部
    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

    // 第三方库
    "github.com/zeromicro/go-zero/core/logx"
)
```

### 9.3 日志规范

```go
// 使用 logx，包含 traceId
logx.WithContext(ctx).Infof("创建任务成功: taskId=%s", taskId)
logx.Errorf("更新字段说明失败: fieldId=%d, error=%v", fieldId, err)
```

### 9.4 错误处理

```go
// 使用预定义错误码
import "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"

if pool == nil {
    return nil, errorx.TaskDataNotExist()
}
```

### 9.5 混合 ID 类型策略

```go
// TaskId: string (36位UUID v7)
type CreateTaskReq struct {
    TaskId string `json:"taskId,optional"`
}

// BusinessTableFieldId: string (36位UUID)
type SubmitDataElementReq struct {
    Id string `json:"id,optional"`
}

// PoolId: int64 (19位BIGINT)
type BusinessTablePool struct {
    Id int64 `db:"f_id"`
}

// DataElementId: int64 (BIGINT)
type BusinessTablePool struct {
    DataElementId int64 `db:"f_data_element_id"`
}
```

### 9.6 Java 对应代码标注

```go
// 对应 Java: BusinessTableStdCreatePoolServiceImpl.createTask()
// 业务流程:
//  1. TaskId格式校验
//  2. 查询池记录
//  3. 批量更新状态
```

---

## 十、参考文档

### 10.1 内部文档

| 文档 | 路径 |
|------|------|
| Java源码 | `specs/src/main/java/com/dsg/standardization/controller/StdCreateTaskController.java` |
| API定义 | `api/doc/task/task.api` |
| 数据模型 | `specs/5-std-task-api/data-model.md` |

### 10.2 外部依赖

| 依赖 | 版本 | 说明 |
|------|------|------|
| go-zero | v1.9+ | 微服务框架 |
| idrm-go-base | v0.1.0+ | 通用基础库 |
| SQLx | - | 高性能数据库访问 |

### 10.3 相关文档

- [项目宪法](../../.specify/memory/constitution.md)
- [SDD 模板](../../.specify/templates/)
- [通用库规范](../../CLAUDE.md)

---

## 十一、后续工作

### 11.1 Mock 服务替换

按优先级替换 Mock 函数：
1. **高优先级**: AfService、DataElementInfoService
2. **中优先级**: StdRecService、Webhook

### 11.2 事务处理

添加事务支持：
```go
// 示例
err := sqlx.Transaction(l.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
    // 数据库操作
    return nil
})
```

### 11.3 测试用例

为每个接口编写测试用例，覆盖：
- 正常流程
- 异常场景
- 边界条件

---

**文档版本**: v1.0
**更新时间**: 2026-02-09
**维护人**: AI Assistant
