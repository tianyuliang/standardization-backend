# 目录管理 (catalog-api) Specification

> **Branch**: `2-catalog-api`
> **Spec Path**: `specs/2-catalog-api/`
> **Created**: 2026-02-06
>
> **Status**: Draft

---

## Overview

目录管理模块用于管理系统中的目录结构（Catalog），支持以下核心功能：
- 目录查询：支持按类型、名称等多维度查询目录树或列表
- 目录创建：支持创建多级目录，继承父目录类型
- 目录编辑：支持修改目录名称和说明
- 目录删除：支持删除目录及其所有子级目录
- 目录移动：支持批量移动数据到指定目录

**重要约束**：从 Java 迁移到 Go，必须100%保持接口兼容。接口路径、请求参数、响应格式、异常信息必须与原Java实现一致。

---

## User Stories

### Story 1: 目录CRUD管理 (P1)

AS a 系统管理员
I WANT 创建、查询、修改、删除目录
SO THAT 能够组织管理数据元、码表、编码规则、标准文件的分类结构

**独立测试**:
- 创建目录后能正确查询到
- 修改目录后信息正确更新
- 删除目录后所有子级也被删除
- 不允许删除根目录

### Story 2: 目录树查询 (P1)

AS a 系统管理员
I WANT 查询目录树结构
SO THAT 能够以层级方式展示目录分类

**独立测试**:
- 按类型查询目录树正确返回层级结构
- 按ID查询指定目录及其子集正确
- 查询结果包含子目录数量统计

### Story 3: 目录检索 (P2)

AS a 系统管理员
I WANT 按关键字检索目录
SO THAT 能够快速找到需要的目录

**独立测试**:
- 支持模糊搜索目录名称
- 返回平铺的目录列表
- 搜索不区分大小写

### Story 4: 目录及文件树查询 (P2)

AS a 系统管理员
I WANT 查询目录及文件的树形结构
SO THAT 能够同时查看目录分类和标准文件

**独立测试**:
- 返回目录树结构
- 每层目录包含关联的标准文件
- 支持按关键字搜索

### Story 5: 目录名称校验 (P1)

AS a 系统管理员
I WANT 创建目录时系统自动校验名称格式
SO THAT 确保目录名称符合规范要求

**独立测试**:
- 名称长度不超过20字符
- 只能包含中英文、数字、下划线、中划线
- 不能以下划线或中划线开头
- 同级目录名称不能重复

---

## Acceptance Criteria (EARS)

### 正常流程

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-01 | 按类型查询目录树成功 | WHEN 用户按type查询目录树 | THE SYSTEM SHALL 返回该类型的完整目录树结构 |
| AC-02 | 按ID查询目录树成功 | WHEN 用户按id查询目录树 | THE SYSTEM SHALL 返回指定目录及其子集的树形结构 |
| AC-03 | 按关键字检索目录成功 | WHEN 用户按keyword检索目录 | THE SYSTEM SHALL 返回匹配的目录平铺列表 |
| AC-04 | 创建目录成功 | WHEN 用户提交有效的目录创建请求 | THE SYSTEM SHALL 保存目录，继承父目录类型和级别+1 |
| AC-05 | 修改目录成功 | WHEN 用户修改目录名称和说明 | THE SYSTEM SHALL 更新目录信息 |
| AC-06 | 删除目录成功 | WHEN 用户删除非根目录且无数据的目录 | THE SYSTEM SHALL 递归删除该目录及所有子级 |
| AC-07 | 查询目录及文件树成功 | WHEN 用户查询目录及文件树 | THE SYSTEM SHALL 返回目录树结构，每层包含关联文件 |

