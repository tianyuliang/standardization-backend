# 数据元管理 (dataelement-api) Specification

> **Branch**: `6-std-dataelement-api`
> **Spec Path**: `specs/6-std-dataelement-api/`
> **Created**: 2026-02-09
>
> **Status**: Draft

---

## Overview

数据元管理模块用于管理系统中的数据元（Data Element），支持以下核心功能：
- **数据元创建**：支持新增单个数据元及批量导入（Excel）
- **数据元查询**：支持按目录、关键字、状态、类型等多维度分页查询
- **数据元编辑**：支持修改数据元信息，含版本控制
- **数据元删除**：支持单个或批量删除
- **启用/停用**：支持批量启用/停用操作
- **目录移动**：支持批量移动数据元到指定目录
- **关联管理**：支持关联码表、编码规则、标准文件
- **导入导出**：支持Excel导入导出
- **重复检查**：支持中英文名称重复性校验

**重要约束**：从 Java 迁移到 Go，必须100%保持接口兼容。接口路径、请求参数、响应格式、异常信息必须与原Java实现一致。

---

## User Stories

### Story 1: 数据元CRUD管理 (P1)

AS a 系统管理员
I WANT 创建、查询、修改、删除数据元
SO THAT 能够管理系统中的数据元定义

**独立测试**:
- 创建数据元后能正确查询到
- 修改数据元时关键属性变更会递增版本号
- 删除数据元后不再出现在列表中
- 批量删除多个数据元成功

### Story 2: 数据元状态管理 (P1)

AS a 系统管理员
I WANT 启用或停用数据元
SO THAT 控制数据元是否可用于业务

**独立测试**:
- 停用数据元时必须填写停用原因
- 停用原因超过800字符时报错
- 启用数据元时清空停用原因

### Story 3: 数据元导入导出 (P2)

AS a 系统管理员
I WANT 批量导入和导出数据元
SO THAT 能够快速迁移大量数据元数据

**独立测试**:
- 支持Excel文件导入（最多5000条）
- 支持按条件导出数据元
- 支持按ID列表导出数据元

### Story 4: 数据元目录管理 (P2)

AS a 系统管理员
I WANT 批量移动数据元到指定目录
SO THAT 能够重新组织数据元的分类结构

**独立测试**:
- 移动后数据元的catalog_id更新
- 移动后数据元版本号自动递增
- 目标目录不存在时报错

### Story 5: 数据元关联查询 (P2)

AS a 系统管理员
I WANT 查询数据元关联的文件及重复性检查
SO THAT 了解数据元的关联情况

**独立测试**:
- 能分页查询数据元关联的标准文件列表
- 能按标准文件/目录查询数据元
- 能检查中英文名称是否重复

---

## Acceptance Criteria (EARS)

### 正常流程

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-01 | 创建数据元成功 | WHEN 用户提交有效的数据元创建请求 | THE SYSTEM SHALL 保存数据元，返回数据元详情 |
| AC-02 | 查询数据元详情成功 | WHEN 用户查询存在的数据元ID | THE SYSTEM SHALL 返回数据元详情（含关联码表、规则、文件、部门信息） |
| AC-03 | 列表查询成功 | WHEN 用户查询数据元列表 | THE SYSTEM SHALL 返回分页列表（含目录、部门、关联信息） |
| AC-04 | 修改数据元成功 | WHEN 用户修改数据元字段 | THE SYSTEM SHALL 更新数据元，关键属性变更时版本号+1 |
| AC-05 | 批量删除成功 | WHEN 用户删除多个数据元ID | THE SYSTEM SHALL 物理删除数据元及关联记录，发送MQ消息 |
| AC-06 | 启用数据元成功 | WHEN 用户启用已停用的数据元 | THE SYSTEM SHALL 更新状态为启用，清空停用原因 |
| AC-07 | 停用数据元成功 | WHEN 用户停用数据元并提供停用原因 | THE SYSTEM SHALL 更新状态为停用，保存原因 |
| AC-08 | 目录移动成功 | WHEN 用户批量移动数据元到目标目录 | THE SYSTEM SHALL 更新catalog_id，版本号+1 |
| AC-09 | 批量导入成功 | WHEN 用户上传有效的Excel文件 | THE SYSTEM SHALL 解析并保存数据元，返回导入结果 |
| AC-10 | 导出数据元成功 | WHEN 用户按条件导出数据元 | THE SYSTEM SHALL 生成Excel文件并返回下载链接 |
| AC-11 | 按文件查询成功 | WHEN 用户按标准文件查询数据元 | THE SYSTEM SHALL 返回关联该文件的数据元列表 |
| AC-12 | 按文件目录查询成功 | WHEN 用户按文件目录查询数据元 | THE SYSTEM SHALL 返回关联该目录下文件的数据元列表 |
| AC-13 | 查询关联文件成功 | WHEN 用户查询数据元关联的标准文件 | THE SYSTEM SHALL 返回标准文件分页列表 |
| AC-14 | 检查名称重复成功 | WHEN 用户检查中英文名称是否重复 | THE SYSTEM SHALL 返回true/false |
| AC-15 | 删除标签成功 | WHEN 用户删除数据元的标签 | THE SYSTEM SHALL 清空label_id字段 |
| AC-16 | 查询列表成功 | WHEN 用户按ID或Code列表查询数据元 | THE SYSTEM SHALL 返回匹配的数据元列表 |
| AC-17 | 内部查询成功 | WHEN 内部接口查询数据元列表 | THE SYSTEM SHALL 返回数据元分页列表 |
| AC-18 | 内部详情查询成功 | WHEN 内部接口查询数据元详情 | THE SYSTEM SHALL 返回数据元详情 |

