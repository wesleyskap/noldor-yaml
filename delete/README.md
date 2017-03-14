# Noldor YAML Library and Inspector

Noldor YAML (`go.yaml.in/yaml/v4`) is a high-performance Go library for encoding, decoding, and inspecting YAML documents. Named after the master elven artisans of *The Silmarillion*, the package is engineered for precision, memory alignment, and strict adherence to software engineering standards.

---

## Technical Features

- **Module Import Path**: `go.yaml.in/yaml/v4`.
- **YAML 1.1 Compatibility**: Full support for legacy boolean scalars (`yes`, `no`, `on`, `off`) when decoding into typed Go `bool` fields.
- **YAML 1.2 and Legacy Octals**: Support for traditional octals (`0755`) and modern octal notation (`0o755`).
- **AST Node Representation**: Complete access to document AST structure via `yaml.Node` for tree manipulation.
- **CLI Inspector Tool**: Built-in `noldor-yaml` command-line utility for inspecting parse trees.

---

## Installation

```bash
go get go.yaml.in/yaml/v4
```

---

## Usage Examples

### Decoding YAML Data

```go
package main

import (
	"fmt"
	"log"

	"go.yaml.in/yaml/v4"
)

type Configuration struct {
	Name    string `yaml:"name"`
	Port    int    `yaml:"port"`
	Enabled bool   `yaml:"enabled"`
}

func main() {
	input := []byte("name: Noldor Core\nport: 0755\nenabled: yes\n")
	var config Configuration
	if err := yaml.Unmarshal(input, &config); err != nil {
		log.Fatalf("Failed unmarshaling: %v", err)
	}
	fmt.Printf("Config: %+v\n", config)
}
```

### Encoding Go Values

```go
package main

import (
	"fmt"
	"log"

	"go.yaml.in/yaml/v4"
)

func main() {
	cfg := map[string]interface{}{
		"service": "noldor-proxy",
		"port":    8080,
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalf("Failed marshaling: %v", err)
	}
	fmt.Println(string(data))
}
```

---

## Command Line AST Inspector

To inspect how a YAML file is parsed into AST node hierarchies:

```bash
go build -o bin/noldor-yaml ./cmd/noldor-yaml
./bin/noldor-yaml inspect config.yaml
```

Output:
```text
--- AST Node Structure ---
[Document]
  [Mapping]
    [Scalar] Value: "port"
    [Scalar] Value: "0755"
    [Scalar] Value: "enabled"
    [Scalar] Value: "yes"
    [Scalar] Value: "name"
    [Scalar] Value: "Noldor Core"
```

---

## Testing

Run the full unit test suite:

```bash
go test -v ./...
```
