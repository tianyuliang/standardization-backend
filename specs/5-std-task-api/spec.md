# 标准任务管理 (std-task-api) Specification

> **Branch**: `5-std-task-api`
> **Spec Path**: `specs/5-std-task-api/`
> **Created**: 2026-02-09
>
> **Status**: Draft

---

## Overview

标准任务管理模块用于管理系统中的标准创建任务，支持以下核心功能：
- **任务管理**：标准创建任务的查询、创建、完成
- **标准推荐**：调用外部推荐算法获取标准推荐
- **规则推荐**：调用外部推荐算法获取编码规则推荐
- **业务表管理**：待新建标准关联的业务表管理

**重要约束**：从 Java 迁移到 Go，必须100%保持接口兼容。接口路径、请求参数、响应格式、异常信息必须与原Java实现一致。

---

## User Stories

### Story 1: 任务查询管理 (P1)

AS a 系统管理员
I WANT 查询未处理和已完成的任务列表及详情
SO THAT 能够跟踪标准创建任务的进度

**独立测试**:
- 查询未处理任务列表能正确分页返回
- 查询已完成任务列表能正确分页返回
- 查询任务详情能返回任务结果

### Story 2: 任务创建与完成 (P1)

AS a 系统管理员
I WANT 创建标准创建任务并完成任务
SO THAT 能够管理标准创建的完整流程

**独立测试**:
- 标准关联暂存成功
- 标准关联提交后异步调用推荐算法
- 完成任务后更新状态并发送回调

### Story 3: 业务表管理 (P1)

AS a 系统管理员
I WANT 管理待新建标准的业务表
SO THAT 能够批量处理业务表字段

**独立测试**:
- 添加业务表到待新建成功
- 查询业务表列表成功
- 查询业务表字段列表成功

### Story 4: 标准推荐 (P2)

AS a 系统管理员
I WANT 调用推荐算法获取标准推荐
SO THAT 能够快速匹配相关标准

**独立测试**:
- 标准推荐接口返回最多3条推荐结果
- 编码规则推荐接口返回推荐结果

### Story 5: 任务状态查询 (P2)

AS a 外部系统
I WANT 查询任务进度和状态
SO THAT 能够获取任务执行情况

**独立测试**:
- 进度查询接口返回任务状态
- 任务状态查询接口返回状态信息

---

## Acceptance Criteria (EARS)

### 正常流程

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-01 | 未处理任务列表查询成功 | WHEN 用户查询未处理任务 | THE SYSTEM SHALL 返回分页列表 |
| AC-02 | 已完成任务列表查询成功 | WHEN 用户查询已完成任务 | THE SYSTEM SHALL 返回分页列表 |
| AC-03 | 任务详情查询成功 | WHEN 用户查询任务ID | THE SYSTEM SHALL 返回任务及结果列表 |
| AC-04 | 标准关联暂存成功 | WHEN 用户提交暂存请求 | THE SYSTEM SHALL 保存任务数据 |
| AC-05 | 标准关联提交成功 | WHEN 用户提交任务 | THE SYSTEM SHALL 保存数据并调用推荐算法 |
| AC-06 | 添加至待新建成功 | WHEN 用户添加业务表 | THE SYSTEM SHALL 保存到待新建表 |
| AC-07 | 业务表列表查询成功 | WHEN 用户查询业务表 | THE SYSTEM SHALL 返回分页列表 |
| AC-08 | 业务表字段列表查询成功 | WHEN 用户查询字段 | THE SYSTEM SHALL 返回字段列表 |
| AC-09 | 移除字段成功 | WHEN 用户删除字段 | THE SYSTEM SHALL 移除指定字段 |
| AC-10 | 新建标准任务成功 | WHEN 用户创建任务 | THE SYSTEM SHALL 生成任务并更新状态 |
| AC-11 | 撤销操作成功 | WHEN 用户执行撤销 | THE SYSTEM SHALL 更新业务表状态 |
| AC-12 | 任务关联业务表查询成功 | WHEN 用户查询任务业务表 | THE SYSTEM SHALL 返回业务表列表 |
| AC-13 | 任务关联字段查询成功 | WHEN 用户查询任务字段 | THE SYSTEM SHALL 返回字段列表 |
| AC-14 | 提交选定数据元成功 | WHEN 用户提交数据元 | THE SYSTEM SHALL 保存数据元选择 |
| AC-15 | 完成任务成功 | WHEN 用户完成任务 | THE SYSTEM SHALL 更新状态并发送回调 |
| AC-16 | 进度查询成功 | WHEN 用户查询进度 | THE SYSTEM SHALL 返回任务状态 |
| AC-17 | 任务状态查询成功 | WHEN 用户查询任务状态 | THE SYSTEM SHALL 返回状态信息 |
| AC-18 | 修改字段说明成功 | WHEN 用户修改字段说明 | THE SYSTEM SHALL 更新说明 |
| AC-19 | 采纳操作成功 | WHEN 用户采纳建议 | THE SYSTEM SHALL 保存采纳结果 |
| AC-20 | 修改表名称成功 | WHEN 用户修改表名 | THE SYSTEM SHALL 更新表名 |
| AC-21 | 标准推荐成功（内部） | WHEN 调用标准推荐 | THE SYSTEM SHALL 返回推荐结果 |
| AC-22 | 标准创建成功（内部） | WHEN 调用标准创建 | THE SYSTEM SHALL 执行创建流程 |
| AC-23 | 标准推荐成功（弹框） | WHEN 调用弹框推荐 | THE SYSTEM SHALL 返回推荐结果 |
| AC-24 | 编码规则推荐成功 | WHEN 调用规则推荐 | THE SYSTEM SHALL 返回规则推荐 |

