# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [4.0.0] - 2017-11-20

### Added
- **Module Path**: Established canonical module import path as `go.yaml.in/yaml/v4`.
- **YAML 1.1 Compatibility**: Added scalar resolution for legacy booleans (`yes`, `no`, `on`, `off`) targeting typed bool fields.
- **Octal Representation**: Support for YAML 1.1 octals (`0755`) and YAML 1.2 octals (`0o755`).
- **AST Node API**: Exposed `yaml.Node` tree hierarchy for manual traversal and decoding.
- **CLI Inspector**: Created `cmd/noldor-yaml` CLI tool to display parse tree structures.
- **Unit Test Suite**: Created F.I.R.S.T. unit test suite and executable Go examples.
