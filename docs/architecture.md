# Architecture

This document describes the structure, design decisions, and non-obvious implementation details of the project. It is kept up to date with the current codebase and is not tied to any specific version.

---

## Layer Map

```
cmd/                          composition roots — entry points only, no business logic
  api/                        REST API (Uber Fx dependency injection)
    modules/                  Fx module definitions (Config, DB, AI, Server)
  daemon/                     voice daemon (manual wiring, no framework)

app/ai/                       application core — framework-free
  entity/                     domain types: Message, Conversation, Memory, Command, Intent
  repository/                 repository interfaces defined here, consumed by use cases
  usecase/                    use cases + input-port interfaces (ports.go)
  service/                    domain services: IntentRouter, ContextBuilder, ActionExecutor
  handler/
    http/                     HTTP adapters — convert HTTP requests to use-case calls
    voice/                    voice adapters — Listener, Transcriber, Speaker + AudioCapture interface

internal/                     infrastructure — implements interfaces defined in app/ai/
  config/                     environment variable loading and validation
  database/                   pgx connection wrapper
  repository/                 PostgreSQL implementations of app/ai/repository interfaces
  httpclient/                 external clients: ClaudeClient, WhisperLocalClient, LocalTTSClient

pkg/
  logger/                     logging utilities
```

Dependencies always point inward. `app/ai/` knows nothing about `internal/` or `cmd/`. Infrastructure implements interfaces; it does not define them.

---

## Domain Entities

### Intent

`Intent` is a string enum representing the category of what the user wants the assistant to do:

| Value | Meaning |
|---|---|
| `chat` | General conversation — no local action required |
| `open_app` | Open an application on the device |
| `set_alarm` | Create an alarm or reminder |
| `save_memory` | Persist a piece of information for future recall |
| `query_memory` | Retrieve previously saved information |
| `unknown` | Could not be determined |

### Command

`Command` is the structured result of parsing raw user text. `Action` is a pointer and is `nil` when the intent is `chat` — the assistant responds with text only, with no local operation to perform. When `Action` is non-nil, its `Payload` carries arbitrary key-value parameters that the `ActionExecutor` service interprets.

### Conversation

A `Conversation` groups all `Message` turns under a `SessionID`. Each message has a `Role` (`user`, `assistant`, or `system`) and `Content`. Conversations are stored as JSONB in PostgreSQL and loaded in full on each request to provide Claude with the full context window.

### Memory

`Memory` is long-term key-value storage independent of any session. Keys are indexed for fast lookup. The `ManageMemoryUseCase` exposes save, search, find-all, and delete operations.

---

## Use Case: ProcessCommand

This is the main request-response flow, shared by both entry points:

1. Validate that input text is non-empty; return `ErrEmptyInput` if not
2. Load the conversation for the given `SessionID`; start a new empty conversation if none exists
3. Append the user message to the conversation
4. Call Claude via `AIClient.Complete` with the full message history
5. Append the assistant response to the conversation
6. Persist the updated conversation — this step is non-fatal: if saving fails, the response is still returned to the caller. Persistence failure is logged but does not degrade the user-facing result
7. Return the response text and the detected intent

---

## Voice Activity Detection (VAD)

The `Listener` reads audio from the default input device via PortAudio at 16 kHz with 1024-frame buffers (16-bit mono samples). On each buffer, it computes the RMS amplitude:

```
RMS = sqrt( sum(sample²) / n )
```

The RMS value is compared against `SilenceThreshold` (configured on a 0–32767 scale). Recording does not stop on silence alone — it first waits for the user to produce speech above the threshold (`hasSpeech = true`), and only then starts counting consecutive silent chunks. When silent chunks reach `silenceChunkThreshold` (derived from `SilenceDurationMs`), the turn ends. `MaxRecordSeconds` is a hard upper limit that fires regardless of speech state.

This two-phase approach (wait for speech, then wait for silence) prevents the daemon from treating initial ambient noise as the end of a turn.

---

## Dependency Injection

### REST API — Uber Fx

The API uses [Uber Fx](https://github.com/uber-go/fx) split into four modules:

| Module | Provides |
|---|---|
| `ConfigModule` | `*config.Config` |
| `DBModule` | `*database.DB`, `ConversationRepository`, `MemoryRepository` |
| `AIModule` | `AIClient`, `CommandProcessor`, `HistoryManager`, `MemoryManager`, HTTP handlers |
| `ServerModule` | `*zap.Logger`, HTTP server with graceful shutdown lifecycle |

Repository and client implementations are registered against their interface types using Fx's `fx.As` annotation, enforcing that no layer depends on a concrete type from another layer.

### Voice Daemon — Manual Wiring

The daemon wires all dependencies manually in `cmd/daemon/main.go`. Its lifecycle is a flat, sequential loop: connect → listen loop → close. There is no benefit to introducing a DI framework for a linear program with no parallel lifecycles or optional modules. Manual wiring is more direct and easier to follow.

---

## Design Decisions

**Interfaces at every boundary.** All cross-layer dependencies go through interfaces defined in the consuming layer or in `usecase/ports.go`. No concrete type from `internal/` is ever referenced in `app/ai/`.

**Local-only AI dependencies.** Transcription uses whisper.cpp running locally (no API key, no network call). TTS uses the OS built-in engine (`say` on macOS, `espeak` on Linux). The only external network call is to the Anthropic Claude API.

**Shared use case layer.** Both the REST API and the voice daemon call the same `ProcessCommandUseCase`. The entry-point protocols (HTTP vs. microphone) are adapters; the business logic is not duplicated.

**System prompt from file.** The Claude system prompt is read from `resources/system_prompt.txt` and injected per-request. This allows the prompt to be modified without recompiling the binary.
