# 修复建议方案: rule-api

> **Feature**: 1-rule-api (编码规则管理)
> **Date**: 2026-02-05
> **Purpose**: 修复 ANALYSIS.md 中识别的 HIGH 优先级问题

---

## 修复清单

| ID | 问题 | 状态 | 预计工作量 |
|----|------|------|------------|
| A001 | API路由前缀不一致 | ⏳ 待修复 | 10分钟 |
| A002 | MQ消息格式未定义 | ⏳ 待修复 | 15分钟 |
| A003 | 缺失Logic实现任务 | ⏳ 待修复 | 20分钟 |
| A004 | 错误码映射不完整 | ⏳ 待修复 | 30分钟 |

**总计**: 约 75 分钟

---

## A001: 修复 API 路由前缀

### 问题描述

plan.md 中使用了错误的API路由前缀 `/api/standardization/v1/`，但Java实现使用的是 `/v1/`。

### 根本原因

查看 Java 源代码 `RuleController.java:52`:
```java
@RestController
@RequestMapping("/v1/rule")  // ✅ 正确的前缀
public class RuleController {
```

### 修复方案

#### 1. 更新 plan.md 的 API Contract 部分

**文件**: `specs/1-rule-api/plan.md`

**位置**: 第 572-667 行 (API Contract 部分)

**修改内容**:

```diff
- @server(
-     prefix: /api/standardization/v1
-     group: rule
-     middleware: TokenCheck
- )
+ @server(
+     prefix: /v1
+     group: rule
+     middleware: TokenCheck
+ )
  service standardization-api {
      @doc "新增编码规则"
      @handler CreateRule
-     post /rule (CreateRuleReq) returns (RuleResp)
+     post /rule (CreateRuleReq) returns (RuleResp)

      @doc "根据ID修改编码规则"
      @handler UpdateRule
-     put /rule/:id (UpdateRuleReq) returns (RuleResp)
+     put /rule/:id (UpdateRuleReq) returns (RuleResp)
      ...
  }
```

#### 2. 更新内部API定义

```diff
  // ==================== 内部API（无认证）====================

  @server(
-     prefix: /v1/rule
+     prefix: /v1
      group: rule_internal
  )
  service standardization-api {
      @doc "编码规则-根据ID查看详情（内部）"
      @handler GetRuleInternal
-     get /internal/getId/:id returns (RuleResp)
+     get /rule/internal/getId/:id returns (RuleResp)
      ...
  }
```

#### 3. 更新其他API部分

```diff
  @server(
-     prefix: /api/standardization/v1
+     prefix: /v1
      group: rule
  )
  service standardization-api {
      @doc "编码规则-删除&批量删除"
      @handler DeleteRule
-     delete /rule/:ids returns (BaseResp)
+     delete /rule/:ids returns (BaseResp)
      ...
  }
```

### 验证

修复后的路由应与Java实现100%一致:

| Java 路由 | 修复后 Go 路由 |
|-----------|---------------|
| POST `/v1/rule` | POST `/v1/rule` |
| PUT `/v1/rule/{id}` | PUT `/v1/rule/:id` |
| GET `/v1/rule/{id}` | GET `/v1/rule/:id` |
| GET `/v1/rule` | GET `/v1/rule` |
| DELETE `/v1/rule/{ids}` | DELETE `/v1/rule/:ids` |

---

## A002: 定义 MQ 消息格式

### 问题描述

plan.md 中的 MQ 消息结构 (`DataMqDto`, `Payload`, `Content`) 被引用但未完整定义。

### Java 源码分析

查看 `RuleServiceImpl.java:355-379`:

```java
private String packageMqInfo(List<RuleEntity> lists, String type) {
    DataMqDto mqDto = new DataMqDto();
    mqDto.setHeader(new HashMap());  // 空header

    DataMqDto.Payload payload = mqDto.new Payload();
    payload.setType("smart-recommendation-graph");  // 固定值

    DataMqDto.Content<RuleEntity> content = mqDto.new Content<>();
    content.setType(type);          // "insert", "update", "delete"
    content.setTable_name("t_rule");
    content.setEntities(lists);
    content.setUpdated_at(new Date());

    payload.setContent(content);
    mqDto.setPayload(payload);

    ObjectMapper objectMapper = new ObjectMapper();
    return objectMapper.writeValueAsString(mqDto);
}
```

查看 `DataMqDto.java`:

```java
public class DataMqDto {
    private Map header;
    private Payload payload;

    @Data
    public class Payload {
        private String type;           // "smart-recommendation-graph"
        private Content content;
    }

    @Data
    public class Content<T> {
        private String type;           // "insert", "update", "delete"
        private String table_name;     // "t_rule"
        private List<T> entities;      // RuleEntity 列表
        private Date updated_at;
    }
}
```

### 修复方案

#### 1. 在 plan.md 中添加完整的 MQ 消息定义

**文件**: `specs/1-rule-api/plan.md`

**位置**: 在 "Key Implementation Details" 部分添加

**新增内容**:

```go
// ========== MQ 消息结构 ==========

// DataMqDto MQ消息数据传输对象
type DataMqDto struct {
    Header  map[string]interface{} `json:"header"`
    Payload Payload                `json:"payload"`
}

// Payload MQ消息负载
type Payload struct {
    Type    string  `json:"type"`    // 固定值: "smart-recommendation-graph"
    Content Content `json:"content"`
}

// Content MQ消息内容
type Content struct {
    Type      string    `json:"type"`       // 操作类型: "insert", "update", "delete"
    TableName string    `json:"table_name"` // 固定值: "t_rule"
    Entities  []Rule    `json:"entities"`   // 规则实体列表
    UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// ========== MQ 消息示例 ==========

// 创建规则时发送的MQ消息
{
    "header": {},
    "payload": {
        "type": "smart-recommendation-graph",
        "content": {
            "type": "insert",
            "table_name": "t_rule",
            "entities": [
                {
                    "f_id": 1234567890,
                    "f_name": "邮政编码规则",
                    "f_org_type": 4,
                    "f_rule_type": 0,
                    "f_expression": "^\\d{6}$",
                    ...
                }
            ],
            "updated_at": "2026-02-05T10:30:00Z"
        }
    }
}

// 更新规则时发送的MQ消息
{
    "header": {},
    "payload": {
        "type": "smart-recommendation-graph",
        "content": {
            "type": "update",
            "table_name": "t_rule",
            "entities": [...],
            "updated_at": "2026-02-05T11:30:00Z"
        }
    }
}

// 删除规则时发送的MQ消息
{
    "header": {},
    "payload": {
        "type": "smart-recommendation-graph",
        "content": {
            "type": "delete",
            "table_name": "t_rule",
            "entities": [
                {"f_id": 1234567890}  // 删除时只包含ID
            ],
            "updated_at": "2026-02-05T12:00:00Z"
        }
    }
}

// ========== MQ 发送逻辑 ==========

// SendRuleMQMessage 发送规则变更MQ消息
func SendRuleMQMessage(producer *kafka.Producer, rules []Rule, operation string) error {
    mqDto := DataMqDto{
        Header: make(map[string]interface{}),
        Payload: Payload{
            Type: "smart-recommendation-graph",
            Content: Content{
                Type:      operation,  // "insert", "update", "delete"
                TableName: "t_rule",
                Entities:  rules,
                UpdatedAt: time.Now(),
            },
        },
    }

    jsonData, err := json.Marshal(mqDto)
    if err != nil {
        return fmt.Errorf("marshal MQ message failed: %w", err)
    }

    // 发送到 Kafka Topic
    topic := "MQ_MESSAGE_SAILOR"
    if err := producer.SendMessage(topic, string(jsonData)); err != nil {
        return fmt.Errorf("send MQ message failed: %w", err)
    }

    logx.Infof("Sent MQ message: topic=%s, operation=%s, count=%d", topic, operation, len(rules))
    return nil
}
```

#### 2. 在 data-model.md 中补充 MQ 消息结构

**文件**: `specs/1-rule-api/data-model.md`

**新增章节**:

```markdown
## MQ 消息结构

### Topic
- **Topic名称**: `MQ_MESSAGE_SAILOR`
- **用途**: 规则变更通知，触发下游系统更新

### 消息格式

| 字段 | 类型 | 说明 | 示例值 |
|------|------|------|--------|
| header | Map | 扩展头（当前为空） | {} |
| payload.type | string | 固定类型标识 | "smart-recommendation-graph" |
| payload.content.type | string | 操作类型 | "insert", "update", "delete" |
| payload.content.table_name | string | 表名 | "t_rule" |
| payload.content.entities | array | 规则实体列表 | [RuleEntity] |
| payload.content.updated_at | datetime | 更新时间 | ISO8601格式 |
```