### 异常处理

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-101 | 目录不存在 | WHEN 指定的catalog_id不存在或类型不符 | THE SYSTEM SHALL 返回 InvalidParameter，提示"数据元对应的目录不存在" |
| AC-102 | 目录类型错误 | WHEN catalog_id对应的目录不是数据元目录 | THE SYSTEM SHALL 返回 InvalidParameter，提示"数据元对应的目录类型不正确" |
| AC-103 | 关联类型为空 | WHEN relationType为空 | THE SYSTEM SHALL 返回 InvalidParameter，提示"数据元关联类型不能为空" |
| AC-104 | 码表关联但dictId为空 | WHEN relationType=codeTable但dictId为空 | THE SYSTEM SHALL 返回 InvalidParameter，提示"数据元关联码表不能为空" |
| AC-105 | 规则关联但ruleId为空 | WHEN relationType=codeRule但ruleId为空 | THE SYSTEM SHALL 返回 InvalidParameter，提示"数据元关联编码规则不能为空" |
| AC-106 | 编码规则不存在 | WHEN ruleId指向的规则不存在或已删除 | THE SYSTEM SHALL 返回 InvalidParameter，提示"编码规则不存在或已删除" |
| AC-107 | 码表不存在 | WHEN dictId指向的码表不存在或已删除 | THE SYSTEM SHALL 返回 InvalidParameter，提示"码表数据不存在或已删除" |
| AC-108 | 中文名称重复 | WHEN nameCn+stdType组合重复 | THE SYSTEM SHALL 返回 InvalidParameter，提示"中文名称+标准分类不能全部重复" |
| AC-109 | 标准分类无效 | WHEN stdType枚举值无效 | THE SYSTEM SHALL 返回 InvalidParameter，提示"[标准分类]:输入错误" |
| AC-110 | 数据类型无效 | WHEN dataType枚举值无效 | THE SYSTEM SHALL 返回 InvalidParameter，提示"[数据类型]:输入错误" |
| AC-111 | 英文名称为空 | WHEN nameEn为空 | THE SYSTEM SHALL 返回 InvalidParameter，提示"[英文名称]:空" |
| AC-112 | 英文名称格式错误 | WHEN nameEn格式不符合要求或长度超过128 | THE SYSTEM SHALL 返回 InvalidParameter，提示"[英文名称]:字符不符合要求或长度超过128" |
| AC-113 | 中文名称为空 | WHEN nameCn为空 | THE SYSTEM SHALL 返回 InvalidParameter，提示"[中文名称]:空" |
| AC-114 | 中文名称过长 | WHEN nameCn长度超过128 | THE SYSTEM SHALL 返回 InvalidParameter，提示"[中文名称]:长度超过128" |
| AC-115 | 数据长度无效 | WHEN Number/Decimal类型长度或精度范围错误 | THE SYSTEM SHALL 返回 InvalidParameter，提示"[数据长度]:输入错误" |
| AC-116 | 字符长度超限 | WHEN Char类型长度超过65535 | THE SYSTEM SHALL 返回 InvalidParameter，提示"[数据长度]:输入错误" |
| AC-117 | 文件ID不存在 | WHEN stdFiles包含不存在的文件ID | THE SYSTEM SHALL 返回 InvalidParameter，提示"一个或多个文件id不存在" |
| AC-118 | 停用原因缺失 | WHEN state=disable时reason为空 | THE SYSTEM SHALL 返回 PARAMETER_EMPTY，提示"停用必须填写停用原因" |
| AC-119 | 停用原因过长 | WHEN reason长度超过800字符 | THE SYSTEM SHALL 返回 InvalidParameter，提示"长度超过800" |
| AC-120 | 修改记录不存在 | WHEN 修改时数据元ID不存在 | THE SYSTEM SHALL 返回 DATA_NOT_EXIST，提示"数据不存在" |

