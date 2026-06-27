# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2026-06-27

First major release. The import path is now `github.com/goloop/log/v2`.

### Migration from v1

- Update imports to `github.com/goloop/log/v2` (and `.../v2/level`,
  `.../v2/layout`). The package name stays `log`.
- Requires Go 1.24 or newer.
- `Panicln` no longer returns `(int, error)` — drop any use of its result.
- `Fxxx(w, ...)` now writes to the configured outputs **and** additionally to
  `w`. Previously it also leaked an extra output into the logger and raced
  under concurrency; if you depended on that behaviour, revisit the call.

### Added

- `log/slog` integration: `Logger.Handler()` returns a `slog.Handler`, and
  `NewSlog(prefixes ...string)` returns a ready `*slog.Logger` backed by the
  logger's outputs. slog levels map onto the logger levels; record attributes
  become typed JSON fields in JSON outputs and `key=value` pairs in text
  outputs.
- `Logger.Enabled(level.Level)` (and the package-level `Enabled`) to guard the
  preparation of expensive arguments for a level no output is interested in.
- `Logger.SetErrorHandler` (and the package-level `SetErrorHandler`) to observe
  output write errors, which were previously discarded.
- `SetDefault(*Logger)` to swap the package-level default logger (paired with
  the existing `Log`).
- Buffer pooling on the hot path, plus benchmarks for the disabled-level and
  no-stack-frame fast paths.
- godoc examples, regression tests, and fuzz tests for `cutFilePath` and `New`.

### Changed

- **BREAKING:** module path is now `github.com/goloop/log/v2`.
- **BREAKING:** minimum Go version raised to 1.24.
- **BREAKING:** `Panicln(a ...any)` no longer returns `(int, error)`.
- `Fxxx(w, ...)` writes to the configured outputs and additionally to `w`
  without mutating the logger; the message body is now rendered once and
  reused across all outputs.
- A single timestamp and at most one stack frame are computed per call, and
  the stack frame is captured only when an output actually needs it.
- The call site for file/function/line layouts is located by walking the stack
  to the first frame outside the package, so it reports the caller regardless
  of the logger's internal call depth. The default `SetSkipStackFrames` value
  is now 0 (extra user wrapper frames only).
- Outputs are rendered and written outside the logger lock; `EditOutputs` and
  `DeleteOutputs` replace the outputs map copy-on-write, so a slow or
  re-entrant writer can no longer block configuration or deadlock the logger.

### Removed

- Dead code: the commented-out `getWriterID` helper and the unused `ioCopy`
  function are gone from the package surface.

### Fixed

- Data race (and possible `concurrent map writes` panic) when several
  goroutines called the `Fxxx` methods concurrently.
- `Fxxx(w, ...)` leaked a `"*"` output into the logger and wrote to every
  configured output instead of the intended targets.
- `go vet` reported 28 format-string diagnostics that stopped `go test` from
  building on a clean toolchain.
- Inconsistent timestamps between outputs within a single logging call.
- Possible panic in stack-frame capture for an oversized skip;
  `SetSkipStackFrames` no longer relies on `recover`.
- A JSON marshalling error no longer silently drops the message — it falls
  back to a plain text line.
- File/function/line layouts reported the logger's own method (or runtime
  internals through the package-level functions) instead of the caller's
  source location; `SetSkipStackFrames` also clamped any value to 4 and so
  could not correct it.
- Write errors were silently discarded; they can now be observed via
  `SetErrorHandler`.
- The logger lock was held for the duration of each output write, so a slow or
  re-entrant writer could block configuration or deadlock.
- `io.Discard` and other zero-size writers were rejected as a "nil writer"
  (typed-nil writers are still rejected).
- slog attributes are now emitted as typed JSON fields in JSON outputs instead
  of being flattened into the message string.

### Performance

- A typical `Info` call costs ~6 allocations (down from ~14); a disabled level
  costs 1. Buffers and the call-site program-counter scratch are pooled, the
  per-`Fxxx` map copy is gone, and panic message bodies are formatted once.

### Documentation

- README and package docs updated for the v2 import path; corrected the JSON
  key names (`filePath`, `lineNumber`, `funcName`) and replaced the
  "minimal allocations" claim with the actual buffer-pooled behaviour.

[2.0.0]: https://github.com/goloop/log/compare/v1.4.3...v2.0.0