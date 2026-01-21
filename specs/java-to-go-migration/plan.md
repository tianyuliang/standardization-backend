# Java to Go Migration Technical Plan

> **Branch**: `feature/java-to-go-migration`
> **Spec Path**: `specs/java-to-go-migration/`
> **Created**: 2026-01-19
> **Status**: Draft

---

## Summary

本技术方案将 Spring Boot 2.7.5 数据标准化管理后台服务迁移到 Go-Zero 微服务架构。采用单体 API 服务模式，使用 GORM 进行数据访问，完整迁移 6 个功能模块。保持现有 MariaDB/达梦数据库，完整支持 Kafka/NSQ 消息队列和 Redis 缓存/分布式锁功能。

---

## Technical Context

| Item | Value |
|------|-------|
| **Language** | Go 1.24+ |
| **Framework** | Go-Zero v1.9+ |
| **Storage** | MariaDB / 达梦数据库 (保持现有) |
| **Cache** | Redis + Redisson 分布式锁 |
| **ORM** | GORM v1.25+ |
| **Message Queue** | Kafka + NSQ |
| **Testing** | go test |
| **Common Lib** | idrm-go-base v0.1.0+ |

---

## 通用库 (idrm-go-base)

**安装**:
```bash
go get github.com/jinguoxing/idrm-go-base@latest
```

### 模块初始化

| 模块 | 初始化方式 |
|------|-----------|
| validator | `validator.Init()` 在 main.go |
| telemetry | `telemetry.Init(cfg)` 在 main.go |
| response | `httpx.SetErrorHandler(response.ErrorHandler)` |
| middleware | `rest.WithMiddlewares(...)` |

### 自定义错误码

| 功能 | 范围 | 位置 |
|------|------|------|
| 数据元素管理 | 30100-30199 | `internal/errorx/codes.go` |
| 目录管理 | 30200-30299 | `internal/errorx/codes.go` |
| 码表字典 | 30300-30399 | `internal/errorx/codes.go` |
| 编码规则 | 30400-30499 | `internal/errorx/codes.go` |
| 标准化任务 | 30500-30599 | `internal/errorx/codes.go` |
| 标准文件 | 30600-30699 | `internal/errorx/codes.go` |

### 第三方库确认

| 库 | 版本 | 用途 | 确认状态 |
|----|------|------|----------|
| github.com/go-sql-driver/mysql | v1.7.1 | MariaDB 驱动 | ✅ 确认 |
| github.com/dtm-labs/driver | 待定 | 达梦数据库驱动 | ⏳ 待确认 |
| github.com/IBM/sarama | v1.42.1 | Kafka 客户端 | ✅ 确认 |
| github.com/nsqio/go-nsq | v1.1.0 | NSQ 客户端 | ✅ 确认 |
| github.com/redis/go-redis/v9 | v9.0.5 | Redis 客户端 | ✅ 确认 |
| github.com/go-redsync/redsync/v4 | v4.11.0 | 分布式锁 | ✅ 确认 |
| github.com/xuri/excelize/v2 | v2.8.0 | Excel 导入导出 | ✅ 确认 |

---

## Go-Zero 开发流程

| Step | 任务 | 方式 | 产出 |
|------|------|------|------|
| 1 | 定义 API 入口文件 | AI 实现 | `api/doc/api.api` |
| 2 | 定义基础类型 | AI 实现 | `api/doc/base.api` |
| 3 | 定义模块 API 文件 | AI 实现 | `api/doc/{module}/*.api` |
| 4 | 生成 Handler/Types | goctl | `api/internal/handler/`, `types/` |
| 5 | 定义 DDL 文件 | AI 实现 | `migrations/{module}/*.sql` |
| 6 | 实现 Model 接口 | AI 实现 | `model/{module}/*/` |
| 7 | 实现 Logic 层 | AI 实现 | `api/internal/logic/` |
| 8 | 实现中间件 | AI 实现 | `api/internal/middleware/` |

