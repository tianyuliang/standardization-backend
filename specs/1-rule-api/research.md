# Research Notes: 编码规则管理 (rule-api)

> **Feature**: rule-api
> **Date**: 2026-02-05
> **Purpose**: Document technical decisions and research findings for Java to Go migration

---

## Migration Strategy

### Decision: 保持 Java 表结构和主键类型

**选择**: 使用 BIGINT 自增主键，而非 UUID v7

**理由**:
- Java 实现使用 `BIGINT AUTO_INCREMENT`，前端已依赖此格式
- 数据库表已存在且有数据，无法更改主键类型
- 这是唯一例外，其他模块仍遵循 UUID v7 规范

**影响**:
- Model 层使用 `int64` 类型而非 `string`
- 生成 ID 使用数据库自增，无需服务端生成

---

## ORM Selection

### Decision: 双 ORM 策略

**选择**:
- **GORM**: 复杂查询（多表 JOIN、动态条件、分页排序）
- **SQLx**: 简单 CRUD（单表操作）

**理由**:
- GORM 提供强大的关联查询能力，适合 `queryByStdFileCatalog` 等多表 JOIN 场景
- SQLx 性能更高，适合 `Insert/Update/Delete` 等简单操作
- Java 实现使用 MyBatis，Go 的 GORM 更接近 MyBatis 的灵活查询能力

**对比**:

| 场景 | Java (MyBatis) | Go (GORM) | Go (SQLx) |
|------|----------------|-----------|-----------|
| 复杂多表查询 | ✅ 动态 SQL | ✅ 链式查询 | ❌ 手写 SQL |
| 简单 CRUD | ✅ Mapper 方法 | ✅ 基础操作 | ✅ 高性能 |
| 事务管理 | ✅ @Transactional | ✅ Tx | ✅ Tx |

---

## API Path Design

### Decision: 保持 Java 接口路径

**选择**: 完全复用 Java 的 API 路径结构

**理由**:
- 前端已硬编码这些路径
- 100% 兼容性要求

**路径映射**:

| Java | Go |
|------|-----|
| `/v1/rule` | `/v1/rule` (内部接口) |
| `/api/standardization/v1/rule` | `/api/standardization/v1/rule` (外部接口) |

---

## Error Handling

### Decision: 使用 idrm-go-base errorx

**选择**: 统一使用通用库的错误处理

**错误码映射**:

| Java 错误码 | Go 错误码 | 含义 |
|-------------|-----------|------|
| `Standardization.DATA_NOT_EXIST` | `30301` | 数据不存在 |
| `Standardization.PARAMETER_EMPTY` | `30302` | 参数为空 |
| `Standardization.InvalidParameter` | `30303` | 参数无效 |

**响应格式**（与 Java 一致）:

```json
{
    "code": "30303",
    "description": "参数值校验不通过",
    "detail": [
        {"Key": "name", "Message": "规则名称已存在"}
    ],
    "solution": "请使用请求参数构造规范化的请求字符串。"
}
```

---

## Regular Expression Validation

### Decision: 使用 Go 标准库 regexp

**选择**: `regexp.Compile()` 校验正则表达式有效性

**理由**:
- Go 标准库已提供完整的正则支持
- 与 Java 的 `Pattern.compile()` 行为一致

**实现**:
```go
if _, err := regexp.Compile(regex); err != nil {
    return errorx.NewWithCode(30303, "regex", "正则表达式非法")
}
```

---

## Custom Rule Date Format

### Decision: 硬编码支持格式列表

**选择**: 与 Java 保持一致的日期格式集合

**支持的格式**:
```go
var CustomDateFormat = []string{
    "yyyy", "yyyyMM", "yyyy-MM-dd", "yyyyMMdd",
    "yyyy-MM-dd HH:mm:ss", "yyyyMMddHHmmss",
    "yyyy-MM-dd'T'HH:mm:ss", "yyyyMMdd'T'HHmmss",
    "HHmmss", "HH:mm:ss",
}
```

---

## Kafka Integration

### Decision: 使用 sarama 库

**选择**: `github.com/IBM/sarama` (Kafka Go 客户端)

**理由**:
- 成熟的 Kafka Go 客户端
- 与 Java KafkaProducer 功能对等

**消息格式**（与 Java 一致）:
```json
{
    "header": {},
    "payload": {
        "type": "smart-recommendation-graph",
        "content": {
            "type": "insert/update/delete",
            "table_name": "t_rule",
            "entities": [...]
        }
    }
}
```

---

## Department ID Processing

### Decision: 路径处理逻辑

**选择**:
- **存储**: 完整路径（如 "a/ab"）
- **查询**: 使用 `LIKE '%dept_id%'` 匹配
- **返回**: 仅最后一段（从 "a/ab" 提取 "ab"）

**理由**: 与 Java `StringUtil.PathSplitAfter()` 行为一致

---

## Version Control

### Decision: 版本递增字段

**选择**: 修改以下字段时版本号 +1

**触发字段**:
- name
- catalog_id
- department_ids
- org_type
- description
- rule_type
- expression
- 关联文件（stdFiles 变化）

**实现**: 在 `Update` Logic 中比较新旧值，有任何变化则递增版本

---

## XSS Prevention

### Decision: Keyword 参数转义

**选择**: 对 `keyword` 查询参数进行 XSS 转义

**理由**: 与 Java `StringUtil.XssEscape()` 行为一致

**实现**: 使用 `html.EscapeString()` 或自定义转义函数

---

## Dependencies Resolution

### Decision: 服务调用方式

**选择**:
- 目录服务: 通过 Model 层直接查询数据库（同库）
- 数据元服务: 通过 Model 层直接查询数据库（同库）
- 字典服务: 通过 Model 层直接查询数据库（同库）
- 标准文件服务: 通过 Model 层直接查询数据库（同库）

**理由**: Java 实现中这些服务都在同一应用内，Go 实现也保持这种方式

---

## Conclusion

所有技术决策均以 **100% 保持 Java 兼容性** 为前提。关键决策点：

1. ✅ BIGINT 主键（例外情况）
2. ✅ 完全复用 Java 表结构
3. ✅ 18 个 API 路径和参数保持一致
4. ✅ 错误码和响应格式完全一致
5. ✅ 业务逻辑与 Java 实现对等
6. ✅ MQ 消息格式保持一致