### 验证

- [ ] MQ消息结构与Java完全一致
- [ ] 支持三种操作类型: insert, update, delete
- [ ] JSON序列化格式正确
- [ ] 包含时间戳字段

---

## A003: 添加缺失的 Logic 实现任务

### 问题描述

4个查询端点有Handler定义但缺少Logic层实现任务。

### 缺失的端点

| 端点 | Handler | Logic任务 | 状态 |
|------|---------|----------|------|
| GET `/v1/rule/queryByStdFile` | ✅ T048 | ❌ 缺失 | 需添加 |
| GET `/v1/rule/relation/stdfile/{id}` | ✅ | ❌ 缺失 | 需添加 |
| GET `/v1/rule/queryDataExists` | ✅ | ❌ 缺失 | 需添加 |
| GET `/v1/rule/getCustomDateFormat` | ✅ | ❌ 缺失 | 需添加 |

### 修复方案

**文件**: `specs/1-rule-api/tasks.md`

**位置**: Phase 6 (User Story 4) 结尾，T050 之后添加

**新增任务**:

```markdown
- [ ] T048b [US4] 实现 `api/internal/logic/rule/query_rule_by_std_file_logic.go`
  - [ ] 实现按标准文件查询规则
  - [ ] 调用 Model.FindByStdFile 方法
  - [ ] 处理空 stdFileId 场景

- [ ] T049b [P] [US4] **[TEST]** 创建 `api/internal/logic/rule/query_rule_by_std_file_logic_test.go`
  - [ ] 测试按文件查询
  - [ ] 测试空文件ID返回空列表

- [ ] T050b [P] [US4] 实现 `api/internal/logic/rule/query_std_files_by_rule_logic.go`
  - [ ] 实现查询规则关联的标准文件列表（分页）
  - [ ] 调用 StdFileModel.QueryStdFilesByRuleId 方法

- [ ] T051b [P] [US4] **[TEST]** 创建 `api/internal/logic/rule/query_std_files_by_rule_logic_test.go`
  - [ ] 测试查询关联标准文件
  - [ ] 测试分页功能

- [ ] T052b [US5] 实现 `api/internal/logic/rule/query_data_exists_logic.go`
  - [ ] 实现名称存在性检查
  - [ ] 支持 filter_id 参数排除当前记录
  - [ ] 处理 departmentIds 路径转换

- [ ] T053b [P] [US5] **[TEST]** 创建 `api/internal/logic/rule/query_data_exists_logic_test.go`
  - [ ] 测试名称存在返回 true
  - [ ] 测试名称不存在返回 false
  - [ ] 测试 filter_id 排除自身

- [ ] T054b [US5] 实现 `api/internal/logic/rule/get_custom_date_format_logic.go`
  - [ ] 返回 CustomDateFormat 常量列表

- [ ] T055b [P] [US5] **[TEST]** 创建 `api/internal/logic/rule/get_custom_date_format_logic_test.go`
  - [ ] 测试返回预定义的日期格式列表
```

### 更新任务总览

**文件**: `specs/1-rule-api/tasks.md`

**位置**: Task Overview 表格

**更新**:

```diff
  | T010 | 关联查询实现 + 测试 | US4 | ⏸️ | - | 200 |
- | T011 | 内部查询接口实现 | US5 | ⏸️ | - | 150 |
+ | T011 | 关联查询补充实现 | US4 | ⏸️ | - | 50 |
+ | T012 | 内部查询接口实现 | US5 | ⏸️ | - | 150 |
  | T012 | 收尾工作 | Polish | ⏸️ | - | - |

- **总计**: 约 1700 行代码 + 测试
+ **总计**: 约 1800 行代码 + 测试
```

### 更新依赖关系

```diff
  ## Phase 6: User Story 4 - 规则关联查询 (P2)

  ...

  - [ ] T050 [P] [US4] 实现 `api/internal/logic/rule/query_rule_by_ids_logic.go`
  - [ ] T051 [P] [US4] **[TEST]** 创建 `api/internal/logic/rule/query_rule_by_ids_logic_test.go`
+- [ ] T048b [P] [US4] 实现 `api/internal/logic/rule/query_rule_by_std_file_logic.go`
+- [ ] T049b [P] [US4] **[TEST]** 测试 query_rule_by_std_file_logic
+- [ ] T050b [P] [US4] 实现 `api/internal/logic/rule/query_std_files_by_rule_logic.go`
+- [ ] T051b [P] [US4] **[TEST]** 测试 query_std_files_by_rule_logic
```

