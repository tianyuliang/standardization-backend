# 目录管理 API (Catalog API) 转写 Specification

> **Branch**: `feature/catalog-api`
> **Spec Path**: `specs/catalog-api/`
> **Created**: 2025-01-21
> **Status**: Draft

---

## Overview

将 source/ 目录下 Java Spring Boot 项目的目录管理 API 转写为 Go-Zero 微服务架构。这是一个**100% 还原转写**任务，要求：

- **接口路由 100% 还原** - 保持所有 HTTP 端点路径、方法、参数不变
- **协议格式 100% 还原** - 保持请求/响应 JSON 结构不变
- **业务逻辑 100% 还原** - 保持所有校验规则、业务流程不变
- **禁止修改业务** - 不新增、不修改、不删除任何业务逻辑

**源码位置**:
- Controller: `source/src/main/java/com/dsg/standardization/controller/DeCatalogInfoController.java`
- Service: `source/src/main/java/com/dsg/standardization/service/IDeCatalogInfoService.java`
- ServiceImpl: `source/src/main/java/com/dsg/standardization/service/impl/DeCatalogInfoServiceImpl.java`
- Entity: `source/src/main/java/com/dsg/standardization/entity/DeCatalogInfo.java`
- Mapper: `source/src/main/java/com/dsg/standardization/mapper/DeCatalogInfoMapper.java`
- VO: `source/src/main/java/com/dsg/standardization/vo/CatalogVo/`

---

## Clarifications

### Session 2025-01-21

- Q: 认证与审计日志策略（Java 源码中有 @AuditLog 注解和 token check 配置，Go 版本是否需要实现？） → A: 仅实现核心功能，暂不实现审计和认证
  - **影响**: 转写范围聚焦于核心 CRUD 业务逻辑，排除认证中间件和审计日志功能
- Q: 已废弃接口是否需要实现（Java 中 GET /v1/catalog/{id} 和 GET /v1/catalog/query_tree_by_file 标记 @Deprecated） → A: 不实现废弃接口，简化转写范围
  - **影响**: 移除废弃接口 #3 和 #7，仅实现 6 个核心接口
- Q: 外部服务依赖的处理策略（删除目录时需检查数据元、码表、编码规则、文件等关联数据） → A: 先实现独立目录功能，外部检查通过预留接口/桩模块
  - **影响**: 当前实现目录的独立 CRUD 功能，关联数据检查通过接口预留，后续通过 RPC 或服务间调用实现

---

## User Stories

### Story 1: 目录树查询 (P1)

AS a 前端开发人员
I WANT 通过类型或ID查询目录树结构
SO THAT 可以在界面上展示层级目录选择器

**独立测试**: 调用 GET /v1/catalog/query_tree?type=1 返回完整的目录树 JSON

### Story 2: 目录检索 (P1)

AS a 前端开发人员
I WANT 通过目录名称模糊检索目录
SO THAT 用户可以快速找到目标目录

**独立测试**: 调用 GET /v1/catalog/query?catalog_name=测试&type=1 返回匹配的目录列表

### Story 3: 创建目录 (P1)

AS a 系统管理员
I WANT 创建新的目录节点
SO THAT 可以组织数据元、码表、编码规则等资源

**独立测试**: 调用 POST /v1/catalog 创建目录，返回成功

### Story 4: 修改目录 (P1)

AS a 系统管理员
I WANT 修改目录名称和父目录
SO THAT 可以调整目录结构

**独立测试**: 调用 PUT /v1/catalog/{id} 更新目录，返回成功

### Story 5: 删除目录 (P1)

AS a 系统管理员
I WANT 删除目录及其所有子目录
SO THAT 可以清理不需要的目录结构

**独立测试**: 调用 DELETE /v1/catalog/{id} 删除目录及其子目录

### Story 6: 目录文件树查询 (P2)

AS a 前端开发人员
I WANT 查询包含文件的目录树结构
SO THAT 可以展示目录和文件的混合视图

**独立测试**: 调用 GET /v1/catalog/query/with_file 返回目录和文件列表

---

## Acceptance Criteria (EARS)

### 正常流程

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-01 | 查询目录树（按类型） | WHEN 请求 GET /v1/catalog/query_tree?type=1&id=空 | THE SYSTEM SHALL 返回该类型根节点开始的完整目录树 |
| AC-02 | 查询目录树（按ID） | WHEN 请求 GET /v1/catalog/query_tree?type=1&id=123 | THE SYSTEM SHALL 返回指定节点及其子树 |
| AC-03 | 按名称检索目录 | WHEN 请求 GET /v1/catalog/query?catalog_name=测试&type=1 | THE SYSTEM SHALL 返回名称包含"测试"的所有目录 |
| AC-04 | 创建目录成功 | WHEN 提交有效的目录数据 | THE SYSTEM SHALL 创建目录并继承父目录的type和level+1 |
| AC-05 | 修改目录成功 | WHEN 提交有效的更新数据 | THE SYSTEM SHALL 更新目录信息，可修改父目录 |
| AC-06 | 删除目录成功 | WHEN 请求 DELETE /v1/catalog/{id} | THE SYSTEM SHALL 删除目录及所有子目录 |
| AC-07 | 查询目录文件树 | WHEN 请求 GET /v1/catalog/query/with_file?catalog_name=测试 | THE SYSTEM SHALL 返回目录列表和文件列表 |

