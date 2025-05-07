# Terminal Developer Guide

- [Terminal Developer Guide](#terminal-developer-guide)
  - [1. 项目概述 (Project Overview)](#1-项目概述-project-overview)
    - [1.1 项目愿景与目标](#11-项目愿景与目标)
    - [1.2 核心架构](#12-核心架构)
    - [1.3 技术栈总览](#13-技术栈总览)
    - [1.4 设计文档](#14-设计文档)
  - [2. 环境搭建与构建 (Setup \& Build)](#2-环境搭建与构建-setup--build)
    - [2.1 开发环境要求](#21-开发环境要求)
    - [2.2 获取源码](#22-获取源码)
    - [2.3 构建与运行指南](#23-构建与运行指南)
      - [2.3.1 开发模式 (Dev Mode)](#231-开发模式-dev-mode)
      - [2.3.2 生产构建](#232-生产构建)
  - [3. 核心功能模块](#3-核心功能模块)
    - [3.1 管理端 (Management Server)](#31-管理端-management-server)
    - [3.2 Agent 端 (Agent)](#32-agent-端-agent)
    - [3.3 客户端 (Client)](#33-客户端-client)
      - [3.3.1 Web 客户端 (Web Client)](#331-web-客户端-web-client)
      - [3.3.2 TUI 客户端 (TUI Client)](#332-tui-客户端-tui-client)
  - [4. 代码结构与规范 (Code Structure \& Conventions)](#4-代码结构与规范-code-structure--conventions)
    - [4.1 主要目录结构说明](#41-主要目录结构说明)
    - [4.2 编码规范](#42-编码规范)
    - [4.3 分支与合并策略](#43-分支与合并策略)
  - [5. 测试策略 (Testing Strategy)](#5-测试策略-testing-strategy)
    - [5.1 单元测试](#51-单元测试)
  - [6. 贡献指南 (Contribution Guide)](#6-贡献指南-contribution-guide)
  - [7. 术语表 (Glossary)](#7-术语表-glossary)

## 1. 项目概述 (Project Overview)

### 1.1 项目愿景与目标

Terminal 致力于打造统一的终端访问与管理平台，支持多种环境（服务器、Kubernetes、容器等），提升运维与开发效率，保障安全与可审计性。

### 1.2 核心架构

项目主要由三个核心组件构成：管理端 (Management Server)、Agent 端 (Agent) 和客户端 (Client)。

- **管理端 (Management Server)**: 作为中央控制平面，负责用户认证、授权管理 (RBAC)、Agent 注册与管理、会话审计日志、配置管理以及作为 Client 与 Agent 之间通信的协调者（或代理）。
- **Agent 端 (Agent)**: 轻量级守护进程，部署在目标系统（物理机、虚拟机、K8s Node 等）。负责与管理端建立安全连接，接收指令，执行命令，管理本地终端会话（如 PTY），并上报操作信息。
- **客户端 (Client)**: 用户与 Terminal 系统交互的界面，目前包括 Web 客户端和 TUI 客户端。客户端通过管理端进行认证并获取可访问的 Agent 列表，然后与目标 Agent 建立会话。

**通信流程概要:**

1.  Agent 启动后向 Management Server 注册。
2.  Client 通过 Management Server 进行认证和授权。
3.  Client 从 Management Server 获取可连接的 Agent 列表。
4.  Client 与 Agent 之间的会话数据流主要通过 NATS Pub/Sub 进行。在某些连接模式下，Management Server 可能充当数据流的代理。
5.  Management Server 记录所有操作和会话元数据以供审计。

```
+-------------+ (API / NATS) +---------------+ (NATS) +---------+
|    Client   | <---> | Management Server | <---> |   Agent     |
| (Web / TUI) |       |(Auth, RBAC, Audit)|       | (on Target) |
+-------------+ +---------------------------------+ +-----------+
```

### 1.3 技术栈总览

- 后端 (Management Server, Agent): Go
  - API 框架: [Gin](https://github.com/gin-gonic/gin)
  - 测试驱动开发 (TDD): [ginkgo](https://github.com/onsi/ginkgo) & [gomega](https://github.com/onsi/gomega)
  - ORM: [Ent](https://github.com/ent/ent)
- 前端：
  - Web 客户端: [HTMX](https://github.com/bigskysoftware/htmx), [xterm.js](https://xtermjs.org/)
  - TUI 客户端: [bubbletea](https://github.com/charmbracelet/bubbletea)
- 数据库：PostgreSQL (生产推荐), SQLite (开发模式默认)
- 消息队列/存储：[Nats](https://github.com/nats-io/nats-server) (Pub/Sub, KV, Object Storage)
- 通信协议：
  - Agent 与 Management Server：基于 Nats Pub/Sub
  - Client 与 Agent：主要基于 Nats Pub/Sub (可由 Management Server 代理)
  - Client 与 Management Server：HTTP/HTTPS API (Gin) 和 Nats Pub/Sub

### 1.4 设计文档

- [API 设计](./designs/api.md) (建议补充或创建此文档)
- [Schema 设计](./designs/schema.md) (建议补充或创建此文档，涵盖数据库和 NATS KV/Object Store 结构)
- [Terminal 设计](./designs/terminal.md) (建议补充或创建此文档，可能包含更详细的会话流程等)

## 2. 环境搭建与构建 (Setup & Build)

### 2.1 开发环境要求

- Go 版本：建议 1.24 及以上
- 操作系统：MacOS、Linux
- 依赖工具：
  - `make`
  - `git`
  - `docker` 和 `docker-compose` (可选，用于运行 PostgreSQL, NATS 等依赖服务)
  - `golangci-lint` (用于代码检查)

### 2.2 获取源码

```bash
git clone https://github.com/<your-username>/terminal.git # 替换为实际仓库地址
cd terminal
```

### 2.3 构建与运行指南

#### 2.3.1 开发模式 (Dev Mode)

Dev 模式旨在简化本地开发和测试，默认使用 SQLite 作为数据库，并启动内嵌的 NATS 服务。管理端与 Agent 功能通常会合并在单一进程中运行，或通过 make 脚本同时启动。

运行命令:

```bash
make run-dev
# 或者根据项目实际情况
# go run cmd/terminal/main.go --dev-mode
```

访问 Web UI (如果已实现): http://localhost:<port> (具体端口查看启动日志)

#### 2.3.2 生产构建

生产构建会生成独立的二进制文件，用于部署。

构建命令:

```bash
make build # 构建所有组件 (Management Server, Agent)
# 或单独构建
# make build-server
# make build-agent
```

构建产物通常在 bin/ 或 \_build/ 目录下。

## 3. 核心功能模块

### 3.1 管理端 (Management Server)

管理端是 Terminal 平台的大脑，主要职责包括：

- 用户认证与授权: 支持多种认证方式 (TOTP, OIDC 等)，实现基于角色的访问控制 (RBAC)。
- Agent 管理: 负责 Agent 的注册、心跳、状态监控和配置下发。
- 会话调度与代理: 协调 Client 与 Agent 之间的连接，并可在需要时代理会话数据流。
- 审计与日志: 记录用户操作、系统事件和会话元数据，支持会话录像存储索引。
- API 服务: 提供 RESTful API 供客户端和外部系统集成。
- 配置管理: 统一管理系统配置、命令模板等。

### 3.2 Agent 端 (Agent)

Agent 端部署在需要被管理的远端目标机器上，主要职责包括：

- 安全通信: 与管理端建立基于 NATS 的安全长连接。
- 终端会话管理: 创建和管理 PTY (伪终端) 会话，支持交互式 Shell。
- 命令执行: 执行来自 Client 的一次性命令或来自管理端的批量命令。
- 本地资源访问: 根据授权访问本地系统资源 (如 Docker, Containerd, Kubernetes API)。
- 会话录制: (可选) 在 Agent 端进行终端会话录制 (如 Asciinema 格式)。
- 状态上报: 定期向管理端上报自身状态和监控信息。

### 3.3 客户端 (Client)

客户端是用户与 Terminal 系统交互的入口。

#### 3.3.1 Web 客户端 (Web Client)

- 技术栈: HTMX, xterm.js, Vanilla JS/CSS。
- 功能: 提供图形化界面，用户可以通过浏览器登录、管理和连接到目标终端，查看审计日志等。
- 交互: 通过 HTTP API 与 Management Server 通信获取数据和配置，通过 WebSocket 或 NATS (经由 Management Server 代理或直连) 与 Agent 进行终端会话交互。

#### 3.3.2 TUI 客户端 (TUI Client)

- 技术栈: Go, [bubbletea](https://github.com/charmbracelet/bubbletea)。
- 功能: 提供文本用户界面，方便习惯命令行的用户快速连接和操作。
- 交互: 通常直接通过 HTTP API 和 NATS 与 Management Server 及 Agent 通信。

## 4. 代码结构与规范 (Code Structure & Conventions)

### 4.1 主要目录结构说明

一个推荐的 Go 项目结构可能如下：

- main.go: (如果使用单一入口) 主程序入口。
- cmd/: 存放各个应用的主程序 (main package)。
  - server/: 管理端 (Management Server) 入口。
  - agent/: Agent 端入口。
  - cli/: TUI 客户端或其他命令行工具入口。
- server/: 管理端核心逻辑。
- agent/: Agent 端核心逻辑。
- client/: TUI 客户端核心逻辑。
- pkg/: 内部共享的平台代码，如数据库交互、NATS 客户端封装、配置加载等。
  - auth/: 认证授权相关模块。
  - store/: 数据存储抽象及实现。
- api/: API 定义 (例如 Protobuf, OpenAPI/Swagger 规范文件)。
- web/: Web 客户端静态资源 (HTML, JS, CSS, images)。
  - htmx/ or static/
- docs/: 项目文档。
  - designs/: 设计文档。
- configs/: 示例配置文件。
- scripts/: 构建、部署等辅助脚本。
- build/: (或 Makefile 中定义的) 构建相关文件，如 Dockerfile。
- test/: 端到端测试或集成测试相关文件。

### 4.2 编码规范

- 遵循 Go 官方风格指南。
- 使用 gofmt 或 goimports 自动格式化代码。
- 使用 golangci-lint 进行静态代码检查，配置文件 (.golangci.yml) 应纳入版本控制。
- 统一格式化命令：make fmt
- 静态检查命令：make lint
- 编写清晰、可维护的注释。
- 错误处理：遵循 Go 的错误处理最佳实践，避免 panic (除非在 main 或不可恢复的场景)。

### 4.3 分支与合并策略

- 主分支: main (或 master)，保持稳定，对应已发布或即将发布的版本。
- 开发分支: develop (可选)，作为日常开发的基础分支，新功能从此分支创建 feature 分支。
- 功能分支: feature/xxx (例如 feature/e2e-encryption)，用于开发新功能。开发完成后合并到 develop (或 main)。
- 修复分支: fix/xxx (例如 fix/login-bug)，用于修复 main 分支上的 bug。修复后通常会合并到 main 和 develop。
- 发布分支: release/vX.Y.Z (可选)，用于准备版本发布，进行最后的测试和文档更新。
- 合并方式:
  - 所有代码变更必须通过 Pull Request (PR) / Merge Request (MR) 提交。
  - PR 需要至少一名其他开发者进行 Code Review。
  - PR 需要通过所有 CI 检查 (linting, tests)。
  - 推荐使用 Squash and Merge 或 Rebase and Merge 以保持 main 分支的提交历史清晰。

## 5. 测试策略 (Testing Strategy)

### 5.1 单元测试

- 范围: 针对独立的函数、方法、模块或包进行测试。
- 工具: 使用 Go 内置的 testing 包，配合 ginkgo 和 gomega 进行 BDD 风格测试。
- 目标: 重点覆盖核心逻辑、边界条件和错误处理。模拟外部依赖 (如数据库、NATS) 以保证测试的隔离性和速度。
- 执行: make test-unit 或 go test ./...

## 6. 贡献指南 (Contribution Guide)

- 欢迎通过 GitHub Issues 报告 Bug、提出功能建议或参与讨论。
- 贡献代码前，建议先创建一个 Issue 描述你的想法或要修复的问题，并与项目维护者讨论方案。
- Fork 项目仓库到你自己的账户。
- 基于 develop (或 main) 分支创建你的 feature/xxx 或 fix/xxx 分支。
- 遵循编码规范和代码结构。
- 为新增或修改的代码编写必要的单元测试和集成测试。
- 确保所有测试通过 (make test) 且代码通过 lint 检查 (make lint)。
- 如果有重要变更或新增功能，请更新相关文档 (README, docs/)。
- 提交 Pull Request 到上游仓库的 develop (或 main) 分支。
- 在 PR 描述中清晰说明变更内容、原因和测试情况。
- 积极参与 PR 的 Code Review 讨论，并根据反馈修改代码。

## 7. 术语表 (Glossary)

- Management Server (管理端): 中央控制服务，负责认证、授权、调度、审计等。
- Agent (代理端): 部署在目标环境的轻量级代理，负责执行命令和管理终端会话。
- Client (客户端): 用户与 Terminal 系统交互的界面，如 Web 客户端或 TUI 客户端。
- Web Client (Web 客户端): 基于浏览器的图形用户界面。
- TUI (Text User Interface, 文本用户界面): 基于命令行的用户界面，本项目使用 bubbletea 实现。
- Dev 模式 (开发模式): 一种简化的运行模式，通常管理端与 Agent 功能合并运行，使用内嵌数据库和消息队列，方便本地开发调试。
- NATS: 高性能消息队列系统，用于组件间的异步通信和数据流传输。
- PTY (Pseudo-Terminal, 伪终端): 模拟物理终端的软件接口，用于实现交互式 Shell 会话。
- E2EE (End-to-End Encryption, 端到端加密): 一种安全通信方式，确保数据在发送方加密，在接收方解密，中间节点（包括 Management Server）无法读取明文内容。
- RBAC (Role-Based Access Control, 基于角色的访问控制): 一种权限管理模型，通过为用户分配角色，再为角色授予权限，来实现灵活的访问控制。
- HTMX: 一种允许你通过 HTML 属性直接访问现代浏览器功能的库，用于构建动态 Web 界面。
- xterm.js: 一个功能强大的 Web 组件，用于在浏览器中实现功能齐全的终端模拟器。
- Gin: Go 语言的高性能 HTTP Web 框架。
- Ent: Go 语言的实体框架，用于数据建模和数据库交互。