**goctl 命令**:
```bash
# 完整生成（在 api.doc/api.api 修改后执行）
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group

# 生成 Swagger 文档
goctl api plugin -plugin goctl-swagger="swagger -filename api.json" -api api/doc/api.api -dir .
```

---

## File Structure

### 项目结构

```
standardization-backend/
├── api/                          # API 服务
│   ├── doc/                      # API 定义
│   │   ├── api.api               # 入口文件
│   │   ├── base.api              # 基础类型
│   │   ├── dataelement/          # 数据元素模块
│   │   ├── catalog/              # 目录模块
│   │   ├── dict/                 # 码表模块
│   │   ├── rule/                 # 规则模块
│   │   ├── task/                 # 任务模块
│   │   └── file/                 # 文件模块
│   ├── etc/
│   │   └── standardization-api.yaml  # 配置文件
│   └── internal/
│       ├── handler/              # 请求处理
│       ├── logic/                # 业务逻辑
│       ├── middleware/           # 中间件
│       ├── svc/                  # 服务上下文
│       └── types/                # 类型定义
├── model/                        # 数据模型
│   ├── dataelement/
│   ├── catalog/
│   ├── dict/
│   ├── rule/
│   ├── task/
│   └── file/
├── migrations/                   # DDL 迁移
│   ├── dataelement/
│   ├── catalog/
│   ├── dict/
│   ├── rule/
│   ├── task/
│   └── file/
├── pkg/                          # 公共包
│   ├── mq/                       # 消息队列
│   │   ├── kafka/                # Kafka 生产者/消费者
│   │   └── nsq/                  # NSQ 生产者/消费者
│   ├── cache/                    # 缓存
│   │   ├── redis/                # Redis 客户端
│   │   └── lock/                 # 分布式锁
│   ├── excel/                    # Excel 处理
│   └── idgen/                    # ID 生成 (UUID v7)
├── specs/java-to-go-migration/   # SDD 文档
└── go.mod
```

---

## Architecture Overview

### 分层架构

```
┌─────────────────────────────────────────────────────────────┐
│                        HTTP Request                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  Handler Layer (参数校验、格式化响应)                         │
│  - 参数绑定                                                  │
│  - validator 校验                                           │
│  - 响应格式化                                                │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  Logic Layer (业务逻辑)                                      │
│  - 业务规则实现                                              │
│  - 事务管理                                                  │
│  - 缓存操作                                                  │
│  - 消息发送                                                  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  Model Layer (数据访问)                                      │
│  - GORM 操作                                                 │
│  - 缓存读写                                                  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  Storage Layer                                               │
│  MariaDB / 达梦 / Redis / Kafka / NSQ                       │
└─────────────────────────────────────────────────────────────┘
```

### 中间件链

```
Request → [TraceID] → [TokenAuth] → [AuditLog] → [RateLimit] → Handler
         ↑              ↑              ↑             ↑
      telemetry    middleware    middleware    middleware
```

---

## Module 1: 数据元素管理 (DataElement)

### API 定义

**位置**: `api/doc/dataelement/dataelement.api`