### 异常处理

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-10 | 类型参数为空 | WHEN type 参数为空 | THE SYSTEM SHALL 返回错误码 MissingParameter |
| AC-11 | 类型参数无效 | WHEN type 不在 [1,2,3,4] 范围内 | THE SYSTEM SHALL 返回错误码 InvalidParameter |
| AC-12 | 目录ID不存在 | WHEN 查询/更新/删除的ID不存在 | THE SYSTEM SHALL 返回错误码 Empty |
| AC-13 | 创建目录-名称为空 | WHEN catalogName 为空 | THE SYSTEM SHALL 返回错误码 InvalidParameter |
| AC-14 | 创建目录-名称格式错误 | WHEN 名称不匹配正则 `^[\\u4e00-\\u9fa5a-zA-Z0-9][\\u4e00-\\u9fa5a-zA-Z0-9_-]{0,19}$` | THE SYSTEM SHALL 返回错误码 InvalidParameter |
| AC-15 | 创建目录-父目录为空 | WHEN parentId 为空 | THE SYSTEM SHALL 返回错误码 InvalidParameter |
| AC-16 | 创建目录-父目录不存在 | WHEN parentId 指向的目录不存在 | THE SYSTEM SHALL 返回错误码 Empty |
| AC-17 | 创建目录-父目录级别超限 | WHEN 父目录 level >= 255 | THE SYSTEM SHALL 返回错误码 OutOfRange |
| AC-18 | 创建目录-同级名称重复 | WHEN 同一父目录下存在同名目录 | THE SYSTEM SHALL 返回错误码 OperationConflict |
| AC-19 | 修改目录-修改根目录 | WHEN 尝试修改 level <= 1 的根目录 | THE SYSTEM SHALL 返回错误码 InvalidParameter |
| AC-20 | 修改目录-父目录是自身子目录 | WHEN 新父目录是当前目录的子目录 | THE SYSTEM SHALL 返回错误码 InvalidParameter |
| AC-21 | 修改目录-父目录类型不一致 | WHEN 新父目录 type 与当前目录不同 | THE SYSTEM SHALL 返回错误码 InvalidParameter |
| AC-22 | 删除目录-删除根目录 | WHEN 尝试删除 level <= 1 的根目录 | THE SYSTEM SHALL 返回错误码 InvalidParameter |
| AC-23 | 删除目录-目录下存在数据元 | WHEN type=1 且目录或子目录下存在数据元 | THE SYSTEM SHALL 返回错误码 DATA_EXIST（通过预留接口检查，当前为桩模块） |
| AC-24 | 删除目录-目录下存在码表 | WHEN type=2 且目录或子目录下存在码表 | THE SYSTEM SHALL 返回错误码 DATA_EXIST（通过预留接口检查，当前为桩模块） |
| AC-25 | 删除目录-目录下存在编码规则 | WHEN type=3 且目录或子目录下存在编码规则 | THE SYSTEM SHALL 返回错误码 DATA_EXIST（通过预留接口检查，当前为桩模块） |
| AC-26 | 删除目录-目录下存在文件 | WHEN type=4 且目录或子目录下存在文件 | THE SYSTEM SHALL 返回错误码 DATA_EXIST（通过预留接口检查，当前为桩模块） |

---

## Edge Cases

| ID | Case | Expected Behavior |
|----|------|-------------------|
| EC-01 | 目录名称含SQL特殊字符 | 自动转义 SQL 特殊字符后查询 |
| EC-02 | 目录名称超过64字符 | 自动截取前64字符 |
| EC-03 | 并发创建同名目录 | 仅一个成功，其他返回 OperationConflict |
| EC-04 | 目录层级超过255层 | 拒绝创建，返回 OutOfRange |
| EC-05 | 循环父子关系检测 | 检测并拒绝，返回 InvalidParameter |
| EC-06 | 删除时关联数据检查 | 递归检查所有子目录是否存在关联数据（通过预留接口实现，当前为桩模块） |

---

## Business Rules

