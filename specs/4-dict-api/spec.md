# 码表管理 (dict-api) Specification

> **Branch**: `4-dict-api`
> **Spec Path**: `specs/4-dict-api/`
> **Created**: 2026-02-06
>
> **Status**: Draft

---

## Overview

码表管理模块用于管理系统中的码表（Dictionary），支持以下核心功能：
- 码表创建：支持新增码表及码值明细
- 码表查询：支持按目录、关键字、状态等多维度分页查询
- 码表编辑：支持修改码表信息，自动版本控制
- 码表删除：支持单个或批量删除
- 启用/停用：支持启用/停用操作，停用时必填原因
- 码值管理：支持查询、添加、修改码值明细
- 关联查询：支持查询引用码表的数据元和标准文件
- 文件关联：支持关联标准文件，按文件目录/文件查询

**重要约束**：从 Java 迁移到 Go，必须100%保持接口兼容。接口路径、请求参数、响应格式、异常信息必须与原Java实现一致。

**废弃接口**（不实现）：
- POST /v1/dataelement/dict/import - 导入（废弃）
- POST /v1/dataelement/dict/export - 导出（废弃）

**实际实现接口数量**：16个

---

## User Stories

### Story 1: 码表CRUD管理 (P1)

AS a 系统管理员
I WANT 创建、查询、修改、删除码表
SO THAT 能够管理系统中的数据字典

**独立测试**:
- 创建码表后能正确查询到
- 修改码表后版本号自动递增（有变更时）
- 删除码表后物理删除记录
- 批量删除多个码表成功

### Story 2: 码值管理 (P1)

AS a 系统管理员
I WANT 管理码表的码值明细
SO THAT 能够定义码表的具体值和含义

**独立测试**:
- 创建码表时同时创建码值
- 码值支持分页查询
- 码值支持列表查询（全量）
- 码值code不能重复

### Story 3: 码表状态管理 (P1)

AS a 系统管理员
I WANT 启用或停用码表
SO THAT 控制码表是否可用于数据校验

**独立测试**:
- 停用码表时必须填写停用原因
- 停用原因超过800字符时报错
- 启用码表时清空停用原因

### Story 4: 码表关联管理 (P2)

AS a 系统管理员
I WANT 查询和管理码表与标准文件的关联关系
SO THAT 了解码表的使用范围

**独立测试**:
- 能查询码表关联的标准文件列表
- 能添加/替换码表的关联文件
- 能查询引用码表的数据元列表
- 能按标准文件目录查询码表
- 能按标准文件查询码表

### Story 5: 码表数据校验 (P2)

AS a 系统管理员
I WANT 检查码表数据是否存在重复
SO THAT 避免数据冲突

**独立测试**:
- 中文名称+orgType组合唯一性检查
- 英文名称+orgType组合唯一性检查
- 支持过滤当前记录ID

---

## Acceptance Criteria (EARS)

### 正常流程

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-01 | 创建码表成功 | WHEN 用户提交有效的码表数据 | THE SYSTEM SHALL 生成码表编码（雪花算法），保存码表和码值，返回详情 |
| AC-02 | 查询码表列表成功 | WHEN 用户查询码表列表 | THE SYSTEM SHALL 返回分页列表（含目录、部门、引用状态） |
| AC-03 | 查询码表详情成功 | WHEN 用户按ID查询码表 | THE SYSTEM SHALL 返回码表详情（含码值列表、目录、部门信息） |
| AC-04 | 按Code查询详情成功 | WHEN 用户按Code查询码表 | THE SYSTEM SHALL 返回码表详情（含码值列表） |
| AC-05 | 修改码表成功 | WHEN 用户修改码表字段 | THE SYSTEM SHALL 检查变更，有变更则更新且版本号+1 |
| AC-06 | 修改无变更 | WHEN 用户修改但所有字段无变化 | THE SYSTEM SHALL 直接返回原数据，不更新版本 |
| AC-07 | 删除码表成功 | WHEN 用户删除单个码表 | THE SYSTEM SHALL 物理删除码表、码值、关联记录 |
| AC-08 | 批量删除成功 | WHEN 用户批量删除码表 | THE SYSTEM SHALL 物理删除所有相关记录 |
| AC-09 | 启用码表成功 | WHEN 用户启用已停用的码表 | THE SYSTEM SHALL 更新状态为启用，清空停用原因 |
| AC-10 | 停用码表成功 | WHEN 用户停用码表并提供停用原因 | THE SYSTEM SHALL 更新状态为停用，保存原因 |
| AC-11 | 码值分页查询成功 | WHEN 用户分页查询码值 | THE SYSTEM SHALL 返回码值分页列表 |
| AC-12 | 码值列表查询成功 | WHEN 用户查询码值列表 | THE SYSTEM SHALL 返回所有码值（不分页） |
| AC-13 | 查询引用数据元成功 | WHEN 用户查询引用码表的数据元 | THE SYSTEM SHALL 返回数据元分页列表 |
| AC-14 | 按文件目录查询成功 | WHEN 用户按文件目录查询码表 | THE SYSTEM SHALL 返回关联该目录文件的码表列表 |
| AC-15 | 按文件查询成功 | WHEN 用户按标准文件查询码表 | THE SYSTEM SHALL 返回关联该文件的码表列表 |
| AC-16 | 查询关联文件成功 | WHEN 用户查询码表关联的标准文件 | THE SYSTEM SHALL 返回标准文件分页列表 |
| AC-17 | 添加关联关系成功 | WHEN 用户添加/替换关联文件 | THE SYSTEM SHALL 更新关联关系，变更时版本号+1 |
| AC-18 | 数据存在检查成功 | WHEN 用户检查数据是否存在 | THE SYSTEM SHALL 返回true/false |