```api
syntax = "v1"

import "../base.api"

type (
    // 创建数据元请求
    CreateDataElementReq {
        ElementName   string `json:"elementName" validate:"required,max=100"`
        ElementCode   string `json:"elementCode" validate:"required,max=50"`
        DepartmentIds string `json:"departmentIds"`
        DataType      int32  `json:"dataType" validate:"required,min=0"`
        DataLength    int32  `json:"dataLength"`
    }

    // 创建数据元响应
    CreateDataElementResp {
        Id string `json:"id"`
    }

    // 查询数据元请求
    GetDataElementReq {
        Id string `json:"id" validate:"required"`
    }

    // 查询数据元响应
    GetDataElementResp {
        Id            string `json:"id"`
        ElementName   string `json:"elementName"`
        ElementCode   string `json:"elementCode"`
        DepartmentIds string `json:"departmentIds"`
        DataType      int32  `json:"dataType"`
        DataLength    int32  `json:"dataLength"`
        Status        int32  `json:"status"`
        CreateTime    string `json:"createTime"`
        UpdateTime    string `json:"updateTime"`
    }

    // 列表查询请求
    ListDataElementReq {
        Page     int32  `json:"page" validate:"required,min=1"`
        PageSize int32  `json:"pageSize" validate:"required,min=1,max=100"`
        Keyword  string `json:"keyword"`
    }

    // 列表查询响应
    ListDataElementResp {
        Total int64                `json:"total"`
        List  []GetDataElementResp `json:"list"`
    }
)

@server(
    prefix: /api/v1/dataelement
    group: dataelement
)
service standardization-api {
    @handler CreateDataElement
    post / (CreateDataElementReq) returns (CreateDataElementResp)

    @handler GetDataElement
    get /:id (GetDataElementReq) returns (GetDataElementResp)

    @handler UpdateDataElement
    put /:id (CreateDataElementReq) returns (BaseResp)

    @handler DeleteDataElement
    delete /:id (GetDataElementReq) returns (BaseResp)

    @handler ListDataElement
    post /list (ListDataElementReq) returns (ListDataElementResp)
}
```

### DDL 设计

**位置**: `migrations/dataelement/data_element_info.sql`

```sql
CREATE TABLE `data_element_info` (
    `id` CHAR(36) NOT NULL COMMENT 'ID (UUID v7)',
    `element_name` VARCHAR(100) NOT NULL COMMENT '数据元名称',
    `element_code` VARCHAR(50) NOT NULL COMMENT '数据元编码',
    `department_ids` VARCHAR(500) DEFAULT NULL COMMENT '部门路径',
    `third_dept_id` CHAR(36) DEFAULT NULL COMMENT '三级部门ID',
    `data_type` INT NOT NULL DEFAULT 0 COMMENT '数据类型',
    `data_length` INT DEFAULT NULL COMMENT '数据长度',
    `status` INT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_element_code` (`element_code`),
    KEY `idx_status` (`status`),
    KEY `idx_department` (`department_ids`(255))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据元信息表';
```

### Model 接口

```go
type DataElementModel interface {
    Insert(ctx context.Context, data *DataElement) (*DataElement, error)
    FindOne(ctx context.Context, id string) (*DataElement, error)
    FindByCode(ctx context.Context, code string) (*DataElement, error)
    Update(ctx context.Context, data *DataElement) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, page, pageSize int32, keyword string) ([]*DataElement, int64, error)
    WithTx(tx *gorm.DB) DataElementModel
}
```

---

## Module 2: 目录管理 (Catalog)

### API 定义

**位置**: `api/doc/catalog/catalog.api`

```api
syntax = "v1"

import "../base.api"

type (
    CreateCatalogReq {
        CatalogName string `json:"catalogName" validate:"required,max=100"`
        ParentId    string `json:"parentId"`
        CatalogType int32  `json:"catalogType" validate:"required,min=0"`
        SortOrder   int32  `json:"sortOrder"`
    }

    CreateCatalogResp {
        Id string `json:"id"`
    }

    GetCatalogTreeReq {
        CatalogType int32 `json:"catalogType"`
    }

    GetCatalogTreeResp {
        Id          string                `json:"id"`
        CatalogName string                `json:"catalogName"`
        ParentId    string                `json:"parentId"`
        Children    []GetCatalogTreeResp  `json:"children"`
    }

    MoveCatalogReq {
        SourceId string `json:"sourceId" validate:"required"`
        TargetId string `json:"targetId" validate:"required"`
    }

    MoveCatalogResp {
        Success bool `json:"success"`
    }
)

