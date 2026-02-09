# 修复建议方案: task-api

> **Feature**: 5-std-task-api (标准任务管理)
> **Date**: 2026-02-09
> **Purpose**: 修复 ANALYSIS.md 中识别的问题

---

## 修复清单

| ID | 问题 | 状态 | 预计工作量 |
|----|------|------|------------|
| A001 | Mock服务需要替换 | ⏳ 待修复 | 120分钟 |
| A002 | Webhook失败处理未定义 | ⏳ 待修复 | 30分钟 |
| A003 | API文档不完整 | ⏳ 待修复 | 20分钟 |
| A004 | 混合ID类型策略 | ✅ 已文档化 | - |
| A005 | 测试用例未实现 | ⏳ 待修复 | 180分钟 |

**总计**: 约 350 分钟

---

## A001: 替换 Mock 服务

### 问题描述

所有外部服务调用当前使用 Mock 实现，需要替换为实际的 RPC/HTTP 调用。

### 需要替换的 Mock 函数

| 函数 | 目标服务 | 优先级 |
|------|----------|--------|
| `GetTaskDetailDto()` | AfService RPC | HIGH |
| `GetDataElementInfo()` | DataElementInfoService RPC | HIGH |
| `GetDataElementDetailVo()` | DataElementInfoService RPC | HIGH |
| `CallStdRecService()` | Recommendation Service HTTP | MEDIUM |
| `SendTaskCallback()` | Webhook HTTP POST | MEDIUM |

### 修复方案

#### 1. 替换 AfService 调用

**新建文件**: `api/internal/logic/task/af/`
```go
// af/client.go
package af

import (
    "context"

    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
)

// AfClient Af服务客户端
type AfClient struct {
    // RPC client configuration
}

// GetTaskDetailDto 获取任务详情
// 对应 Java: AfService.getTaskDetailDto(String taskNo)
func (c *AfClient) GetTaskDetailDto(ctx context.Context, taskNo string) (*types.TaskDetailDto, error) {
    // TODO: 实现RPC调用
    return &types.TaskDetailDto{
        TaskNo:  taskNo,
        OrgType: "",
    }, nil
}
```

**更新 Logic**:
```go
// get_task_by_id_logic.go
import "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/task/af"

func (l *GetTaskByIdLogic) GetTaskById(id int64) (*types.TaskDetailResp, error) {
    // 替换 Mock 调用
    taskDto, err := l.svcCtx.AfClient.GetTaskDetailDto(l.ctx, task.TaskNo)
    if err != nil {
        return nil, err
    }
    // ...
}
```

#### 2. 替换 DataElementInfoService 调用

**新建文件**: `api/internal/logic/task/dataelement/`
```go
// dataelement/client.go
package dataelement

import (
    "context"

    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
)

// DataElementClient 数据元服务客户端
type DataElementClient struct {
    // RPC client configuration
}

// GetById 根据ID获取数据元信息
// 对应 Java: DataElementInfoService.getById(Long id)
func (c *DataElementClient) GetById(ctx context.Context, id int64) (*types.DataElementInfo, error) {
    // TODO: 实现RPC调用
    return &types.DataElementInfo{
        FieldId: id,
        StdCode: "",
    }, nil
}

// GetDetailVo 获取数据元详情
// 对应 Java: DataElementInfoService.getDetailVo(Long id)
func (c *DataElementClient) GetDetailVo(ctx context.Context, id int64) (*types.DataElementDetailVo, error) {
    // TODO: 实现RPC调用
    return &types.DataElementDetailVo{
        Id: id,
    }, nil
}
```

#### 3. 实现推荐服务 HTTP 调用

**新建文件**: `api/internal/logic/task/recommendation/`
```go
// recommendation/client.go
package recommendation

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
)

// RecommendationClient 推荐服务客户端
type RecommendationClient struct {
    httpClient *http.Client
    baseURL    string
}

// StdRecRequest 标准推荐请求
type StdRecRequest struct {
    BusinessTable        string `json:"businessTable"`
    TableField           string `json:"tableField"`
    TableFieldDescription string `json:"tableFieldDescription"`
}

// StdRecResponse 标准推荐响应
type StdRecResponse struct {
    Code        string          `json:"code"`
    Description string          `json:"description"`
    Data        []StdRecItem    `json:"data"`
}

// StdRecItem 标准推荐项
type StdRecItem struct {
    StdCode   string  `json:"stdCode"`
    StdName   string  `json:"stdName"`
    MatchRate float64 `json:"matchRate"`
}

// CallStdRecService 调用标准推荐服务
// 对应 Java: HTTP POST to recommendation service
func (c *RecommendationClient) CallStdRecService(ctx context.Context, req *StdRecRequest) (*StdRecResponse, error) {
    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }

    httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/std-rec/rec", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }

    httpReq.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(httpReq)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result StdRecResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &result, nil
}
```

