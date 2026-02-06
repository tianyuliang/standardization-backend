# 标准文件管理 (std-file-api) Specification

> **Branch**: `3-std-file-api`
> **Spec Path**: `specs/3-std-file-api/`
> **Created**: 2026-02-06
>
> **Status**: Draft

---

## Overview

标准文件管理模块用于管理系统中的标准文件（Standard File），支持以下核心功能：
- **文件管理**：支持标准文件的上传、下载、查询
- **文件查询**：支持按目录、关键字、状态等多维度分页查询
- **文件编辑**：支持修改文件信息
- **文件删除**：支持批量删除
- **启用/停用**：支持启用/停用操作
- **关联管理**：支持关联数据元、码表、编码规则

**重要约束**：从 Java 迁移到 Go，必须100%保持接口兼容。接口路径、请求参数、响应格式、异常信息必须与原Java实现一致。

---

## User Stories

### Story 1: 标准文件CRUD管理 (P1)

AS a 系统管理员
I WANT 创建、查询、修改、删除标准文件
SO THAT 能够管理系统中的标准文件规范

**独立测试**:
- 创建文件后能正确查询到
- 修改文件后版本号自动递增
- 删除文件后不再出现在列表中
- 批量删除多个文件成功

### Story 2: 文件状态管理 (P1)

AS a 系统管理员
I WANT 启用或停用标准文件
SO THAT 控制文件是否可用于系统

**独立测试**:
- 停用文件时必须填写停用原因
- 停用原因超过800字符时报错
- 启用文件时清空停用原因

### Story 3: 文件下载 (P1)

AS a 系统用户
I WANT 下载标准文件附件
SO THAT 能够查看和使用标准文件内容

**独立测试**:
- FILE类型：从OSS下载文件
- URL类型：重定向到链接地址
- 批量下载：打包成ZIP

### Story 4: 文件目录移动 (P2)

AS a 系统管理员
I WANT 批量移动文件到指定目录
SO THAT 能够重新组织文件的分类结构

**独立测试**:
- 移动后文件的catalog_id更新
- 目标目录不存在时报错

### Story 5: 关联关系管理 (P2)

AS a 系统管理员
I WANT 管理文件与数据元、码表、编码规则的关联关系
SO THAT 能够追踪标准的引用关系

**独立测试**:
- 能添加关联关系
- 能查询关联的数据元、码表、编码规则
- 能批量启用/停用文件

---

## Acceptance Criteria (EARS)

