client register:

```json
{
    "action": "register",
    "client-name": "...",
    "secret-key": "..."
}
```

server ack:

```json
{
    "action": "register-ack"
}
```

client->server message

```json
{
    "action": "client-message",
    "content": "..."
}
```

server->client forwarding message

```json
{
    "action": "forwarding-message",
    "source-client-name": "...",
    "content": "..."
}
```