@server(
    prefix: /api/v1/catalog
    group: catalog
)
service standardization-api {
    @handler CreateCatalog
    post / (CreateCatalogReq) returns (CreateCatalogResp)

    @handler GetCatalogTree
    get /tree (GetCatalogTreeReq) returns ([]GetCatalogTreeResp)

    @handler MoveCatalog
    post /move (MoveCatalogReq) returns (MoveCatalogResp)

    @handler DeleteCatalog
    delete /:id (BaseIdReq) returns (BaseResp)
}
```

### DDL 设计

**位置**: `migrations/catalog/de_catalog_info.sql`

```sql
CREATE TABLE `de_catalog_info` (
    `id` CHAR(36) NOT NULL COMMENT 'ID (UUID v7)',
    `catalog_name` VARCHAR(100) NOT NULL COMMENT '目录名称',
    `parent_id` CHAR(36) DEFAULT NULL COMMENT '父目录ID',
    `catalog_type` INT NOT NULL DEFAULT 0 COMMENT '目录类型',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `tree_path` VARCHAR(1000) DEFAULT NULL COMMENT '树路径',
    `level` INT NOT NULL DEFAULT 1 COMMENT '层级',
    `status` INT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_parent` (`parent_id`),
    KEY `idx_type` (`catalog_type`),
    KEY `idx_tree_path` (`tree_path`(255))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='目录信息表';
```

---

## Module 3-6: 其他模块

其他模块（码表字典、编码规则、标准化任务、标准文件）的 API 和 DDL 设计遵循相同模式，详细内容在 tasks.md 中展开。

---

## Message Queue Integration

### Kafka 配置

```go
// pkg/mq/kafka/producer.go
type KafkaProducer struct {
    producer sarama.SyncProducer
    config   *KafkaConfig
}

func (k *KafkaProducer) SendMessage(ctx context.Context, topic string, key string, value []byte) error {
    msg := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.StringEncoder(key),
        Value: sarama.ByteEncoder(value),
    }
    partition, offset, err := k.producer.SendMessage(msg)
    logx.Infof("Message sent to partition %d at offset %d", partition, offset)
    return err
}

// Topic 定义
const (
    TopicDataElementChanged = "data.element.changed"
    TopicCatalogChanged     = "catalog.changed"
    TopicDictChanged        = "dict.changed"
    TopicRuleChanged        = "rule.changed"
    TopicTaskChanged        = "task.changed"
    TopicFileChanged        = "file.changed"
)
```

### NSQ 配置

```go
// pkg/mq/nsq/producer.go
type NSQProducer struct {
    producer *nsq.Producer
}

func (n *NSQProducer) Publish(ctx context.Context, topic string, value []byte) error {
    return n.producer.Publish(topic, value)
}
```

---

## Redis & Distributed Lock

### Redis 配置

```go
// pkg/cache/redis/redis.go
type RedisClient struct {
    client *redis.Client
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
    return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
    return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
    return r.client.Del(ctx, keys...).Err()
}
```

### 分布式锁

```go
// pkg/cache/lock/distributed_lock.go
type DistributedLock interface {
    Lock(ctx context.Context, key string, expiration time.Duration) (bool, error)
    Unlock(ctx context.Context, key string) error
}

type RedsyncLock struct {
    rs *redsync.Redsync
}

func (r *RedsyncLock) Lock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
    mutex := r.rs.NewMutex(
        fmt.Sprintf("lock:%s", key),
        redsync.WithExpiry(expiration),
        redsync.WithTries(3),
    )
    return mutex.TryLock()
}

func (r *RedsyncLock) Unlock(ctx context.Context, key string) error {
    // 实现解锁逻辑
    return nil
}
```

---

## Testing Strategy

| 类型 | 方法 | 覆盖率 |
|------|------|--------|
| 单元测试 | 表驱动测试，Mock Model | > 80% |
| 集成测试 | testcontainers | 核心流程 |
| 对比测试 | 并行调用 Java/Go API | 100% API 覆盖 |

### 测试目录结构

```
api/
├── internal/
│   ├── logic/
│   │   ├── dataelement/
│   │   │   ├── create_dataelement_logic.go
│   │   │   └── create_dataelement_logic_test.go
│   │   └── ...
└── test/
    ├── integration/
    │   └── api_test.go
    └── comparison/
        └── java_vs_go_test.go
```

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-01-19 | - | 初始版本 - Java to Go 迁移技术方案 |
