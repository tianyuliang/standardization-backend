# Java to Go Migration Specification

> **Branch**: `feature/java-to-go-migration`
> **Spec Path**: `specs/java-to-go-migration/`
> **Created**: 2026-01-19
> **Status**: Draft

---

## Overview

将现有基于 Spring Boot 2.7.5 + Java 1.8 的数据标准化管理后台服务迁移到 Go-Zero 微服务架构。迁移范围包括全部 6 个功能模块（数据元素管理、目录管理、码表字典、编码规则、标准化任务、标准文件管理），保持业务逻辑不变，同时保持现有 MariaDB/达梦数据库，完整迁移 Kafka/NSQ 消息队列和 Redis 缓存功能。

---

## User Stories

### Story 1: 数据元素管理模块 (P1)

AS a 系统管理员
I WANT 使用 Go 版本的数据元素管理 API
SO THAT 可以管理数据元素的 CRUD 操作、历史记录和业务表字段管理

**独立测试**: 对比 Java 版本和 Go 版本的 API 响应，确保业务逻辑一致

### Story 2: 目录管理模块 (P1)

AS a 系统管理员
I WANT 使用 Go 版本的目录管理 API
SO THAT 可以管理目录树结构、目录移动和统计功能

**独立测试**: 对比 Java 版本和 Go 版本的 API 响应，确保业务逻辑一致

### Story 3: 码表字典模块 (P1)

AS a 系统管理员
I WANT 使用 Go 版本的码表字典 API
SO THAT 可以管理码表数据、枚举，以及 Excel 导入导出功能

**独立测试**: 对比 Java 版本和 Go 版本的 API 响应，确保业务逻辑一致

### Story 4: 编码规则管理模块 (P1)

AS a 系统管理员
I WANT 使用 Go 版本的编码规则 API
SO THAT 可以管理编码规则定义、自定义规则和规则推荐

**独立测试**: 对比 Java 版本和 Go 版本的 API 响应，确保业务逻辑一致

### Story 5: 标准化任务模块 (P1)

AS a 系统管理员
I WANT 使用 Go 版本的标准化任务 API
SO THAT 可以管理任务创建、任务状态和任务生命周期

**独立测试**: 对比 Java 版本和 Go 版本的 API 响应，确保业务逻辑一致

### Story 6: 标准文件管理模块 (P1)

AS a 系统管理员
I WANT 使用 Go 版本的标准文件 API
SO THAT 可以管理标准文件和文件附件

**独立测试**: 对比 Java 版本和 Go 版本的 API 响应，确保业务逻辑一致

### Story 7: 消息队列集成 (P1)

AS a 系统架构师
I WANT Go 版本完整支持 Kafka 和 NSQ 消息队列
SO THAT 可以保持与现有系统的异步通信能力

**独立测试**: 发送测试消息，验证 Go 版本能正确生产和消费消息

### Story 8: Redis 缓存和分布式锁 (P1)

AS a 系统架构师
I WANT Go 版本完整支持 Redis 缓存和 Redisson 分布式锁
SO THAT 可以保持与现有系统的缓存一致性和并发控制能力

**独立测试**: 验证缓存读写和分布式锁功能

---

## Acceptance Criteria (EARS)

### 正常流程 - 数据元素管理

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-01 | 创建数据元成功 | WHEN 用户提交有效数据元信息 | THE SYSTEM SHALL 保存数据元并返回 201 |
| AC-02 | 查询数据元成功 | WHEN 用户查询存在的数据元 ID | THE SYSTEM SHALL 返回数据元详情 |
| AC-03 | 更新数据元成功 | WHEN 用户更新有效数据元信息 | THE SYSTEM SHALL 保存更新并返回 200 |
| AC-04 | 删除数据元成功 | WHEN 用户删除存在的数据元 | THE SYSTEM SHALL 删除数据元并返回 204 |
| AC-05 | 数据元历史记录 | WHEN 用户查询数据元历史 | THE SYSTEM SHALL 返回完整变更历史 |
| AC-06 | 业务表字段管理 | WHEN 用户管理业务表字段 | THE SYSTEM SHALL 支持字段 CRUD 和状态管理 |

### 正常流程 - 目录管理

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-10 | 创建目录成功 | WHEN 用户提交有效目录信息 | THE SYSTEM SHALL 保存目录并返回 201 |
| AC-11 | 目录树查询 | WHEN 用户查询目录树 | THE SYSTEM SHALL 返回完整树形结构 |
| AC-12 | 目录移动成功 | WHEN 用户移动目录到新位置 | THE SYSTEM SHALL 更新目录位置并返回 200 |
| AC-13 | 目录统计 | WHEN 用户查询目录统计信息 | THE SYSTEM SHALL 返回统计数据（按目录分组） |

