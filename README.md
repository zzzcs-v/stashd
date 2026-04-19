# stashd

A lightweight key-value store daemon with TTL support and a simple HTTP API.

## Installation

```bash
go install github.com/yourusername/stashd@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/stashd.git && cd stashd && go build ./...
```

## Usage

Start the daemon:

```bash
stashd --port 8080
```

**Set a key** (with optional TTL in seconds):

```bash
curl -X PUT http://localhost:8080/keys/mykey \
  -d '{"value": "hello", "ttl": 60}'
```

**Get a key:**

```bash
curl http://localhost:8080/keys/mykey
```

**Delete a key:**

```bash
curl -X DELETE http://localhost:8080/keys/mykey
```

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/keys/:key` | Retrieve a value |
| `PUT` | `/keys/:key` | Set a value with optional TTL |
| `DELETE` | `/keys/:key` | Delete a key |
| `GET` | `/health` | Health check |

## Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | `8080` | Port to listen on |
| `--ttl` | `0` | Default TTL in seconds (0 = no expiry) |

## License

MIT © [yourusername](https://github.com/yourusername)