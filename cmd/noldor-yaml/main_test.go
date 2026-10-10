package main

import (
	"strings"
	"testing"

	"go.yaml.in/yaml/v4"
)

func TestGetKindName(t *testing.T) {
	tests := []struct {
		kind     yaml.Kind
		expected string
	}{
		{yaml.DocumentNode, "Document"},
		{yaml.MappingNode, "Mapping"},
		{yaml.SequenceNode, "Sequence"},
		{yaml.ScalarNode, "Scalar"},
		{yaml.AliasNode, "Alias"},
		{yaml.Kind(999), "Unknown"},
	}
	for _, tt := range tests {
		got := getKindName(tt.kind)
		if got != tt.expected {
			t.Errorf("getKindName(%v) = %q, expected %q", tt.kind, got, tt.expected)
		}
	}
}

func TestPrintNodeTree(t *testing.T) {
	node := &yaml.Node{
		Kind:  yaml.MappingNode,
		Value: "",
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "key"},
			{Kind: yaml.ScalarNode, Value: "value"},
		},
	}
	// Verify nil safety
	printNodeTree(nil, 0)

	// Verify tree execution
	printNodeTree(node, 0)
}

func TestInspectReader(t *testing.T) {
	input := strings.NewReader("name: noldor\n")
	inspectReader(input)
}

func TestPrintPrettyNode(t *testing.T) {
	node := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: "test",
	}
	printPrettyNode(node)
}

func TestOpenInput(t *testing.T) {
	r, err := openInput("-")
	if err != nil {
		t.Fatalf("unexpected error for stdin: %v", err)
	}
	_ = r.Close()

	// Test non-existent file
	_, err = openInput("non_existent_file_noldor.yaml")
	if err == nil {
		t.Fatalf("expected error for non-existent file, got nil")
	}
}

func TestPrintUsage(t *testing.T) {
	printUsage()
}
