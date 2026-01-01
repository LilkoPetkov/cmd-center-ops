# CMD-Center (ops)

 **CMD-Center** is a powerful, user-friendly Command Line Interface (CLI) application built for developers, system administrators, network engineers, and cybersecurity professionals. It provides precise and flexible tools for domain name resolution, DNS diagnostics, network services, YAML manipulation, and more, directly from your terminal. 🌟

## Features

- **DNS Tools**:
  - Resolve domain names for various record types (A, AAAA, CNAME, NS, TXT).
  - Start a local DNS server for testing and diagnostics.
- **Network Utilities**:
  - Start a simple TCP server.
  - Test TCP connections to any host and port (a `telnet`-like utility).
- **YAML Management**:
  - Edit YAML files by updating key-value pairs globally.
  - Perform fine-grained, scoped edits within specific YAML blocks.
- **Data Generation & Management**:
  - Generate various versions of UUIDs (v4, v6, v7).
  - Orchestrate Open-Metadata ingestion pipelines.
- **Styled Output**:
  - Uses `lipgloss` for a clean, readable, and colorful user experience.

## Installation

1.  **Clone the repository:**

    ```sh
    git clone git@github.com:LilkoPetkov/cmd-center-ops.git
    cd cmd-center-ops
    ```

2.  **Build the binary:**

    ```sh
    go build -o ops .
    ```

3.  **(Optional) Move the binary to a directory in your PATH:**
    ```sh
    sudo mv ops ~/go/bin/
    ```

## Usage

The root command is `ops`. You can see a list of all available commands by running `ops --help`.

### DNS Commands

#### Resolve a Domain

Get information about the resolution of a domain name.

- **Usage:** `ops dns resolve [flags]`
- **Examples:**

  ```sh
  # Resolve the AAAA record for example.com (default)
  ops dns resolve -d example.com

  # Resolve the CNAME record for example.com
  ops dns resolve -d example.com -q cname

  # Resolve all main record types (A, AAAA, CNAME, NS, TXT) for example.com
  ops dns resolve -d example.com -a
  ```

### Server Commands

#### Start a DNS Server

Start a DNS server on a specified port for testing purposes.

- **Usage:** `ops server dns [flags]`
- **Examples:**

  ```sh
  # Start a DNS server on the default port 8888
  ops server dns

  # Start a DNS server on the standard DNS port 53
  ops server dns -p 53
  ```

#### Start a TCP Server

Start a TCP server on a specified port for network testing.

- **Usage:** `ops server tcp [flags]`
- **Examples:**

  ```sh
  # Start a TCP server on default port 8888
  ops server tcp

  # Start a TCP server on port 9000
  ops server tcp -p 9000
  ```

#### Test a Connection (Telnet)

Test the connection to a server on a specific port, similar to `telnet`.

- **Usage:** `ops server telnet [flags]`
- **Examples:**

  ```sh
  # Test connection to localhost on default port 443
  ops server telnet

  # Test HTTPS port on a public server
  ops server telnet -n google.com -p 443

  # Check if an internal service is reachable on a custom port
  ops server telnet -n 10.0.0.15 -p 8080
  ```

### YAML Commands

#### General YAML Edit

Update a YAML file's general structure without fine-grained control. This command will update all keys that match.

- **Usage:** `ops yaml yedit [flags]`
- **Examples:**

  ```sh
  # Update the value of a top-level key in a YAML file
  ops yaml yedit -p ./config.yaml -k port -v 8080

  # Update the service name across all YAML files in a directory
  ops yaml yedit -p ./configs/ -k name -v my-updated-service
  ```

#### Scoped YAML Edit

Update a YAML file's scoped value for more fine-grained control. This is useful for updating a key within a specific block or parent key.

- **Usage:** `ops yaml yedit_scoped [flags]`
- **Examples:**

  ```sh
  # Update the value of a specific key inside a matched YAML block
  ops yaml yedit_scoped -p ./config.yaml -s "database" -k "host" -v "localhost"

  # Apply scoped update across all YAML files in a directory
  ops yaml yedit_scoped -p ./configs/ -s "authService" -k "port" -v "8081"
  ```

### UUID Generation

Generate a new UUID. Supported types are `uuid4`, `uuid6`, `uuid7`, and `clock`.

- **Usage:** `ops uuid [flags]`
- **Examples:**

  ```sh
  # Generate a default UUID (uuid4)
  ops uuid

  # Generate a UUID version 6
  ops uuid -t uuid6

  # Generate a UUID version 7
  ops uuid --type=uuid7
  ```

### Open-Metadata Commands

Allows execution of Open-Metadata commands. This requires a Python virtual environment and `uv`.

#### Ingest Metadata

Ingest metadata using a connector and a YAML configuration file.

- **Usage:** `ops ometadata ingest [flags]`
- **Examples:**

  ```sh
  # Ingest metadata using default pyproject.toml and config.yaml files
  ops ometadata ingest

  # Ingest metadata specifying a custom pyproject.toml file and config.yaml
  ops ometadata ingest -p path/to/pyproject.toml -c path/to/config.yaml
  ```

## Dependencies

This project is built with several key Go libraries:

- **[github.com/spf13/cobra](https://github.com/spf13/cobra)**: A powerful library for creating modern CLI applications.
- **[github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)**: A library for style-based text rendering in the terminal.
- **[github.com/miekg/dns](https://github.com/miekg/dns)**: A comprehensive DNS library for Go.
- **[github.com/goccy/go-yaml](https://github.com/goccy/go-yaml)**: A robust YAML parser for Go.
- **[github.com/google/uuid](https://github.com/google/uuid)**: A library for generating and working with UUIDs.

## License

This project is licensed under the terms of the [LICENSE](LICENSE) file.
