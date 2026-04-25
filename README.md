# pulsar-watch

A lightweight CLI tool to monitor and replay Apache Pulsar topic messages with filtering and export support.

---

## Installation

```bash
go install github.com/yourname/pulsar-watch@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/pulsar-watch.git
cd pulsar-watch
go build -o pulsar-watch .
```

---

## Usage

```bash
# Watch messages on a topic
pulsar-watch --broker pulsar://localhost:6650 --topic persistent://public/default/my-topic

# Filter messages by key
pulsar-watch --topic persistent://public/default/my-topic --filter "user-id=42"

# Replay messages from a specific timestamp and export to JSON
pulsar-watch --topic persistent://public/default/my-topic \
  --replay --since "2024-01-15T10:00:00Z" \
  --export output.json
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--broker` | Pulsar broker URL | `pulsar://localhost:6650` |
| `--topic` | Full topic name | *(required)* |
| `--filter` | Key-value filter expression | — |
| `--replay` | Enable replay mode | `false` |
| `--since` | Replay start timestamp (RFC3339) | — |
| `--export` | Export messages to a file | — |

---

## Requirements

- Go 1.21+
- Apache Pulsar 2.x or later

---

## License

This project is licensed under the [MIT License](LICENSE).