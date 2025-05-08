[English](../README.md) | [中文](README_zh.md) | [开发指南](./develop_guide_zh.md) | [开发计划](./plans_zh.md)

# Terminal

**统一、安全、可审计的命令行访问**

Terminal 是一个全面的平台，旨在简化和保护跨各种基础设施的命令行操作，包括物理/虚拟服务器（Linux、MacOS）、Kubernetes 集群（内部/外部集群）、Docker 和 Containerd 环境。

它提供了一个集中的网关用于：

- **统一访问：** 通过单一界面（Web 和 TUI）连接到各种后端。
- **增强安全性：** 强大的身份验证（TOTP、OIDC 等）、细粒度的授权（基于 JWT 的 RBAC）以及可选的端到端会话加密。
- **全面审计：** 详细的操作日志和会话记录，实现完全可见性和合规性。
- **运营效率：** 命令模板、批量命令执行以及代理到内部服务的 API。

## 核心组件

1.  **客户端 (Client)：** 用户交互界面，包括基于 HTMX + xterm.js 的 Web 界面和基于 Bubbletea 的 TUI (文本用户界面)。
2.  **管理服务器 (Management Server)：** 中央控制平面，使用 Go 和 NATS 构建。负责身份验证、授权、Agent 管理、日志记录、会话记录和配置管理。
3.  **Agent：** 轻量级守护进程，使用 Go 编写，部署在目标系统上，用于安全地暴露终端访问、执行命令，并与管理服务器通信。

## 主要功能

- **多种终端类型：** 本地 Shell、SSH 连接、Kubernetes Pod、Docker 容器、Containerd 容器。
- **灵活的操作模式：** 交互式 PTY 会话、一次性命令执行、服务器端流式输出。
- **身份验证：** 支持 NoAuth (不推荐生产)、用户名/密码、TOTP、Email 验证、OIDC (如 Keycloak, Okta)、非对称密钥对。
- **授权：** 基于 JWT 的角色访问控制 (RBAC)，支持精细化权限配置。
- **会话记录：** Agent 端可记录终端会话（例如，Asciinema 格式），记录元数据由管理服务器存储和索引。
- **命令模板：** 定义和重用常见的命令模式，提高操作效率。
- **API 代理：** (规划中) 通过 Terminal 代理安全地暴露内部 API。
- **批量命令执行：** 同时在多个目标 Agent 上运行命令。
- **端到端加密 (E2EE)：** (规划中) 可选的会话加密保护，防止管理服务器窥探会话内容。
- **连接模式：**
  - Client <-> Agent (通过 NATS 直连，需网络可达)
  - Client <-> Management Server (代理) <-> Agent (通过 NATS)
  - 支持 WebSocket, Server-Sent Events (SSE), HTTP 进行客户端通信。
- **开发模式：** 捆绑管理服务器 + Agent，内嵌 SQLite 和 NATS，简化本地开发和测试。

## 技术栈

- **后端和 Agent：** Go
  - API 框架: [Gin](https://github.com/gin-gonic/gin)
  - 测试驱动开发 (TDD): [ginkgo](https://github.com/onsi/ginkgo) & [gomega](https://github.com/onsi/gomega)
  - ORM: [Ent](https://github.com/ent/ent)
- **Web 前端：** [HTMX](https://github.com/bigskysoftware/htmx), [xterm.js](https://xtermjs.org/)
- **TUI 客户端：** [bubbletea](https://github.com/charmbracelet/bubbletea)
- **数据库：** PostgreSQL (生产推荐), SQLite (开发模式默认)
- **消息队列/存储：** [Nats](https://github.com/nats-io/nats-server) (Pub/Sub, Key-Value Store, Object Storage)
- **通信协议：**
  - Agent 与 Management Server：基于 Nats Pub/Sub
  - Client 与 Agent：主要基于 Nats Pub/Sub (可由 Management Server 代理)
  - Client 与 Management Server：HTTP/HTTPS API (Gin) 和 Nats Pub/Sub

## 快速入门

1.  **获取源码：**

    ```bash
    git clone https://github.com/<your-username>/terminal.git # 替换为实际仓库地址
    cd terminal
    ```

2.  **运行开发模式：**
    此模式将启动一个集成的服务，包含管理端、Agent 功能，并使用内嵌的 SQLite 和 NATS。

    ```bash
    make run-dev
    ```

    (如果 `make run-dev` 不可用，请查阅[开发指南](./develop_guide.md#23-构建与运行指南)中关于启动开发模式的说明，例如 `go run cmd/terminal/main.go --dev-mode`)

3.  **访问客户端：**

    - Web 客户端：打开浏览器访问 `http://localhost:<port>` (具体端口请查看启动日志，默认为 8080 或类似)。
    - TUI 客户端：(如果已实现) `./bin/terminal-cli` 或类似命令。

4.  **更多信息：**
    - 详细的开发环境搭建、构建和运行指南，请参考 [开发指南](./develop_guide.md#2-环境搭建与构建-setup--build)。
    - 了解项目组件和代码结构，请参考 [开发指南](./develop_guide.md)。

## 路线图

项目的详细开发计划和各版本目标，请参见 [开发计划](./plans_zh.md)。

## 贡献

我们欢迎各种形式的贡献！无论是提交 Bug 报告、提出功能建议，还是贡献代码。
请参考我们的 [贡献指南](./develop_guide_zh.md#6-贡献指南-contribution-guide) 了解如何参与。

## 许可证

本项目采用 [MIT License](../LICENSE) 授权。
