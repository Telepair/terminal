[English](README.md) | [中文](docs/README_zh.md) | [Development Guide](docs/develop_guide.md) | [Development Plan](docs/plans.md)

# Terminal

**Unified, Secure, and Auditable Command-Line Access**

Terminal is a comprehensive platform designed to simplify and secure command-line operations across various infrastructures, including physical/virtual servers (Linux, MacOS), Kubernetes clusters (internal/external), Docker, and Containerd environments.

It provides a centralized gateway for:

- **Unified Access:** Connect to various backends through a single interface (Web and TUI).
- **Enhanced Security:** Robust authentication (TOTP, OIDC, etc.), fine-grained authorization (JWT-based RBAC), and optional end-to-end session encryption.
- **Comprehensive Auditing:** Detailed operational logs and session recordings for complete visibility and compliance.
- **Operational Efficiency:** Command templates, bulk command execution, and an API proxy to internal services.

## Core Components

1.  **Client:** User interaction interface, including a Web UI based on HTMX + xterm.js and a TUI (Text User Interface) based on Bubbletea.
2.  **Management Server:** Central control plane built with Go and NATS. Responsible for authentication, authorization, Agent management, logging, session recording, and configuration management.
3.  **Agent:** A lightweight daemon written in Go, deployed on target systems to securely expose terminal access, execute commands, and communicate with the Management Server.

## Key Features

- **Multiple Terminal Types:** Local Shell, SSH connections, Kubernetes Pods, Docker containers, Containerd containers.
- **Flexible Operation Modes:** Interactive PTY sessions, one-off command execution, server-side streaming output.
- **Authentication:** Supports NoAuth (not recommended for production), Username/Password, TOTP, Email verification, OIDC (e.g., Keycloak, Okta), Asymmetric key pairs.
- **Authorization:** JWT-based Role-Based Access Control (RBAC) supporting fine-grained permission configuration.
- **Session Recording:** Agents can record terminal sessions (e.g., in Asciinema format), with metadata stored and indexed by the Management Server.
- **Command Templates:** Define and reuse common command patterns to improve operational efficiency.
- **API Proxy:** (Planned) Securely expose internal APIs through Terminal.
- **Bulk Command Execution:** Run commands simultaneously on multiple target Agents.
- **End-to-End Encryption (E2EE):** (Planned) Optional session encryption to prevent the Management Server from snooping on session content.
- **Connection Modes:**
  - Client <-> Agent (Direct connection via NATS, requires network reachability)
  - Client <-> Management Server (Proxy) <-> Agent (via NATS)
  - Supports WebSocket, Server-Sent Events (SSE), HTTP for client communication.
- **Development Mode:** Bundles Management Server + Agent, with embedded SQLite and NATS, to simplify local development and testing.

## Technology Stack

- **Backend and Agent:** Go
  - API Framework: [Gin](https://github.com/gin-gonic/gin)
  - Test-Driven Development (TDD): [ginkgo](https://github.com/onsi/ginkgo) & [gomega](https://github.com/onsi/gomega)
  - ORM: [Ent](https://github.com/ent/ent)
- **Web Frontend:** [HTMX](https://github.com/bigskysoftware/htmx), [xterm.js](https://xtermjs.org/)
- **TUI Client:** [bubbletea](https://github.com/charmbracelet/bubbletea)
- **Database:** PostgreSQL (recommended for production), SQLite (default in development mode)
- **Message Queue/Storage:** [Nats](https://github.com/nats-io/nats-server) (Pub/Sub, Key-Value Store, Object Storage)
- **Communication Protocols:**
  - Agent to Management Server: Nats Pub/Sub based
  - Client to Agent: Primarily Nats Pub/Sub based (can be proxied by Management Server)
  - Client to Management Server: HTTP/HTTPS API (Gin) and Nats Pub/Sub

## Quick Start

1.  **Get the Source Code:**

    ```bash
    git clone https://github.com/<your-username>/terminal.git # Replace with actual repository address
    cd terminal
    ```

2.  **Run in Development Mode:**
    This mode will start an integrated service, including management and Agent functionalities, using embedded SQLite and NATS.

    ```bash
    make run-dev
    ```

    (If `make run-dev` is not available, please refer to the instructions for starting development mode in the [Development Guide](./develop_guide.md#23-build-and-run-guide), e.g., `go run cmd/terminal/main.go --dev-mode`)

3.  **Access the Client:**

    - Web Client: Open a browser and navigate to `http://localhost:<port>` (check startup logs for the specific port, defaults to 8080 or similar).
    - TUI Client: (If implemented) `./bin/terminal-cli` or a similar command.

4.  **More Information:**
    - For detailed development environment setup, build, and run guides, please refer to the [Development Guide](./develop_guide.md#2-environment-setup-and-build-setup--build).
    - To understand project components and code structure, please refer to the [Development Guide](./develop_guide.md).

## Roadmap

For the detailed development plan and version targets of the project, please see the [Development Plan](docs/plans.md).

## Contributing

We welcome contributions of all kinds! Whether it's submitting bug reports, proposing feature suggestions, or contributing code.
Please refer to our [Contribution Guide](docs/develop_guide.md#6-contribution-guide-contribution-guide) to learn how to participate.

## License

This project is licensed under the [Apache License 2.0](./LICENSE).
