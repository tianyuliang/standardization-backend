# 编码规则管理 (rule-api) Specification

> **Branch**: `1-rule-api`
> **Spec Path**: `specs/1-rule-api/`
> **Created**: 2026-02-05
>
> **Status**: Draft

---

## Overview

编码规则管理模块用于管理系统中的编码规则，支持以下核心功能：
- 规则创建：支持正则表达式（REGEX）和自定义配置（CUSTOM）两种类型
- 规则查询：支持按目录、关键字、状态、类型等多维度查询
- 规则修改：支持版本控制，修改时自动递增版本号
- 规则停用/启用：支持停用规则，停用时必须填写停用原因
- 规则删除：支持批量删除，删除后自动推送MQ消息
- 目录移动：支持将规则移动到指定目录
- 关联查询：支持查询规则被哪些数据元引用
- 关联文件：支持关联标准文件

**重要约束**：从 Java 迁移到 Go，必须100%保持接口兼容。接口路径、请求参数、响应格式、异常信息必须与原Java实现一致。

---

## User Stories

### Story 1: 编码规则CRUD管理 (P1)

AS a 系统管理员
I WANT 创建、查询、修改、删除编码规则
SO THAT 能够管理系统中的数据编码规范

**独立测试**:
- 创建规则后能正确查询到
- 修改规则后版本号自动递增
- 删除规则后不再出现在列表中
- 批量删除多个规则成功

### Story 2: 规则状态管理 (P1)

AS a 系统管理员
I WANT 启用或停用编码规则
SO THAT 控制规则是否可用于数据校验

**独立测试**:
- 停用规则时必须填写停用原因
- 停用原因超过800字符时报错
- 启用规则时清空停用原因
- 停用的规则不影响已引用它的数据元

### Story 3: 规则目录移动 (P2)

AS a 系统管理员
I WANT 批量移动规则到指定目录
SO THAT 能够重新组织规则的分类结构

**独立测试**:
- 移动后规则的catalog_id更新
- 移动后规则版本号自动递增
- 目标目录不存在时报错

### Story 4: 规则关联查询 (P2)

AS a 系统管理员
I WANT 查询规则被哪些数据元引用
SO THAT 了解规则的影响范围

**独立测试**:
- 能分页查询引用规则的数据元列表
- 能查询规则关联的标准文件列表
- 能根据标准文件目录查询规则

### Story 5: 内部查询接口 (P3)

AS a 内部服务
I WANT 根据数据元ID或编码查询关联的编码规则
SO THAT 能够获取数据元的编码规则信息

**独立测试**:
- 根据数据元ID查询到正确的规则
- 根据数据元编码查询到正确的规则
- 数据元不存在或无关联规则时返回null

---

## Acceptance Criteria (EARS)

