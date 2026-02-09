# standardization-backend

> 基于 Go-Zero 框架构建的微服务标准化管理后端系统。

## 概述

这是一个标准化管理后端服务，提供全面的数据标准化能力，包括：

- **编码规则管理** (`rule-api`) - 编码规则管理
- **目录管理** (`catalog-api`) - 标准目录层级管理
- **标准文件管理** (`stdfile-api`) - 标准文档/文件管理
- **码表管理** (`dict-api`) - 代码字典管理
- **数据元管理** (`dataelement-api`) - 数据元全生命周期管理
- **标准任务管理** (`task-api`) - 标准任务管理

## 技术栈

| 组件 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.24+ |
| 框架 | Go-Zero | v1.9+ |
| 数据库 | MySQL | 8.0 |
| ORM | GORM + SQLx | latest |
| 消息队列 | Kafka | 3.0 |
| API 文档 | Swagger/OpenAPI | 3.0 |

## 项目结构

```
standardization-backend/
├── api/                          # API 网关服务
│   ├── doc/                      # API 定义 (.api 文件)
│   │   ├── api.api              # 主 API 入口
│   │   ├── base.api             # 通用类型定义
│   │   ├── rule/                # 规则模块 API
│   │   ├── catalog/             # 目录模块 API
│   │   ├── stdfile/             # 标准文件模块 API
│   │   ├── dict/                # 码表模块 API
│   │   ├── dataelement/         # 数据元模块 API
│   │   └── task/                # 任务模块 API
│   ├── internal/
│   │   ├── handler/             # HTTP 请求处理器
│   │   ├── logic/               # 业务逻辑
│   │   ├── svc/                 # 服务上下文
│   │   ├── types/               # 类型定义
│   │   └── errorx/              # 错误处理
│   └── etc/                     # 配置文件
├── model/                        # 数据模型
├── migrations/                   # 数据库迁移脚本
├── specs/                        # SDD 规格文档
│   ├── 1-rule-api/
│   ├── 2-catalog-api/
│   ├── 3-std-file-api/
│   ├── 4-dict-api/
│   ├── 5-std-task-api/
│   └── 6-std-dataelement-api/
├── deploy/                       # 部署配置
│   ├── docker/
│   └── k8s/
├── Makefile                      # 构建自动化
├── CLAUDE.md                     # 项目指南
└── go.mod                        # Go 模块定义
```

## 快速开始

### 前置要求

- Go 1.24 或更高版本
- MySQL 8.0
- Kafka 3.0 (可选，用于异步消息)
- `goctl` CLI 工具

### 安装步骤

```bash
# 克隆仓库
git clone https://github.com/tianyuliang/standardization-backend.git
cd standardization-backend

# 安装依赖
make deps

# 生成 API 代码
make api

# 运行服务
make run
```

API 服务将在 `http://localhost:8888` 启动

### 使用 Docker

```bash
# 构建 Docker 镜像
make docker-build

# 运行容器
make docker-run
```

### 使用 Kubernetes

```bash
# 部署到开发环境
make k8s-deploy-dev

# 部署到生产环境
make k8s-deploy-prod

# 检查状态
make k8s-status
```

## 可用命令

```bash
# 开发命令
make init          # 初始化项目
make api           # 使用 goctl 生成 API 代码
make swagger       # 生成 Swagger 文档
make gen           # 生成 API 代码 + Swagger 文档
make fmt           # 格式化代码
make lint          # 运行代码检查
make test          # 运行测试
make build         # 构建二进制文件
make run           # 运行服务器
make clean         # 清理构建产物
make deps          # 安装依赖

# Docker 命令
make docker-build  # 构建 Docker 镜像
make docker-run    # 运行 Docker 容器
make docker-stop   # 停止 Docker 容器
make docker-push   # 推送 Docker 镜像

# Kubernetes 命令
make k8s-deploy    # 部署到 K8s (默认: dev)
make k8s-manifest  # 查看生成的清单
make k8s-delete    # 删除 K8s 部署
make k8s-status    # 检查 K8s 状态
```

## API 文档

### Swagger UI

启动服务后，访问 Swagger 文档：
- **JSON**: http://localhost:8888/swagger/swagger.json
- **YAML**: http://localhost:8888/swagger/swagger.yaml

### 模块 API

| 模块 | 前缀 | 描述 |
|------|------|------|
| Rule | `/api/v1/rule` | 编码规则 CRUD 操作 |
| Catalog | `/api/v1/catalog` | 目录层级管理 |
| StdFile | `/api/v1/stdfile` | 标准文件管理 |
| Dict | `/api/v1/dict` | 码表管理 |
| DataElement | `/api/v1/dataelement` | 数据元操作 |
| Task | `/api/v1/task` | 任务管理 |

## 开发指南

### 代码生成

始终使用 `goctl` 生成 API 代码：

```bash
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

### 测试

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 代码检查

```bash
# 格式化代码
make fmt

# 运行 linter
make lint
```

## 规格文档

每个模块遵循规格驱动开发 (SDD) 工作流：

1. **Context** - 项目章程和指南
2. **Specify** - 需求规格说明 (EARS 格式)
3. **Design** - 技术设计 (API/DDL/Model)
4. **Tasks** - 实现任务
5. **Implement** - 代码实现

详见 `specs/` 目录。

## 贡献指南

1. 遵循 [SDD 工作流](.specify/workflows/)
2. 从 `master` 分支创建功能分支
3. 提交 PR 前确保所有测试通过
4. 遵循 Go 编码规范

## 开源协议

MIT License - 详见 LICENSE 文件

## 联系方式

- 项目负责人: 田玉亮
- 代码仓库: https://github.com/tianyuliang/standardization-backend
