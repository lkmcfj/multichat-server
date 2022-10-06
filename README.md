# multichat

multichat is a JSON protocol over WebSocket for cross-platform chat forwarding. It consists of a server for routing messages, and several clients for each chat source. To connect a chat platform (like some IM app) into this framework, you just need to implement a client for that platform in this protocol.

This repository is an implementation of the server in Go.

## Usage

1. Clone this repo.
2. `go build`
3. Rename `config.json.example` to `config.json` and edit it on your needs.
4. Run the executable 

### Configuration File Fields

| name         | type                     | value                                                        |
| ------------ | ------------------------ | ------------------------------------------------------------ |
| `ws-path`    | string                   | The http path where the connection is upgraded to WebSocket. Usually you want to leave it to `"/"`. |
| `host`       | string                   | Where to bind. A typical value may look like `0.0.0.0:1234`. |
| `secret-key` | string                   | A secret key shared between the client and the server. Note that if the connection is not protected by WSS (maybe a reverse proxy or the native WSS support), this authentication isn't really secure. |
| `wss`        | WSS configuration object | Optional. If this field is omitted, WSS is disabled. The fields in this object are illustrated in the next sheet. |

WSS configuration object:

| name       | type   | value                         |
| ---------- | ------ | ----------------------------- |
| `keyfile`  | string | The path to your private key. |
| `certfile` | string | The path to your certificate. |

## Protocol

All WebSocket packets should be in text mode.

After a client connects to a server, the connection enters the authentication stage. The first packet sent to the server should be a registration message:

```
{
    "action": "register",
    "client-name": "...",
    "secret-key": "..."
}
```

If the authentication is successful, the server will send an ACK message to the client:

```
{
    "action": "register-ack"
}
```

After this, the normal stage begins, and messages of type `client-message` and `forwarding-message` may occur.

On an authentication fail, the server closes the WebSocket connection. For both the client and the server, if any critical error is encountered, you can just close the connection.

When a chat message occurs in the chat source, the client may send a `client-message` to the server:

```
{
    "action": "client-message",
    "content": "..."
}
```

And the server broadcasts the message to everyone connected:

```
{
    "action": "forwarding-message",
    "source-client-name": "...",
    "content": "..."
}
```