### 正常流程

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-01 | 创建REGEX规则成功 | WHEN 用户提交有效的正则表达式规则 | THE SYSTEM SHALL 保存规则，返回规则详情 |
| AC-02 | 创建CUSTOM规则成功 | WHEN 用户提交有效的自定义配置规则 | THE SYSTEM SHALL 保存规则，返回规则详情 |
| AC-03 | 查询规则详情成功 | WHEN 用户查询存在的规则ID | THE SYSTEM SHALL 返回规则详情（含关联文件、部门信息） |
| AC-04 | 列表查询成功 | WHEN 用户查询规则列表 | THE SYSTEM SHALL 返回分页列表（含目录、部门、引用状态） |
| AC-05 | 修改规则成功 | WHEN 用户修改规则字段 | THE SYSTEM SHALL 更新规则，版本号+1 |
| AC-06 | 修改无变更 | WHEN 用户修改但所有字段无变化 | THE SYSTEM SHALL 直接返回原数据，不更新版本 |
| AC-07 | 批量删除成功 | WHEN 用户删除多个规则ID | THE SYSTEM SHALL 物理删除规则及关联记录，发送MQ消息 |
| AC-08 | 启用规则成功 | WHEN 用户启用已停用的规则 | THE SYSTEM SHALL 更新状态为启用，清空停用原因 |
| AC-09 | 停用规则成功 | WHEN 用户停用规则并提供停用原因 | THE SYSTEM SHALL 更新状态为停用，保存原因 |
| AC-10 | 目录移动成功 | WHEN 用户批量移动规则到目标目录 | THE SYSTEM SHALL 更新catalog_id，版本号+1 |
| AC-11 | 批量查询成功 | WHEN 用户按ID列表查询规则 | THE SYSTEM SHALL 返回匹配的规则列表 |
| AC-12 | 按标准文件查询成功 | WHEN 用户按标准文件查询规则 | THE SYSTEM SHALL 返回关联该文件的规则列表 |
| AC-13 | 按文件目录查询成功 | WHEN 用户按文件目录查询规则 | THE SYSTEM SHALL 返回关联该目录下文件的规则 |
| AC-14 | 查询引用数据元成功 | WHEN 用户查询规则被引用情况 | THE SYSTEM SHALL 返回引用该规则的数据元分页列表 |
| AC-15 | 查询关联文件成功 | WHEN 用户查询规则关联的标准文件 | THE SYSTEM SHALL 返回标准文件分页列表 |
| AC-16 | 获取日期格式成功 | WHEN 用户查询自定义规则日期格式 | THE SYSTEM SHALL 返回支持的日期格式集合 |
| AC-17 | 检查数据存在成功 | WHEN 用户检查规则名称是否存在 | THE SYSTEM SHALL 返回true/false |
| AC-18 | 内部查询规则成功 | WHEN 内部接口按数据元ID查询 | THE SYSTEM SHALL 返回关联的规则详情 |

### 异常处理

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-101 | 名称重复 | WHEN 创建时规则名称在同一orgType下已存在 | THE SYSTEM SHALL 返回 InvalidParameter，提示"规则名称已存在" |
| AC-102 | 目录不存在 | WHEN 指定的catalog_id不存在或类型不符 | THE SYSTEM SHALL 返回 InvalidParameter，提示"目录id[{id}]对应的目录不存在" |
| AC-103 | 正则为空 | WHEN ruleType=REGEX时regex为空 | THE SYSTEM SHALL 返回 InvalidParameter，提示"正则表达式为空" |
| AC-104 | 正则非法 | WHEN 正则表达式格式错误 | THE SYSTEM SHALL 返回 InvalidParameter，提示"正则表达式非法" |
| AC-105 | 自定义配置为空 | WHEN ruleType=CUSTOM时custom为空 | THE SYSTEM SHALL 返回 InvalidParameter，提示"不能为空" |
| AC-106 | 分段长度无效 | WHEN custom[].segment_length <= 0 | THE SYSTEM SHALL 返回 InvalidParameter，提示"值必须为正整数" |
| AC-107 | 码表不存在 | WHEN custom类型为DICT但码表ID不存在 | THE SYSTEM SHALL 返回 InvalidParameter，提示"码表不存在" |
| AC-108 | 日期格式不支持 | WHEN custom类型为DATE但格式不支持 | THE SYSTEM SHALL 返回 InvalidParameter，提示"不支持的日期格式" |
| AC-109 | 修改记录不存在 | WHEN 修改时规则ID不存在 | THE SYSTEM SHALL 返回 DATA_NOT_EXIST，提示"数据不存在" |
| AC-110 | 修改名称重复 | WHEN 修改后名称与其他记录重复 | THE SYSTEM SHALL 返回 InvalidParameter，提示"规则名称已存在" |
| AC-111 | 删除IDs为空 | WHEN 批量删除时ids参数为空 | THE SYSTEM SHALL 返回 InvalidParameter，提示"ids 不能为空" |
| AC-112 | 停用原因缺失 | WHEN state=disable时reason为空 | THE SYSTEM SHALL 返回 PARAMETER_EMPTY，提示"停用必须填写停用原因" |
| AC-113 | 停用原因过长 | WHEN reason长度超过800字符 | THE SYSTEM SHALL 返回 InvalidParameter，提示"长度超过800" |
| AC-114 | 查询记录不存在 | WHEN 查询引用数据元时规则ID不存在 | THE SYSTEM SHALL 返回 DATA_NOT_EXIST，提示"记录不存在" |
| AC-115 | 批量查询IDs为空 | WHEN 批量查询时ids为空 | THE SYSTEM SHALL 返回 InvalidParameter，提示"ids 不能为空" |

