# Research Notes: 码表管理 (dict-api)

> **Feature**: dict-api
> **Date**: 2026-02-06
> **Purpose**: Document technical decisions and research findings for Java to Go migration

---

## Migration Strategy

### Decision: 保持 Java 表结构和主键类型

**选择**: 使用 BIGINT 自增主键，而非 UUID v7

**理由**:
- Java 实现使用 `BIGINT AUTO_INCREMENT`，前端已依赖此格式
- 数据库表已存在且有数据，无法更改主键类型
- 与 rule-api、stdfile-api 保持一致
- 这是唯一例外，其他模块仍遵循 UUID v7 规范

**影响**:
- Model 层使用 `int64` 类型而非 `string`
- 生成 ID 使用数据库自增，无需服务端生成
- 码表编码（code）使用雪花算法生成

---

## ORM Selection

### Decision: 纯 SQLx 策略

**选择**: 所有数据库操作使用 SQLx，手工编写 SQL 查询

**理由**:
- 与 stdfile-api 保持一致（参考其实现模式）
- SQLx 性能更高，适合高频操作
- Java 实现使用 MyBatis，手工 SQL 更接近 MyBatis 的灵活性
- 便于优化查询性能

**对比**:

| 场景 | Java (MyBatis) | Go (SQLx) |
|------|----------------|-----------|
| 复杂多表查询 | ✅ 动态 SQL | ✅ 手写 SQL |
| 简单 CRUD | ✅ Mapper 方法 | ✅ 高性能 |
| 事务管理 | ✅ @Transactional | ✅ Tx |

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
| `/v1/dataelement/dict` | `/api/standardization/v1/dataelement/dict` (外部接口，带 TokenCheck) |
| 内部接口（无认证） | `/api/standardization/v1/dataelement/dict` |

**废弃接口**:
- POST /v1/dataelement/dict/import - 导入（不实现）
- POST /v1/dataelement/dict/export - 导出（不实现）

---

## Error Handling

### Decision: 使用 idrm-go-base errorx

**选择**: 统一使用通用库的错误处理，单独文件

**错误码映射**:

| Java 错误码 | Go 错误码 | 含义 | 单独文件 |
|-------------|-----------|------|----------|
| `Standardization.DATA_NOT_EXIST` | `30401` | 数据不存在 | dict.go |
| `Standardization.PARAMETER_EMPTY` | `30402` | 参数为空 | dict.go |
| `Standardization.InvalidParameter` | `30403` | 参数无效 | dict.go |
| - | `30404` | 目录不存在 | dict.go |
| - | `30405` | 中文名称重复 | dict.go |
| - | `30406` | 英文名称重复 | dict.go |
| - | `30407` | 码值为空 | dict.go |
| - | `30408` | 码值描述为空 | dict.go |
| - | `30409` | 码值重复 | dict.go |
| - | `30410` | 停用原因过长 | dict.go |

**响应格式**（与 Java 一致）:

```json
{
    "code": "30405",
    "description": "码表中文名称、标准分类不能全部重复"
}
```

---

## Snowflake Algorithm

### Decision: 雪花算法生成码表编码

**选择**: 使用雪花算法生成全局唯一的 Long 类型 code

**理由**:
- Java 实现使用雪花算法生成码表编码
- 需要保证全局唯一性，避免分布式环境冲突
- 与 Java 实现保持一致

**实现**:

```go
import "github.com/sony/sonyflake"

var snowflake *sonyflake.Snowflake

func init() {
    snowflake, _ = sonyflake.New(0) // 使用默认配置
}

func generateCode() int64 {
    return snowflake.NextID()
}
```

---

## Department ID Processing

### Decision: 路径处理逻辑

**选择**:
- **存储**: 完整路径（如 "a/ab"）
- **查询**: 使用 `LIKE '%dept_id%'` 匹配
- **返回**: 仅最后一段（从 "a/ab" 提取 "ab"）

**理由**: 与 Java `StringUtils.PathSplitAfter()` 行为一致

**实现**:

```go
// 存储完整路径
departmentIds := "a/ab"

// 查询时使用 LIKE
query := "SELECT * FROM t_dict WHERE f_department_ids LIKE ?"
db.Query(query, "%"+deptId+"%")

// 返回时提取最后一段
lastSegment := PathSplitAfter(departmentIds) // "ab"
```

---

## Version Control

### Decision: 版本递增字段

**选择**: 修改以下字段时版本号 +1

**触发字段**:
- chName（中文名称）
- enName（英文名称）
- catalogId（所属目录）
- departmentIds（部门ID）
- orgType（组织类型）
- description（业务含义）
- enums（码值列表）
- stdFiles（关联文件）

**实现**: 在 `Update` Logic 中比较新旧值，有任何变化则递增版本

---

## Code Uniqueness Validation

### Decision: 名称唯一性校验

**选择**:
- 同一 orgType 下中文名称必须唯一
- 同一 orgType 下英文名称必须唯一
- 修改时排除当前记录自身

**实现**:

```go
// 创建时校验
func CheckChNameUnique(model DictModel, chName string, orgType int32) error {
    existing, _ := model.FindByChNameAndOrgType(ctx, chName, orgType)
    if existing != nil {
        return errorx.NewWithMsg(ErrCodeDictChNameDuplicate, "码表中文名称、标准分类不能全部重复")
    }
    return nil
}

// 修改时校验（排除自身）
func CheckChNameUniqueExcludeSelf(model DictModel, id, chName string, orgType int32) error {
    existing, _ := model.FindByChNameAndOrgType(ctx, chName, orgType)
    if existing != nil && existing.Id != id {
        return errorx.NewWithMsg(ErrCodeDictChNameDuplicate, "码表中文名称、标准分类不能全部重复")
    }
    return nil
}
```

---

## Enum Code Validation

### Decision: 码值唯一性校验

**选择**:
- 同一码表内 code 必须唯一
- 创建和修改时都需要校验

**实现**:

```go
func CheckEnumCodesUnique(enums []DictEnumVo) error {
    codeMap := make(map[string]bool)
    for _, enum := range enums {
        if enum.Code == "" {
            return errorx.NewWithMsg(ErrCodeDictEnumCodeEmpty, "码值输入不能为空")
        }
        if codeMap[enum.Code] {
            return errorx.NewWithMsg(ErrCodeDictEnumCodeDuplicate, "码值出现重复记录")
        }
        codeMap[enum.Code] = true
    }
    return nil
}
```

---

## Disable Reason Validation

### Decision: 停用原因校验

**选择**:
- 停用（state=disable）时原因必填
- 原因长度不能超过800字符
- 启用时清空原因

**实现**:

```go
if state == StateDisable {
    if strings.TrimSpace(req.Reason) == "" {
        return errorx.NewWithMsg(ErrCodeDictParamEmpty, "停用原因不能为空")
    }
    if len([]rune(req.Reason)) > 800 {
        return errorx.NewWithMsg(ErrCodeDictReasonTooLong, "停用原因长度不能超过800字符")
    }
}
```

---

## Dependencies Resolution

### Decision: 服务调用方式

**选择**:
- 目录服务: 使用 Mock，后续补充 RPC
- 数据元服务: 使用 Mock，后续补充 RPC
- 标准文件服务: 使用 Mock，后续补充 RPC

**理由**:
- 这些服务目前未实现 Go 版本
- 使用 Mock 可以保证开发进度
- 后续通过 build tag 控制是否使用 Mock

**Mock 目录结构**:

```
api/internal/logic/dict/mock/
├── catalog.go       // 目录服务 Mock
├── dataelement.go   // 数据元服务 Mock
├── stdfile.go       // 标准文件服务 Mock
└── token.go         // Token 服务 Mock
```

---

## Conclusion

所有技术决策均以 **100% 保持 Java 兼容性** 为前提。关键决策点：

1. ✅ BIGINT 主键（例外情况）
2. ✅ 完全复用 Java 表结构
3. ✅ 16 个 API 路径和参数保持一致（排除2个废弃接口）
4. ✅ 错误码和响应格式完全一致（30400-30499）
5. ✅ 业务逻辑与 Java 实现对等
6. ✅ 雪花算法生成码表编码
7. ✅ 纯 SQLx 数据访问策略
