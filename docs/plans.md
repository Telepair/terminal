# Terminal Development Plan

- [Terminal Development Plan](#terminal-development-plan)
  - [v0.1 MVP Implementation](#v01-mvp-implementation)
  - [v0.2 Introduce Web Client \& Interactive Sessions](#v02-introduce-web-client--interactive-sessions)
  - [v0.3 Extend Terminal Type Support](#v03-extend-terminal-type-support)
  - [v0.4 Core Authentication, Authorization, and Auditing](#v04-core-authentication-authorization-and-auditing)
  - [v0.5 Introduce Command Templates](#v05-introduce-command-templates)
  - [v0.6 Architecture Evolution: Separate Management and Agent Deployment](#v06-architecture-evolution-separate-management-and-agent-deployment)
  - [v0.7 Feature and UX Enhancements](#v07-feature-and-ux-enhancements)
  - [v0.8 Security Enhancement: End-to-End Encryption](#v08-security-enhancement-end-to-end-encryption)
  - [v0.9 Extended Integrations and Interaction Methods](#v09-extended-integrations-and-interaction-methods)
  - [V1.0: Production Ready](#v10-production-ready)

## v0.1 MVP Implementation

**Goal**: Build the Minimum Viable Product to validate core local terminal connection (one-off command execution) and basic infrastructure.

- [ ] **Project Initialization & Infrastructure**
  - [ ] Project structure setup (`cmd/`, `internal/`, `pkg/`, `api/`, etc.)
  - [ ] Introduce core tech stack (Go, Gin, Ent, NATS client)
  - [ ] Test-driven development environment (ginkgo, gomega)
  - [ ] Basic Makefile (build, lint, fmt, test-unit, run-dev)
- [ ] **Core API & Data Model Design**
  - [ ] Initial database schema design (SQLite & PostgreSQL)
  - [ ] Draft Management Server core APIs (Client <-> Management, Agent <-> Management)
  - [ ] Draft Agent core APIs (Management <-> Agent)
  - [ ] Initial NATS topics and message structure design
- [ ] **Dev Mode Core Features**
  - [ ] Management Server basic startup, embedded NATS and SQLite
  - [ ] Agent basic startup (can run merged with Server in Dev mode)
  - [ ] Agent registration to Management Server
  - [ ] **One-off Command Execution**: Client -> Management Server -> Agent (via NATS/API) -> return result
    - Support Linux local shell commands
- [ ] **Basic TUI Client**
  - [ ] Implement basic TUI client framework (Bubbletea)
  - [ ] TUI connects to Management Server, lists registered Agents
  - [ ] TUI supports executing one-off commands on selected Agent and displaying results

## v0.2 Introduce Web Client & Interactive Sessions

**Goal**: Introduce a web client, support interactive PTY sessions, and implement basic operation logging.

- [ ] **Basic Web Client**
  - [ ] Implement basic web client framework (HTMX, xterm.js)
  - [ ] Web client communicates with Management Server via API
  - [ ] Web client displays Agent list and supports one-off command execution
- [ ] **Interactive PTY Sessions**
  - [ ] Agent supports PTY (pseudo-terminal) creation and management (Linux)
  - [ ] Implement PTY data streaming: Client <-> Management Server (proxy) <-> Agent (via NATS & WebSocket/SSE)
  - [ ] Web client (xterm.js) supports interactive terminal sessions
  - [ ] TUI client supports interactive terminal sessions
- [ ] **Management Server Enhancements**
  - [ ] Basic operation log storage (command execution history)
  - [ ] Web client displays operation logs

## v0.3 Extend Terminal Type Support

**Goal**: Extend Agent support for more mainstream terminal types.

- [ ] **Agent Capability Expansion**
  - [ ] Support connecting to local Docker containers (exec, attach)
  - [ ] Support connecting to local Containerd containers (exec, attach)
  - [ ] Support connecting to containers in Kubernetes Pods (exec, attach, requires Kubeconfig or ServiceAccount)
  - [ ] Support MacOS local shell (PTY and one-off commands)
  - [ ] (Optional) Support SSH connections to other hosts via Agent

## v0.4 Core Authentication, Authorization, and Auditing

**Goal**: Introduce formal user authentication, role-based authorization (RBAC), and key security auditing features.

- [ ] **Management Server Security Module**
  - [ ] **Authentication Design & Implementation**:
    - [ ] Username/password authentication (store encrypted passwords)
    - [ ] JWT (JSON Web Tokens) as session tokens
    - [ ] (Optional initial) TOTP two-factor authentication
  - [ ] **Authorization Design & Implementation (RBAC)**:
    - [ ] Define user, role, and permission models
    - [ ] Implement role-based API access control
    - [ ] Implement role-based Agent access control (which users/roles can access which Agents)
  - [ ] **Key Security Auditing Features**:
    - [ ] Log user login, authentication failures, permission changes, etc.
    - [ ] Log detailed operation audit records (who, when, where, what target, what operation)
- [ ] **Agent Security Cooperation**
  - [ ] (Optional) Agent session recording (e.g., Asciinema format), upload to Management Server-controlled storage (NATS Object Store)
  - [ ] Management Server provides session playback (Web client)

## v0.5 Introduce Command Templates

**Goal**: Introduce command template functionality to simplify common operations, improve efficiency and standardization.

- [ ] **Management Server Features**
  - [ ] Command template CRUD API (support parameterized templates)
  - [ ] Web client manages and uses command templates
- [ ] **Agent Features**
  - [ ] Support executing parameterized command templates from Management Server
- [ ] **Client Features**
  - [ ] TUI and Web client support listing, selecting, and executing command templates

## v0.6 Architecture Evolution: Separate Management and Agent Deployment

**Goal**: Optimize architecture to support fully independent deployment of Management Server and Agent, laying the foundation for advanced security features.

- [ ] **Deployment Mode Enhancement**
  - [ ] Documentation and configuration examples for deploying Management Server and Agent as independent services
  - [ ] Improve robustness of Agent auto-registration to remote Management Server
- [ ] **Connection Proxy & Routing**
  - [ ] Enhance Management Server's ability to proxy communication between Client and Agent (for complex network environments)
- [ ] **Authentication Enhancement**
  - [ ] Support OIDC (OpenID Connect) authentication integration (e.g., Keycloak, Okta)

## v0.7 Feature and UX Enhancements

**Goal**: Improve user experience and add convenient features.

- [ ] **Notification Mechanism**:
  - [ ] Agent status changes, important event notifications (e.g., via Web UI, Email)
- [ ] **Batch Operations**:
  - [ ] Support batch command/template execution on multiple selected Agents
- [ ] **File Transfer**:
  - [ ] Support secure file upload/download between Client and Agent (via Management Server relay or NATS)
- [ ] **Background & Scheduled Tasks**:
  - [ ] (Exploratory) Support long-running background tasks on Agent
  - [ ] (Exploratory) Support configuring scheduled tasks from Management Server to Agent
- [ ] **Infinite Canvas/Dashboard**:
  - [ ] (Exploratory) Web client provides customizable dashboard to display Agent status, common connections, etc.

## v0.8 Security Enhancement: End-to-End Encryption

**Goal**: Introduce stronger security mechanisms, especially for session data end-to-end encryption.

- [ ] **Key Pair Authentication**:
  - [ ] Support Client and/or Agent authentication using asymmetric key pairs
- [ ] **End-to-End Encryption (E2EE) Design & Implementation**:
  - [ ] Client and Agent negotiate session keys
  - [ ] PTY data stream is end-to-end encrypted between Client and Agent, Management Server cannot decrypt session content
  - [ ] Adjust session recording scheme for E2EE mode (e.g., client-side recording or encrypted recording)

## v0.9 Extended Integrations and Interaction Methods

**Goal**: Provide more ways to integrate and automate interactions with other systems.

- [ ] **Bot Support**:
  - [ ] (Exploratory) Provide API or SDK for chatbots (e.g., Slack, Teams) to interact with Terminal and execute commands
- [ ] **Webhook Support**:
  - [ ] Management Server supports configuring webhooks to notify external systems on specific events (e.g., Agent online/offline)
- [ ] **Event Subscription**:
  - [ ] Client or external systems can subscribe to events published by Management Server via NATS or API

## V1.0: Production Ready

**Goal**: Polish the product to meet production deployment requirements.

- [ ] **Documentation Completion**:
  - [ ] User manual, admin manual, API documentation, deployment guide
- [ ] **Multi-Tenancy Support**: (if applicable)
  - [ ] Data isolation and permission control, support for multiple independent organizations or teams using the same system
- [ ] **Performance & Stability**:
  - [ ] Stress testing, long-term stability testing
  - [ ] Performance bottleneck analysis and optimization (high concurrency, large-scale Agent management)
- [ ] **Operations Support**:
  - [ ] Complete monitoring metrics (Prometheus metrics)
  - [ ] Structured log output for easy collection and analysis
  - [ ] Alert mechanism integration
- [ ] **Security Hardening**:
  - [ ] Third-party security audit or penetration testing
  - [ ] Dependency security scanning and updates
- [ ] **Production Deployment Solutions**:
  - [ ] Provide Docker-compose, Kubernetes Helm Chart, etc. deployment solutions