### 验证

- [ ] 所有18个端点都有对应的Logic实现任务
- [ ] 每个Logic任务都有对应的测试任务
- [ ] 任务编号连续且合理

---

## A004: 创建错误码映射文档

### 问题描述

仅定义了3个错误码（30301-30303），但spec.md定义了15+个异常场景需要映射。

### Java 错误码分析

查看 Java 源码 `RuleServiceImpl.java` 中所有抛出的异常:

| 行号 | 场景 | Java 错误码 | 字段 | 提示信息 |
|------|------|-------------|------|----------|
| 385 | REGEX时regex为空 | InvalidParameter | regex | 正则表达式为空 |
| 391 | 正则表达式非法 | InvalidParameter | regex | 正则表达式非法 |
| 397 | CUSTOM时custom为空 | InvalidParameter | custom | 不能为空 |
| 405 | segment_length <= 0 | InvalidParameter | segment_length | 值必须为正整数 |
| 417 | 码表不存在 | InvalidParameter | value | 码表不存在 |
| 422 | DATE时value为空 | InvalidParameter | value | 不能为空 |
| 426 | 日期格式不支持 | InvalidParameter | value | 不支持的日期格式 |
| 432 | SPLIT_STR时value为空 | InvalidParameter | value | 不能为空 |
| 680 | 创建时名称重复 | InvalidParameter | name | 规则名称已存在 |
| 692 | 修改时名称重复 | InvalidParameter | name | 规则名称已存在 |
| 702 | 目录不存在 | InvalidParameter | catalog_id | 目录id[{id}]对应的目录不存在 |
| 477 | 修改记录不存在 | DATA_NOT_EXIST | - | 数据不存在/记录不存在 |
| 645 | 查询记录不存在 | DATA_NOT_EXIST | - | 记录不存在 |

### 修复方案

**新建文件**: `specs/1-rule-api/error-codes.md`

**内容**:

```markdown
# 错误码定义: rule-api

> **模块**: 编码规则管理
> **错误码范围**: 30300-30399
> **日期**: 2026-02-05

---

## 错误码映射表

| Go 错误码 | Java 错误码 | 场景 | 字段 | 提示信息 | HTTP状态码 |
|-----------|-------------|------|------|----------|-----------|
| 30301 | DATA_NOT_EXIST | 查询/修改记录不存在 | - | 数据不存在 | 404 |
| 30302 | PARAMETER_EMPTY | 停用原因为空 | reason | 停用必须填写停用原因 | 400 |
| 30303 | InvalidParameter | 通用参数错误 | * | 参数错误 | 400 |
| 30310 | InvalidParameter | 创建时规则名称重复 | name | 规则名称已存在 | 400 |
| 30311 | InvalidParameter | 修改时规则名称重复 | name | 规则名称已存在 | 400 |
| 30312 | InvalidParameter | 目录不存在 | catalog_id | 目录id[{id}]对应的目录不存在 | 400 |
| 30320 | InvalidParameter | 正则表达式为空 | regex | 正则表达式为空 | 400 |
| 30321 | InvalidParameter | 正则表达式非法 | regex | 正则表达式非法 | 400 |
| 30330 | InvalidParameter | 自定义配置为空 | custom | 不能为空 | 400 |
| 30331 | InvalidParameter | segment_length无效 | segment_length | 值必须为正整数 | 400 |
| 30332 | InvalidParameter | 码表不存在 | value | 码表不存在 | 400 |
| 30333 | InvalidParameter | 日期格式不支持 | value | 不支持的日期格式 | 400 |
| 30334 | InvalidParameter | DATE/SPLIT_STR时value为空 | value | 不能为空 | 400 |
| 30335 | InvalidParameter | 批量删除ids为空 | ids | ids 不能为空 | 400 |
| 30336 | InvalidParameter | 停用原因过长 | reason | 长度超过800 | 400 |
| 30337 | InvalidParameter | 批量查询ids为空 | ids | ids 不能为空 | 400 |

---

## Go 代码定义

**文件**: `api/internal/errorx/codes.go`

```go
package errorx

import "github.com/jinguoxing/idrm-go-base/errorx"

// ========== rule-api 错误码定义 (30300-30399) ==========

