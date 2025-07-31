# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [4.8.0] - 2025-07-31

### Added
- **Pretty Printer & Formatter Engine**: Formatted AST tree emission into indented YAML in `pretty.go`.
- **CLI Pretty Command**: Added `noldor-yaml pretty` subcommand in `cmd/noldor-yaml/main.go`.

## [4.7.0] - 2025-07-10

### Added
- **Explicit Tag Handling**: Support for explicit custom and standard YAML type tags (`!!binary`, `!custom`) in `tags.go`.

## [4.6.0] - 2025-06-20

### Added
- **AST Comment Preserver**: Preservation of header, inline, and footer comments during parsing and AST traversal in `comment.go`.

## [4.5.0] - 2025-05-30

### Added
- **Block Scalar Formatting**: Support for Literal (`|`) and Folded (`>`) multiline block text styles in `block_scalar.go`.

## [4.4.0] - 2025-05-10

### Added
- **Multi-Document Stream Decoder**: Sequential decoding for multi-document streams separated by `---` in `stream.go`.

## [4.3.0] - 2025-04-18

### Added
- **ISO-8601 Timestamp Parser**: Native parsing for timestamps and dates into Go `time.Time` values in `timestamp.go`.

## [4.2.0] - 2025-03-26

### Added
- **Merge Key Resolver**: Support for YAML `<<: *anchor` dictionary merging semantics in `merge.go`.

## [4.1.0] - 2025-03-05

### Added
- **Anchor and Alias Resolution**: Native anchor definition (`&anchor`) and alias dereferencing (`*anchor`) in `anchor.go`.

## [1.0.0] - 2025-02-10

### Added
- **Dynamic Reverse Proxy & Router Engine**: Full support for Host, Path, and PathPrefix request matching.
- **Round Robin Load Balancer**: Atomic target selection across healthy backend pools.
- **Resilience Middlewares**: Consecutive failure Circuit Breaker (CLOSED, OPEN, HALF-OPEN states), token-bucket Rate Limiter, and exponential backoff Retries.
- **Multi-Protocol Support**: Full proxying for HTTP/1.1, HTTP/2 (h2c), WebSocket hijacking, and gRPC trailers.
- **Observability Exporters**: Prometheus metrics exposition endpoint (/metrics) and structured JSON access logging.
- **Embedded Web Dashboard**: Sleek dark glassmorphism management UI built into the single binary.
- **Docker Auto-Discovery**: Automatic label discovery (`palantir.enable=true`) from local Docker socket.