### 正常流程 - 码表字典

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-20 | 创建码表成功 | WHEN 用户提交有效码表信息 | THE SYSTEM SHALL 保存码表并返回 201 |
| AC-21 | 码表枚举管理 | WHEN 用户管理码表枚举 | THE SYSTEM SHALL 支持枚举 CRUD |
| AC-22 | 码表 Excel 导入 | WHEN 用户上传 Excel 文件 | THE SYSTEM SHALL 解析并导入码表数据 |
| AC-23 | 码表 Excel 导出 | WHEN 用户导出码表 | THE SYSTEM SHALL 生成并返回 Excel 文件 |

### 正常流程 - 编码规则

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-30 | 创建规则成功 | WHEN 用户提交有效规则信息 | THE SYSTEM SHALL 保存规则并返回 201 |
| AC-31 | 自定义规则管理 | WHEN 用户管理自定义规则 | THE SYSTEM SHALL 支持自定义规则 CRUD |
| AC-32 | 规则推荐 | WHEN 用户请求规则推荐 | THE SYSTEM SHALL 返回推荐规则列表 |

### 正常流程 - 标准化任务

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-40 | 创建任务成功 | WHEN 用户提交有效任务信息 | THE SYSTEM SHALL 创建任务并返回 201 |
| AC-41 | 任务状态流转 | WHEN 用户更新任务状态 | THE SYSTEM SHALL 按状态机规则流转状态 |
| AC-42 | 任务列表查询 | WHEN 用户查询任务列表 | THE SYSTEM SHALL 返回分页任务列表 |

### 正常流程 - 标准文件

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-50 | 文件上传成功 | WHEN 用户上传标准文件 | THE SYSTEM SHALL 保存文件并返回 201 |
| AC-51 | 文件下载 | WHEN 用户下载标准文件 | THE SYSTEM SHALL 返回文件内容 |
| AC-52 | 附件管理 | WHEN 用户管理文件附件 | THE SYSTEM SHALL 支持附件 CRUD |

### 正常流程 - 消息队列

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-60 | Kafka 消息发送 | WHEN 系统需要发送 Kafka 消息 | THE SYSTEM SHALL 成功发送到指定 topic |
| AC-61 | NSQ 消息发送 | WHEN 系统需要发送 NSQ 消息 | THE SYSTEM SHALL 成功发送到指定 topic |
| AC-62 | 消息消费 | WHEN 消息队列有新消息 | THE SYSTEM SHALL 正确处理消息 |

### 正常流程 - Redis

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-70 | 缓存写入 | WHEN 数据需要缓存 | THE SYSTEM SHALL 写入 Redis 并设置过期时间 |
| AC-71 | 缓存读取 | WHEN 查询缓存数据 | THE SYSTEM SHALL 从 Redis 读取并返回 |
| AC-72 | 分布式锁获取 | WHEN 需要分布式锁 | THE SYSTEM SHALL 成功获取锁并在用完后释放 |

### 异常处理

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-100 | 参数为空 | WHEN 必填参数为空 | THE SYSTEM SHALL 返回 400 和错误信息 |
| AC-101 | 资源不存在 | WHEN 查询不存在的 ID | THE SYSTEM SHALL 返回 404 |
| AC-102 | 名称重复 | WHEN 名称与已有资源重复 | THE SYSTEM SHALL 返回 409 |
| AC-103 | 权限不足 | WHEN 用户无操作权限 | THE SYSTEM SHALL 返回 403 |
| AC-104 | 数据库连接失败 | WHEN 数据库不可用 | THE SYSTEM SHALL 返回 503 并记录日志 |
| AC-105 | 消息队列连接失败 | WHEN 消息队列不可用 | THE SYSTEM SHALL 降级处理并记录日志 |
| AC-106 | Redis 连接失败 | WHEN Redis 不可用 | THE SYSTEM SHALL 降级到直连数据库并记录日志 |

---

## Edge Cases

| ID | Case | Expected Behavior |
|----|------|-------------------|
| EC-01 | 并发创建同名资源 | 仅一个成功，其他返回 409 |
| EC-02 | 删除被引用的资源 | 返回 400，提示存在关联 |
| EC-03 | 批量操作部分失败 | 返回失败项列表和成功项列表 |
| EC-04 | Excel 导入格式错误 | 返回详细的行级错误信息 |
| EC-05 | 分布式锁超时 | 释放锁并记录超时日志 |
| EC-06 | 消息消费失败 | 重试 3 次，失败后转入死信队列 |
| EC-07 | 大文件上传 | 支持分片上传和断点续传 |

---

## Business Rules

