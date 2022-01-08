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
    "content": "..."
}
```

server->client forwarding message

```json
{
    "source-client-name": "...",
    "content": "..."
}
```