### 异常处理

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-101 | 任务不存在 | WHEN 查询不存在的任务ID | THE SYSTEM SHALL 返回 DATA_NOT_EXIST |
| AC-102 | 参数为空 | WHEN 必填参数为空 | THE SYSTEM SHALL 返回 PARAMETER_EMPTY |
| AC-103 | 业务表不存在 | WHEN 操作不存在的业务表 | THE SYSTEM SHALL 返回 DATA_NOT_EXIST |
| AC-104 | 推荐服务调用失败 | WHEN 推荐服务不可用 | THE SYSTEM SHALL 记录日志并返回空结果 |
| AC-105 | 回调失败 | WHEN AF系统回调失败 | THE SYSTEM SHALL 记录日志不阻塞主流程 |

---

## Edge Cases

| ID | Case | Expected Behavior |
|----|------|-------------------|
| EC-01 | 推荐结果为空 | 返回空列表，不报错 |
| EC-02 | 任务已完成再次完成 | 返回成功，状态保持完成 |
| EC-03 | 业务表状态为已完成 | 跳过处理已完成的业务表 |
| EC-04 | Webhook为空 | 完成任务时不发送回调 |

---

## Business Rules

| ID | Rule | Description |
|----|------|-------------|
| BR-01 | 任务编号 | 使用UUID或自增ID生成唯一任务编号 |
| BR-02 | 任务状态 | 0-未处理，1-处理中，2-处理完成 |
| BR-03 | 推荐算法 | 调用外部推荐算法服务，最多返回3条结果 |
| BR-04 | 异步处理 | 标准关联提交后异步调用推荐算法 |
| BR-05 | 回调机制 | 任务完成后调用Webhook发送结果 |
| BR-06 | 逻辑删除 | 使用deleted字段标记删除 |

---

## Data Considerations

### t_task_std_create (标准创建任务表)

| Field | Description | Constraints |
|-------|-------------|-------------|
| f_id | 主键 | BIGINT(20) AUTO_INCREMENT |
| f_task_no | 任务编号 | VARCHAR(64) NOT NULL |
| f_table | 业务表名称 | VARCHAR(128) |
| f_table_description | 业务表描述 | VARCHAR(256) |
| f_table_field | 表字段名称 | VARCHAR(1024) |
| f_status | 任务状态 | INT(2) NOT NULL DEFAULT 0 |
| f_create_time | 创建时间 | DATETIME |
| f_create_user | 创建用户ID | VARCHAR(128) |
| f_create_user_phone | 创建用户联系方式 | VARCHAR(32) |
| f_update_time | 修改时间 | DATETIME |
| f_update_user | 修改用户ID | VARCHAR(128) |
| f_webhook | AF回调地址 | VARCHAR(256) |
| f_deleted | 逻辑删除标记 | BIGINT(20) DEFAULT 0 |

### t_task_std_create_result (任务结果表)

| Field | Description | Constraints |
|-------|-------------|-------------|
| f_id | 主键 | BIGINT(20) AUTO_INCREMENT |
| f_task_id | 标准推荐任务ID | BIGINT(20) NOT NULL |
| f_table_field | 表字段名称 | VARCHAR(64) |
| f_table_field_description | 表字段描述 | VARCHAR(256) |
| f_std_ref_file | 参考标准文件 | VARCHAR(256) |
| f_std_code | 标准编码 | VARCHAR(64) |
| f_rec_std_codes | 推荐算法结果标准编码 | VARCHAR(512) |
| f_std_ch_name | 标准中文名称 | VARCHAR(128) |
| f_std_en_name | 标准英文名称 | VARCHAR(256) |

### t_business_table_std_create_pool (待新建标准表)

| Field | Description | Constraints |
|-------|-------------|-------------|
| f_id | 主键 | BIGINT(20) AUTO_INCREMENT |
| f_table_name | 业务表名称 | VARCHAR(128) NOT NULL |
| f_table_description | 业务表描述 | VARCHAR(256) |
| f_table_field | 表字段 | VARCHAR(2048) |
| f_status | 状态 | INT(2) DEFAULT 0 |
| f_create_user | 创建用户 | VARCHAR(128) |
| f_create_time | 创建时间 | DATETIME |
| f_deleted | 逻辑删除标记 | BIGINT(20) DEFAULT 0 |

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

1. **主键类型**：使用BIGINT自增ID（与Java实现一致）
2. **推荐服务**：推荐算法服务URL从配置文件读取
3. **回调机制**：Webhook回调使用HTTP POST
4. **分页参数**：offset从1开始
5. **用户信息**：从Token中解析用户信息
6. **推荐结果限制**：标准推荐最多返回3条结果

---

## Dependencies

- 推荐算法服务（recServiceUrl）：标准推荐
- 规则推荐服务（recRuleServiceUrl）：规则推荐
- AF系统：任务回调
- Kafka：异步消息处理

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-02-09 | - | 初始版本 |
