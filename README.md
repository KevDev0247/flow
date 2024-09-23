# Flow
PLEASE CHECK AGAGAIN AT THE END OF 2024-09-23

Take home project for flow interview! Work in progress... Currently testing out concurrency features and edge cases.

Also, check out [Proofster](https://github.com/KevDev0247/proofster), a first order logic computation project where I applied knowledge of Golang, Python, algorithms, and microservices.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Environment Setup](#environment-setup)
3. [Dependencies](#dependencies)
4. [Starting the Server](#starting-the-server)
5. [Interacting with the Program](#interacting-with-the-program)
6. [Assumptions](#assumptions)

## Prerequisites

- Go (version 1.18 or higher)
- Curl

## Environment Setup

1. **Install Go**: Follow the instructions on the [official Go website](https://golang.org/doc/install) to install Go on your machine.
2. **Set Up Your Workspace**:
   - Create a directory for the project:
     ```bash
     mkdir flow
     cd flow
     ```

3. **Clone the Repository**:
   - Clone this repository (or copy the code files) into the workspace.
     ```bash
     git clone git@github.com:KevDev0247/flow.git
     ```

## Dependencies

None so far, will be added when tests are introduced

## Starting the Server

To start the server, navigate to the directory containing the `main.go` file and run:

```bash
go run main.go
```

# Interacting With The Program
## Add a Block and Vote
```bash
curl -X POST http://localhost:8080/process \
-H "Content-Type: application/json" \
-d '{"block": {"id": "block1", "view": 1}, "vote": {"block_id": "block1"}}'
```

## Add Another Block and Vote
```bash
curl -X POST http://localhost:8080/process \
-H "Content-Type: application/json" \
-d '{"block": {"id": "block1", "view": 2}, "vote": {"block_id": "block1"}}'
```

## Add a Third Block and Vote
```bash
curl -X POST http://localhost:8080/process \
-H "Content-Type: application/json" \
-d '{"block": {"id": "block1", "view": 3}, "vote": {"block_id": "block1"}}'
```
More to be added...
You can check the server logs to see the output of accepted blocks.

# Assumptions
- Each block has a unique ID and an associated view count.
- Votes are only valid for existing blocks.
- The program assumes that the input data is correctly formatted as JSON. Invalid or improperly formatted requests will result in a `400 Bad Request` response.
- The implementation does not persist data across server restarts; all blocks and votes are lost when the server stops.