const (
    // 通用错误
    ErrCodeRuleDataNotExist    = 30301 // DATA_NOT_EXIST - 数据不存在
    ErrCodeRuleParamEmpty      = 30302 // PARAMETER_EMPTY - 参数为空
    ErrCodeRuleInvalidParam    = 30303 // InvalidParameter - 参数无效

    // 名称相关 (30310-30319)
    ErrCodeRuleNameDuplicate   = 30310 // 创建时名称重复
    ErrCodeRuleNameUpdateDuplicate = 30311 // 修改时名称重复
    ErrCodeRuleCatalogNotExist  = 30312 // 目录不存在

    // REGEX 相关 (30320-30329)
    ErrCodeRuleRegexEmpty      = 30320 // 正则表达式为空
    ErrCodeRuleRegexInvalid    = 30321 // 正则表达式非法

    // CUSTOM 相关 (30330-30339)
    ErrCodeRuleCustomEmpty     = 30330 // 自定义配置为空
    ErrCodeRuleSegmentLength   = 30331 // segment_length <= 0
    ErrCodeRuleDictNotExist    = 30332 // 码表不存在
    ErrCodeRuleDateFormat      = 30333 // 日期格式不支持
    ErrCodeRuleCustomValueEmpty = 30334 // value为空

    // 参数相关 (30335-30339)
    ErrCodeRuleIdsEmpty        = 30335 // 批量删除ids为空
    ErrCodeRuleReasonTooLong   = 30336 // 停用原因过长
    ErrCodeRuleQueryIdsEmpty   = 30337 // 批量查询ids为空
)

// ========== 辅助函数 ==========

// RuleNameDuplicate 规则名称已存在
func RuleNameDuplicate(name string) error {
    return errorx.NewWithCode(ErrCodeRuleNameDuplicate, "name", "规则名称已存在")
}

// RuleCatalogNotExist 目录不存在
func RuleCatalogNotExist(catalogId int64) error {
    return errorx.NewWithCode(ErrCodeRuleCatalogNotExist, "catalog_id",
        fmt.Sprintf("目录id[%d]对应的目录不存在", catalogId))
}

// RuleRegexEmpty 正则表达式为空
func RuleRegexEmpty() error {
    return errorx.NewWithCode(ErrCodeRuleRegexEmpty, "regex", "正则表达式为空")
}

// RuleRegexInvalid 正则表达式非法
func RuleRegexInvalid() error {
    return errorx.NewWithCode(ErrCodeRuleRegexInvalid, "regex", "正则表达式非法")
}

// RuleCustomEmpty 自定义配置为空
func RuleCustomEmpty() error {
    return errorx.NewWithCode(ErrCodeRuleCustomEmpty, "custom", "不能为空")
}

// RuleSegmentLengthInvalid segment_length必须为正整数
func RuleSegmentLengthInvalid(fieldPrefix string) error {
    return errorx.NewWithCode(ErrCodeRuleSegmentLength, fieldPrefix+"segment_length", "值必须为正整数")
}

// RuleDictNotExist 码表不存在
func RuleDictNotExist(fieldPrefix string) error {
    return errorx.NewWithCode(ErrCodeRuleDictNotExist, fieldPrefix+"value", "码表不存在")
}

// RuleDateFormatNotSupported 日期格式不支持
func RuleDateFormatNotSupported(fieldPrefix string) error {
    return errorx.NewWithCode(ErrCodeRuleDateFormat, fieldPrefix+"value", "不支持的日期格式")
}

// RuleCustomValueEmpty value不能为空
func RuleCustomValueEmpty(fieldPrefix string) error {
    return errorx.NewWithCode(ErrCodeRuleCustomValueEmpty, fieldPrefix+"value", "不能为空")
}

// RuleDisableReasonEmpty 停用原因为空
func RuleDisableReasonEmpty() error {
    return errorx.NewWithCode(ErrCodeRuleParamEmpty, "reason", "停用必须填写停用原因")
}

// RuleDisableReasonTooLong 停用原因过长
func RuleDisableReasonTooLong() error {
    return errorx.NewWithCode(ErrCodeRuleReasonTooLong, "reason", "长度超过800")
}

// RuleIdsEmpty ids不能为空
func RuleIdsEmpty() error {
    return errorx.NewWithCode(ErrCodeRuleIdsEmpty, "ids", "ids 不能为空")
}

