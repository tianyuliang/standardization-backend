# Research Notes: 数据元管理 (dataelement-api)

> **Feature**: dataelement-api
> **Date**: 2026-02-09
> **Purpose**: Document technical decisions and research findings for Java to Go migration

---

## Migration Strategy

### Decision: 保持 Java 表结构和主键类型

**选择**: 使用 BIGINT 自增主键，而非 UUID v7

**理由**:
- Java 实现使用 `BIGINT AUTO_INCREMENT`，前端已依赖此格式
- 数据库表已存在且有数据，无法更改主键类型
- 与已完成的 rule-api、dict-api、task-api 保持一致

**影响**:
- Model 层使用 `int64` 类型而非 `string`
- 生成 ID 使用数据库自增，无需服务端生成
- `f_de_id` 等于 `f_id`

---

## ORM Selection

### Decision: 纯 SQLx 策略

**选择**: 所有数据访问使用 SQLx，手工编写 SQL 查询

**理由**:
- SQLx 性能更高，适合高并发场景
- 手写 SQL 更容易优化和调试
- Java 实现使用 MyBatis，SQLx 的手动 SQL 更接近 MyBatis 模式
- 与 rule-api、task-api 保持一致的技术栈

**对比**:

| 场景 | Java (MyBatis) | Go (SQLx) |
|------|----------------|-----------|
| 简单 CRUD | ✅ Mapper 方法 | ✅ 高性能 |
| 动态查询 | ✅ 动态 SQL | ✅ 手写 SQL |
| 事务管理 | ✅ @Transactional | ✅ Tx |
| 分页查询 | ✅ Mapper | ✅ 手写 LIMIT |

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
| `/v1/dataelement` | `/v1/dataelement` (内部接口) |
| `/api/standardization/v1/dataelement` | `/api/standardization/v1/dataelement` (外部接口) |

---

## Error Handling

### Decision: 使用 idrm-go-base errorx

**选择**: 统一使用通用库的错误处理

**错误标识格式**（不使用数字码值）:

```go
// api/internal/errorx/dataelement.go

// DataElementNotExist 数据不存在
func DataElementNotExist() error {
    return errorx.New("DataElementNotExist", "数据不存在")
}

// ParameterEmpty 参数为空
func ParameterEmpty(field string) error {
    return errorx.New("ParameterEmpty", fmt.Sprintf("[%s]:参数为空", field))
}

// InvalidParameter 参数无效
func InvalidParameter(field, message string) error {
    return errorx.New("InvalidParameter", fmt.Sprintf("[%s]:%s", field, message))
}
```

**响应格式**（与 Java 一致）:

```json
{
    "code": "InvalidParameter",
    "description": "参数值校验不通过",
    "detail": [
        {"Key": "catalog_id", "Message": "数据元对应的目录不存在"}
    ]
}
```

---

## Data Type Handling

### Decision: 数据类型枚举

| 枚举值 | code | 说明 |
|--------|------|------|
| Number | 0 | 数字型 |
| Char | 1 | 字符型 |
| Date | 2 | 日期型 |
| DateTime | 3 | 日期时间型 |
| Boolean | 5 | 布尔型 |
| Decimal | 7 | 高精度型 |
| Binary | 8 | 小数型 |
| Time | 9 | 时间型 |
| Integer | 10 | 整数型 |
| Unknown | 99 | 未知 |

**Go 常量定义**:

```go
const (
    DataTypeNumber     int32 = 0  // 数字型
    DataTypeChar        int32 = 1  // 字符型
    DataTypeDate        int32 = 2  // 日期型
    DataTypeDateTime    int32 = 3  // 日期时间型
    DataTypeBoolean     int32 = 5  // 布尔型
    DataTypeDecimal     int32 = 7  // 高精度型
    DataTypeBinary      int32 = 8  // 小数型
    DataTypeTime        int32 = 9  // 时间型
    DataTypeInteger     int32 = 10 // 整数型
    DataTypeUnknown     int32 = 99 // 未知
)
```

---

## Value Range Calculation

### Decision: 值域计算规则

```go
// CalculateValueRange 计算值域
func CalculateValueRange(de *DataElement, dictData string) string {
    // 关联码表时返回码表值
    if de.RelationType == "codeTable" && dictData != "" {
        return dictData // 格式: [M,F]
    }

    // Number/Decimal 类型计算范围
    if de.DataType == DataTypeNumber || de.DataType == DataTypeDecimal {
        length := de.DataLength ?? 38
        precision := de.DataPrecision ?? 0
        min := -int(math.Pow10(float64(length-precision))) + 1
        max := int(math.Pow10(float64(length-precision))) - 1
        return fmt.Sprintf("(%d,%d)", min, max)
    }

    return ""
}
```

---

## Version Control

### Decision: 版本控制规则

**版本递增场景**:
- 关联类型（relationType）变更
- 关联对象（dictCode/ruleId）变更
- 名称（nameCn/nameEn）变更

**版本不递增场景**:
- 目录ID（catalogId）变更
- 数据类型（dataType）变更
- 数据长度/精度变更
- 说明（description）变更

---

## Import/Export Strategy

### Decision: Excel 导入导出

**导入限制**:
- 单次最多 5000 条
- 文件大小限制 10M
- 支持 .xlsx, .xls 格式

**导出策略**:
- catalog_id 为空时导出模板
- 目录不存在或类型不对时导出模板

---

## MQ Message Format

### Decision: MQ 消息结构

**Topic**: `MQ_MESSAGE_SAILOR`

```go
type DataMqDto struct {
    Header  map[string]interface{} `json:"header"`
    Payload Payload                  `json:"payload"`
}

type Payload struct {
    Type    string  `json:"type"`    // 固定值: "smart-recommendation-graph"
    Content Content `json:"content"`
}

type Content struct {
    Type       string      `json:"type"`       // insert/update/delete
    TableName string      `json:"table_name"` // t_data_element_info
    Entities  []DataElement `json:"entities"`
}
```

---

## Batch Processing

### Decision: 批量操作处理

```go
// 批量删除
func (l *DeleteDataElementLogic) DeleteDataElement(ids string) error {
    // 将字符串ID转换为int64数组
    idList := strings.Split(ids, ",")
    var idsInt []int64
    for _, idStr := range idList {
        id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
        if err != nil || id <= 0 {
            return errorx.InvalidParameter("ids", "ID格式错误")
        }
        idsInt = append(idsInt, id)
    }

    // 批量物理删除
    // 批量删除关联关系
    // 发送MQ消息
}
```

---

## Catalog ID Processing

### Decision: 目录ID 处理逻辑

```go
// GetCatalogIds 获取当前目录及所有子目录ID
func GetCatalogIds(catalogId int64) ([]int64, error) {
    // TODO: 调用 Catalog RPC 获取子目录ID列表
    // 当前返回 mock 数据
    return []int64{catalogId}, nil
}
```

---

## Conclusion

所有技术决策均以 **100% 保持 Java 兼容性** 为前提。关键决策点：

1. ✅ BIGINT 主键（与 rule-api、dict-api、task-api 一致）
2. ✅ 纯 SQLx ORM 策略（高性能，手写 SQL）
3. ✅ 完全复用 Java 表结构（f_ 前缀字段）
4. ✅ 19个 API 路径和参数保持一致
5. ✅ 错误标识使用英文+消息（不使用数字码值）
6. ✅ 批量操作处理逻辑一致
7. ✅ 版本控制规则与 Java 一致

---

**文档版本**: v1.0
**更新时间**: 2026-02-09
**维护人**: AI Assistant
