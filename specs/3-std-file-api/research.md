# Research Notes: 标准文件管理 (std-file-api)

> **Feature**: std-file-api
> **Date**: 2026-02-06
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

### Decision: 纯 SQLx 策略

**选择**: 所有数据库操作使用 SQLx，手工编写 SQL 查询

**理由**:
- 规范要求使用 sqlx 驱动方式
- SQLx 性能更高，适合所有 CRUD 操作
- 手写 SQL 可以精确控制查询逻辑，与 Java MyBatis 保持一致
- 避免引入 GORM 依赖，保持代码轻量

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
| `/v1/std-file` | `/api/standardization/v1/std-file` (外部接口) |

---

## Error Handling

### Decision: 使用 idrm-go-base errorx

**选择**: 统一使用通用库的错误处理

**错误码映射**:

| Java 错误码 | Go 错误码 | 含义 |
|-------------|-----------|------|
| `Standardization.DATA_NOT_EXIST` | `30201` | 数据不存在 |
| `Standardization.PARAMETER_EMPTY` | `30202` | 参数为空 |
| `Standardization.InvalidParameter` | `30203` | 参数无效 |
| `Standardization.DATA_EXIST` | `30204` | 数据已存在 |
| `Standardization.FileDownloadFailed` | `30205` | 文件下载失败 |

**响应格式**（与 Java 一致）:

```json
{
    "code": "30203",
    "description": "参数值校验不通过",
    "detail": [
        {"Key": "number", "Message": "标准编号重复"}
    ],
    "solution": "请使用请求参数构造规范化的请求字符串。"
}
```

---

## File Upload/Download Strategy

### Decision: OSS 集成方式

**选择**: 使用阿里云 OSS SDK 进行文件存储

**理由**:
- Java 实现使用 OSS 文件上传下载工具 `ossFileUploadDownloadUtil`
- 需要保持文件存储位置一致

**支持的文件类型**:
- doc, docx
- pdf
- txt
- ppt, pptx
- xls, xlsx

**文件大小限制**: 最大 30M

---

## Attachment Type Handling

### Decision: 两种附件类型处理

**选择**:
1. **FILE (0)**: 文件附件，存储在 OSS
2. **URL (1)**: 外置链接，仅存储 URL 字符串

**下载逻辑**:
- FILE 类型：从 OSS 下载文件流
- URL 类型：返回错误提示（Java 实现行为）

---

## Version Control

### Decision: 版本递增字段

**选择**: 修改以下字段时版本号 +1

**触发字段**:
- number (标准编号)
- name (文件名称)
- catalog_id (目录ID)
- department_ids (部门ID)
- org_type (标准组织类型)
- description (说明)
- act_date (实施日期)
- attachment_type (附件类型)
- attachment_url (链接地址)
- file_name (文件名)

**实现**: 在 `Update` Logic 中比较新旧值，有任何变化则递增版本

---

## Department ID Processing

### Decision: 路径处理逻辑

**选择**:
- **存储**: 完整路径（如 "a/ab"）
- **查询**: 使用 `LIKE '%dept_id%'` 匹配
- **返回**: 仅最后一段（从 "a/ab" 提取 "ab"）

**理由**: 与 Java `StringUtil.PathSplitAfter()` 行为一致

---

## XSS Prevention

### Decision: Keyword 参数转义

**选择**: 对 `keyword`、`name` 查询参数进行 XSS 转义

**理由**: 与 Java `StringUtil.XssEscape()` 行为一致

**实现**: 使用自定义转义函数 `escapeSqlSpecialChars()`

---

## Batch Download Strategy

### Decision: ZIP 压缩下载

**选择**: 多文件批量下载时打包成 ZIP

**理由**: 与 Java 实现行为一致

**文件名处理规则**:
- 如果文件名不重复：使用原文件名
- 如果文件名重复：`{文件前缀}({标准分类})({标准文件名称}).{后缀}`

**支持的文件类型**:
- 仅 FILE 类型的文件可以下载
- URL 类型的文件会被过滤掉（返回错误提示）

---

## Dependencies Resolution

### Decision: 服务调用方式

**选择**:
- 目录服务: 通过 Model 层直接查询数据库（同库）
- 数据元服务: 通过 Mock 函数实现，后续补充 RPC
- 字典服务: 通过 Mock 函数实现，后续补充 RPC
- 编码规则服务: 通过 Mock 函数实现，后续补充 RPC

**理由**: Java 实现中这些服务可能在不同应用，Go 实现阶段使用 Mock 模拟

**Mock 目录结构**:
```
api/internal/logic/stdfile/mock/
├── catalog.go      # 目录服务 Mock
├── dataelement.go  # 数据元服务 Mock
├── dict.go         # 码表服务 Mock
└── rule.go         # 编码规则服务 Mock
```

---

## State Management

### Decision: 启用/停用状态处理

**选择**:
- **启用**: `state = 1`, 清空 `disable_date` 和 `disable_reason`
- **停用**: `state = 0`, 设置 `disable_date` 和 `disable_reason`

**停用原因校验**:
- 停用时必须填写原因，否则返回错误
- 停用原因长度不能超过 800 字符

---

## Date Handling

### Decision: 日期格式

**选择**: 使用 `yyyy-MM-dd` 格式

**支持的日期字段**:
- `act_date` (实施日期)
- `publish_date` (发布日期)
- `disable_date` (停用时间)

**实现**: 使用自定义函数 `parseDate()` 解析日期字符串

---

## Unique Constraints

### Decision: 唯一性校验

**标准编号唯一性**:
- 全局唯一，不允许重复

**文件名称唯一性**:
- 同一 `org_type` 下唯一
- 不同 `org_type` 可以有相同的文件名称

**实现**: 在 Model 层提供 `FindByNumber()` 和 `FindByNameAndOrgType()` 方法

---

## Catalog ID Default Value

### Decision: 默认目录 ID

**选择**: 默认值为 `44`（全部目录）

**理由**: 与 Java 实现保持一致

---

## Conclusion

所有技术决策均以 **100% 保持 Java 兼容性** 为前提。关键决策点：

1. ✅ BIGINT 主键（例外情况）
2. ✅ 完全复用 Java 表结构
3. ✅ 16 个 API 路径和参数保持一致
4. ✅ 错误码和响应格式完全一致
5. ✅ 业务逻辑与 Java 实现对等
6. ✅ 文件上传下载使用 OSS
7. ✅ 批量下载使用 ZIP 压缩
8. ✅ Mock 函数统一管理在 `logic/stdfile/mock/` 目录