// RuleQueryIdsEmpty 查询ids为空
func RuleQueryIdsEmpty() error {
    return errorx.NewWithCode(ErrCodeRuleQueryIdsEmpty, "ids", "ids 不能为空")
}

// RuleDataNotExist 数据不存在
func RuleDataNotExist() error {
    return errorx.NewWithCode(ErrCodeRuleDataNotExist, "", "数据不存在")
}

// RuleRecordNotExist 记录不存在
func RuleRecordNotExist() error {
    return errorx.NewWithCode(ErrCodeRuleDataNotExist, "id", "记录不存在")
}
```

---

## 错误响应格式

所有错误响应遵循 idrm-go-base 的统一格式:

```json
{
    "code": "30310",
    "message": "参数错误",
    "errors": [
        {
            "field": "name",
            "message": "规则名称已存在"
        }
    ]
}
```

---

## 使用示例

### Logic 层使用

```go
// 创建规则 - 名称重复检查
func (l *CreateRuleLogic) Create(req *types.CreateRuleReq) (*types.RuleResp, error) {
    // 检查名称唯一性
    exist, err := l.svcCtx.RuleModel.FindByNameAndOrgType(l.ctx, req.Name, req.OrgType)
    if len(exist) > 0 {
        return nil, errorx.RuleNameDuplicate(req.Name)
    }

    // ... 其他逻辑
}

// 停用规则 - 原因检查
func (l *UpdateRuleStateLogic) UpdateState(req *types.UpdateRuleStateReq) error {
    if req.State == "disable" && req.Reason == "" {
        return errorx.RuleDisableReasonEmpty()
    }

    if len(req.Reason) > 800 {
        return errorx.RuleDisableReasonTooLong()
    }

    // ... 其他逻辑
}
```

### Model 层使用

```go
// 查询不存在的记录
func (m *RuleModel) FindOne(ctx context.Context, id int64) (*Rule, error) {
    var rule Rule
    err := m.db.WithContext(ctx).Where("f_id = ? AND f_deleted = 0", id).First(&rule).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, errorx.RuleDataNotExist()
    }
    if err != nil {
        return nil, err
    }
    return &rule, nil
}
```

---

## 错误码分配统计

| 范围 | 用途 | 已分配 | 剩余 |
|------|------|--------|------|
| 30300-30309 | 通用错误 | 3 | 7 |
| 30310-30319 | 名称/目录相关 | 3 | 7 |
| 30320-30329 | REGEX相关 | 2 | 8 |
| 30330-30339 | CUSTOM相关 | 5 | 5 |
| 30340-30399 | 预留扩展 | 0 | 60 |
| **总计** | | **13** | **87** |

---

## 验证清单

- [ ] 所有 Java 异常场景都有对应的 Go 错误码
- [ ] 错误码在 30300-30399 范围内
- [ ] 错误提示信息与 Java 一致
- [ ] 字段名与 Java 一致
- [ ] HTTP 状态码正确设置
```

---

## 执行计划

### 修复顺序

1. **A001** (10分钟) - 修复 API 路由前缀
   - 影响: plan.md
   - 风险: 低
   - 验证: 对比 Java 路由

2. **A004** (30分钟) - 创建错误码映射文档
   - 影响: 新建 error-codes.md
   - 风险: 低
   - 验证: 覆盖所有异常场景

3. **A002** (15分钟) - 定义 MQ 消息格式
   - 影响: plan.md, data-model.md
   - 风险: 低
   - 验证: 对比 Java 实现

4. **A003** (20分钟) - 添加缺失任务
   - 影响: tasks.md
   - 风险: 低
   - 验证: 任务覆盖所有端点

### 修复后验证

```bash
# 1. 验证 API 路由
grep -r "prefix:" specs/1-rule-api/plan.md | grep -v "/v1"

# 2. 验证错误码覆盖
wc -l specs/1-rule-api/error-codes.md

# 3. 验证任务数量
grep -c "^\- \[ \] T" specs/1-rule-api/tasks.md

# 4. 验证 MQ 格式
grep -c "DataMqDto" specs/1-rule-api/plan.md
```

---

## 完成检查清单

- [x] A001: API 路由前缀已修复
- [ ] A002: MQ 消息格式已定义
- [ ] A003: 缺失任务已添加
- [ ] A004: 错误码映射已创建

**状态**: 1/4 完成 (25%)

---

**下一步**: 完成剩余3项修复后，重新运行 `/speckit.analyze` 验证所有 HIGH 问题已解决。