#### 4. 实现 Webhook 回调

**新建文件**: `api/internal/logic/task/webhook/`
```go
// webhook/client.go
package webhook

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "time"
)

// WebhookClient Webhook客户端
type WebhookClient struct {
    httpClient *http.Client
    maxRetries int
}

// TaskCompletionPayload 任务完成回调负载
type TaskCompletionPayload struct {
    TaskId   string `json:"taskId"`
    Status   string `json:"status"`
    FinishTime int64 `json:"finishTime"`
}

// SendTaskCallback 发送任务完成回调
// 对应 Java: Webhook HTTP POST on task completion
func (c *WebhookClient) SendTaskCallback(ctx context.Context, webhookURL string, taskId string) error {
    payload := TaskCompletionPayload{
        TaskId:   taskId,
        Status:   "completed",
        FinishTime: time.Now().Unix(),
    }

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    // 实现重试逻辑
    var lastErr error
    for attempt := 0; attempt <= c.maxRetries; attempt++ {
        httpReq, err := http.NewRequestWithContext(ctx, "POST", webhookURL, bytes.NewBuffer(jsonData))
        if err != nil {
            return err
        }

        httpReq.Header.Set("Content-Type", "application/json")

        resp, err := c.httpClient.Do(httpReq)
        if err != nil {
            lastErr = err
            time.Sleep(time.Duration(attempt+1) * time.Second) // 指数退避
            continue
        }
        defer resp.Body.Close()

        if resp.StatusCode >= 200 && resp.StatusCode < 300 {
            return nil // 成功
        }

        lastErr = fmt.Errorf("webhook returned status %d", resp.StatusCode)
        time.Sleep(time.Duration(attempt+1) * time.Second)
    }

    // 记录失败但不影响任务完成
    logx.Errorf("Webhook callback failed after %d attempts: %v", c.maxRetries, lastErr)
    return nil // 返回nil，不阻塞任务完成
}
```

### 更新 ServiceContext

**文件**: `api/internal/svc/service_context.go`

```go
type ServiceContext struct {
    Config config.Config

    // Task Models
    TaskStdCreateModel         model.TaskStdCreateModel
    TaskStdCreateResultModel   model.TaskStdCreateResultModel
    BusinessTablePoolModel     model.BusinessTablePoolModel

    // RPC Clients (新增)
    AfClient              *af.AfClient
    DataElementClient     *dataelement.DataElementClient
    RecommendationClient  *recommendation.RecommendationClient
    WebhookClient         *webhook.WebhookClient
}
```

### 验证清单

- [ ] 所有 Mock 函数已替换为实际服务调用
- [ ] RPC 客户端正确配置
- [ ] HTTP 客户端包含超时和重试逻辑
- [ ] 错误处理正确实现
- [ ] 单元测试已更新

---

## A002: 定义 Webhook 失败处理策略

### 问题描述

Webhook 回调失败时的处理策略未定义。

### 修复方案

#### 1. 定义策略

**文件**: `specs/5-std-task-api/research.md`

**新增内容**:

```markdown
## Webhook 失败处理策略

### 策略选择

**选择**: 异步重试 + 降级处理

**理由**:
- Webhook 是通知性质，不应阻塞业务流程
- 任务完成应立即返回成功，webhook 异步发送
- 失败后记录日志，不重试（或有限重试）

### 实现

1. **非阻塞**: Webhook 发送失败不影响任务完成
2. **重试**: 最多重试 3 次，每次间隔递增（1s, 2s, 4s）
3. **降级**: 重试失败后记录日志，发送告警
4. **监控**: 统计 webhook 成功率

### 代码示例

```go
// 异步发送 webhook
go func() {
    err := l.svcCtx.WebhookClient.SendTaskCallback(context.Background(), webhook, taskId)
    if err != nil {
        logx.Errorf("Webhook callback failed: %v", err)
        // 发送告警
    }
}()
```
```

---

## A003: 完善 API 文档

### 问题描述

部分 API 端点缺少详细的文档注释。

### 修复方案

**文件**: `api/doc/task/task.api`

**更新**:

