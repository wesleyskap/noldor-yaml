# Noldor YAML

Noldor YAML (`go.yaml.in/yaml/v4`) is a high-performance Go library for encoding, decoding, and inspecting YAML documents. Inspired by the master elven craftspeople of *The Silmarillion*, Noldor YAML is designed for algorithmic efficiency, memory-aligned data structures, and compliance with modern software engineering standards.

---

## Overview and Features

Noldor YAML provides a complete suite for working with YAML data structures in Go applications:

- **Module Import Path**: `go.yaml.in/yaml/v4`
- **Anchor & Alias Resolution (`yaml.AnchorRegistry`)**: Native support for defining anchors (`&anchor`) and resolving references (`*anchor`) across document trees.
- **Merge Key Processing (`yaml.MergeKeyResolver`)**: Complete support for YAML `<<: *anchor` dictionary merging semantics.
- **ISO-8601 Timestamp Parser (`yaml.ParseTimestamp`)**: Automatic detection and conversion of date and timestamp scalars into typed `time.Time` fields.
- **Multi-Document Streaming (`yaml.StreamDecoder`)**: Sequential streaming and decoding of multi-document YAML payloads separated by `---`.
- **Block Scalar Formatter (`yaml.BlockScalarFormatter`)**: Full handling of literal (`|`) and folded (`>`) multiline block text styles.
- **Comment Preserver (`yaml.CommentPreserver`)**: Preservation and association of header, line, and footer comments within AST nodes.
- **Explicit Tag Resolver (`yaml.ExplicitTagResolver`)**: Custom type tag handling and URI namespace expansion (`!!binary`, `!custom`).
- **Pretty Printer Formatter (`yaml.PrettyPrinter`)**: Clean indented YAML emission and CLI formatting (`noldor-yaml pretty`).
- **YAML 1.1 Scalar Compatibility**: Seamless decoding of legacy boolean scalars (`yes`, `no`, `on`, `off`) into typed Go `bool` fields.
- **Octal Number Handling**: Support for legacy YAML 1.1 octals (`0755`) and YAML 1.2 notation (`0o755`).
- **AST Node Inspection**: Tree representation (`yaml.Node`) for low-level document traversal and manipulation.
- **CLI Inspector**: Command line utility (`cmd/noldor-yaml`) to parse and display AST node hierarchies.

---

## Architecture and Guidelines

This library strictly adheres to the ecosystem development guidelines:

1. **Memory Alignment**: All Go structs order fields from largest to smallest byte footprint (`pointers` -> `int64`/`uint64` -> `int32`/`float32` -> `bool`) to minimize memory padding and maximize CPU cache efficiency.
2. **Modular Design**: Functions are kept between 4 and 20 lines of code to enforce single responsibility.
3. **Explicit Errors**: Error messages include the offending invalid input value and the expected target shape.
4. **Clean Documentation**: Zero emojis, zero icons, and standard heading titles without CamelCase.

---

## Installation

```bash
go get go.yaml.in/yaml/v4
```

---

## Usage Examples

### Unmarshaling YAML into Structs

```go
package main

import (
	"fmt"
	"log"

	"go.yaml.in/yaml/v4"
)

type ServerConfig struct {
	Name    string `yaml:"name"`
	Port    int    `yaml:"port"`
	Enabled bool   `yaml:"enabled"`
}

func main() {
	input := []byte("name: Noldor App\nport: 0755\nenabled: yes\n")
	var config ServerConfig
	if err := yaml.Unmarshal(input, &config); err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}
	fmt.Printf("Config: %+v\n", config)
}
```

### Marshaling Go Values to YAML

```go
package main

import (
	"fmt"
	"log"

	"go.yaml.in/yaml/v4"
)

func main() {
	payload := map[string]interface{}{
		"service": "noldor-core",
		"port":    8080,
	}
	output, err := yaml.Marshal(payload)
	if err != nil {
		log.Fatalf("Marshal error: %v", err)
	}
	fmt.Println(string(output))
}
```

### Streaming Multi-Document YAML

```go
package main

import (
	"fmt"
	"strings"

	"go.yaml.in/yaml/v4"
)

func main() {
	stream := "doc: 1\n---\ndoc: 2\n"
	dec := yaml.NewStreamDecoder(strings.NewReader(stream))
	for dec.More() {
		var doc map[string]int
		if err := dec.Decode(&doc); err == nil {
			fmt.Printf("Doc: %+v\n", doc)
		}
	}
}
```

---

## Command Line AST Inspector

The package includes a command-line tool `noldor-yaml` located in `cmd/noldor-yaml` for inspecting and formatting YAML syntax trees.

### Running the Inspector

```bash
go run ./cmd/noldor-yaml inspect config.yaml
```

Output:
```text
--- AST Node Structure ---
[Document]
  [Mapping]
    [Scalar] Value: "name"
    [Scalar] Value: "Noldor App"
    [Scalar] Value: "port"
    [Scalar] Value: "0755"
    [Scalar] Value: "enabled"
    [Scalar] Value: "yes"
```

---

## Running Tests

Execute the F.I.R.S.T. unit test suite across all packages:

```bash
go test -v ./...
```

---

## Versioning and Governance

This library adheres to [Semantic Versioning 2.0.0](https://semver.org/). See the following governance files for details:

- [`CHANGELOG.md`](CHANGELOG.md): Release history and backlog tracking.
- [`CONTRIBUTING.md`](CONTRIBUTING.md): Contribution guidelines and code style policies.
- [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md): Community participation standards.