| ID | Rule | Description |
|----|------|-------------|
| BR-01 | ID 生成 | 使用 UUID v7 替代雪花算法生成主键 |
| BR-02 | 数据库兼容 | Go 版本需同时支持 MariaDB 和达梦数据库 |
| BR-03 | 响应格式 | 使用 idrm-go-base 标准响应格式 |
| BR-04 | 错误码 | 使用 idrm-go-base errorx 模块统一错误码 |
| BR-05 | 审计日志 | 关键操作必须记录审计日志（复用 Java 格式） |
| BR-06 | 分页规范 | 使用 go-zero PageResult 标准分页 |
| BR-07 | 树形结构 | 目录树支持无限层级，使用递归查询优化 |
| BR-08 | 状态机 | 任务状态必须按指定顺序流转 |
| BR-09 | 消息幂等 | 消息消费必须支持幂等性处理 |
| BR-10 | 缓存一致性 | 数据更新时同步更新/删除缓存 |

---

## Data Considerations

### 数据元素 (data_element_info)

| Field | Description | Constraints |
|-------|-------------|-------------|
| id | ID (UUID v7) | CHAR(36), PK |
| element_name | 数据元名称 | VARCHAR(100), NOT NULL |
| element_code | 数据元编码 | VARCHAR(50), UNIQUE |
| department_ids | 部门路径 | VARCHAR(500) |
| data_type | 数据类型 | INT (枚举) |
| data_length | 数据长度 | INT |
| status | 状态 | INT (0:禁用, 1:启用) |
| create_time | 创建时间 | DATETIME |
| update_time | 更新时间 | DATETIME |

### 目录信息 (de_catalog_info)

| Field | Description | Constraints |
|-------|-------------|-------------|
| id | ID (UUID v7) | CHAR(36), PK |
| catalog_name | 目录名称 | VARCHAR(100), NOT NULL |
| parent_id | 父目录 ID | CHAR(36) |
| catalog_type | 目录类型 | INT (枚举) |
| sort_order | 排序 | INT, DEFAULT 0 |
| tree_path | 树路径 | VARCHAR(1000) |
| create_time | 创建时间 | DATETIME |

### 码表字典 (dict_info)

| Field | Description | Constraints |
|-------|-------------|-------------|
| id | ID (UUID v7) | CHAR(36), PK |
| dict_code | 字典编码 | VARCHAR(50), UNIQUE |
| dict_name | 字典名称 | VARCHAR(100), NOT NULL |
| dict_type | 字典类型 | VARCHAR(50) |
| status | 状态 | INT (0:禁用, 1:启用) |
| create_time | 创建时间 | DATETIME |

### 编码规则 (rule_info)

| Field | Description | Constraints |
|-------|-------------|-------------|
| id | ID (UUID v7) | CHAR(36), PK |
| rule_name | 规则名称 | VARCHAR(100), NOT NULL |
| rule_type | 规则类型 | INT (枚举) |
| rule_expression | 规则表达式 | TEXT |
| status | 状态 | INT (0:禁用, 1:启用) |
| create_time | 创建时间 | DATETIME |

### 标准化任务 (std_create_task)

| Field | Description | Constraints |
|-------|-------------|-------------|
| id | ID (UUID v7) | CHAR(36), PK |
| task_name | 任务名称 | VARCHAR(100), NOT NULL |
| task_type | 任务类型 | INT (枚举) |
| task_status | 任务状态 | INT (枚举) |
| create_source | 创建来源 | INT (枚举) |
| create_time | 创建时间 | DATETIME |

### 标准文件 (std_file_mgr)

| Field | Description | Constraints |
|-------|-------------|-------------|
| id | ID (UUID v7) | CHAR(36), PK |
| file_name | 文件名称 | VARCHAR(200), NOT NULL |
| file_path | 文件路径 | VARCHAR(500) |
| file_size | 文件大小 | BIGINT |
| file_status | 文件状态 | INT (枚举) |
| create_time | 创建时间 | DATETIME |

---

## Success Metrics

| ID | Metric | Target |
|----|--------|--------|
| SC-01 | API 接口覆盖率 | 100% (所有 Java API 都有对应的 Go 实现) |
| SC-02 | 业务逻辑一致性 | 100% (单元测试通过) |
| SC-03 | 接口响应时间 | P99 < 200ms (与 Java 版本持平或更优) |
| SC-04 | 测试覆盖率 | > 80% |
| SC-05 | 消息队列功能 | Kafka 和 NSQ 双队列正常工作 |
| SC-06 | Redis 功能 | 缓存和分布式锁正常工作 |

---

## Open Questions

- [ ] 达梦数据库的 Go 驱动是否可用？如果不可用，是否有替代方案？
- [ ] 现有 Kafka/NSQ 消息格式是否需要调整？
- [ ] 是否需要支持双写（Java 和 Go 同时写入）作为过渡？
- [ ] 审计日志格式是否需要调整以适配 Go 版本？
- [ ] Excel 导入导出功能使用哪个 Go 库？

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-01-19 | - | 初始版本 - Java to Go 迁移规格 |