```api
@doc "未处理任务列表 - 查询所有状态为待处理或处理中的标准创建任务"
@handler GetUncompletedTasks
get /std-create/uncompleted (PageInfoWithKeyword) returns (TaskDataListResp)

@doc "已完成任务列表 - 查询所有状态为已提交的标准创建任务"
@handler GetCompletedTasks
get /std-create/completed (PageInfoWithKeyword) returns (TaskDataListResp)

@doc "任务详情 - 根据任务ID查询任务详情及关联的标准创建结果"
@handler GetTaskById
get /std-create/completed/:id returns (TaskDetailResp)

@doc "创建标准任务 - 将业务表字段关联到任务，状态变更为处理中"
@handler CreateTask
post /createTask (CreateTaskReq) returns (TaskBaseResp)

@doc "完成任务 - 验证所有字段都关联数据元后，批量更新状态为已完成"
@handler FinishTask
post /finishTask/:task_id returns (TaskBaseResp)

@doc "进度查询 - 查询任务的总字段数和已完成关联的字段数"
@handler QueryTaskProcess
post /queryTaskProcess returns (ProcessResp)

@doc "任务状态查询 - 根据业务表ID和状态查询业务表字段状态"
@handler QueryTaskState
post /queryTaskState returns (BusinessTableStateListResp)

@doc "采纳 - 批量采纳业务表字段关联的标准，状态变更为已采纳"
@handler Accept
put /accept (AcceptReq) returns (TaskBaseResp)

@doc "撤销 - 批量撤销业务表字段的任务关联，清空任务ID"
@handler CancelField
put /cancelBusinessTableField (CancelFieldReq) returns (TaskBaseResp)
```

---

## A005: 编写测试用例

### 问题描述

24 个 API 端点缺少测试用例。

### 修复方案

#### 测试文件结构

```
api/internal/logic/task/
├── accept_logic_test.go              # 采纳测试
├── add_to_pending_logic_test.go      # 添加至待新建测试
├── cancel_field_logic_test.go        # 撤销测试
├── create_task_logic_test.go         # 创建任务测试
├── delete_field_logic_test.go        # 移除字段测试
├── finish_task_logic_test.go         # 完成任务测试
├── get_business_table_field_logic_test.go  # 业务表字段列表测试
├── get_business_table_logic_test.go  # 业务表列表测试
├── get_completed_tasks_logic_test.go # 已完成任务列表测试
├── get_field_from_task_logic_test.go # 任务关联字段测试
├── get_table_from_task_logic_test.go # 任务关联业务表测试
├── get_task_by_id_logic_test.go      # 任务详情测试
├── get_uncompleted_tasks_logic_test.go # 未处理任务列表测试
├── query_task_process_logic_test.go  # 进度查询测试
├── query_task_state_logic_test.go    # 任务状态查询测试
├── rule_rec_logic_test.go            # 编码规则推荐测试
├── stand_rec_logic_test.go           # 标准采纳推荐测试
├── std_create_logic_test.go          # 标准关联暂存测试
├── std_rec_logic_test.go             # 标准推荐测试
├── submit_data_element_logic_test.go # 提交选定数据元测试
├── submit_relation_logic_test.go     # 标准关联提交测试
├── update_description_logic_test.go  # 修改字段说明测试
└── update_table_name_logic_test.go   # 修改表名称测试
```

#### 测试用例模板

**示例**: `create_task_logic_test.go`

```go
package task

import (
    "context"
    "testing"

    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
    "github.com/stretchr/testify/assert"
)

func TestCreateTaskLogic_CreateTask(t *testing.T) {
    // 准备测试数据
    taskId := "550e8400-e29b-41d4-a716-446655440000"
    ids := []string{"550e8400-e29b-41d4-a716-446655440001"}

    // 创建测试请求
    req := &types.CreateTaskReq{
        TaskId: taskId,
        Ids:    ids,
        Webhook: "",
    }

    // 执行测试
    resp, err := logic.CreateTask(req)

    // 断言
    assert.NoError(t, err)
    assert.Equal(t, "0", resp.Code)
}
```

### 验证清单

- [ ] 所有 24 个 Logic 文件都有对应的测试文件
- [ ] 测试覆盖正常流程
- [ ] 测试覆盖异常场景
- [ ] 测试覆盖边界条件
- [ ] 测试覆盖率 ≥ 80%

---

## 执行计划

### 修复顺序

1. **A003** (20分钟) - 完善 API 文档
   - 影响: task.api
   - 风险: 低
   - 验证: 文档完整性

2. **A002** (30分钟) - 定义 Webhook 策略
   - 影响: research.md, webhook/client.go
   - 风险: 低
   - 验证: 策略文档

3. **A001** (120分钟) - 替换 Mock 服务
   - 影响: 新增 RPC/HTTP 客户端
   - 风险: 中
   - 验证: 服务调用

4. **A005** (180分钟) - 编写测试用例
   - 影响: 新增测试文件
   - 风险: 低
   - 验证: 测试覆盖率

---

## 完成检查清单

- [ ] A001: Mock 服务已替换
- [ ] A002: Webhook 失败处理已定义
- [ ] A003: API 文档已完善
- [ ] A004: 混合 ID 类型已文档化
- [ ] A005: 测试用例已编写

**状态**: 1/5 完成 (20%)

---

**下一步**: 按照执行计划逐步完成剩余修复工作