### 异常处理

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-101 | 目录名称为空 | WHEN 创建时catalogName为空 | THE SYSTEM SHALL 返回 InvalidParameter，提示"目录名称不能为空" |
| AC-102 | 目录名称过长 | WHEN catalogName长度超过20字符 | THE SYSTEM SHALL 返回 InvalidParameter，提示"目录名称长度不能超过20个字符" |
| AC-103 | 目录名称格式错误 | WHEN catalogName包含非法字符 | THE SYSTEM SHALL 返回 InvalidParameter，提示"目录名称只能由中英文、数字、下划线、中划线组成" |
| AC-104 | 目录名称以下划线开头 | WHEN catalogName以_开头 | THE SYSTEM SHALL 返回 InvalidParameter，提示"目录名称不能以下划线和中划线开头" |
| AC-105 | 目录名称以中划线开头 | WHEN catalogName以-开头 | THE SYSTEM SHALL 返回 InvalidParameter，提示"目录名称不能以下划线和中划线开头" |
| AC-106 | 父目录不存在 | WHEN parentId指向不存在的目录 | THE SYSTEM SHALL 返回 Empty，提示"无法找到对应的父目录" |
| AC-107 | 目录级别超限 | WHEN 父目录级别>=255 | THE SYSTEM SHALL 返回 OutOfRange，提示"目录级别取值范围(1-255)" |
| AC-108 | 同级目录名称重复 | WHEN 同一父目录下存在同名目录 | THE SYSTEM SHALL 返回 OperationConflict，提示"同级目录名称不能重复" |
| AC-109 | 修改时名称重复 | WHEN 修改后名称与其他同级目录重复 | THE SYSTEM SHALL 返回 OperationConflict，提示"同级目录名称不能重复" |
| AC-110 | 删除根目录 | WHEN 用户删除level=1的目录 | THE SYSTEM SHALL 返回 InvalidParameter，提示"不允许删除根目录" |
| AC-111 | 删除包含数据的目录 | WHEN 目录下存在数据元/码表/规则/文件 | THE SYSTEM SHALL 返回 DATA_EXIST，提示"目录或子目录下已存在数据，不允许删除" |
| AC-112 | 目录类型无效 | WHEN type不在1-4范围内 | THE SYSTEM SHALL 返回 InvalidParameter，提示"此类型不在有效值范围内" |
| AC-113 | 修改根目录 | WHEN 用户修改level=1的目录 | THE SYSTEM SHALL 返回 InvalidParameter，提示"不能修改根目录" |
| AC-114 | 新父目录是子目录 | WHEN 移动时新父目录是原目录的子目录 | THE SYSTEM SHALL 返回 InvalidParameter，提示"新的父目录不能是自身及其子目录" |

---

## Edge Cases

| ID | Case | Expected Behavior |
|----|------|-------------------|
| EC-01 | 查询时id为空 | 返回指定类型的完整目录树 |
| EC-02 | 查询时type为空 | 返回 InvalidParameter 错误 |
| EC-03 | 根目录的特殊处理 | level=1的目录为根目录，不允许修改和删除 |
| EC-04 | 目录级别继承 | 子目录level=父目录level+1 |
| EC-05 | 目录类型继承 | 子目录type继承父目录type |
| EC-06 | 递归删除 | 删除目录时递归删除所有子级目录 |
| EC-07 | 空关键字检索 | keyword为空时返回所有level>1的目录 |

---

## Business Rules

| ID | Rule | Description |
|----|------|-------------|
| BR-01 | 目录名称格式 | 最大20字符，中英文、数字、_、-，不以_-开头 |
| BR-02 | 目录级别范围 | 1-255，子目录=父目录+1 |
| BR-03 | 目录类型枚举 | 1-数据元，2-码表，3-编码规则，4-文件 |
| BR-04 | 同级名称唯一性 | 同一父目录下目录名称不能重复 |
| BR-05 | 根目录保护 | level=1的根目录不允许修改和删除 |
| BR-06 | 删除前校验 | 删除前需检查目录及子目录下是否存在数据 |
| BR-07 | 递归删除 | 删除目录时递归删除所有子级 |
| BR-08 | 继承规则 | 子目录继承父目录的type，level=父目录level+1 |
| BR-09 | 数据校验 | 根据type校验对应数据：数据元/码表/规则/文件 |

---

## Data Considerations

| Field | Description | Constraints |
|-------|-------------|-------------|
| id | 主键 | Long类型，自增 |
| catalog_name | 目录名称 | 必填，最大20字符，同级唯一 |
| description | 目录说明 | 可选，最大300字符 |
| level | 目录级别 | 必填，1-255，继承父目录+1 |
| parent_id | 父目录ID | 必填，0或存在的目录ID |
| type | 目录类型 | 必填，1-4枚举值 |
| authority_id | 权限域 | 预留字段 |
| deleted | 逻辑删除标记 | 0-未删除，非0-已删除 |

---

## Success Metrics

| ID | Metric | Target |
|----|--------|--------|
| SC-01 | 接口响应时间 | < 500ms (P95) |
| SC-02 | 测试覆盖率 | > 80% |
| SC-03 | 接口兼容性 | 100% 与Java实现一致 |
| SC-04 | 异常信息准确性 | 错误码、字段、提示信息与Java实现完全一致 |

---

## Open Questions

无

---

## Assumptions

1. **主键类型**：使用Long类型的自增ID（与Java实现一致）
2. **逻辑删除**：使用deleted字段标记删除状态，0表示未删除
3. **根目录**：level=1的目录为根目录，由系统初始化创建
4. **目录类型**：1-数据元，2-码表，3-编码规则，4-文件
5. **数据校验服务**：依赖dataelement-api、dict-api、rule-api、stdfile-api进行数据校验
6. **目录统计**：查询目录树时可选择是否包含数据统计
7. **SQLx驱动**：使用sqlx作为数据库驱动
8. **错误码**：参考Java实现的错误码体系

---

## Dependencies

- dataelement-api：数据元查询服务（删除校验）
- dict-api：字典查询服务（删除校验）
- rule-api：编码规则查询服务（删除校验）
- stdfile-api：标准文件查询服务（删除校验、文件树查询）

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-02-06 | - | 初始版本 |
