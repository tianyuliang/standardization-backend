# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# standardization-backend

> **项目状态**: Java → Go-Zero 迁移中 (Migration in Progress)
>
> 当前代码库包含 **Java Spring Boot** 源代码和 **Go-Zero** 目标架构的规范文档。新功能开发应遵循 SDD 流程，最终实现为 Go 代码。

基于 Go-Zero 微服务架构的项目，采用 AI 辅助的规范驱动开发 (SDD) 模式。

## ⚠️ 强制工作流 (必读)

**任何功能开发必须遵循 SDD 流程，不允许跳过！**

### 开发请求识别

当用户请求涉及以下内容时，视为"功能开发请求"：
- 新增功能、接口、模块
- 修改现有功能逻辑
- Bug 修复（非简单配置修改）
- 重构代码

### 强制行为

1. **必须使用 speckit 命令启动**
   ```
   /speckit.start <功能描述>    # 智能场景匹配
   /speckit.specify <功能描述>  # 创建规格文档
   ```

2. **禁止跳过 SDD 阶段**
   - ❌ 直接编写代码
   - ❌ 未创建 spec.md 就开始实现
   - ❌ 未经用户确认就进入下一阶段

3. **如果用户要求直接编码**
   - 停止并提醒："此请求涉及功能开发，请使用 `/speckit.start <功能描述>` 启动 SDD 流程"
   - 解释 SDD 流程的价值（可追溯、可测试、减少返工）

4. **specs 文件格式**
   - 所有 `specs/` 下的文件必须使用 `.specify/templates/` 中的模板
   - spec.md 使用 EARS 格式 (WHEN...THE SYSTEM SHALL)
   - plan.md 包含 API/DDL/Model 设计

5. **测试先行 (Test-First) 🧪**
   - ❌ 没有测试的代码不允许提交
   - ✅ 实现代码必须同时编写测试用例
   - ✅ 测试覆盖率 > 80%
   - ✅ 使用 `[TEST]` 标记测试任务

### 例外情况

以下请求可以直接处理，无需 SDD 流程：
- 配置文件修改
- 依赖版本更新
- 代码格式化
- 简单问答和解释

---

## 技术栈

### 当前 (Java Spring Boot)
- **语言**: Java 1.8
- **框架**: Spring Boot 2.7.5
- **数据库**: MariaDB / 达梦数据库 (DM8)
- **ORM**: MyBatis Plus 3.5.2
- **服务器**: Undertow
- **API 文档**: Swagger 3.0 + Knife4j
- **消息队列**: Kafka + NSQ
- **缓存**: Redis + Redisson (分布式锁)

### 目标 (Go-Zero)
- **语言**: Go 1.24+
- **框架**: Go-Zero v1.9+
- **数据库**: MySQL 8.0
- **ORM**: GORM (复杂查询) + SQLx (高性能)
- **架构**: 微服务 (API/RPC/Job/Consumer)
- **通用库**: idrm-go-base v0.1.0+

## 通用库规范 (idrm-go-base)

### 必须使用

| 场景 | 使用模块 | 禁止行为 |
|------|----------|----------|
| 错误处理 | `errorx` | ❌ 自定义 error struct |
| HTTP 响应 | `response` | ❌ 自定义响应格式 |
| API 中间件 | `middleware` | ❌ 重复实现认证/日志 |
| 参数校验 | `validator` | ❌ 手写校验逻辑 |
| 日志追踪 | `telemetry` | ❌ 直接使用 fmt/log |

### Import 路径

```go
import "github.com/jinguoxing/idrm-go-base/{module}"
```

### 引入其他库规则

如需使用通用库以外的第三方库：
- **停止** 并询问：该库是否可以使用？
- 等待确认后再继续

### 主键规范 (UUID v7)

所有表使用 UUID v7 作为主键：

```sql
`id` CHAR(36) NOT NULL COMMENT 'ID (UUID v7)'
```

```go
Id string `gorm:"primaryKey;size:36"`  // UUID v7
```

## 项目结构

### 当前 (Java Spring Boot)
```
standardization-backend/
├── src/main/java/com/dsg/standardization/
│   ├── controller/           # REST 控制器 (6个)
│   ├── service/              # 业务逻辑层
│   ├── mapper/               # MyBatis 数据访问层
│   ├── entity/               # JPA 实体类
│   ├── dto/                  # 数据传输对象
│   ├── vo/                   # 视图对象
│   ├── common/               # 公共组件
│   │   ├── annotation/       # @AuditLog 等注解
│   │   ├── constant/         # 常量定义
│   │   ├── enums/            # 枚举类
│   │   ├── exception/        # 异常处理
│   │   └── util/             # 工具类
│   ├── aspect/               # AOP 切面
│   └── config/               # 配置类
├── src/main/resources/
│   └── application.yml       # Spring Boot 配置
├── pom.xml                   # Maven 依赖
├── docker/                   # Docker 配置
├── helm/                     # Kubernetes Helm Chart
└── migrations/               # DDL 迁移脚本
```

### 目标 (Go-Zero - 待实现)
```
standardization-backend/
├── api/                      # API 服务
│   ├── doc/                  # API 定义 (.api 文件)
│   ├── etc/                  # 配置文件
│   └── internal/             # 内部实现
│       ├── handler/          # 请求处理 (参数校验)
│       ├── logic/            # 业务逻辑
│       ├── svc/              # 服务上下文
│       └── types/            # 类型定义
├── model/                    # 数据模型
├── pkg/                      # 公共包
│   ├── mq/                   # 消息队列 (Kafka/NSQ)
│   ├── cache/                # Redis 缓存
│   └── excel/                # Excel 处理
├── migrations/               # DDL 迁移
├── specs/                    # SDD 规格文档
└── .specify/                 # Spec Kit 配置
```