| ID | Rule | Description |
|----|------|-------------|
| BR-01 | 目录类型枚举 | 1=数据元, 2=码表, 3=编码规则, 4=文件, 0=根目录(保留), 99=其他(保留) |
| BR-02 | 目录名称格式 | 长度1-20字符，中文/英文/数字/符号_-, 且_-不能作为首字符 |
| BR-03 | 目录层级 | level=1 为根目录，level=2 为一级子目录，最大 255 |
| BR-04 | 父目录继承 | 创建时继承父目录的 type，level = 父目录 level + 1 |
| BR-05 | 根目录保护 | 不允许创建、修改、删除 level <= 1 的根目录 |
| BR-06 | 同级名称唯一 | 同一父目录下的目录名称不能重复 |
| BR-07 | 类型一致性 | 目录移动时，新父目录的 type 必须与当前目录一致 |
| BR-08 | 级联删除 | 删除目录时递归删除所有子目录 |
| BR-09 | 关联数据检查 | 删除前必须检查目录及其子目录下是否存在关联数据（通过预留接口实现，当前阶段为桩模块） |

---

## Data Considerations

### 目录实体 (t_de_catalog_info)

| Field | Java Type | Go Type | Description | Constraints |
|-------|-----------|---------|-------------|-------------|
| f_id | Long | string | 目录唯一标识 | 主键，雪花算法生成 |
| f_catalog_name | String | string | 目录名称 | 必填，1-20字符，正则校验 |
| f_description | String | string | 目录说明 | 可选 |
| f_level | Integer | int32 | 目录级别 | 1-255，1为根目录 |
| f_parent_id | Long | string | 父级标识 | 必填，指向父目录ID |
| f_type | CatalogTypeEnum | int32 | 目录类型 | 1=数据元, 2=码表, 3=编码规则, 4=文件 |
| f_authority_id | String | string | 权限域 | 预留字段 |

### 请求/响应结构

**CatalogTreeNodeVo** (目录树节点):
```json
{
  "id": "123456789",
  "catalogName": "目录名称",
  "description": "目录说明",
  "level": 2,
  "parentId": "123456788",
  "type": 1,
  "authorityId": null,
  "count": 10,
  "children": [...]
}
```

**CatalogInfoVo** (目录信息):
```json
{
  "id": "123456789",
  "catalogName": "目录名称",
  "description": "目录说明",
  "level": 2,
  "parentId": "123456788",
  "type": 1,
  "authorityId": null,
  "count": 10
}
```

**CatalogListByFileVo** (目录文件列表):
```json
{
  "catalogs": [...],
  "files": [...]
}
```

### 错误码映射 (Java -> Go)

| Java ErrorCodeEnum | 含义 | Go errorx.Error Code |
|-------------------|------|---------------------|
| MissingParameter | 必填参数丢失 | errorx.ErrMissingParam |
| InvalidParameter | 参数值校验不通过 | errorx.ErrInvalidParam |
| Empty | 数据不存在 | errorx.ErrNotFound |
| OutOfRange | 超出范围 | errorx.ErrOutOfRange |
| OperationConflict | 冲突的操作 | errorx.ErrConflict |
| DATA_EXIST | 数据已存在 | errorx.ErrDataExist |
| CatalogServiceError | 目录模块服务异常 | errorx.ErrInternal (30100) |

---

## API 接口清单

### 1. 查询目录树
```
GET /v1/catalog/query_tree
Query Params:
  - type: int (required) - 目录类型 1=数据元, 2=码表, 3=编码规则, 4=文件
  - id: long (optional) - 目录ID，指定时返回该节点子树
Response: Result<CatalogTreeNodeVo>
```

### 2. 按名称检索目录
```
GET /v1/catalog/query
Query Params:
  - type: int (required) - 目录类型
  - catalog_name: string (optional) - 目录名称（模糊查询）
Response: Result<List<CatalogInfoVo>>
```

### 3. 创建目录
```
POST /v1/catalog
Request Body: DeCatalogInfo
{
  "catalogName": "目录名称",
  "parentId": 123456789
}
Response: Result<?>
```

### 4. 修改目录
```
PUT /v1/catalog/{id}
Path Params:
  - id: long - 目录ID
Request Body: DeCatalogInfo
{
  "catalogName": "新名称",
  "parentId": 123456789
}
Response: Result<?>
```

### 5. 删除目录
```
DELETE /v1/catalog/{id}
Path Params:
  - id: long - 目录ID
Response: Result<?>
```

### 6. 查询目录与文件树列表
```
GET /v1/catalog/query/with_file
Query Params:
  - catalog_name: string (optional) - 目录名称（模糊查询）
Response: Result<CatalogListByFileVo>
```

---

## Success Metrics

| ID | Metric | Target |
|----|--------|--------|
| SM-01 | API 契约测试通过率 | 100% (与 Java 版本响应完全一致) |
| SM-02 | 单元测试覆盖率 | > 80% |
| SM-03 | 接口响应时间 | < 100ms (P95) |
| SM-04 | 业务逻辑还原度 | 100% (所有校验规则与 Java 版本一致) |

---

## Open Questions

无

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-01-21 | - | 初始版本 - Java 目录管理 API 转写规格 |
