# ai-assistant

A personal assistant powered by Claude (Anthropic). Exposes two independent entry points: a REST API for text-based commands and a voice daemon that listens, transcribes, reasons, and speaks — all locally except for the Claude call.

---

## Entry Points

### REST API (`cmd/api`)

Accepts text commands over HTTP and returns Claude-powered responses. Any client (web, mobile, curl) can use it. No voice hardware required.

```
POST /ai/command
Content-Type: application/json

{ "text": "What is the weather today?", "session_id": "user-1" }
```

```
GET /ai/history?session_id=user-1
GET /ai/history?limit=5
GET /health
```

### Voice Daemon (`cmd/daemon`)

Runs a continuous loop: captures microphone audio → transcribes locally via whisper.cpp → processes with Claude → speaks the response using the OS TTS engine. Requires PortAudio and a running whisper.cpp server.

---

## Prerequisites

| Dependency | Required by | Notes |
|---|---|---|
| Go 1.25+ | both | |
| PostgreSQL 14+ | both | |
| Anthropic API key | both | Claude model access |
| whisper.cpp server | daemon only | local transcription |
| PortAudio | daemon + `go build ./...` | `brew install portaudio` / `apt install portaudio19-dev` |
| `say` / `espeak` | daemon only | built-in on macOS; `apt install espeak` on Linux |

---

## Setup

```bash
# 1. Clone and install dependencies
git clone https://github.com/KevenAbraham/ai-assistant
cd ai-assistant
go mod tidy

# 2. Configure environment
cp .env.example .env
# Edit .env with real values

# 3. Start PostgreSQL
make docker-up

# 4. Run migrations
make migrate-up

# 5. Run the REST API
make run-api

# 6. (Optional) Run the voice daemon
make run-daemon
```

---

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `DATABASE_URL` | yes | — | PostgreSQL connection string |
| `ANTHROPIC_API_KEY` | yes | — | Anthropic API key |
| `CLAUDE_MODEL` | no | `claude-haiku-4-5-20251001` | Claude model ID |
| `HTTP_ADDR` | no | `:3000` | HTTP server listen address |
| `SYSTEM_PROMPT_PATH` | no | `resources/system_prompt.txt` | Path to the Claude system prompt file |
| `WHISPER_URL` | no | `http://localhost:9000` | whisper.cpp server URL (daemon only) |
| `RECORD_SECONDS` | no | `30` | Hard cap on recording duration in seconds (daemon only) |
| `SILENCE_THRESHOLD` | no | `500` | RMS amplitude below which audio counts as silence, on a 0–32767 scale (daemon only) |
| `SILENCE_DURATION_MS` | no | `1000` | Consecutive silence in milliseconds that ends a recording turn (daemon only) |
| `LOG_LEVEL` | no | `info` | Log level (`debug`, `info`, `warn`, `error`) |

---

## Voice Daemon Setup (whisper.cpp)

```bash
# macOS
brew install whisper-cpp
whisper-download-ggml-model base.en
whisper-server -m ~/.cache/whisper/ggml-base.en.bin --port 9000
```

```bash
# Linux
git clone https://github.com/ggerganov/whisper.cpp
cd whisper.cpp && make
bash models/download-ggml-model.sh base.en
./server -m models/ggml-base.en.bin --port 9000
```

Set `WHISPER_URL=http://localhost:9000` in `.env`.

---

## Makefile Targets

| Target | Description |
|---|---|
| `make build` | Compile API and daemon binaries |
| `make run-api` | Run REST API with `.env` loaded |
| `make run-daemon` | Run voice daemon with `.env` loaded |
| `make migrate-up` | Apply SQL migrations |
| `make migrate-down` | Roll back SQL migrations |
| `make docker-up` | Start PostgreSQL container |
| `make docker-down` | Stop PostgreSQL container |
| `make test` | Run tests |
| `make vet` | Run `go vet` |
| `make tidy` | Tidy `go.mod` |

---

## Documentation

- [Architecture](docs/architecture.md) — layer map, design decisions, voice activity detection, dependency injection
- [Changelog](docs/CHANGELOG.md) — release history
