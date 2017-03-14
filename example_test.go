package yaml_test

import (
	"fmt"
	"log"

	"go.yaml.in/yaml/v4"
)

func ExampleUnmarshal() {
	input := []byte("name: Noldor\nactive: yes\nport: 0755\n")
	var spec struct {
		Name   string `yaml:"name"`
		Port   int    `yaml:"port"`
		Active bool   `yaml:"active"`
	}
	if err := yaml.Unmarshal(input, &spec); err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}
	fmt.Printf("Name: %s, Port: %d, Active: %t\n", spec.Name, spec.Port, spec.Active)
	// Output:
	// Name: Noldor, Port: 493, Active: true
}

func ExampleMarshal() {
	spec := struct {
		Service string `yaml:"service"`
		Port    int    `yaml:"port"`
	}{
		Service: "Noldor YAML",
		Port:    8080,
	}
	out, err := yaml.Marshal(spec)
	if err != nil {
		log.Fatalf("Marshal error: %v", err)
	}
	fmt.Print(string(out))
	// Output:
	// service: Noldor YAML
	// port: 8080
}
