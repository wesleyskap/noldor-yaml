package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"go.yaml.in/yaml/v4"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	cmd := os.Args[1]
	switch cmd {
	case "inspect", "-inspect":
		runInspect(os.Args[2:])
	case "pretty", "-pretty":
		runPretty(os.Args[2:])
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Noldor YAML Inspector tool")
	fmt.Println("Usage: noldor-yaml inspect <file.yaml>")
	fmt.Println("       noldor-yaml pretty <file.yaml>")
	fmt.Println("       cat file.yaml | noldor-yaml inspect -")
}

func runInspect(args []string) {
	if len(args) == 0 {
		fmt.Println("error: missing input file or '-' for stdin")
		os.Exit(1)
	}
	r, err := openInput(args[0])
	if err != nil {
		fmt.Printf("error opening input: %v\n", err)
		os.Exit(1)
	}
	defer r.Close()
	inspectReader(r)
}

func runPretty(args []string) {
	if len(args) == 0 {
		fmt.Println("error: missing input file or '-' for stdin")
		os.Exit(1)
	}
	r, err := openInput(args[0])
	if err != nil {
		fmt.Printf("error opening input: %v\n", err)
		os.Exit(1)
	}
	defer r.Close()
	p := yaml.NewParser(r)
	tree, err := p.ParseTree()
	if err != nil {
		fmt.Printf("YAML parsing error: %v\n", err)
		os.Exit(1)
	}
	pp := yaml.NewPrettyPrinter(2)
	out, err := pp.PrintNode(tree)
	if err != nil {
		fmt.Printf("Pretty printing error: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(out)
}

func openInput(target string) (io.ReadCloser, error) {
	if target == "-" {
		return io.NopCloser(os.Stdin), nil
	}
	return os.Open(target)
}

func inspectReader(r io.Reader) {
	p := yaml.NewParser(r)
	tree, err := p.ParseTree()
	if err != nil {
		fmt.Printf("YAML parsing error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("--- AST Node Structure ---")
	printNodeTree(tree, 0)
}

func printNodeTree(n *yaml.Node, depth int) {
	if n == nil {
		return
	}
	indent := strings.Repeat("  ", depth)
	kindStr := getKindName(n.Kind)
	val := n.Value
	if val != "" {
		fmt.Printf("%s[%s] Value: %q\n", indent, kindStr, val)
	} else {
		fmt.Printf("%s[%s]\n", indent, kindStr)
	}
	for _, child := range n.Content {
		printNodeTree(child, depth+1)
	}
}

func getKindName(k yaml.Kind) string {
	switch k {
	case yaml.DocumentNode:
		return "Document"
	case yaml.MappingNode:
		return "Mapping"
	case yaml.SequenceNode:
		return "Sequence"
	case yaml.ScalarNode:
		return "Scalar"
	case yaml.AliasNode:
		return "Alias"
	default:
		return "Unknown"
	}
}

