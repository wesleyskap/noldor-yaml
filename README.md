# Noldor YAML

Noldor YAML (`go.yaml.in/yaml/v4`) is a high-performance Go library for encoding, decoding, and inspecting YAML documents.

---

## Overview

- **Module Import Path**: `go.yaml.in/yaml/v4`
- **YAML 1.1 Scalar Compatibility**: Legacy boolean and octal resolution.
- **AST Node Inspection**: Tree representation (`yaml.Node`).
- **CLI Inspector**: Command line utility (`cmd/noldor-yaml`).
- **Anchor & Alias Resolution (`yaml.AnchorRegistry`)**: Native support for defining anchors and resolving aliases.
- **Merge Key Processing (`yaml.MergeKeyResolver`)**: Complete support for YAML <<: *anchor dictionary merging.
- **ISO-8601 Timestamp Parser (`yaml.ParseTimestamp`)**: Automatic detection and conversion of date and timestamp scalars into time.Time.
- **Multi-Document Streaming (`yaml.StreamDecoder`)**: Sequential streaming and decoding of multi-document YAML payloads separated by `---`.
- **Block Scalar Formatter (`yaml.BlockScalarFormatter`)**: Full handling of literal (|) and folded (>) multiline block text styles.