## 快速命令

### Java (当前代码库)

```bash
# 构建
mvn clean package                # 编译打包
mvn clean package -DskipTests    # 跳过测试打包

# 运行
java -jar target/standardization-web-0.0.1-SNAPSHOT.jar

# 测试
mvn test                         # 运行测试

# API 文档
# Swagger: http://localhost:8888/swagger-ui/
# Knife4j: http://localhost:8888/doc.html
```

### Go (目标架构 - 待实现)

```bash
# 开发
goctl api go -api api/doc/api.api -dir api/ --style=go_zero  # 生成 API 代码
goctl api plugin -plugin goctl-swagger="swagger -filename api.json" -api api/doc/api.api -dir .  # 生成 Swagger
go run api/main.go                                                         # 运行服务
go test ./...                                                              # 运行测试

# 部署
docker build -t standardization-backend:latest -f docker/Dockerfile .      # 构建镜像
helm upgrade --install standardization ./helm/standardization               # 部署到 K8s
```

## SDD 工作流程

本项目遵循 Spec-Driven Development 5 阶段工作流:

1. **Context** - 阅读 `.specify/memory/constitution.md` 理解项目规范
2. **Specify** - 使用 EARS 格式定义需求 → `specs/<feature>/spec.md`
3. **Design** - 创建技术方案 → `specs/<feature>/plan.md`
4. **Tasks** - 拆分任务 (每个 <50 行) → `specs/<feature>/tasks.md`
5. **Implement** - 按任务顺序编码实现

**重要**: 每个阶段完成后等待用户确认，再进入下一阶段。

## 架构规范

### 分层职责 (严格遵守)

| 层 | 职责 | 禁止 |
|---|------|------|
| Handler | 参数绑定、校验、调用 Logic | 包含业务逻辑 |
| Logic | 业务逻辑、事务管理 | 直接操作 HTTP |
| Model | 数据访问 (GORM/SQLx) | 包含业务逻辑 |

### API 设计

- 入口文件: `api/doc/api.api`
- 基础类型: `api/doc/base.api`
- 模块 API: `api/doc/<module>/<module>.api`
- 使用 `goctl api go` 生成代码

## 编码约定

### 命名规范

```
通用文件名: snake_case.go
包名:       lowercase
结构体:     PascalCase
方法:       PascalCase
变量:       camelCase
常量:       UPPER_SNAKE_CASE
```

### ⚠️ Go-Zero 文件命名特殊规则

**重要**: Go-Zero 框架生成的代码文件名必须与结构体名称完全匹配 (使用小写形式,不加下划线)

| 结构体名称 | 正确文件名 | ❌ 错误文件名 |
|-----------|-----------|--------------|
| `DeleteCatalogLogic` | `deletecataloglogic.go` | `delete_catalog_logic.go` |
| `CreateCatalogHandler` | `createcataloghandler.go` | `create_catalog_handler.go` |
| `QueryCatalogLogic` | `querycataloglogic.go` | `query_catalog_logic.go` |

**规则**:
- ✅ **保持 `goctl` 生成的文件名不变**
- ❌ **禁止重命名为带下划线的形式**
- 文件名 = 结构体名转小写 (无分隔符)

**原因**: `goctl` 工具严格按照此约定生成代码,重命名会导致后续生成代码时出现重复文件。

### 错误处理

```go
import "github.com/jinguoxing/idrm-go-base/errorx"

// 使用预定义错误码
if user == nil {
    return nil, errorx.NewWithCode(errorx.ErrCodeNotFound)
}

// 自定义业务错误码 (在 internal/errorx/codes.go 定义)
if user.Status == 0 {
    return nil, errorx.New(30102, "用户已禁用")
}
```

### 日志规范

```go
// 使用 logx，包含 traceId
logx.WithContext(ctx).Infof("user login: %s", phone)
```

## 重要约束

### 必须

- ✅ Handler 使用 validator 校验参数
- ✅ Logic 层管理事务边界
- ✅ 使用配置文件管理环境变量
- ✅ 错误信息使用 errors.Wrapf 包装

### 禁止

- ❌ Handler 直接操作数据库
- ❌ Model 层包含业务判断
- ❌ 硬编码配置值
- ❌ 使用 fmt.Println 替代 logx

## 相关文档

### SDD 规范
- 项目宪法: `.specify/memory/constitution.md`
- SDD 模板: `.specify/templates/`
- Spec Kit 命令: `.claude/commands/speckit.*.md`

### 迁移文档
- Java → Go 迁移规格: `specs/java-to-go-migration/spec.md`
- 迁移技术方案: `specs/java-to-go-migration/plan.md`
- 迁移任务清单: `specs/java-to-go-migration/tasks.md`

## 常见操作

### 新增 API 接口

1. 在 `api/doc/<module>/` 创建 `.api` 文件
2. 运行 `make api` 生成代码
3. 在 `api/internal/logic/` 实现业务逻辑
4. 运行 `make swagger` 更新文档

### 新增数据表

1. 在 `migrations/` 创建 DDL 文件
2. 执行 DDL 创建表
3. 使用 goctl 生成 Model 或手写 GORM Model
4. 在 Logic 层调用 Model

### 部署服务

#### Java (当前)
```bash
# Docker
docker build -t standardization-backend:latest -f docker/Dockerfile .

# Kubernetes
helm install standardization ./helm/standardization
helm upgrade standardization ./helm/standardization
helm status standardization
```

#### Go (目标 - 待实现)
```bash
# Docker
docker build -t standardization-backend:latest -f docker/Dockerfile .

# Kubernetes
helm upgrade --install standardization ./helm/standardization
```
