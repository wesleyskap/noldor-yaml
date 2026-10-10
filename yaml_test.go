package yaml_test

import (
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

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

func TestAnchorRegistryAndAliases(t *testing.T) {
	reg := yaml.NewAnchorRegistry()
	targetNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "target_value"}
	if err := reg.Register("base", targetNode); err != nil {
		t.Fatalf("failed registering anchor: %v", err)
	}
	resolved, err := reg.Resolve("base")
	if err != nil || resolved.Value != "target_value" {
		t.Fatalf("failed resolving anchor: %v", err)
	}

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

func TestMergeKeyResolver(t *testing.T) {
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
	if len(targetMap.Content) != 4 {
		t.Errorf("expected 4 elements in merged map content, got %d", len(targetMap.Content))
	}
}

func TestTimestampParsing(t *testing.T) {
	ts, err := yaml.ParseTimestamp("2025-06-01T12:30:00Z")
	if err != nil {
		t.Fatalf("failed parsing timestamp: %v", err)
	}
	if ts.Year() != 2025 || ts.Month() != time.June || ts.Day() != 1 {
		t.Errorf("unexpected parsed date: %v", ts)
	}
}

func TestStreamDecoder(t *testing.T) {
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

func TestBlockScalarFormatter(t *testing.T) {
	formatter := yaml.NewBlockScalarFormatter(2, "clip")
	res := formatter.FormatLiteral("first line\nsecond line")
	if !strings.Contains(res, "|\n  first line\n  second line") {
		t.Errorf("unexpected formatted literal: %q", res)
	}
}

func TestCommentPreserver(t *testing.T) {
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

func TestExplicitTagResolver(t *testing.T) {
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

func TestSequenceAndInterfaceDecode(t *testing.T) {
	input := []byte("- apple\n- banana\n- cherry\n")
	var fruits []string
	if err := yaml.Unmarshal(input, &fruits); err != nil {
		t.Fatalf("failed decoding slice: %v", err)
	}
	if len(fruits) != 3 || fruits[1] != "banana" {
		t.Errorf("unexpected slice content: %v", fruits)
	}

	var anyVal interface{}
	if err := yaml.Unmarshal(input, &anyVal); err != nil {
		t.Fatalf("failed decoding into interface: %v", err)
	}
	slice, ok := anyVal.([]interface{})
	if !ok || len(slice) != 3 {
		t.Errorf("expected []interface{} of length 3, got: %T (%v)", anyVal, anyVal)
	}
}

func TestMapMarshalAndSequenceMarshal(t *testing.T) {
	m := map[string]string{
		"alpha": "1",
		"beta":  "2",
	}
	out, err := yaml.Marshal(m)
	if err != nil {
		t.Fatalf("failed marshaling map: %v", err)
	}
	if !strings.Contains(string(out), "alpha: 1") {
		t.Errorf("expected map output containing 'alpha: 1', got: %s", string(out))
	}

	items := []string{"val1", "val2"}
	seqOut, err := yaml.Marshal(items)
	if err != nil {
		t.Fatalf("failed marshaling slice: %v", err)
	}
	if !strings.Contains(string(seqOut), "- val1") {
		t.Errorf("expected sequence output containing '- val1', got: %s", string(seqOut))
	}
}

func TestNodeHelpers(t *testing.T) {
	var n yaml.Node
	if !n.IsZero() {
		t.Errorf("expected uninitialized node to be zero")
	}

	docNode := &yaml.Node{Kind: yaml.DocumentNode}
	var nilNode *yaml.Node
	if !nilNode.IsZero() {
		t.Errorf("expected nil node to be zero")
	}

	if err := nilNode.Decode(&docNode); err == nil {
		t.Errorf("expected error decoding nil node, got nil")
	}

	err := n.Encode(map[string]int{"a": 1})
	if err != nil {
		t.Fatalf("failed encoding node: %v", err)
	}
	if n.Kind != yaml.MappingNode {
		t.Errorf("expected MappingNode, got %v", n.Kind)
	}

	var res map[string]int
	if err := n.Decode(&res); err != nil {
		t.Fatalf("failed decoding node: %v", err)
	}
	if res["a"] != 1 {
		t.Errorf("expected res['a'] == 1, got %v", res["a"])
	}
}

func TestScalarAndUtilityHelpers(t *testing.T) {
	if !yaml.IsNullScalar("null") || !yaml.IsNullScalar("~") || !yaml.IsNullScalar("") {
		t.Errorf("expected null scalar detection to be true")
	}
	if yaml.IsNullScalar("value") {
		t.Errorf("expected non-null scalar to be false")
	}

	if !yaml.IsTimestampScalar("2026-10-11") {
		t.Errorf("expected timestamp detection to be true")
	}
	if yaml.IsTimestampScalar("not-a-date") {
		t.Errorf("expected invalid timestamp detection to be false")
	}

	b := yaml.NewBlockScalarFormatter(2, "clip")
	parsed, err := b.ParseLiteralBlock([]string{"  line1", "  line2"})
	if err != nil || parsed != "line1\nline2" {
		t.Errorf("unexpected parsed literal block: %q (err: %v)", parsed, err)
	}
	if _, err := b.ParseLiteralBlock(nil); err == nil {
		t.Errorf("expected error for empty literal block lines")
	}

	tr := yaml.NewExplicitTagResolver()
	tr.RegisterTag("!h", "uri:custom")
	if tr.ResolveTag("!h") != "uri:custom" {
		t.Errorf("unexpected resolved tag: %s", tr.ResolveTag("!h"))
	}
	if tr.ResolveTag("!other") != "!other" {
		t.Errorf("expected fallback tag handle for unknown tags")
	}

	cp := yaml.NewCommentPreserver()
	cp.AttachComment(5, "# comment 5")
	if cp.GetComment(5) != "# comment 5" {
		t.Errorf("unexpected comment text: %s", cp.GetComment(5))
	}

	r := strings.NewReader("key: val\n")
	dec := yaml.NewDecoderStream(r)
	var m map[string]string
	if err := dec.Decode(&m); err != nil || m["key"] != "val" {
		t.Errorf("unexpected stream decode: %+v, err: %v", m, err)
	}

	var buf strings.Builder
	enc := yaml.NewEncoderStream(&buf)
	if err := enc.Encode(m); err != nil {
		t.Fatalf("failed encoder stream: %v", err)
	}

	nodeParsed, err := yaml.ParseBytes([]byte("hello: world"))
	if err != nil || nodeParsed == nil {
		t.Fatalf("failed ParseBytes: %v", err)
	}
}

func BenchmarkUnmarshalStruct(b *testing.B) {
	input := []byte("title: Noldor App\nport: 8080\nenabled: yes\n")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var cfg Config
		if err := yaml.Unmarshal(input, &cfg); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalStruct(b *testing.B) {
	cfg := Config{
		Title:   "Elven Core",
		Port:    9000,
		Enabled: true,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := yaml.Marshal(cfg); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseBool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := yaml.ParseBool("yes"); err != nil {
			b.Fatal(err)
		}
	}
}
