# gRPC Chat

A small gRPC chat server with two RPCs:

- `CreateStream` registers a client and keeps a server-to-client message stream open.
- `BroadcastMessage` sends a message to every registered active client.

The server listens on port `8080`. The RPC and message definitions are in [`proto/chat.proto`](proto/chat.proto).

## Run the server

From the repository root:

```bash
go run ./lld/grpc-chat
```

The server should print:

```text
Server started at port :8080
```

## Try it with Postman

Create two gRPC requests in Postman and set the server address to:

```text
localhost:8080
```

Load `lld/grpc-chat/proto/chat.proto` as the service definition. The server does not enable gRPC reflection, so Postman needs the proto file to discover the service and RPCs.

### 1. Open a client stream

Select `chat.Broadcast/CreateStream` and send:

```json
{
  "user": {
    "ID": "user-1",
    "name": "Alice"
  },
  "active": true
}
```

`CreateStream` is server-streaming: it accepts one `Connect` request and keeps the RPC open so the server can send `Message` responses. Leave this request running in Postman to receive broadcasts. Open another `CreateStream` request with a different user ID to test multiple clients.

### 2. Broadcast a message

While at least one stream is open, create a second Postman gRPC request for `chat.Broadcast/BroadcastMessage` and send:

```json
{
  "ID": "message-1",
  "content": "Hello, everyone!",
  "timestamp": "2026-11-19T12:00:00Z"
}
```

Each open client stream should receive the message. `BroadcastMessage` returns an empty `Close` response.

## How the handlers work

### `CreateStream`

The handler creates a `Connection` containing the user's ID and gRPC stream, appends it to the connection pool, and waits on that connection's error channel. This keeps the RPC active while the client is connected.

### `BroadcastMessage`

The handler loops over connections in the pool and calls `stream.Send(msg)` for each active connection. It waits for those sends to finish, then returns a `Close` response.

## Current implementation notes

- `CreateStream` reads `pconn.User.ID`; include a `user` object in the request to avoid a nil pointer panic.
- The current handler does not remove normally disconnected streams from the pool or synchronize concurrent access to the pool. It is a learning example rather than a production-ready chat server.
- The proto declares fields with uppercase names such as `ID`; use the names Postman displays after loading the proto.