---

## Edge Cases

| ID | Case | Expected Behavior |
|----|------|-------------------|
| EC-01 | 标准文件目录ID=-1 | 返回未关联任何标准文件的数据元 |
| EC-02 | 查询时排除自身 | isRepeat接口支持id参数排除当前记录 |
| EC-03 | 导入文件格式错误 | 返回模板文件和错误提示 |
| EC-04 | 导入超过5000条 | 返回错误，提示"单次导入限制5000条" |
| EC-05 | 值域计算 | 关联码表时返回"[码值1,码值2,...]"，否则按数据类型计算 |
| EC-06 | 版本不递增场景 | 修改目录ID、数据类型、长度/精度、说明时版本号不变 |

---

## Business Rules

| ID | Rule | Description |
|----|------|-------------|
| BR-01 | 名称唯一性 | 同一stdType下nameCn必须唯一，nameEn在同一部门下唯一 |
| BR-02 | 版本控制 | 修改以下字段时版本号+1：relationType、dictCode/ruleId、nameCn/nameEn |
| BR-03 | 停用原因 | 停用时必须填写原因且<=800字符，启用时清空原因 |
| BR-04 | 关联类型限制 | relationType为codeTable时必须有dictId，为codeRule时必须有ruleId |
| BR-05 | 数据长度校验 | Number/Decimal类型：长度1-38，精度0-38且精度<长度；Char类型：长度0-65535 |
| BR-06 | MQ消息 | 创建、修改、删除数据元时发送MQ消息到MQ_MESSAGE_SAILOR |
| BR-07 | 物理删除 | 删除数据元时物理删除t_data_element_info和t_relation_de_file记录 |
| BR-08 | 值域计算 | 关联码表时值域为码表值JSON数组，否则按数据类型公式计算 |
| BR-09 | 导入限制 | 单次导入最多5000条，文件大小限制10M |
| BR-10 | 目录ID处理 | 查询时catalog_id获取当前目录及所有子目录的数据元 |

---

## Data Considerations

| Field | Description | Constraints |
|-------|-------------|-------------|
| id | 主键 | BIGINT(20)，自增 |
| code | 关联标识 | BIGINT(20)，等于id |
| name_en | 英文名称 | 必填，最大128字符，英文_-组成 |
| name_cn | 中文名称 | 必填，最大128字符，同stdType下唯一 |
| synonym | 同义词 | 可选，最大300字符，逗号分隔 |
| std_type | 标准分类 | 必填，0-99枚举值 |
| data_type | 数据类型 | 必填，0-10枚举值 |
| data_length | 数据长度 | 可选，数字型1-38，字符型0-65535 |
| data_precision | 数据精度 | 可选，0-38且<长度 |
| dict_code | 关联码表编码 | 可选，关联码表时必填 |
| rule_id | 关联编码规则ID | 可选，关联规则时必填 |
| relation_type | 关联类型 | 必填，no/codeTable/codeRule |
| catalog_id | 目录ID | 必填，关联数据元目录 |
| label_id | 数据分级标签ID | 可选 |
| description | 说明 | 可选，最大300字符 |
| version | 版本号 | 从1开始，修改特定字段时+1 |
| state | 状态 | ENABLE(1)或DISABLE(0) |
| department_ids | 部门ID | 可选，存储完整路径 |
| third_dept_id | 第三方部门ID | 可选，CHAR(36) |
| std_files | 关联文件ID数组 | 可选，最多无限制 |

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

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-02-09 | - | 初始版本 |
