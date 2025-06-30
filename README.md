# Godis

`godis` is a Redis-like in-memory data store implemented in Go. It supports basic key-value operations, list operations, and follows the RESP (REdis Serialization Protocol) for communication.

## Features

`godis` currently supports the following Redis commands:

- **SET** `<key> <value>`: Sets the string `value` of `key`.
- **GET** `<key>`: Get the value of `key`. If the key does not exist, it returns nil.
- **LPUSH** `<key> <value> [value ...]`: Inserts all the specified `values` at the head of the list stored at `key`. If `key` does not exist, it is created as an empty list before pre-pending these values.
- **LPOP** `<key>`: Removes and returns the first element of the list stored at `key`.
- **PING** `[message]`: Returns `PONG` if no `message` is provided, otherwise returns `message`.
- **QUIT**: Closes the connection.

## Project Structure

The project is organized into the following directories and files:

- `main.go`: The main entry point of the server. It handles incoming connections, parses commands, and interacts with the data store.
- `parser/parser.go`: Contains the implementation for parsing the RESP protocol, including types for simple strings, bulk strings, integers, arrays, and errors.
- `store/store.go`: Implements the in-memory data storage, including thread-safe operations for setting/getting key-value pairs and performing list operations (LPUSH, LPOP).
- `README.md`: This file.
- `.gitpod.yml`: Configuration file for Gitpod, setting up the development environment.

## Getting Started

### Prerequisites

- Go (version 1.18 or higher recommended)

### Installation and Running Locally

1. **Clone the repository:**
   ```bash
   git clone https://github.com/malayt04/godis.git
   cd godis
   ```

2. **Build the project:**
   ```bash
   go build -o godis .
   ```

3. **Run the server:**
   ```bash
   ./godis
   ```

The server will start listening on port `6379`.

### Using with Gitpod

This project includes a `.gitpod.yml` configuration, allowing you to easily set up a development environment in Gitpod.

1. Open the project in Gitpod
2. Gitpod will automatically run `go get`, `go build`, `go test`, and `make` (if a Makefile is present) as part of the `init` task, and then execute `go run .` to start the server.

## Usage Examples

Once the `godis` server is running (e.g., on `localhost:6379`), you can interact with it using a Redis client or `netcat`.

### Using `netcat` (nc)

1. **Connect to the server:**
   ```bash
   nc localhost 6379
   ```

2. **Enter commands in RESP format. For example:**

   **SET command:**
   ```
   *3
   $3
   SET
   $4
   mykey
   $5
   value
   ```
   Expected output: `+OK`

   **GET command:**
   ```
   *2
   $3
   GET
   $4
   mykey
   ```
   Expected output:
   ```
   $5
   value
   ```

   **LPUSH command:**
   ```
   *4
   $5
   LPUSH
   $5
   mylist
   $1
   A
   $1
   B
   ```
   Expected output: `:2`

   **LPOP command:**
   ```
   *2
   $4
   LPOP
   $5
   mylist
   ```
   Expected output:
   ```
   $1
   B
   ```

   **PING command:**
   ```
   *1
   $4
   PING
   ```
   Expected output: `+PONG`

   ```
   *2
   $4
   PING
   $5
   Hello
   ```
   Expected output: `+Hello`

   **QUIT command:**
   ```
   *1
   $4
   QUIT
   ```
   Expected output: `+OK`

## Contributing

Feel free to fork the repository, open issues, and submit pull requests.
