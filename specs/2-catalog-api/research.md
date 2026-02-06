# Research Notes: 目录管理 (catalog-api)

> **Feature**: catalog-api
> **Date**: 2026-02-06
> **Purpose**: Document technical decisions and research findings for Java to Go migration

---

## Migration Strategy

### Decision: 保持 Java 表结构和主键类型

**选择**: 使用 BIGINT 自增主键

**理由**:
- Java 实现使用 `BIGINT AUTO_INCREMENT`
- 数据库表已存在且有数据，无法更改主键类型
- 与 rule-api 保持一致

**影响**:
- Model 层使用 `int64` 类型
- 生成 ID 使用数据库自增

---

## ORM Selection

### Decision: 使用纯 SQLx 策略

**选择**: 统一使用 SQLx 进行所有数据库操作

**理由**:
- catalog-api 的操作相对简单（单表 CRUD + 树形结构构建）
- SQLx 性能更高，代码更直观
- 避免引入 GORM 的复杂性
- 与 rule-api 的双 ORM 策略不同，catalog-api 选择纯 SQLx

**对比**:

| 场景 | Java (MyBatis) | Go (SQLx) |
|------|----------------|-----------|
| 单表 CRUD | ✅ Mapper 方法 | ✅ 直接 SQL |
| 树形结构查询 | ✅ 递归查询 | ✅ 应用层构建 |
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
| `/v1/catalog/query_tree` | `/api/standardization/v1/catalog/query_tree` |
| `/v1/catalog/query` | `/api/standardization/v1/catalog/query` |
| `/v1/catalog` | `/api/standardization/v1/catalog` (POST/PUT/DELETE) |
| `/v1/catalog/query/with_file` | `/api/standardization/v1/catalog/query/with_file` |

---

## Error Handling

### Decision: 使用 idrm-go-base errorx

**选择**: 统一使用通用库的错误处理

**错误码映射**:

| Java 错误码 | Go 错误码 | 含义 |
|-------------|-----------|------|
| `Standardization.DATA_NOT_EXIST` | `30101` | 数据不存在 |
| `Standardization.PARAMETER_EMPTY` | `30102` | 参数缺失 |
| `Standardization.InvalidParameter` | `30103` | 参数无效 |
| `Standardization.OutOfRange` | `30104` | 超出范围 |
| `Standardization.OperationConflict` | `30105` | 操作冲突 |
| `Standardization.DATA_EXIST` | `30106` | 数据已存在 |

**响应格式**（与 Java 一致）:

```json
{
    "code": "30103",
    "description": "参数值校验不通过",
    "detail": [
        {"Key": "catalogName", "Message": "目录名称长度不能超过20个字符"}
    ],
    "solution": "请使用请求参数构造规范化的请求字符串。"
}
```

---

## Catalog Name Validation

### Decision: 使用 Go 标准库 regexp

**选择**: `regexp.MatchString()` 校验目录名称格式

**理由**:
- Go 标准库已提供完整的正则支持
- 与 Java 的 `Pattern.matches()` 行为一致

**正则规则**:
```go
// 中英文、数字、下划线、中划线，不以_-开头，最大20字符
regex := `^[\p{Han}a-zA-Z0-9][\p{Han}a-zA-Z0-9_-]{0,19}$`
```

---

## Tree Structure Building

### Decision: 应用层递归构建

**选择**: 在应用层（Logic）构建树形结构，而非数据库层

**理由**:
- Java 实现在应用层构建树
- 支持按需包含数据统计
- 更灵活的控制（带/不带 Count）

**实现**:
```go
func BuildTree(catalogs []*Catalog, rootLevel int32) []*CatalogResp {
    // 按level分组
    levelMap := make(map[int64][]*Catalog)
    for _, c := range catalogs {
        levelMap[int64(c.Level)] = append(levelMap[int64(c.Level)], c)
    }

    // 获取根节点
    roots := make([]*CatalogResp, 0)
    for _, c := range catalogs {
        if c.Level == rootLevel {
            roots = append(roots, modelToResp(c))
        }
    }

    // 递归设置子节点
    for _, root := range roots {
        setChildren(root, levelMap)
    }

    return roots
}
```

---

## Logical Deletion

### Decision: 逻辑删除标记

**选择**: 使用 `f_deleted = f_id + 1` 标记删除

**理由**:
- 与 Java 实现完全一致
- 避免使用 0/1 布尔值（可能与业务字段冲突）
- 删除后的 ID 可以唯一标识被删除的记录

**实现**:
```sql
UPDATE t_de_catalog_info SET f_deleted = f_id + 1 WHERE f_id = ?
```

---

## Catalog Type Enumeration

### Decision: 四种目录类型

**类型定义**:
```go
const (
    CatalogTypeDataElement = 1 // 数据元
    CatalogTypeDict          = 2 // 码表
    CatalogTypeValueRule     = 3 // 编码规则
    CatalogTypeFile          = 4 // 标准文件
)
```

**规则**:
- 子目录继承父目录的 Type
- Type 不能修改（创建时确定）
- 查询时按 Type 过滤

---

## Recursive Operations

### Decision: 递归获取所有子级 ID

**选择**: 先查询所有同级及下级目录，再应用层过滤

**理由**:
- 避免数据库递归查询
- Java 实现使用此方式
- 支持批量操作

**实现**:
```go
func GetAllChildIds(ctx context.Context, catalogId int64) ([]int64, error) {
    // 1. 获取当前目录
    catalog, err := model.FindOne(ctx, catalogId)
    if err != nil {
        return nil, err
    }

    // 2. 获取所有 level >= catalog.level 的目录
    allCatalogs, err := model.FindByTypeAndLevel(ctx, catalog.Type, catalog.Level)

    // 3. 过滤出后代目录
    childIds := make([]int64, 0)
    for _, cat := range allCatalogs {
        if isDescendant(cat, catalog) && cat.Id != catalogId {
            childIds = append(childIds, cat.Id)
        }
    }

    return childIds, nil
}
```

---

## Delete Validation

### Decision: 删除前多层校验

**校验层级**:
1. 根目录检查（level <= 1 不允许删除）
2. 数据存在性检查（根据 Type 查询对应数据）
3. 递归删除所有子目录

**理由**: 与 Java `checkCatalogDelete()` 逻辑完全一致

---

## Dependencies Resolution

### Decision: Model 层直接查询数据库

**选择**:
- 目录管理: 直接查询 `t_de_catalog_info`
- 数据元检查: TODO - 后续通过 RPC 调用 dataelement-api
- 码表检查: TODO - 后续通过 RPC 调用 dict-api
- 规则检查: TODO - 后续通过 RPC 调用 rule-api
- 文件检查: TODO - 后续通过 RPC 调用 stdfile-api

**当前状态**: 使用 Mock 函数占位，后续补充 RPC 调用

---

## Conclusion

所有技术决策均以 **100% 保持 Java 兼容性** 为前提。关键决策点：

1. ✅ BIGINT 主键
2. ✅ 完全复用 Java 表结构
3. ✅ 6 个 API 路径和参数保持一致
4. ✅ 错误码和响应格式完全一致
5. ✅ 业务逻辑与 Java 实现对等
6. ✅ 树形结构构建逻辑一致