### 正常流程

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-01 | 创建文件成功(FILE类型) | WHEN 用户提交有效的文件附件 | THE SYSTEM SHALL 保存文件，返回文件详情 |
| AC-02 | 创建文件成功(URL类型) | WHEN 用户提交有效的外置链接 | THE SYSTEM SHALL 保存文件，返回文件详情 |
| AC-03 | 查询文件详情成功 | WHEN 用户查询存在的文件ID | THE SYSTEM SHALL 返回文件详情（含目录、部门信息） |
| AC-04 | 列表查询成功 | WHEN 用户查询文件列表 | THE SYSTEM SHALL 返回分页列表（含目录、部门、状态信息） |
| AC-05 | 修改文件成功 | WHEN 用户修改文件字段 | THE SYSTEM SHALL 更新文件，版本号+1 |
| AC-06 | 修改无变更 | WHEN 用户修改但所有字段无变化 | THE SYSTEM SHALL 直接返回原数据，不更新版本 |
| AC-07 | 批量删除成功 | WHEN 用户删除多个文件ID | THE SYSTEM SHALL 物理删除文件 |
| AC-08 | 启用文件成功 | WHEN 用户启用已停用的文件 | THE SYSTEM SHALL 更新状态为启用，清空停用原因 |
| AC-09 | 停用文件成功 | WHEN 用户停用文件并提供停用原因 | THE SYSTEM SHALL 更新状态为停用，保存原因和停用时间 |
| AC-10 | 目录移动成功 | WHEN 用户批量移动文件到目标目录 | THE SYSTEM SHALL 更新catalog_id |
| AC-11 | 下载文件成功(FILE类型) | WHEN 用户下载FILE类型文件 | THE SYSTEM SHALL 从OSS返回文件流 |
| AC-12 | 下载文件失败(URL类型) | WHEN 用户下载URL类型文件 | THE SYSTEM SHALL 返回错误提示 |
| AC-13 | 批量下载成功 | WHEN 用户批量下载多个文件 | THE SYSTEM SHALL 返回ZIP压缩包 |
| AC-14 | 查询关联数据元成功 | WHEN 用户查询文件关联的数据元 | THE SYSTEM SHALL 返回数据元分页列表 |
| AC-15 | 查询关联码表成功 | WHEN 用户查询文件关联的码表 | THE SYSTEM SHALL 返回码表分页列表 |
| AC-16 | 查询关联编码规则成功 | WHEN 用户查询文件关联的编码规则 | THE SYSTEM SHALL 返回编码规则分页列表 |
| AC-17 | 添加关联关系成功 | WHEN 用户添加文件关联关系 | THE SYSTEM SHALL 替换原有关联关系 |
| AC-18 | 查询关联关系成功 | WHEN 用户查询文件关联关系 | THE SYSTEM SHALL 返回所有关联的ID列表 |
| AC-19 | 检查数据存在成功 | WHEN 用户检查文件数据是否存在 | THE SYSTEM SHALL 返回true/false |
| AC-20 | 批量启用/停用成功 | WHEN 用户批量操作文件状态 | THE SYSTEM SHALL 批量更新状态 |

### 异常处理

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-101 | 标准编号重复 | WHEN 创建时标准编号已存在 | THE SYSTEM SHALL 返回 InvalidParameter，提示"标准编号重复" |
| AC-102 | 文件名称重复 | WHEN 创建时文件名称在同一orgType下已存在 | THE SYSTEM SHALL 返回 InvalidParameter，提示"标准文件名称重复" |
| AC-103 | 目录不存在 | WHEN 指定的catalog_id不存在或类型不符 | THE SYSTEM SHALL 返回 InvalidParameter，提示"目录id[{id}]对应的目录不存在" |
| AC-104 | 文件为空 | WHEN attachmentType=FILE时未上传文件 | THE SYSTEM SHALL 返回 InvalidParameter，提示"文件不能为空" |
| AC-105 | 文件类型不支持 | WHEN 上传不支持的文件类型 | THE SYSTEM SHALL 返回 InvalidParameter，提示"不支持的文件类型" |
| AC-106 | 文件超过大小限制 | WHEN 上传文件超过30M | THE SYSTEM SHALL 返回 InvalidParameter，提示"文件不能超过30M" |
| AC-107 | 链接地址为空 | WHEN attachmentType=URL时链接为空 | THE SYSTEM SHALL 返回 InvalidParameter，提示"链接不能为空" |
| AC-108 | 链接地址过长 | WHEN 链接地址超过2048字符 | THE SYSTEM SHALL 返回 InvalidParameter，提示"链接长度超过2048" |
| AC-109 | 修改记录不存在 | WHEN 修改时文件ID不存在 | THE SYSTEM SHALL 返回 DATA_NOT_EXIST，提示"数据不存在" |
| AC-110 | 停用原因缺失 | WHEN state=disable时reason为空 | THE SYSTEM SHALL 返回 PARAMETER_EMPTY，提示"停用必须填写停用原因" |
| AC-111 | 停用原因过长 | WHEN reason长度超过800字符 | THE SYSTEM SHALL 返回 InvalidParameter，提示"长度超过800" |
| AC-112 | 下载URL类型文件 | WHEN 下载URL类型文件 | THE SYSTEM SHALL 返回 FileDownloadFailed，提示"[URL]类型没有文件附件" |
| AC-113 | 文件不存在 | WHEN 下载/查询的文件ID不存在 | THE SYSTEM SHALL 返回 DATA_NOT_EXIST，提示"数据不存在" |

