# Go MCP Client

This is a command-line interface (CLI) client written in Go to interact with an MCP (My Control Plane - Placeholder Name) server. It allows performing CRUD (Create, Read, Update, Delete) operations on MCP configurations.

**Note:** This client is currently a placeholder and assumes the MCP server provides a RESTful API over HTTP with JSON payloads. The `MCPConfig` structure in `config.go` and the client logic in `client.go` will need to be updated based on the actual MCP server's API specification and data model.

## Prerequisites

- Go (version 1.18 or later recommended)

## Setup

1.  **Clone the repository (if applicable) or ensure you have the source files.**
2.  **Define Configuration Structure:**
    Update the `MCPConfig` struct in `config.go` to match the actual data structure of your MCP server's configurations.
3.  **Update Client Logic:**
    Modify `client.go` to correctly interact with your MCP server:
    *   Adjust API endpoints (e.g., `/configs`, `/configs/{id}`).
    *   Implement any necessary authentication mechanisms within the `makeRequest` function or `MCPClient` struct.
    *   Change the request/response handling if the server uses a different protocol (e.g., gRPC) or data format (e.g., YAML).

## Building

To build the client, navigate to the project's root directory and run:

```bash
go build
```

This will produce an executable file named `mcp-client` (or `mcp-client.exe` on Windows).

## Usage

The client uses command-line flags to specify the action and necessary parameters.

```bash
./mcp-client -server <MCP_SERVER_URL> -action <ACTION> [OPTIONS]
```

### Common Flags:

-   `-server <URL>`: (Required) Base URL of the MCP server (e.g., `http://localhost:8080`).
-   `-action <ACTION>`: (Required) The operation to perform. Supported actions:
    -   `create`
    -   `get`
    -   `update`
    -   `delete`
    -   `list`

### Action-Specific Options:

**1. Create a new configuration:**

```bash
./mcp-client -server <URL> -action create -name "My New Config" -version 1 -enabled=true
```

-   `-name <STRING>`: Name of the configuration.
-   `-version <INT>`: Version number.
-   `-enabled <BOOL>`: `true` or `false`.
-   *(You may need to add more flags to `main.go` for other `MCPConfig` fields)*

**2. Get a configuration by ID:**

```bash
./mcp-client -server <URL> -action get -id "config-123"
```

-   `-id <STRING>`: The ID of the configuration to retrieve.

**3. Update an existing configuration:**

```bash
./mcp-client -server <URL> -action update -id "config-123" -name "Updated Name" -version 2
```

-   `-id <STRING>`: The ID of the configuration to update.
-   *(Flags for fields to update, e.g., `-name`, `-version`)*

**4. Delete a configuration by ID:**

```bash
./mcp-client -server <URL> -action delete -id "config-123"
```

-   `-id <STRING>`: The ID of the configuration to delete.

**5. List all configurations:**

```bash
./mcp-client -server <URL> -action list
```

### Examples:

```bash
# List all configurations from a server at http://mcp.example.com
./mcp-client -server http://mcp.example.com -action list

# Get configuration with ID "alpha-001"
./mcp-client -server http://mcp.example.com -action get -id "alpha-001"

# Create a new configuration
./mcp-client -server http://mcp.example.com -action create -name "Primary Settings" -version 1 -enabled=true
```

## Running Tests

To run the unit tests:

```bash
go test
```

## Project Structure

-   `main.go`: Contains the CLI logic (flag parsing, command dispatch).
-   `client.go`: Implements the `MCPClient` for interacting with the MCP server.
-   `config.go`: Defines the `MCPConfig` struct (placeholder, needs customization).
-   `client_test.go`: Unit tests for `MCPClient`.
-   `go.mod`: Go module definition file.
-   `README.md`: This file.

## TODO / Further Development

-   **Implement actual MCP server communication:** The current client logic is a placeholder. It needs to be adapted to the specific MCP server's API (protocol, endpoints, authentication, data models).
-   **More robust CLI:** For complex applications, consider using a library like Cobra for the CLI.
-   **Configuration file support:** Allow passing parameters via a configuration file (e.g., using Viper).
-   **Enhanced error handling and logging:** Implement more detailed and structured logging.
-   **More comprehensive tests:** Add more test cases, including various error scenarios and edge cases.