### 异常处理

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-101 | 目录不存在 | WHEN 指定的catalog_id不存在或已删除 | THE SYSTEM SHALL 返回 InvalidParameter，提示"目录不存在或已删除" |
| AC-102 | 中文名称重复 | WHEN 创建时chName+orgType已存在 | THE SYSTEM SHALL 返回 InvalidParameter，提示"码表中文名称、标准分类不能全部重复" |
| AC-103 | 英文名称重复 | WHEN 创建时enName+orgType已存在 | THE SYSTEM SHALL 返回 InvalidParameter，提示"码表英文名称、标准分类不能全部重复" |
| AC-104 | 修改中文名称重复 | WHEN 修改时chName+orgType与其他记录重复 | THE SYSTEM SHALL 返回 InvalidParameter，提示"码表中文名称、标准分类不能全部重复" |
| AC-105 | 修改英文名称重复 | WHEN 修改时enName+orgType与其他记录重复 | THE SYSTEM SHALL 返回 InvalidParameter，提示"码表英文名称、标准分类不能全部重复" |
| AC-106 | 码值为空 | WHEN 创建时enums[].code为空 | THE SYSTEM SHALL 返回 InvalidParameter，提示"码值输入不能为空" |
| AC-107 | 码值描述为空 | WHEN 创建时enums[].value为空 | THE SYSTEM SHALL 返回 InvalidParameter，提示"码值描述输入不能为空" |
| AC-108 | 码值重复 | WHEN 同一码表内code重复 | THE SYSTEM SHALL 返回 InvalidParameter，提示"码值出现重复记录" |
| AC-109 | 修改记录不存在 | WHEN 修改时码表ID不存在 | THE SYSTEM SHALL 返回 DATA_NOT_EXIST，提示"数据不存在" |
| AC-110 | 停用原因缺失 | WHEN state=disable时reason为空 | THE SYSTEM SHALL 返回 PARAMETER_EMPTY，提示"停用原因不能为空" |
| AC-111 | 停用原因过长 | WHEN reason长度超过800字符 | THE SYSTEM SHALL 返回 InvalidParameter，提示"停用原因长度不能超过800字符" |
| AC-112 | 码表不存在 | WHEN 查询/关联操作时码表ID不存在 | THE SYSTEM SHALL 返回 DATA_NOT_EXIST，提示"数据不存在" |

---

## Edge Cases

| ID | Case | Expected Behavior |
|----|------|-------------------|
| EC-01 | 标准文件目录ID=-1 | 返回未关联任何标准文件的码表 |
| EC-02 | 标准文件顶级目录 | 返回所有码表 |
| EC-03 | 修改时排除自身 | queryDataExists接口支持filter_id参数排除当前记录 |
| EC-04 | 码表已删除但被查询 | 返回null（code=0），不抛异常 |
| EC-05 | 目录ID为44 | 默认值为"全部目录" |
| EC-06 | 码值列表查询 | 不分页，返回所有码值 |

---

## Business Rules

| ID | Rule | Description |
|----|------|-------------|
| BR-01 | 名称唯一性 | 同一orgType下中文名称唯一、英文名称唯一 |
| BR-02 | 版本控制 | 修改以下字段时版本号+1：chName、enName、catalogId、departmentIds、orgType、description、enums、stdFiles |
| BR-03 | 停用原因 | 停用时必须填写原因且<=800字符，启用时清空原因 |
| BR-04 | 码表编码 | 使用雪花算法生成，全局唯一 |
| BR-05 | 码值唯一性 | 同一码表内code必须唯一 |
| BR-06 | 物理删除 | 删除码表时物理删除t_dict、t_dict_enum、t_relation_dict_file记录 |
| BR-07 | 目录处理 | 处理目录ID时获取当前目录及所有子目录ID列表 |
| BR-08 | 引用状态 | 查询列表时判断是否被数据元引用 |
| BR-09 | 关联文件变更 | 添加/替换关联文件时更新码表版本号 |
| BR-10 | 部门ID处理 | department_ids存储完整路径，查询/返回时使用最后一段 |

---

## Data Considerations

| Field | Description | Constraints |
|-------|-------------|-------------|
| id | 主键 | Long类型，自增 |
| code | 码表编码 | Long类型，雪花算法生成，全局唯一 |
| ch_name | 中文名称 | 必填，最大128字符 |
| en_name | 英文名称 | 必填，最大128字符 |
| description | 业务含义 | 可选，最大300字符 |
| catalog_id | 所属目录ID | 必填，默认44 |
| org_type | 标准组织类型 | 必填，0-99枚举值 |
| version | 版本号 | 从1开始，修改特定字段时+1 |
| state | 状态 | ENABLE(1)或DISABLE(0) |
| disable_reason | 停用原因 | 停用时必填，最大800字符 |
| department_ids | 部门ID | 可选，字符串格式 |
| enums | 码值列表 | 创建时至少1条 |
| std_files | 关联标准文件ID数组 | 可选 |

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
2. **码表编码**：使用雪花算法生成Long类型code
3. **分页参数**：offset从1开始（不是0）
4. **目录服务**：依赖catalog-api模块的目录服务
5. **数据元服务**：依赖dataelement-api模块的数据元服务
6. **标准文件服务**：依赖stdfile-api模块的文件服务
7. **部门信息**：从Token中解析部门信息
8. **用户信息**：从Token中解析用户信息
9. **废弃接口**：import/export接口不在本次实现范围

---

## Dependencies

- catalog-api：目录查询服务
- dataelement-api：数据元查询服务
- stdfile-api：标准文件查询服务

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-02-06 | - | 初始版本 |
