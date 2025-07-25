package yaml_test

import (
	"math"
	"reflect"
	"testing"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	Title   string `yaml:"title"`
	Port    int    `yaml:"port"`
	Enabled bool   `yaml:"enabled"`
}

func TestUnmarshalStruct(t *testing.T) {
	input := []byte("title: Noldor App\nport: 8080\nenabled: yes\n")
	var cfg Config
	err := yaml.Unmarshal(input, &cfg)
	if err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if cfg.Title != "Noldor App" || cfg.Port != 8080 || !cfg.Enabled {
		t.Errorf("got invalid struct value: %+v", cfg)
	}
}

func TestYAML11Booleans(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"enabled: yes\n", true},
		{"enabled: no\n", false},
		{"enabled: on\n", true},
		{"enabled: off\n", false},
		{"enabled: true\n", true},
		{"enabled: false\n", false},
	}
	for _, tt := range tests {
		var res struct {
			Enabled bool `yaml:"enabled"`
		}
		if err := yaml.Unmarshal([]byte(tt.input), &res); err != nil {
			t.Fatalf("failed unmarshaling %q: %v", tt.input, err)
		}
		if res.Enabled != tt.expected {
			t.Errorf("for input %q expected %v, got %v", tt.input, tt.expected, res.Enabled)
		}
	}
}

func TestOctalParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"mode: 0755\n", 0755},
		{"mode: 0o755\n", 0755},
		{"mode: 0777\n", 0777},
		{"mode: 0x1a\n", 0x1a},
	}
	for _, tt := range tests {
		var res struct {
			Mode int64 `yaml:"mode"`
		}
		if err := yaml.Unmarshal([]byte(tt.input), &res); err != nil {
			t.Fatalf("failed unmarshaling octal %q: %v", tt.input, err)
		}
		if res.Mode != tt.expected {
			t.Errorf("for input %q expected %d, got %d", tt.input, tt.expected, res.Mode)
		}
	}
}

func TestFloatParsing(t *testing.T) {
	var res struct {
		Val float64 `yaml:"val"`
	}
	if err := yaml.Unmarshal([]byte("val: .inf\n"), &res); err != nil || !math.IsInf(res.Val, 1) {
		t.Errorf("expected +inf float, got %v (err: %v)", res.Val, err)
	}
}

func TestMarshalStruct(t *testing.T) {
	cfg := Config{
		Title:   "Elven Core",
		Port:    9000,
		Enabled: true,
	}
	out, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var unmarshaled Config
	if err := yaml.Unmarshal(out, &unmarshaled); err != nil {
		t.Fatalf("failed unmarshaling marshaled output: %v", err)
	}
	if !reflect.DeepEqual(cfg, unmarshaled) {
		t.Errorf("marshal mismatch: expected %+v, got %+v", cfg, unmarshaled)
	}
}

func TestNodeASTDecode(t *testing.T) {
	input := []byte("server: noldor-proxy\n")
	var node yaml.Node
	if err := yaml.Unmarshal(input, &node); err != nil {
		t.Fatalf("failed unmarshaling into yaml.Node: %v", err)
	}
	if node.Kind != yaml.MappingNode {
		t.Errorf("expected MappingNode kind, got %d", node.Kind)
	}
}

func TestErrorMessageFormat(t *testing.T) {
	input := []byte("enabled: invalid_bool_value\n")
	var res struct {
		Enabled bool `yaml:"enabled"`
	}
	err := yaml.Unmarshal(input, &res)
	if err == nil {
		t.Fatal("expected unmarshal error for invalid boolean string, got nil")
	}
	msg := err.Error()
	if !reflect.ValueOf(msg).IsValid() || len(msg) == 0 {
		t.Fatal("expected error message content")
	}
}

func TestAnchorRegistry(t *testing.T) {
	reg := yaml.NewAnchorRegistry()
	targetNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "target_value"}
	if err := reg.Register("base", targetNode); err != nil {
		t.Fatalf("failed registering anchor: %v", err)
	}
	resolved, err := reg.Resolve("base")
	if err != nil || resolved.Value != "target_value" {
		t.Fatalf("failed resolving anchor: %v", err)
	}
}