---

## Edge Cases

| ID | Case | Expected Behavior |
|----|------|-------------------|
| EC-01 | 标准文件目录ID=-1 | 返回未关联任何标准文件的规则 |
| EC-02 | 标准文件顶级目录 | 返回所有编码规则 |
| EC-03 | 自定义规则value为空 | DATE/SPLIT_STR类型时报错"不能为空" |
| EC-04 | 修改时排除自身 | queryDataExists接口支持filter_id参数排除当前记录 |
| EC-05 | 规则已删除但被查询 | 返回null（code=0），不抛异常 |
| EC-06 | 内部接口数据元无关联规则 | 返回null，不抛异常 |
| EC-07 | 目录ID为33 | 默认值为"全部目录" |

---

## Business Rules

| ID | Rule | Description |
|----|------|-------------|
| BR-01 | 名称唯一性 | 同一orgType下规则名称必须唯一（含部门ID） |
| BR-02 | 版本控制 | 修改以下字段时版本号+1：name、catalog_id、department_ids、org_type、description、rule_type、expression、关联文件 |
| BR-03 | 停用原因 | 停用时必须填写原因且<=800字符，启用时清空原因 |
| BR-04 | 表达式存储 | REGEX类型存储regex字符串，CUSTOM类型存储custom的JSON |
| BR-05 | 关联文件限制 | 最多关联10个标准文件 |
| BR-06 | MQ消息 | 创建、修改、删除规则时发送MQ消息到MQ_MESSAGE_SAILOR |
| BR-07 | 物理删除 | 删除规则时物理删除t_rule和t_relation_rule_file记录 |
| BR-08 | 默认catalog_id | 不传时默认为33（全部目录） |
| BR-09 | 部门ID处理 | department_ids存储完整路径，查询/返回时使用最后一段 |
| BR-10 | XSS转义 | keyword参数需进行XSS转义 |

---

## Data Considerations

| Field | Description | Constraints |
|-------|-------------|-------------|
| id | 主键 | Long类型 |
| name | 规则名称 | 必填，最大128字符，同orgType下唯一 |
| catalog_id | 所属目录ID | 必填，默认33，关联目录表 |
| org_type | 标准组织类型 | 必填，0-99枚举值 |
| rule_type | 规则类型 | 必填，REGEX(0)或CUSTOM(1) |
| expression | 表达式 | 必填，REGEX时为正则，CUSTOM时为JSON |
| version | 版本号 | 从1开始，修改特定字段时+1 |
| state | 状态 | ENABLE(1)或DISABLE(0) |
| disable_reason | 停用原因 | 停用时必填，最大800字符 |
| description | 规则说明 | 可选，最大300字符 |
| std_files | 关联标准文件ID数组 | 可选，最多10个 |
| department_ids | 部门ID | 可选，字符串格式 |
| custom | 自定义配置 | ruleType=CUSTOM时必填，最多1024条 |

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
2. **日期格式**：支持 yyyy、yyyyMM、yyyy-MM-dd、yyyyMMdd、yyyy-MM-dd HH:mm:ss、yyyyMMddHHmmss、yyyy-MM-dd'T'HH:mm:ss、yyyyMMdd'T'HHmmss、HHmmss、HH:mm:ss
3. **MQ消息格式**：与Java实现保持一致的JSON格式
4. **分页参数**：offset从1开始（不是0）
5. **目录服务**：依赖catalog-api模块的目录服务
6. **数据元服务**：依赖dataelement-api模块的数据元服务
7. **字典服务**：依赖dict-api模块的字典服务（CUSTOM类型DICT校验）
8. **标准文件服务**：依赖stdfile-api模块的文件服务
9. **部门信息**：从Token中解析部门信息
10. **用户信息**：从Token中解析用户信息

---

## Dependencies

- catalog-api：目录查询服务
- dataelement-api：数据元查询服务
- dict-api：字典查询服务（CUSTOM类型DICT校验）
- stdfile-api：标准文件查询服务
- MQ服务：Kafka消息队列

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-02-05 | - | 初始版本 |
