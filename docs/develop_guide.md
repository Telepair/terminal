# Terminal Developer Guide

- [Terminal Developer Guide](#terminal-developer-guide)
  - [1. Project Overview](#1-project-overview)
    - [1.1 Vision and Goals](#11-vision-and-goals)
    - [1.2 Core Architecture](#12-core-architecture)
    - [1.3 Technology Stack Overview](#13-technology-stack-overview)
    - [1.4 Design Documents](#14-design-documents)
  - [2. Setup \& Build](#2-setup--build)
    - [2.1 Development Environment Requirements](#21-development-environment-requirements)
    - [2.2 Getting the Source Code](#22-getting-the-source-code)
    - [2.3 Build \& Run Guide](#23-build--run-guide)
      - [2.3.1 Development Mode](#231-development-mode)
      - [2.3.2 Production Build](#232-production-build)
  - [3. Core Functional Modules](#3-core-functional-modules)
    - [3.1 Management Server](#31-management-server)
    - [3.2 Agent](#32-agent)
    - [3.3 Client](#33-client)
      - [3.3.1 Web Client](#331-web-client)
      - [3.3.2 TUI Client](#332-tui-client)
  - [4. Code Structure \& Conventions](#4-code-structure--conventions)
    - [4.1 Main Directory Structure](#41-main-directory-structure)
    - [4.2 Coding Conventions](#42-coding-conventions)
    - [4.3 Branching \& Merging Strategy](#43-branching--merging-strategy)
  - [5. Testing Strategy](#5-testing-strategy)
    - [5.1 Unit Testing](#51-unit-testing)
  - [6. Contribution Guide](#6-contribution-guide)
  - [7. Glossary](#7-glossary)

## 1. Project Overview

### 1.1 Vision and Goals

Terminal aims to build a unified platform for terminal access and management, supporting various environments (servers, Kubernetes, containers, etc.), improving DevOps efficiency, and ensuring security and auditability.

### 1.2 Core Architecture

The project consists of three core components: Management Server, Agent, and Client.

- **Management Server**: Acts as the central control plane, responsible for user authentication, RBAC, agent registration and management, session audit logs, configuration management, and as a coordinator (or proxy) for communication between Client and Agent.
- **Agent**: A lightweight daemon deployed on target systems (physical machines, VMs, K8s nodes, etc.). Responsible for establishing secure connections with the management server, receiving instructions, executing commands, managing local terminal sessions (e.g., PTY), and reporting operation information.
- **Client**: The user interface for interacting with the Terminal system, currently including Web and TUI clients. Clients authenticate via the management server, obtain the list of accessible agents, and then establish sessions with target agents.

**Communication Flow Overview:**

1.  Agent registers with the Management Server upon startup.
2.  Client authenticates and is authorized via the Management Server.
3.  Client obtains the list of connectable Agents from the Management Server.
4.  Session data between Client and Agent mainly flows through NATS Pub/Sub. In some modes, the Management Server may act as a proxy for the data stream.
5.  The Management Server records all operations and session metadata for auditing.

```
+-------------+ (API / NATS) +---------------+ (NATS) +---------+
|    Client   | <---> | Management Server | <---> |   Agent     |
| (Web / TUI) |       |(Auth, RBAC, Audit)|       | (on Target) |
+-------------+ +---------------------------------+ +-----------+
```

### 1.3 Technology Stack Overview

- Backend (Management Server, Agent): Go
  - API Framework: [Gin](https://github.com/gin-gonic/gin)
  - Test Driven Development (TDD): [ginkgo](https://github.com/onsi/ginkgo) & [gomega](https://github.com/onsi/gomega)
  - ORM: [Ent](https://github.com/ent/ent)
- Frontend:
  - Web Client: [HTMX](https://github.com/bigskysoftware/htmx), [xterm.js](https://xtermjs.org/)
  - TUI Client: [bubbletea](https://github.com/charmbracelet/bubbletea)
- Database: PostgreSQL (recommended for production), SQLite (default for development)
- Message Queue/Storage: [Nats](https://github.com/nats-io/nats-server) (Pub/Sub, KV, Object Storage)
- Communication Protocols:
  - Agent & Management Server: Based on Nats Pub/Sub
  - Client & Agent: Mainly based on Nats Pub/Sub (can be proxied by Management Server)
  - Client & Management Server: HTTP/HTTPS API (Gin) and Nats Pub/Sub

### 1.4 Design Documents

- [API Design](./designs/api.md) (Recommended to supplement or create this document)
- [Schema Design](./designs/schema.md) (Recommended to supplement or create this document, covering database and NATS KV/Object Store structures)
- [Terminal Design](./designs/terminal.md) (Recommended to supplement or create this document, may include more detailed session flows, etc.)

## 2. Setup & Build

### 2.1 Development Environment Requirements

- Go version: 1.24 or above recommended
- OS: MacOS, Linux
- Dependencies:
  - `make`
  - `git`
  - `docker` and `docker-compose` (optional, for running PostgreSQL, NATS, etc.)
  - `golangci-lint` (for code linting)

### 2.2 Getting the Source Code

```bash
git clone https://github.com/<your-username>/terminal.git # Replace with actual repo URL
cd terminal
```

### 2.3 Build & Run Guide

#### 2.3.1 Development Mode

Dev mode simplifies local development and testing, using SQLite as the default database and starting an embedded NATS service. Management and Agent functions are usually merged into a single process or started together via make scripts.

Run command:

```bash
make run-dev
# Or depending on the project
# go run cmd/terminal/main.go --dev-mode
```

Access Web UI (if implemented): http://localhost:<port> (check startup logs for the actual port)

#### 2.3.2 Production Build

Production builds generate standalone binaries for deployment.

Build command:

```bash
make build # Build all components (Management Server, Agent)
# Or build separately
# make build-server
# make build-agent
```

Build artifacts are usually in the bin/ or \_build/ directory.

## 3. Core Functional Modules

### 3.1 Management Server

The Management Server is the brain of the Terminal platform, mainly responsible for:

- User authentication and authorization: Supports multiple authentication methods (TOTP, OIDC, etc.), implements RBAC.
- Agent management: Handles agent registration, heartbeat, status monitoring, and configuration delivery.
- Session scheduling and proxy: Coordinates connections between Client and Agent, and can proxy session data streams if needed.
- Audit and logging: Records user operations, system events, and session metadata, supports session recording index.
- API service: Provides RESTful APIs for clients and external integration.
- Configuration management: Unified management of system configs, command templates, etc.

### 3.2 Agent

The Agent is deployed on remote target machines to be managed, mainly responsible for:

- Secure communication: Establishes secure long connections with the management server via NATS.
- Terminal session management: Creates and manages PTY (pseudo-terminal) sessions, supports interactive shells.
- Command execution: Executes one-off commands from the Client or batch commands from the management server.
- Local resource access: Accesses local system resources (e.g., Docker, Containerd, Kubernetes API) as authorized.
- Session recording: (Optional) Records terminal sessions on the Agent side (e.g., Asciinema format).
- Status reporting: Periodically reports its status and monitoring info to the management server.

### 3.3 Client

The Client is the entry point for users to interact with the Terminal system.

#### 3.3.1 Web Client

- Tech stack: HTMX, xterm.js, Vanilla JS/CSS.
- Features: Provides a graphical interface for users to log in, manage, and connect to target terminals, view audit logs, etc.
- Interaction: Communicates with the Management Server via HTTP API to get data and configs, and interacts with the Agent for terminal sessions via WebSocket or NATS (proxied by Management Server or direct).

#### 3.3.2 TUI Client

- Tech stack: Go, [bubbletea](https://github.com/charmbracelet/bubbletea).
- Features: Provides a text user interface for command-line users to quickly connect and operate.
- Interaction: Usually communicates directly with the Management Server and Agent via HTTP API and NATS.

## 4. Code Structure & Conventions

### 4.1 Main Directory Structure

A recommended Go project structure may look like:

- main.go: (if using a single entry) Main program entry.
- cmd/: Main programs (main package) for each application.
  - server/: Management Server entry.
  - agent/: Agent entry.
  - cli/: TUI client or other CLI tool entry.
- server/: Core logic of the management server.
- agent/: Core logic of the agent.
- client/: Core logic of the TUI client.
- pkg/: Shared platform code, e.g., DB interaction, NATS client wrappers, config loading, etc.
  - auth/: Auth-related modules.
  - store/: Data storage abstraction and implementation.
- api/: API definitions (e.g., Protobuf, OpenAPI/Swagger specs).
- web/: Web client static resources (HTML, JS, CSS, images).
  - htmx/ or static/
- docs/: Project documentation.
  - designs/: Design documents.
- configs/: Example config files.
- scripts/: Build, deployment, and helper scripts.
- build/: (or as defined in Makefile) Build-related files, e.g., Dockerfile.
- test/: End-to-end or integration test files.

### 4.2 Coding Conventions

- Follow the official Go style guide.
- Use gofmt or goimports for automatic code formatting.
- Use golangci-lint for static code checks; the config file (.golangci.yml) should be version controlled.
- Unified formatting command: make fmt
- Static check command: make lint
- Write clear, maintainable comments.
- Error handling: Follow Go best practices, avoid panic (except in main or unrecoverable scenarios).

### 4.3 Branching & Merging Strategy

- Main branch: main (or master), kept stable, corresponds to released or soon-to-be-released versions.
- Development branch: develop (optional), base for daily development, new features branch from here.
- Feature branch: feature/xxx (e.g., feature/e2e-encryption), for new features. Merge to develop (or main) after completion.
- Fix branch: fix/xxx (e.g., fix/login-bug), for bug fixes on main. Usually merged to both main and develop after fixing.
- Release branch: release/vX.Y.Z (optional), for preparing releases, final testing, and docs update.
- Merge policy:
  - All code changes must be submitted via Pull Request (PR) / Merge Request (MR).
  - PRs require at least one other developer for code review.
  - PRs must pass all CI checks (linting, tests).
  - Squash and Merge or Rebase and Merge is recommended to keep main branch history clean.

## 5. Testing Strategy

### 5.1 Unit Testing

- Scope: Test individual functions, methods, modules, or packages.
- Tools: Use Go's built-in testing package, with ginkgo and gomega for BDD-style tests.
- Goal: Focus on core logic, edge cases, and error handling. Mock external dependencies (DB, NATS) for isolation and speed.
- Run: make test-unit or go test ./...

## 6. Contribution Guide

- Bug reports, feature suggestions, and discussions are welcome via GitHub Issues.
- Before contributing code, it's recommended to create an Issue describing your idea or the problem to be fixed, and discuss with maintainers.
- Fork the repo to your own account.
- Create your feature/xxx or fix/xxx branch based on develop (or main).
- Follow coding conventions and code structure.
- Write necessary unit and integration tests for new or modified code.
- Ensure all tests pass (make test) and code passes lint checks (make lint).
- Update relevant docs (README, docs/) for significant changes or new features.
- Submit a Pull Request to the upstream develop (or main) branch.
- Clearly describe the changes, reasons, and test status in the PR description.
- Actively participate in PR code review and revise code based on feedback.

## 7. Glossary

- Management Server: Central control service, responsible for authentication, authorization, scheduling, auditing, etc.
- Agent: Lightweight proxy deployed in the target environment, responsible for command execution and terminal session management.
- Client: User interface for interacting with the Terminal system, such as Web or TUI client.
- Web Client: Browser-based graphical user interface.
- TUI (Text User Interface): Command-line based user interface, implemented with bubbletea in this project.
- Dev Mode: A simplified run mode, usually merging management and agent functions, using embedded DB and message queue for local development and debugging.
- NATS: High-performance message queue system for asynchronous communication and data streaming between components.
- PTY (Pseudo-Terminal): Software interface simulating a physical terminal, used for interactive shell sessions.
- E2EE (End-to-End Encryption): Secure communication method ensuring data is encrypted by the sender and decrypted by the receiver, with intermediaries (including Management Server) unable to read plaintext.
- RBAC (Role-Based Access Control): Permission management model assigning roles to users and granting permissions to roles for flexible access control.
- HTMX: Library enabling modern browser features via HTML attributes for dynamic web interfaces.
- xterm.js: Powerful web component for fully functional terminal emulation in browsers.
- Gin: High-performance HTTP web framework for Go.
- Ent: Entity framework for Go, used for data modeling and DB interaction.