---

## Edge Cases

| ID | Case | Expected Behavior |
|----|------|-------------------|
| EC-01 | 批量下载文件名重复 | ZIP中文件名格式：{文件前缀}({标准分类})({标准文件名称}).{后缀} |
| EC-02 | 修改时排除自身 | queryDataExists接口支持filter_id参数排除当前记录 |
| EC-03 | 实施日期格式不正确 | 返回"实施日期格式不正确"错误 |
| EC-04 | 目录ID为44 | 默认值为"全部目录" |
| EC-05 | 修改时文件为空(FILE类型) | 保持原文件不变 |

---

## Business Rules

| ID | Rule | Description |
|----|------|-------------|
| BR-01 | 标准编号唯一性 | 标准编号必须唯一 |
| BR-02 | 文件名称唯一性 | 同一orgType下文件名称必须唯一 |
| BR-03 | 版本控制 | 修改特定字段时版本号+1 |
| BR-04 | 文件类型限制 | 支持doc、pdf、docx、txt、ppt、pptx、xls、xlsx |
| BR-05 | 文件大小限制 | 最大30M |
| BR-06 | 物理删除 | 删除文件时物理删除记录 |
| BR-07 | 默认catalog_id | 不传时默认为44（全部目录） |
| BR-08 | 停用原因 | 停用时必须填写原因且<=800字符，启用时清空原因 |
| BR-09 | XSS转义 | keyword、name参数需进行XSS转义 |
| BR-10 | 部门ID处理 | department_ids存储完整路径，查询/返回时使用最后一段 |

---

## Data Considerations

| Field | Description | Constraints |
|-------|-------------|-------------|
| id | 主键 | Long类型 |
| number | 标准编号 | 可选，最大300字符，唯一 |
| name | 文件名称 | 必填，最大300字符，同orgType下唯一 |
| catalog_id | 所属目录ID | 必填，默认44，关联目录表 |
| org_type | 标准组织类型 | 必填，0-99枚举值 |
| attachment_type | 附件类型 | 必填，FILE(0)或URL(1) |
| attachment_url | 链接地址 | attachmentType=URL时必填，最大2048字符 |
| file_name | 文件名 | attachmentType=FILE时必填 |
| act_date | 实施日期 | 可选，格式yyyy-MM-dd |
| publish_date | 发布日期 | 可选，格式yyyy-MM-dd |
| state | 状态 | ENABLE(1)或DISABLE(0)，默认ENABLE |
| disable_date | 停用时间 | 停用时记录 |
| disable_reason | 停用原因 | 停用时必填，最大800字符 |
| version | 版本号 | 从1开始，修改特定字段时+1 |
| description | 说明 | 可选，最大300字符 |
| department_ids | 部门ID | 可选，字符串格式 |
| third_dept_id | 第三方部门ID | 可选，UUID格式 |

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
2. **日期格式**：yyyy-MM-dd
3. **文件存储**：OSS（阿里云对象存储）
4. **分页参数**：offset从1开始（不是0）
5. **目录服务**：依赖catalog-api模块的目录服务
6. **数据元服务**：依赖dataelement-api模块的数据元服务
7. **字典服务**：依赖dict-api模块的字典服务
8. **编码规则服务**：依赖rule-api模块的编码规则服务
9. **部门信息**：从Token中解析部门信息
10. **用户信息**：从Token中解析用户信息

---

## Dependencies

- catalog-api：目录查询服务
- dataelement-api：数据元查询服务
- dict-api：字典查询服务
- rule-api：编码规则查询服务
- OSS服务：文件存储服务
- 部门服务：部门信息查询

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-02-06 | - | 初始版本 |
