# Terminal 开发计划

- [Terminal 开发计划](#terminal-开发计划)
  - [v0.1 实现 MVP (最小可行产品)](#v01-实现-mvp-最小可行产品)
  - [v0.2 引入 Web 客户端与交互式会话](#v02-引入-web-客户端与交互式会话)
  - [v0.3 扩展终端类型支持](#v03-扩展终端类型支持)
  - [v0.4 核心认证、授权和审计](#v04-核心认证授权和审计)
  - [v0.5 引入命令模板](#v05-引入命令模板)
  - [v0.6 架构演进：管理端与 Agent 分离部署](#v06-架构演进管理端与-agent-分离部署)
  - [v0.7 增强功能与用户体验](#v07-增强功能与用户体验)
  - [v0.8 安全性增强：端到端加密](#v08-安全性增强端到端加密)
  - [v0.9 扩展集成与交互方式](#v09-扩展集成与交互方式)
  - [V1.0: 生产可用](#v10-生产可用)

## v0.1 实现 MVP (最小可行产品)

**目标**: 实现最小可行产品，验证核心的本地终端连接（一次性命令执行）和基础架构。

- [ ] **项目初始化与基础架构**
  - [ ] 项目结构搭建 (`cmd/`, `internal/`, `pkg/`, `api/` 等)
  - [ ] 引入核心技术栈 (Go, Gin, Ent, NATS client)
  - [ ] 测试驱动开发环境配置 (ginkgo, gomega)
  - [ ] Makefile 基础 (build, lint, fmt, test-unit, run-dev)
- [ ] **核心 API 与数据模型设计**
  - [ ] 初步数据库 Schema 设计 (SQLite & PostgreSQL)
  - [ ] Management Server 核心 API 初稿 (Client <-> Management, Agent <-> Management)
  - [ ] Agent 端核心 API 初稿 (Management <-> Agent)
  - [ ] NATS 主题与消息结构设计初稿
- [ ] **Dev 模式核心功能**
  - [ ] Management Server 基础启动，内嵌 NATS 与 SQLite
  - [ ] Agent 基础启动 (可与 Server 合并运行在 Dev 模式)
  - [ ] Agent 向 Management Server 注册
  - [ ] **一次性命令执行**: Client -> Management Server -> Agent (通过 NATS/API) -> 返回结果
    - 支持 Linux 本地 Shell 命令
- [ ] **TUI 客户端基础**
  - [ ] 实现 TUI 客户端基本框架 (Bubbletea)
  - [ ] TUI 连接 Management Server，列出已注册 Agent
  - [ ] TUI 支持在选定 Agent 上执行一次性命令并显示结果

## v0.2 引入 Web 客户端与交互式会话

**目标**: 引入 Web 客户端，支持交互式 PTY 会话，初步实现操作记录。

- [ ] **Web 客户端基础**
  - [ ] 实现 Web 客户端基础框架 (HTMX, xterm.js)
  - [ ] Web 客户端通过 API 与 Management Server 通信
  - [ ] Web 客户端展示 Agent 列表，支持执行一次性命令
- [ ] **交互式 PTY 会话**
  - [ ] Agent 端支持 PTY (伪终端) 创建与管理 (Linux)
  - [ ] 实现 PTY 数据流传输：Client <-> Management Server (代理) <-> Agent (通过 NATS 及 WebSocket/SSE)
  - [ ] Web 客户端 (xterm.js) 支持交互式终端会话
  - [ ] TUI 客户端支持交互式终端会话
- [ ] **Management Server 增强**
  - [ ] 基础操作记录存储 (命令执行历史)
  - [ ] Web 客户端展示操作记录

## v0.3 扩展终端类型支持

**目标**: Agent 端扩展对更多主流终端类型的支持。

- [ ] **Agent 端能力扩展**
  - [ ] 支持连接本地 Docker 容器 (exec, attach)
  - [ ] 支持连接本地 Containerd 容器 (exec, attach)
  - [ ] 支持连接 Kubernetes Pod 内的容器 (exec, attach, 需要 Kubeconfig 或 ServiceAccount)
  - [ ] 支持 MacOS 本地 Shell (PTY 和一次性命令)
  - [ ] (可选) 支持通过 Agent 建立 SSH 连接到其他主机

## v0.4 核心认证、授权和审计

**目标**: 引入正式的用户认证、基于角色的授权 (RBAC) 和关键安全审计功能。

- [ ] **Management Server 安全模块**
  - [ ] **认证机制设计与实现**:
    - [ ] 用户名/密码认证 (存储加密密码)
    - [ ] JWT (JSON Web Tokens) 作为会话令牌
    - [ ] (可选初期) TOTP 双因素认证
  - [ ] **授权机制设计与实现 (RBAC)**:
    - [ ] 定义用户、角色、权限模型
    - [ ] 实现基于角色的 API 访问控制
    - [ ] 实现基于角色的 Agent 访问控制 (哪些用户/角色可以访问哪些 Agent)
  - [ ] **关键安全审计功能**:
    - [ ] 记录用户登录、认证失败、权限变更等安全事件
    - [ ] 记录详细的操作审计日志 (谁、何时、何地、对何目标、执行了何操作)
- [ ] **Agent 端安全配合**
  - [ ] (可选) Agent 端会话录制 (如 Asciinema 格式)，上传至 Management Server 控制的存储 (NATS Object Store)
  - [ ] Management Server 提供会话录像回放功能 (Web 客户端)

## v0.5 引入命令模板

**目标**: 引入命令模板功能，简化常用操作，提高效率和标准化。

- [ ] **Management Server 功能**
  - [ ] 命令模板 CRUD API (支持参数化模板)
  - [ ] Web 客户端管理和使用命令模板
- [ ] **Agent 端功能**
  - [ ] 支持执行来自 Management Server 的参数化命令模板
- [ ] **Client 端功能**
  - [ ] TUI 和 Web 客户端支持列出、选择和执行命令模板

## v0.6 架构演进：管理端与 Agent 分离部署

**目标**: 优化架构，支持管理端与 Agent 完全分离独立部署，并为未来高级安全特性奠定基础。

- [ ] **部署模式强化**
  - [ ] Management Server 与 Agent 作为独立服务部署的文档和配置示例
  - [ ] Agent 自动注册到远程 Management Server 的健壮性提升
- [ ] **连接代理与路由**
  - [ ] 强化 Management Server 作为 Client 与 Agent 通信代理的能力 (应对复杂网络环境)
- [ ] **认证增强**
  - [ ] 支持 OIDC (OpenID Connect) 认证集成 (如 Keycloak, Okta)

## v0.7 增强功能与用户体验

**目标**: 提升用户体验，增加便捷功能。

- [ ] **通知机制**:
  - [ ] Agent 状态变更、重要事件通知 (如通过 Web UI, Email)
- [ ] **批量操作**:
  - [ ] 支持在多个选定 Agent 上批量执行命令或命令模板
- [ ] **文件传输**:
  - [ ] 支持在 Client 和 Agent 之间安全上传/下载文件 (通过 Management Server 中转或 NATS)
- [ ] **后台与定时任务**:
  - [ ] (探索性) 支持在 Agent 上执行长时间运行的后台任务
  - [ ] (探索性) 支持通过 Management Server 配置定时任务下发到 Agent 执行
- [ ] **无限画布/仪表盘**:
  - [ ] (探索性) Web 客户端提供可定制的仪表盘，展示 Agent 状态、常用连接等

## v0.8 安全性增强：端到端加密

**目标**: 引入更强的安全机制，特别是针对会话数据的端到端加密。

- [ ] **密钥对认证**:
  - [ ] 支持 Client 和/或 Agent 使用非对称密钥对进行身份验证
- [ ] **端到端加密 (E2EE) 设计与实现**:
  - [ ] Client 与 Agent 协商会话密钥
  - [ ] PTY 数据流在 Client 和 Agent 之间端到端加密，Management Server 无法解密会话内容
  - [ ] E2EE 模式下会话录制方案调整 (如 Client 端录制或加密录制)

## v0.9 扩展集成与交互方式

**目标**: 提供更多与其他系统集成和自动化交互的方式。

- [ ] **Bot 支持**:
  - [ ] (探索性) 提供 API 或 SDK 支持聊天机器人 (如 Slack, Teams) 与 Terminal 交互执行命令
- [ ] **Webhook 支持**:
  - [ ] Management Server 支持配置 Webhook，在特定事件发生时 (如 Agent 上线/离线) 通知外部系统
- [ ] **事件订阅**:
  - [ ] Client 或外部系统可通过 NATS 或 API 订阅 Management Server 发布的事件流

## V1.0: 生产可用

**目标**: 打磨产品，使其达到生产环境部署的要求。

- [ ] **文档完善**:
  - [ ] 用户手册、管理员手册、API 文档、部署指南
- [ ] **多租户支持**: (如果适用)
  - [ ] 数据隔离与权限控制，支持多个独立的组织或团队使用同一套系统
- [ ] **性能与稳定性**:
  - [ ] 进行压力测试、长时间稳定性测试
  - [ ] 性能瓶颈分析与优化 (高并发连接、大量 Agent 管理)
- [ ] **运维支持**:
  - [ ] 完善的监控指标 (Prometheus metrics)
  - [ ] 结构化日志输出，方便收集和分析
  - [ ] 告警机制集成
- [ ] **安全加固**:
  - [ ] 第三方安全审计或渗透测试
  - [ ] 依赖库安全扫描与更新
- [ ] **生产环境部署方案**:
  - [ ] 提供 Docker-compose, Kubernetes Helm Chart 等部署方案