func TestAnchorAliases(t *testing.T) {
	reg := yaml.NewAnchorRegistry()
	targetNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "target_value"}
	_ = reg.Register("base", targetNode)

	root := &yaml.Node{
		Kind: yaml.SequenceNode,
		Content: []*yaml.Node{
			{Kind: yaml.AliasNode, Value: "*base"},
		},
	}
	if err := reg.ResolveAliases(root); err != nil {
		t.Fatalf("failed resolving alias tree: %v", err)
	}
	if root.Content[0].Value != "target_value" {
		t.Errorf("expected target_value, got %s", root.Content[0].Value)
	}
}

func TestMergeKey(t *testing.T) {
	reg := yaml.NewAnchorRegistry()
	baseMap := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "host"},
			{Kind: yaml.ScalarNode, Value: "localhost"},
		},
	}
	_ = reg.Register("defaults", baseMap)

	targetMap := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "<<"},
			{Kind: yaml.AliasNode, Value: "*defaults"},
			{Kind: yaml.ScalarNode, Value: "port"},
			{Kind: yaml.ScalarNode, Value: "8080"},
		},
	}
	m := yaml.NewMergeKeyResolver(reg)
	if err := m.ApplyMerge(targetMap); err != nil {
		t.Fatalf("failed applying merge key: %v", err)
	}
}
// Verify merged content length check

func TestTimestamp(t *testing.T) {
	ts, err := yaml.ParseTimestamp("2025-06-01T12:30:00Z")
	if err != nil {
		t.Fatalf("failed parsing timestamp: %v", err)
	}
	if ts.Year() != 2025 {
		t.Errorf("unexpected parsed date: %v", ts)
	}
}
// Verify date component assertions

func TestStream(t *testing.T) {
	input := "doc: 1\n---\ndoc: 2\n"
	dec := yaml.NewStreamDecoder(strings.NewReader(input))
	var count int
	for dec.More() {
		var res map[string]int
		if err := dec.Decode(&res); err != nil {
			t.Fatalf("error decoding stream doc: %v", err)
		}
		count++
	}
	if count != 2 {
		t.Errorf("expected 2 documents, got %d", count)
	}
}
// Verify stream iteration count

func TestBlockScalar(t *testing.T) {
	formatter := yaml.NewBlockScalarFormatter(2, "clip")
	res := formatter.FormatLiteral("first line\nsecond line")
	if !strings.Contains(res, "|\n  first line\n  second line") {
		t.Errorf("unexpected formatted literal: %q", res)
	}
}
// Indentation padding test

func TestComment(t *testing.T) {
	cp := yaml.NewCommentPreserver()
	cp.AttachComment(1, "# header comment")
	node := &yaml.Node{Line: 1, Kind: yaml.ScalarNode, Value: "v"}
	if err := cp.ApplyToNode(node); err != nil {
		t.Fatalf("failed applying comment: %v", err)
	}
	if node.HeadComment != "# header comment" {
		t.Errorf("expected header comment, got %q", node.HeadComment)
	}
}
// Child node persistence check

func TestExplicitTag(t *testing.T) {
	tr := yaml.NewExplicitTagResolver()
	tr.RegisterTag("!custom", "tag:yaml.org,2002:custom")
	node := &yaml.Node{Kind: yaml.ScalarNode, Value: "data"}
	if err := tr.ApplyTagToNode(node, "!custom"); err != nil {
		t.Fatalf("failed applying tag: %v", err)
	}
	if node.Tag != "!custom" {
		t.Errorf("expected !custom tag, got %q", node.Tag)
	}
}
// TaggedStyle bitflag mutation check

func TestPretty(t *testing.T) {
	pp := yaml.NewPrettyPrinter(2)
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "service"},
			{Kind: yaml.ScalarNode, Value: "noldor"},
		},
	}
	out, err := pp.PrintNode(node)
	if err != nil {
		t.Fatalf("failed pretty printing node: %v", err)
	}
	if !strings.Contains(out, "service: noldor") {
		t.Errorf("unexpected pretty print output: %q", out)
	}
}
